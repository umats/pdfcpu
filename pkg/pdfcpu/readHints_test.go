package pdfcpu

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/umats/pdfcpu/pkg/filter"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func hintAuditTestContext(t *testing.T) (*model.Context, *hintReadState) {
	t.Helper()
	ctx, err := model.NewContext(bytes.NewReader(nil), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	ctx.Table[10] = model.NewXRefTableEntryGen0(types.Dict{"Linearized": types.Integer(1)})
	sd := types.StreamDict{Dict: types.Dict{"Length": types.Integer(142), "S": types.Integer(36)}}
	ctx.Table[11] = model.NewXRefTableEntryGen0(sd)
	return ctx, &hintReadState{linearNr: 10, streams: map[int]types.StreamDict{11: sd}}
}

func TestHintReferenceCoverage(t *testing.T) {
	for _, origin := range []string{"dict", "array", "stream", "single-ref", "trailer", "trailer-nested", "root", "info", "encrypt", "linearization-ref", "lazy-compressed", "lazy-container", "unparsed", "ambiguous-linearization"} {
		t.Run(origin, func(t *testing.T) {
			ctx, state := hintAuditTestContext(t)
			ref := *types.NewIndirectRef(11, 0)
			var o types.Object
			switch origin {
			case "dict":
				o = types.Dict{"Target": ref}
			case "array":
				o = types.Array{types.Dict{"Target": ref}}
			case "stream":
				o = types.StreamDict{Dict: types.Dict{"Target": ref}}
			case "single-ref":
				o = ref
			case "trailer":
				ctx.Read.TrailerDicts = []types.Dict{{"Uncommon": ref}}
			case "trailer-nested":
				ctx.Read.TrailerDicts = []types.Dict{{"Uncommon": types.Array{types.Dict{"Target": ref}}}}
			case "root":
				ctx.Root = &ref
			case "info":
				ctx.Info = &ref
			case "encrypt":
				ctx.Encrypt = &ref
			case "linearization-ref":
				o = *types.NewIndirectRef(10, 0)
			case "lazy-compressed", "lazy-container":
				// Real existing object-stream prolog/lazy parser, not a refcount
				// stub. The reference is invisible until lazy materialization.
				content := []byte("12 0 << /Target 11 0 R >>")
				length := int64(len(content))
				osd := types.ObjectStreamDict{StreamDict: types.NewStreamDict(types.Dict{}, 0, &length, nil, nil), ObjCount: 1, FirstObjOffset: 5}
				osd.Raw = content
				if err := parseObjectStream(t.Context(), &osd, ctx.Configuration.Limits); err != nil {
					t.Fatal(err)
				}
				if origin == "lazy-container" {
					o = osd
				} else {
					o = osd.ObjArray[0]
				}
			case "ambiguous-linearization":
				o = types.Dict{"Linearized": types.Integer(1)}
			case "unparsed":
				ctx.Table[12] = model.NewXRefTableEntryGen0(nil)
			}
			if o != nil {
				ctx.Table[12] = model.NewXRefTableEntryGen0(o)
			}
			got, err := hintReferences(t.Context(), ctx, state)
			if err != nil || !got {
				t.Fatalf("reference audit = %t, %v", got, err)
			}
		})
	}
	t.Run("unreferenced", func(t *testing.T) {
		ctx, state := hintAuditTestContext(t)
		got, err := hintReferences(t.Context(), ctx, state)
		if err != nil || got {
			t.Fatalf("audit = %t, %v", got, err)
		}
	})
}

func TestHintDictionaryShapeGuards(t *testing.T) {
	ctx, _ := hintAuditTestContext(t)
	base := types.Dict{"Length": types.Integer(142), "S": types.Integer(36), "Filter": types.Name("FlateDecode")}
	if !supportedHintDict(ctx, base, true) {
		t.Fatal("basic hint shape rejected")
	}
	for _, tc := range []struct {
		name, key string
		value     types.Object
	}{
		{"metadata", "Type", types.Name("Metadata")},
		{"image", "Subtype", types.Name("Image")},
		{"unknown-key", "Unknown", types.Integer(1)},
		{"negative-S", "S", types.Integer(-1)},
		{"indirect-Length", "Length", *types.NewIndirectRef(1, 0)},
		{"indirect-DecodeParms", "DecodeParms", types.Dict{"Target": *types.NewIndirectRef(1, 0)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := base.Clone().(types.Dict)
			d[tc.key] = tc.value
			if supportedHintDict(ctx, d, true) {
				t.Fatal("ordinary/ambiguous shape classified as hint")
			}
		})
	}
}

func TestHintReferenceAuditLimits(t *testing.T) {
	ctx, state := hintAuditTestContext(t)
	ctx.Configuration.Limits.MaxRecursionDepth = 2
	ctx.Table[12] = model.NewXRefTableEntryGen0(types.Array{types.Array{types.Array{types.Integer(1)}}})
	if _, err := hintReferences(t.Context(), ctx, state); !errors.Is(err, model.ErrMaxRecursionDepthExceeded) {
		t.Fatalf("depth limit: %v", err)
	}
	c, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := hintReferences(c, ctx, state); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel: %v", err)
	}
}

