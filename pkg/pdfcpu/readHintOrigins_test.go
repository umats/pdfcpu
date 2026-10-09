package pdfcpu

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintPhysicalGapFraming(t *testing.T) {
	for _, tc := range []struct {
		name, gap string
		valid     bool
	}{
		{"separators", "\n% harmless\n ", true},
		{"clipped-comment", "% unfinished ", false},
		{"epilogue", "startxref\n12\n%%EOF\n", true},
		{"epilogue-comment", "startxref\n12\n%%EOF\n% harmless\n", true},
		{"hidden-value", "startxref\n12\n%%EOF\n<< /Target 8 0 R >>\nstartxref\n12\n%%EOF\n", false},
		{"false-offset", "startxref\n99\n%%EOF\n", false},
		{"false-eof", "startxref\n12\n%%EOFJUNK\n", false},
		{"clipped-eof", "startxref\n12\n%%EOF", false},
		{"clipped-eof-space", "startxref\n12\n%%EOF ", false},
		{"clipped-header-comment", "startxref\n12\n%%EOF\n% ", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := model.NewContext(bytes.NewReader([]byte(tc.gap+"next")), model.NewDefaultConfiguration())
			if err != nil {
				t.Fatal(err)
			}
			ctx.Read.FileSize = int64(len(tc.gap) + len("next"))
			ok, err := hintPhysicalGap(t.Context(), ctx, 0, int64(len(tc.gap)), map[int64]bool{12: true}, false)
			if err != nil || ok != tc.valid {
				t.Fatalf("valid=%t err=%v", ok, err)
			}
		})
	}
}

func TestHintPhysicalGapRejectsUnvisitedXRef(t *testing.T) {
	gap := "xref\n0 1\n0000000000 65535 f \ntrailer\n<< /Size 12 /Target 8 0 R >>\nstartxref\n0\n%%EOF\n"
	ctx, err := model.NewContext(bytes.NewReader([]byte(gap+"next")), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := hintPhysicalGap(t.Context(), ctx, 0, int64(len(gap)), map[int64]bool{}, true); err != nil || ok {
		t.Fatalf("unvisited xref accepted: valid=%t err=%v", ok, err)
	}
	p := &contextutil.ParseProvenance{}
	p.RecordXRefOrigin(0)
	c := contextutil.WithParseProvenance(t.Context(), p)
	if ok, err := hintPhysicalGap(c, ctx, 0, int64(len(gap)), map[int64]bool{}, true); err != nil || !ok {
		t.Fatalf("visited xref framing denied: valid=%t err=%v", ok, err)
	}
}

func TestHintPhysicalMissingLength(t *testing.T) {
	file := []byte("%PDF-1.7\n1 0 obj\n<< >>\nstream\nx\nendstream\nendobj\n")
	ctx, err := model.NewContext(bytes.NewReader(file), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok, err := hintPhysicalEnd(t.Context(), ctx, int64(len("%PDF-1.7\n")), 1, 0); err != nil || ok {
		t.Fatalf("missing length authorized an envelope: valid=%t err=%v", ok, err)
	}
}

func TestHintPhysicalCachedXRefStream(t *testing.T) {
	for _, tc := range []struct {
		name            string
		cached, visited bool
	}{
		{"parsed-and-visited", true, true},
		{"unvisited", true, false},
		{"name-only", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			offset := int64(len("%PDF-1.7\n"))
			file := fmt.Sprintf("%%PDF-1.7\n1 0 obj\n<< /Type /XRef /Length 0 >>\nstream\n\nendstream\nendobj\nstartxref\n%d\n%%%%EOF\n", offset)
			ctx, err := model.NewContext(bytes.NewReader([]byte(file)), model.NewDefaultConfiguration())
			if err != nil {
				t.Fatal(err)
			}
			length := int64(0)
			sd := types.NewStreamDict(types.Dict{"Type": types.Name("XRef"), "Length": types.Integer(0)}, 0, &length, nil, nil)
			e := model.NewXRefTableEntryGen0(sd)
			if tc.cached {
				e.Object = types.XRefStreamDict{StreamDict: sd}
			}
			e.Offset = &offset
			ctx.Table[1] = e
			p := &contextutil.ParseProvenance{}
			if tc.visited {
				p.RecordXRefOrigin(offset)
			}
			c := contextutil.WithParseProvenance(t.Context(), p)
			ok, err := proveHintOrigins(c, ctx, offset)
			if err != nil || ok != (tc.cached && tc.visited) {
				t.Fatalf("cached xref origin: valid=%t err=%v", ok, err)
			}
		})
	}
}

func TestHintPhysicalOrigins(t *testing.T) {
	ctx := hintDiscoveryTestContext(t)
	fileBytes := make([]byte, ctx.Read.FileSize)
	if _, err := ctx.Read.RS.(*bytes.Reader).ReadAt(fileBytes, 0); err != nil {
		t.Fatal(err)
	}
	p := &contextutil.ParseProvenance{}
	p.RecordXRefOrigin(int64(bytes.Index(fileBytes, []byte("\nxref\n")) + 1))
	proofContext := contextutil.WithParseProvenance(t.Context(), p)
	ok, err := proveHintOrigins(proofContext, ctx, int64(len("%PDF-1.7\n")))
	if err != nil || !ok {
		t.Fatalf("valid physical chain: valid=%t err=%v", ok, err)
	}
	c, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := proveHintOrigins(c, ctx, int64(len("%PDF-1.7\n"))); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}

	payload := "2 0 obj << /Length 0 >> endobj"
	file := fmt.Sprintf("%%PDF-1.7\n1 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(payload), payload)
	ctx, err = model.NewContext(bytes.NewReader([]byte(file)), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	ctx.Read.FileSize = int64(len(file))
	for nr, offset := range map[int]int{1: len("%PDF-1.7\n"), 2: strings.Index(file, payload)} {
		e := model.NewXRefTableEntryGen0(nil)
		v := int64(offset)
		e.Offset = &v
		ctx.Table[nr] = e
	}
	if ok, err := proveHintOrigins(t.Context(), ctx, int64(len("%PDF-1.7\n"))); err != nil || ok {
		t.Fatalf("stream-contained header accepted: valid=%t err=%v", ok, err)
	}
}