func TestHintProbePreservesOriginalAndEncryptedBytes(t *testing.T) {
	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	if _, err := zw.Write(bytes.Repeat([]byte("hint-table-bytes"), 20)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	original := compressed.Bytes()
	for _, encrypted := range []bool{false, true} {
		t.Run(map[bool]string{false: "plaintext", true: "encrypted"}[encrypted], func(t *testing.T) {
			ctx, state := hintAuditTestContext(t)
			ctx.EncKey = bytes.Repeat([]byte{1}, 16)
			ctx.AES4Streams = true
			ctx.E = &model.Enc{R: 4, Emd: true}
			raw := bytes.Clone(original)
			if encrypted {
				var err error
				raw, err = encryptStream(raw, 11, 0, ctx.EncKey, true, 4)
				if err != nil {
					t.Fatal(err)
				}
			}
			before := bytes.Clone(raw)
			length := int64(len(raw))
			sd := types.NewStreamDict(types.Dict{"S": types.Integer(0), "Length": types.Integer(length)}, 0, &length, nil, []types.PDFFilter{{Name: "FlateDecode"}})
			sd.Raw = raw
			ctx.Table[11].Object = sd
			state.streams[11] = sd
			if err := finalizeHintStreamContents(t.Context(), ctx, state, true); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(raw, before) {
				t.Fatal("trial decryption mutated original bytes")
			}
			got := ctx.Table[11].Object.(types.StreamDict)
			if !bytes.Equal(got.Raw, original) {
				t.Fatal("encoded bytes not preserved")
			}
			if ctx.Read.BinaryTotalSize != int64(len(original)) {
				t.Fatal("binary size counted incorrectly")
			}
		})
	}
}

func TestHintUnreadableCompressedOriginDisqualifiesOnlyRepair(t *testing.T) {
	for _, mode := range []int{model.ValidationStrict, model.ValidationRelaxed} {
		for _, encrypted := range []bool{false, true} {
			ctx, state := hintAuditTestContext(t)
			ctx.XRefTable.ValidationMode = mode
			ctx.EncKey = bytes.Repeat([]byte{1}, 16)
			ctx.AES4Streams = true
			ctx.E = &model.Enc{R: 4, Emd: true}
			var compressed bytes.Buffer
			zw := zlib.NewWriter(&compressed)
			if _, err := zw.Write([]byte("hint-table-bytes")); err != nil {
				t.Fatal(err)
			}
			if err := zw.Close(); err != nil {
				t.Fatal(err)
			}
			raw := compressed.Bytes()
			if encrypted {
				var err error
				raw, err = encryptStream(raw, 11, 0, ctx.EncKey, true, 4)
				if err != nil {
					t.Fatal(err)
				}
			}
			length := int64(len(raw))
			sd := types.NewStreamDict(types.Dict{"Length": types.Integer(length), "S": types.Integer(0)}, 0, &length, nil, []types.PDFFilter{{Name: "FlateDecode"}})
			sd.Raw = raw
			ctx.Table[11].Object = sd
			osd := types.ObjectStreamDict{StreamDict: types.StreamDict{Dict: types.Dict{}, Raw: []byte("12 0 << /Broken")}, ObjCount: 1, FirstObjOffset: 5}
			if err := parseObjectStream(t.Context(), &osd, ctx.Configuration.Limits); err != nil {
				t.Fatal(err)
			}
			ctx.Table[12] = model.NewXRefTableEntryGen0(osd.ObjArray[0])
			err := finalizeHintStreams(t.Context(), ctx, state)
			if (err == nil) != encrypted {
				t.Fatalf("mode=%d encrypted=%t err=%v", mode, encrypted, err)
			}
			if len(ctx.LinearizationObjs) != 0 {
				t.Fatal("unreadable origin authorized identity")
			}
		}
	}
}

func hintDiscoveryTestContext(t *testing.T) *model.Context {
	t.Helper()
	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	if _, err := zw.Write([]byte("synthetic-hint-table")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	var file []byte
	var hintOffset, span int
	for range 3 {
		var b bytes.Buffer
		b.WriteString("%PDF-1.7\n")
		fmt.Fprintf(&b, "10 0 obj\n<< /Linearized 1 /L %010d /H [%010d %010d] /O 1 /E %010d /N 1 /T %010d >>\nendobj\n", len(file), hintOffset, span, len(file)-1, len(file)-1)
		hintOffset = b.Len()
		fmt.Fprintf(&b, "11 0 obj\n<< /Length %d /Filter /FlateDecode /S 0 >>\nstream\n%s\nendstream\nendobj\n", compressed.Len(), compressed.Bytes())
		span = b.Len() - hintOffset
		pageOffset := b.Len()
		b.WriteString("1 0 obj << /Type /Page >> endobj\n")
		xrefOffset := b.Len()
		fmt.Fprintf(&b, "xref\n0 2\n0000000000 65535 f \n%010d 00000 n \n10 2\n%010d 00000 n \n%010d 00000 n \ntrailer\n<< /Size 12 >>\nstartxref\n%d\n%%%%EOF\n", pageOffset, len("%PDF-1.7\n"), hintOffset, xrefOffset)
		file = b.Bytes()
	}
	ctx, err := model.NewContext(bytes.NewReader(file), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	for nr, offset := range map[int]int{10: len("%PDF-1.7\n"), 11: hintOffset, 1: hintOffset + span} {
		e := model.NewXRefTableEntryGen0(nil)
		v := int64(offset)
		e.Offset = &v
		ctx.Table[nr] = e
	}
	return ctx
}

func TestHintDiscoveryIdentityGuards(t *testing.T) {
	for _, name := range []string{"valid", "alias", "free", "compressed-target", "repaired", "wrong-generation", "wrong-object-tag", "metadata", "image", "non-first", "stale"} {
		t.Run(name, func(t *testing.T) {
			ctx := hintDiscoveryTestContext(t)
			file := ctx.Read.RS.(*bytes.Reader)
			bb := make([]byte, file.Size())
			if _, err := file.ReadAt(bb, 0); err != nil {
				t.Fatal(err)
			}
			switch name {
			case "alias":
				e := model.NewXRefTableEntryGen0(nil)
				e.Offset = ctx.Table[11].Offset
				ctx.Table[12] = e
			case "free":
				ctx.Table[11].Free = true
			case "compressed-target":
				ctx.Table[11].Compressed = true
			case "repaired":
				ctx.Read.RepairedXRef = true
			case "wrong-generation":
				v := 1
				ctx.Table[11].Generation = &v
			case "wrong-object-tag":
				bb = bytes.Replace(bb, []byte("11 0 obj"), []byte("12 0 obj"), 1)
			case "metadata", "image":
				// Type entries are ineligible before any envelope/decoding probe.
				marker := []byte("/S 0")
				replacement := []byte("/Type /Metadata")
				if name == "image" {
					replacement = []byte("/Type /XObject /Subtype /Image")
				}
				bb = bytes.Replace(bb, marker, replacement, 1)
				// Discovery must reject without decoding regardless of spans.
			case "non-first":
				bb = append([]byte("%PDF-1.7\n99 0 obj << >> endobj\n"), bb[len("%PDF-1.7\n"):]...)
			case "stale":
				ctx.Read.FileSize++
			}
			ctx.Read.RS = bytes.NewReader(bb)
			state, err := discoverHintStreams(t.Context(), ctx)
			if err != nil {
				t.Fatal(err)
			}
			if (state != nil) != (name == "valid") {
				t.Fatalf("candidate=%t", state != nil)
			}
		})
	}
}

func TestHintCompatibilityLeavesExplicitCryptUntouched(t *testing.T) {
	ctx, _ := hintAuditTestContext(t)
	ctx.EncKey = bytes.Repeat([]byte{1}, 16)
	ctx.AES4Streams = true
	ctx.E = &model.Enc{R: 4, Emd: true}
	raw := []byte("explicit-identity-plaintext")
	sd := types.StreamDict{Dict: types.Dict{"Filter": types.Name("Crypt"), "DecodeParms": types.Dict{"Name": types.Name("Identity")}}, Raw: raw, FilterPipeline: []types.PDFFilter{{Name: "Crypt", DecodeParms: types.Dict{"Name": types.Name("Identity")}}}}
	if err := saveDecodedStreamContent(ctx, &sd, 11, 0, true); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(sd.Raw, raw) || !bytes.Equal(sd.Content, raw) {
		t.Fatal("explicit Identity path changed")
	}
}

func TestOrdinaryAESAndTerminalProbeErrors(t *testing.T) {
	ctx, _ := hintAuditTestContext(t)
	ctx.EncKey = bytes.Repeat([]byte{1}, 16)
	ctx.AES4Streams = true
	ctx.E = &model.Enc{R: 4, Emd: true}
	sd := types.StreamDict{Dict: types.Dict{}, Raw: bytes.Repeat([]byte{1}, 33)}
	if err := decryptStreamContent(ctx, &sd, 1, 0); !errors.Is(err, errAESCiphertextUnaligned) {
		t.Fatalf("ordinary malformed AES: %v", err)
	}
	for _, err := range []error{context.Canceled, context.DeadlineExceeded, filter.ErrDecodeLimitExceeded, model.ErrMaxRecursionDepthExceeded, errObjectBufferLimit} {
		if !terminalHintError(err) {
			t.Fatalf("terminal probe error not recognized: %v", err)
		}
	}
}
