package api

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/umats/pdfcpu/pkg/filter"
	"github.com/umats/pdfcpu/pkg/pdfcpu"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
)

// plaintextPrimaryHintPDF preserves the encrypted object's numbers and encoded
// bytes, and adds a physically first linearization dictionary, early xref and
// primary hint followed by the document objects and main xref. The early trailer
// points forward to the main xref, and EOF startxref points to the early xref.
// This is a reader-order reproduction, NOT a fully conforming Annex F fixture:
// both xrefs cover all objects, object numbering is inherited from Encrypt, and
// the inflated hint payload is synthetic, not a usable page/shared-object table.
type primaryHintFixtureOptions struct {
	overflow            bool
	aligned             bool
	encrypted           bool
	referenced          bool
	keyLength           int
	lowHint             bool
	objectStream        bool
	compressedReference bool
	unencrypted         bool
	corruptOffset       bool
	ambiguousClaim      string
	ambiguousOrigin     string
	unreliableXRef      string
	shortFilter         bool
	malformedHeader     string
	boundaryOrigin      string
	trailingOrigin      string
	boundaryDepth       int
	review6             string
	review7             string
	review8             string
	review9             string
	review10            string
	review13            string
}

func plaintextPrimaryHintPDF(t *testing.T, options ...primaryHintFixtureOptions) ([]byte, []byte, []byte) {
	t.Helper()
	source, content := encryptedXRefSourcePDF(t)
	var opt primaryHintFixtureOptions
	if len(options) > 0 {
		opt = options[0]
	}
	tc := encryptedXRefTestCase{aes: true, keyLength: 128, userPW: "user"}
	if opt.keyLength != 0 {
		tc.keyLength = opt.keyLength
	}
	clean := source
	if !opt.unencrypted {
		var encrypted bytes.Buffer
		if err := Encrypt(t.Context(), bytes.NewReader(source), &encrypted, encryptedXRefConfig(tc)); err != nil {
			t.Fatal(err)
		}
		clean = bytes.Clone(encrypted.Bytes())
	}
	ctx, err := ReadContext(t.Context(), bytes.NewReader(clean), encryptedXRefReadConfig(tc.userPW))
	if err != nil {
		t.Fatal(err)
	}
	wantCFM := "Identity"
	if !opt.unencrypted {
		enc, err := ctx.EncryptDict()
		if err != nil {
			t.Fatal(err)
		}
		wantCFM = "AESV2"
		if tc.keyLength == 256 {
			wantCFM = "AESV3"
		}
		if !ctx.AES4Streams || enc.NameEntry("StmF") == nil || *enc.NameEntry("StmF") != "StdCF" ||
			enc.DictEntry("CF").DictEntry("StdCF").NameEntry("CFM") == nil ||
			*enc.DictEntry("CF").DictEntry("StdCF").NameEntry("CFM") != wantCFM {
			t.Fatal("fixture did not select AESV2 via StmF/StdCF")
		}
	}

	// Keep the original serialized ciphertext, not the decrypted context objects.
	var numbers []int
	for nr, e := range ctx.Table {
		if !e.Free && e.Offset != nil && *e.Offset > 0 {
			numbers = append(numbers, nr)
		}
	}
	sort.Slice(numbers, func(i, j int) bool { return *ctx.Table[numbers[i]].Offset < *ctx.Table[numbers[j]].Offset })
	trailerStart := bytes.LastIndex(clean, []byte("trailer\n"))
	startXRef := bytes.LastIndex(clean, []byte("startxref\n"))
	if trailerStart < 0 || startXRef < trailerStart {
		t.Fatal("missing classic trailer")
	}
	trailer := bytes.TrimSpace(clean[trailerStart+len("trailer\n") : startXRef])
	objects := map[int][]byte{}
	maxNr := 0
	for i, nr := range numbers {
		start := int(*ctx.Table[nr].Offset)
		end := trailerStart
		if i+1 < len(numbers) {
			end = int(*ctx.Table[numbers[i+1]].Offset)
		}
		endObj := bytes.LastIndex(clean[start:end], []byte("endobj"))
		if endObj < 0 {
			t.Fatal("missing serialized endobj")
		}
		objects[nr] = bytes.Clone(clean[start : start+endObj+len("endobj")])
		maxNr = max(maxNr, nr)
	}
	linearNr, hintNr := maxNr+1, maxNr+2
	size := hintNr + 1
	if opt.lowHint {
		linearNr, hintNr = hintNr, linearNr
	}
	if opt.overflow {
		size++
	}
	compressedNr, containerNr, xrefNr := maxNr+3, maxNr+4, maxNr+5
	secondaryNr := xrefNr + 1
	secondMember := opt.review9 == "compressed-comment" || opt.review9 == "compressed-name"
	if opt.objectStream {
		size = xrefNr + 1
		if secondMember {
			size++
		}
	}
	if opt.unreliableXRef == "skipped-entry" || opt.unreliableXRef == "short-offset" {
		size = compressedNr + 1
		objects[compressedNr] = fmt.Appendf(nil, "%d 0 obj\n<< /Unused 42 >>\nendobj", compressedNr)
		numbers = append(numbers, compressedNr)
	}
	if opt.review10 == "comment-endobj" {
		size++
		objects[size-1] = fmt.Appendf(nil, "%d 0 obj\n/A\\%% endobj\n<< /Target %d 0 R >>\nendobj", size-1, hintNr)
		numbers = append(numbers, size-1)
	}
	if opt.trailingOrigin == "ordinary" || opt.trailingOrigin == "ordinary-comment" {
		tail := fmt.Sprintf(" << /Target %d 0 R >> ", hintNr)
		if opt.trailingOrigin == "ordinary-comment" {
			tail = " % harmless comment\n "
		}
		objects[1] = bytes.Replace(objects[1], []byte("endobj"), []byte(tail+"endobj"), 1)
	}
	if strings.HasSuffix(opt.review8, "ordinary") && strings.HasPrefix(opt.review8, "keyword-") {
		objects[1] = bytes.Replace(objects[1], []byte("/Pages"), []byte(round8KeywordSyntax(opt.review8)+" /Pages"), 1)
	}
	if strings.HasSuffix(opt.review8, "trailer") && strings.HasPrefix(opt.review8, "keyword-") {
		trailer = bytes.Replace(trailer, []byte("<<"), []byte("<< "+round8KeywordSyntax(opt.review8)), 1)
	}
	if strings.HasSuffix(opt.review7, "ordinary") {
		syntax := round7LexicalSyntax(opt.review7)
		objects[1] = bytes.Replace(objects[1], []byte("/Pages"), []byte(syntax+" /Pages"), 1)
	}
	if strings.HasSuffix(opt.review7, "trailer") {
		trailer = bytes.Replace(trailer, []byte("<<"), []byte("<< "+round7LexicalSyntax(opt.review7)), 1)
	}
	if opt.review6 == "hex-ordinary" {
		objects[1] = bytes.Replace(objects[1], []byte("/Pages"), []byte("/Ignored <GG> /Pages"), 1)
	}
	if opt.review13 == "comment-ordinary" {
		changed := false
		for nr, obj := range objects {
			stream := bytes.Index(obj, []byte("stream\n"))
			if stream < 0 {
				continue
			}
			end := bytes.LastIndex(obj[:stream], []byte(">>"))
			if end < 0 {
				t.Fatal("ordinary stream lacks dictionary")
			}
			objects[nr] = append(bytes.Clone(obj[:end]), append([]byte(" /DecodeParms << /Ignored /stream >> >> % >>stream\n"), obj[stream+len("stream\n"):]...)...)
			changed = true
			break
		}
		if !changed {
			t.Fatal("no ordinary stream mutated")
		}
	}
	if opt.review7 == "stream-space-ordinary" {
		for nr, obj := range objects {
			if bytes.Contains(obj, []byte("stream\n")) {
				objects[nr] = bytes.Replace(obj, []byte("stream\n"), []byte("stream "), 1)
				break
			}
		}
	}
	if opt.review7 == "postlude-ordinary" || opt.review7 == "postlude-comment" {
		for nr, obj := range objects {
			if bytes.Contains(obj, []byte("endstream\nendobj")) {
				tail := fmt.Sprintf("endstream\n<< /Target %d 0 R >>\nendobj", hintNr)
				if opt.review7 == "postlude-comment" {
					tail = "endstream\n% harmless NULL\\v\nendobj"
				}
				objects[nr] = bytes.Replace(obj, []byte("endstream\nendobj"), []byte(tail), 1)
				break
			}
		}
	}
	if opt.review6 == "single-ID" {
		trailer = regexp.MustCompile(`(/ID\s*\[\s*<[^>]*>)\s*<[^>]*>\s*\]`).ReplaceAll(trailer, []byte("${1}]"))
	}
	if opt.ambiguousOrigin == "ordinary" {
		objects[1] = bytes.Replace(objects[1], []byte("/Pages"), fmt.Appendf(nil, "/Unused %d 0 R /Unused 42 /Pages", hintNr), 1)
	}
	if strings.HasPrefix(opt.ambiguousOrigin, "ordinary-null") {
		unused := "/Unused null /Unused 42 "
		if strings.HasSuffix(opt.ambiguousOrigin, "last") {
			unused = "/Unused 42 /Unused null "
		}
		objects[1] = bytes.Replace(objects[1], []byte("/Pages"), []byte(unused+"/Pages"), 1)
	}
	if strings.HasPrefix(opt.ambiguousOrigin, "trailer-null") {
		unused := "<< /Unused null /Unused 42 "
		if strings.HasSuffix(opt.ambiguousOrigin, "last") {
			unused = "<< /Unused 42 /Unused null "
		}
		trailer = bytes.Replace(trailer, []byte("<<"), []byte(unused), 1)
	}
	if opt.ambiguousOrigin == "ordinary-fallback" {
		objects[1] = bytes.Replace(objects[1], []byte("<<"), []byte("<< /Missing\n"), 1)
	}
	if opt.ambiguousOrigin == "trailer" {
		trailer = bytes.Replace(trailer, []byte("<<"), fmt.Appendf(nil, "<< /Unused %d 0 R /Unused 42", hintNr), 1)
	}
	if opt.ambiguousOrigin == "trailer-fallback" {
		trailer = bytes.Replace(trailer, []byte("<<"), []byte("<< /Missing\n"), 1)
	}
	if opt.boundaryOrigin == "stream" {
		for nr, obj := range objects {
			if bytes.Contains(obj, []byte("stream\n")) {
				objects[nr] = bytes.Replace(obj, []byte("<<"), []byte("<< /Boundary [[[42]]] "), 1)
				break
			}
		}
	}
	if opt.referenced {
		objects[1] = bytes.Replace(objects[1], []byte("/Pages"), fmt.Appendf(nil, "/Unused %d 0 R /Pages", hintNr), 1)
	}
	declaredSize := size
	if opt.unreliableXRef == "classic-size" {
		declaredSize--
	}
	trailer = regexp.MustCompile(`/Size\s+\d+`).ReplaceAll(trailer, fmt.Appendf(nil, "/Size %d", declaredSize))

	// A deterministic nonaligned Flate stream, with no per-stream Crypt filter.
	payload := make([]byte, 36)
	for i := range 128 {
		payload = append(payload, byte(i))
	}
	var hint []byte
	for i := 0; i < 256; i++ {
		var compressed bytes.Buffer
		zw := zlib.NewWriter(&compressed)
		if _, err := zw.Write(payload); err != nil {
			t.Fatal(err)
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		hint = compressed.Bytes()
		if !opt.aligned || len(hint)%16 == 0 {
			break
		}
		payload = append(payload, byte(i*73+157))
	}
	if len(hint) <= 16 || (!opt.aligned && len(hint)%16 == 0) || (opt.aligned && len(hint)%16 != 0) {
		t.Fatalf("unexpected encoded hint length %d (aligned=%t)", len(hint), opt.aligned)
	}
	t.Logf("primary encoded hint length=%d, inflated length=%d; stream encryption method=%s", len(hint), len(payload), wantCFM)
	hintNumbers := []int{hintNr}
	if opt.overflow {
		hintNumbers = append(hintNumbers, hintNr+1)
	}
	for _, nr := range hintNumbers {
		encoded := hint
		if opt.encrypted {
			encoded = encryptPrimaryHintBytes(t, hint, nr, ctx)
		}
		prefix := ""
		if opt.ambiguousClaim == "Length" {
			prefix = "/Length 1 "
		}
		if opt.ambiguousClaim == "S" {
			prefix = "/S -1 "
		}
		suffix := ""
		for _, key := range []string{"Length", "S"} {
			if opt.ambiguousClaim == key+"-null-first" {
				prefix = "/" + key + " null "
			}
			if opt.ambiguousClaim == key+"-null-last" {
				suffix = " /" + key + " null"
			}
		}
		if opt.ambiguousClaim == "numeric-S" {
			suffix += " /O 0+0"
		}
		if strings.HasSuffix(opt.review7, "hint") {
			suffix += " " + round7LexicalSyntax(opt.review7)
		}
		if opt.review6 == "hex-hint" {
			suffix += " /O <GG>"
		}
		filterName := "FlateDecode"
		if opt.shortFilter {
			filterName = "Fl"
		}
		header := fmt.Sprintf("%d 0 obj", nr)
		if opt.malformedHeader == "hint-junk" {
			header = "junk" + header
		}
		if opt.malformedHeader == "hint-bj" {
			header = strings.Replace(header, "obj", "bj", 1)
		}
		if opt.trailingOrigin == "hint" {
			suffix += fmt.Sprintf(" >> << /Target %d 0 R", hintNr)
		}
		objects[nr] = fmt.Appendf(nil, "%s\n<< %s/Length %d /Filter /%s /S 36%s >>\nstream\n%s\nendstream\nendobj", header, prefix, len(encoded), filterName, suffix, encoded)
		if opt.review13 == "comment-hint" || opt.review13 == "name-hint" {
			replacement := " /DecodeParms << /Ignored /stream >> >>\nstream\n"
			if opt.review13 == "comment-hint" {
				replacement = " /DecodeParms << /Ignored /stream >> >> % >>stream\n"
			}
			objects[nr] = bytes.Replace(objects[nr], []byte(" >>\nstream\n"), []byte(replacement), 1)
		}
		if opt.review7 == "stream-space" {
			objects[nr] = bytes.Replace(objects[nr], []byte("stream\n"), []byte("stream "), 1)
		}
		if opt.review7 == "clipped-endobj" {
			objects[nr] = append(objects[nr], []byte("JUNK")...)
		}
	}

	if opt.objectStream {
		prolog := fmt.Sprintf("%d 0 ", compressedNr)
		body := "<< /Unused 42 >>"
		if opt.boundaryOrigin == "compressed" {
			depth := opt.boundaryDepth
			if depth == 0 {
				depth = 3
			}
			body = "<< /Boundary " + strings.Repeat("[", depth) + "42" + strings.Repeat("]", depth) + " >>"
		}
		if opt.trailingOrigin == "compressed" || opt.trailingOrigin == "compressed-comment" {
			tail := fmt.Sprintf(" << /Target %d 0 R >>", hintNr)
			if opt.trailingOrigin == "compressed-comment" {
				tail = " % harmless comment\n "
			}
			body += tail
		}
		if opt.compressedReference {
			body = fmt.Sprintf("<< /Unused %d 0 R >>", hintNr)
		}
		if opt.review10 == "compressed-length" {
			lengthPattern := regexp.MustCompile(`/Length\s+([0-9]+)`)
			found := false
			for _, nr := range numbers {
				stream := bytes.Index(objects[nr], []byte("stream\n"))
				if stream < 0 {
					continue
				}
				match := lengthPattern.FindSubmatch(objects[nr][:stream])
				if len(match) != 2 {
					continue
				}
				body = string(match[1])
				objects[nr] = bytes.Replace(objects[nr], match[0], fmt.Appendf(nil, "/Length %d 0 R", compressedNr), 1)
				found = true
				break
			}
			if !found {
				t.Fatal("no ordinary stream with a direct length")
			}
		}
		if opt.ambiguousOrigin == "compressed" {
			body = fmt.Sprintf("<< /Unused %d 0 R /Unused 42 >>", hintNr)
		}
		if strings.HasPrefix(opt.ambiguousOrigin, "compressed-null") {
			body = "<< /Unused null /Unused 42 >>"
			if strings.HasSuffix(opt.ambiguousOrigin, "last") {
				body = "<< /Unused 42 /Unused null >>"
			}
		}
		if opt.ambiguousOrigin == "compressed-fallback" {
			body = "<< /Missing\n/Unused 42 >>"
		}
		if strings.HasSuffix(opt.review8, "compressed") && strings.HasPrefix(opt.review8, "keyword-") {
			body = "<< " + round8KeywordSyntax(opt.review8) + " >>"
		}
		if strings.HasSuffix(opt.review7, "compressed") {
			body = "<< " + round7LexicalSyntax(opt.review7) + " >>"
		}
		if opt.review6 == "hex-compressed" {
			body = "<< /Unused <GG> >>"
		}
		prefix := ""
		switch opt.review6 {
		case "member-prefix":
			prefix = fmt.Sprintf("%d 0 R ", hintNr)
		case "member-comment":
			prefix = "% harmless\n "
		case "prolog-tail":
			prolog += fmt.Sprintf("%% comment\n%d 0 R\n", hintNr)
		case "prolog-comment":
			prolog += "% harmless\n "
		case "index-space":
			prolog = fmt.Sprintf("%d\v0 ", compressedNr)
		case "index-junk":
			prolog = "junk 0 "
		case "index-mismatch":
			prolog = fmt.Sprintf("%d 0 ", compressedNr+20)
		}
		if prefix != "" {
			prolog = fmt.Sprintf("%d %d ", compressedNr, len(prefix))
		}
		memberCount := 1
		switch opt.review9 {
		case "compressed-first-one":
			prolog, body = fmt.Sprintf("%d 0", compressedNr), "1"
		case "compressed-first-42":
			prolog, body = fmt.Sprintf("%d 0", compressedNr), "42"
		case "compressed-comment":
			prolog, body, memberCount = fmt.Sprintf("%d 0 %d 6 ", compressedNr, secondaryNr), "42 %x 43", 2
		case "compressed-name":
			prolog, body, memberCount = fmt.Sprintf("%d 0 %d 1 ", compressedNr, secondaryNr), "/1", 2
		case "compressed-gap-comment":
			prolog, prefix, body = fmt.Sprintf("%d 4 ", compressedNr), "%x  ", "43"
		}
		encoded := []byte(prolog + prefix + body)
		if !opt.unencrypted {
			encoded = encryptPrimaryHintBytes(t, encoded, containerNr, ctx)
		}
		objects[containerNr] = fmt.Appendf(nil, "%d 0 obj\n<< /Type /ObjStm /N %d /First %d /Length %d >>\nstream\n%s\nendstream\nendobj", containerNr, memberCount, len(prolog), len(encoded), encoded)
		if opt.review7 == "stream-space-container" {
			objects[containerNr] = bytes.Replace(objects[containerNr], []byte("stream\n"), []byte("stream "), 1)
		}
		if opt.review7 == "postlude-container" {
			objects[containerNr] = bytes.Replace(objects[containerNr], []byte("endstream\nendobj"), fmt.Appendf(nil, "endstream\n<< /Target %d 0 R >>\nendobj", hintNr), 1)
		}
		numbers = append(numbers, containerNr)
		trailer = regexp.MustCompile(`/Size\s+\d+`).ReplaceAll(trailer, fmt.Appendf(nil, "/Size %d", size))
	}

	// All offset fields have fixed width, allowing layout to converge without
	// rewriting any original encrypted string/stream or changing its object key.
	olderNativeNr := 0
	if strings.HasPrefix(opt.review13, "older-native-") {
		olderNativeNr = size
		size++
		trailer = regexp.MustCompile(`/Size\s+\d+`).ReplaceAll(trailer, fmt.Appendf(nil, "/Size %d", size))
	}
	offsets := make([]int, size)
	var fileLength, mainXRef, hintLength, overflowLength, firstPageEnd int
	var result []byte
	for range 3 {
		var b bytes.Buffer
		b.WriteString("%PDF-1.7\n")
		if opt.corruptOffset {
			b.WriteString("% bad bad obj\n")
		} else {
			b.WriteString("%\xE2\xE3\xCF\xD3\n")
		}
		offsets[linearNr] = b.Len()
		h := fmt.Sprintf("%010d %010d", offsets[hintNr], hintLength)
		if opt.ambiguousClaim == "numeric-H" {
			h = fmt.Sprintf("0+%08d %010d", offsets[hintNr], hintLength)
		}
		if opt.overflow {
			h += fmt.Sprintf(" %010d %010d", offsets[hintNr+1], overflowLength)
		}
		prefix := ""
		if opt.ambiguousClaim == "H" {
			prefix = "/H [0000000000 0000000001] "
		}
		if opt.ambiguousClaim == "H-null-first" {
			prefix = "/H null "
		}
		suffix := ""
		if opt.ambiguousClaim == "H-null-last" {
			suffix = " /H null"
		}
		if strings.HasSuffix(opt.review8, "linear") && strings.HasPrefix(opt.review8, "keyword-") {
			suffix += " " + round8KeywordSyntax(opt.review8)
		}
		if opt.trailingOrigin == "linear" {
			suffix += fmt.Sprintf(" >> << /Target %d 0 R", hintNr)
		}
		if opt.ambiguousClaim == "fallback" {
			suffix = " /Missing\n/Extra /IgnoredName"
		}
		linearHeader := fmt.Sprintf("%d 0 obj", linearNr)
		if opt.malformedHeader == "linear-junk" {
			linearHeader = "junk" + linearHeader
		}
		if opt.malformedHeader == "linear-bj" {
			linearHeader = strings.Replace(linearHeader, "obj", "bj", 1)
		}
		fmt.Fprintf(&b, "%s\n<< %s/Linearized 1 /L %010d /H [%s] /O 3 /E %010d /N 1 /T %010d%s >>\nendobj\n", linearHeader, prefix, fileLength, h, firstPageEnd, mainXRef, suffix)
		if opt.review8 == "endobj-linear" {
			b.Bytes()[b.Len()-1] = '\v'
		}
		earlyXRef := b.Len()
		writeXRef := func(prev bool) {
			sectionStart := b.Len()
			defer func() {
				kind := opt.review8
				if !strings.HasPrefix(kind, "xref-") || (strings.HasSuffix(kind, "-older") && prev) {
					return
				}
				section := b.Bytes()[sectionStart:]
				switch {
				case strings.Contains(kind, "header"):
					section[len("xref\n0")] = '\v'
				case strings.Contains(kind, "record"):
					start := bytes.IndexByte(section[5:], '\n') + 6
					section[start+10] = '\v'
				}
			}()

			count := size
			if !prev && opt.unreliableXRef == "older-surplus" {
				count--
			}
			xrefMarker := "xref"
			if opt.review8 == "xref-marker" {
				xrefMarker += "\v"
			}
			if opt.unreliableXRef == "missing-zero" {
				fmt.Fprintf(&b, "%s\n1 %d\n", xrefMarker, size-1)
			} else if opt.unreliableXRef == "free-list" {
				fmt.Fprintf(&b, "%s\n0 %d\n0000000001 00000 f \n", xrefMarker, size)
			} else {
				fmt.Fprintf(&b, "%s\n0 %d\n0000000000 65535 f \n", xrefMarker, count)
			}
			for nr := 1; nr < size; nr++ {
				if offsets[nr] == 0 {
					b.WriteString("0000000000 00000 f \n")
				} else {
					entryOffset := offsets[nr]
					if nr == compressedNr {
						if opt.unreliableXRef == "skipped-entry" {
							entryOffset = 0
						}
						if opt.unreliableXRef == "short-offset" {
							entryOffset = 1
						}
					}
					if opt.corruptOffset && nr == 2 {
						entryOffset = bytes.Index(b.Bytes(), []byte("bad bad obj"))
					}
					if (opt.unreliableXRef == "prefixed-offset" || (!prev && opt.unreliableXRef == "older-shape")) && nr == hintNr {
						fmt.Fprintf(&b, "x%09d 00000 n \n", entryOffset)
					} else {
						fmt.Fprintf(&b, "%010d 00000 n \n", entryOffset)
					}
				}
			}
			if opt.unreliableXRef == "duplicate-classic" {
				// The retained origin is the hint; the conflicting assignment is ordinary content.
				fmt.Fprintf(&b, "%d 1\n%010d 00000 n \n", hintNr, offsets[1])
			}
			b.WriteString("trailer\n")
			if opt.review6 == "trailer-prefix" {
				fmt.Fprintf(&b, "%d 0 R\n", hintNr)
			}
			if opt.review6 == "trailer-inline" {
				fmt.Fprintf(&b, "%d 0 R ", hintNr)
			}
			if opt.review6 == "trailer-comment" {
				b.WriteString("% harmless\n ")
			}
			if prev {
				b.Write(trailer[:len(trailer)-2])
				fmt.Fprintf(&b, " /Prev %010d >>\n", mainXRef)
			} else {
				if opt.unreliableXRef == "prev-zero" || opt.unreliableXRef == "repeated-offset" {
					b.Write(trailer[:len(trailer)-2])
					prevOffset := mainXRef
					if opt.unreliableXRef == "prev-zero" {
						prevOffset = 0
					}
					fmt.Fprintf(&b, " /Prev %010d >>", prevOffset)
				} else {
					b.Write(trailer)
				}
				b.WriteByte('\n')
			}
			if !prev && opt.review10 == "complete-epilogue" {
				fmt.Fprintf(&b, "startxref\n%d\n%%%%EOF\n<< /Target %d 0 R >>\n", mainXRef, hintNr)
			}
			if !prev && strings.HasPrefix(opt.review9, "trailer-") {
				marker := "% startxref"
				if opt.review9 == "trailer-junk" {
					marker = "startxrefJUNK"
				}
				if opt.review9 == "trailer-spaced-junk" {
					marker = "startxref JUNK"
				}
				if opt.review9 == "trailer-fake-frame" {
					marker = "startxref\n0"
				}
				fmt.Fprintln(&b, marker)
				if opt.review9 != "trailer-harmless" {
					fmt.Fprintf(&b, "<< /Target %d 0 R >>\n", hintNr)
				}
			}
		}
		if !opt.objectStream {
			writeXRef(true)
			if opt.unreliableXRef == "repeated-offset" {
				// Keep the earlier EOF beyond the first partial 512-byte scan
				// window, regardless of randomized encrypted string lengths.
				b.WriteByte('%')
				b.WriteString(strings.Repeat(" ", 1024))
				b.WriteByte('\n')
				fmt.Fprintf(&b, "startxref\n%d\n%%%%EOF\n", earlyXRef)
			} else {
				fmt.Fprintf(&b, "startxref\n0\n%%%%EOF\n")
			}
		}
		if opt.review10 == "comment-origin" {
			b.WriteString("% ")
		}
		if opt.review10 == "stream-origin" {
			fmt.Fprintf(&b, "%d 0 obj\n<< /Length %d >>\nstream\n", size+10, len(objects[hintNr])+1)
		}
		offsets[hintNr] = b.Len()
		b.Write(objects[hintNr])
		b.WriteByte('\n')
		hintLength = b.Len() - offsets[hintNr]
		if opt.review10 == "stream-origin" {
			b.WriteString("endstream\nendobj\n")
		}
		if opt.review7 == "clipped-endobj" {
			hintLength -= len("JUNK\n")
		}
		if opt.overflow {
			offsets[hintNr+1] = b.Len()
			b.Write(objects[hintNr+1])
			b.WriteByte('\n')
			overflowLength = b.Len() - offsets[hintNr+1]
		}
		if opt.review10 == "counterfeit-xref" {
			fakeOffset := b.Len()
			fmt.Fprintf(&b, "xref\n0 1\n0000000000 65535 f \ntrailer\n<< /Size %d /Target %d 0 R >>\nstartxref\n%d\n%%%%EOF\n", size, hintNr, fakeOffset)
		}
		for _, nr := range numbers {
			offsets[nr] = b.Len()
			b.Write(objects[nr])
			if nr == 1 && opt.review8 == "endobj-ordinary" {
				b.WriteByte('\v')
			} else {
				b.WriteByte('\n')
			}
		}
		if olderNativeNr != 0 {
			offsets[olderNativeNr] = b.Len()
			entry := make([]byte, 11)
			entry[0] = 1
			target := hintNr
			if strings.Contains(opt.review13, "generation") {
				entry[6] = 1
			} else {
				entry[0], target = 2, compressedNr
				if strings.Contains(opt.review13, "container") {
					entry[1] = 1
				} else if strings.Contains(opt.review13, "index") {
					entry[6] = 1
				}
			}
			fmt.Fprintf(&b, "%d 0 obj\n", olderNativeNr)
			b.Write(trailer[:len(trailer)-2])
			fmt.Fprintf(&b, " /Type /XRef /Index [%d 1] /W [1 5 5] /Length 11 >>\nstream\n", target)
			b.Write(entry)
			b.WriteString("\nendstream\nendobj\n")
		}
		firstPageEnd = b.Len()
		mainXRef = b.Len()
		if opt.objectStream {
			offsets[xrefNr] = mainXRef
			data := make([]byte, size*7)
			for nr := range size {
				entry := data[nr*7 : (nr+1)*7]
				switch {
				case nr == 0:
					binary.BigEndian.PutUint16(entry[5:], 65535)
				case nr == compressedNr || (secondMember && nr == secondaryNr):
					entry[0] = 2
					binary.BigEndian.PutUint32(entry[1:5], uint32(containerNr))
					if nr == secondaryNr {
						binary.BigEndian.PutUint16(entry[5:], 1)
					}
				default:
					if offsets[nr] != 0 {
						entry[0] = 1
						binary.BigEndian.PutUint32(entry[1:5], uint32(offsets[nr]))
					}
				}
			}
			if opt.unreliableXRef == "surplus-stream-entry" {
				data = append(data, make([]byte, 7)...)
			}
			xrefHeader := fmt.Sprintf("%d 0 obj", xrefNr)
			if opt.unreliableXRef == "stream-header-junk" {
				xrefHeader = "junk" + xrefHeader
			}
			if opt.unreliableXRef == "stream-header-bj" {
				xrefHeader = strings.Replace(xrefHeader, "obj", "bj", 1)
			}
			fmt.Fprintf(&b, "%s\n", xrefHeader)
			xrefTrailer := trailer
			if opt.unreliableXRef == "stream-size" {
				xrefTrailer = regexp.MustCompile(`/Size\s+\d+`).ReplaceAll(trailer, fmt.Appendf(nil, "/Size %d", size-1))
			}
			b.Write(xrefTrailer[:len(xrefTrailer)-2])
			index := fmt.Sprintf("0 %d", size)
			if opt.unreliableXRef == "duplicate-stream" {
				index += fmt.Sprintf(" %d 1", hintNr)
				data = append(data, data[7:14]...)
			}
			w := "1 4 2"
			if strings.HasPrefix(opt.review7, "wide-") {
				stride := 12
				if strings.Contains(opt.review7, "generation") {
					stride = 14
				}
				wide := make([]byte, size*stride)
				for nr := range size {
					old, entry := data[nr*7:(nr+1)*7], wide[nr*stride:(nr+1)*stride]
					entry[0] = old[0]
					if strings.Contains(opt.review7, "offset") {
						w = "1 9 2"
						copy(entry[6:10], old[1:5])
						copy(entry[10:], old[5:])
						if nr == hintNr {
							entry[1] = 1
						}
					} else {
						w = "1 4 9"
						copy(entry[1:5], old[1:5])
						copy(entry[12:], old[5:])
						if nr == hintNr {
							entry[5] = 1
						}
					}
				}
				data = wide
			}
			if strings.HasPrefix(opt.review13, "native-") {
				wide := make([]byte, size*11)
				for nr := range size {
					old, entry := data[nr*7:(nr+1)*7], wide[nr*11:(nr+1)*11]
					entry[0] = old[0]
					copy(entry[2:6], old[1:5])
					copy(entry[9:11], old[5:])
					if opt.review13 == "native-generation" && nr == hintNr {
						entry[6] = 1
					}
					if opt.review13 == "native-container" && nr == compressedNr {
						entry[1] = 1
					}
					if opt.review13 == "native-index" && nr == compressedNr {
						entry[6] = 1
					}
				}
				data, w = wide, "1 5 5"
			}
			declaredLength := len(data)
			if opt.review6 == "xref-length" {
				declaredLength = 0
			}
			marker := "stream\n"
			if opt.review7 == "stream-space-xref" {
				marker = "stream "
			}
			xrefSuffix := ""
			if olderNativeNr != 0 {
				xrefSuffix = fmt.Sprintf(" /Prev %010d", offsets[olderNativeNr])
			}
			if opt.review13 == "comment-xref" {
				xrefSuffix, marker = " /DecodeParms << /Ignored /stream >>", "% >>stream\n"
			}
			fmt.Fprintf(&b, " /Type /XRef /Index [%s] /W [%s] /Length %d%s >>\n%s", index, w, declaredLength, xrefSuffix, marker)
			b.Write(data)
			if opt.review7 == "postlude-xref" {
				fmt.Fprintf(&b, "\nendstream\n<< /Target %d 0 R >>\nendobj\n", hintNr)
			} else {
				b.WriteString("\nendstream\nendobj\n")
			}
			fmt.Fprintf(&b, "startxref\n%d\n%%%%EOF\n", mainXRef)
		} else {
			writeXRef(false)
			fmt.Fprintf(&b, "startxref\n%d\n%%%%EOF\n", earlyXRef)
		}
		fileLength = b.Len()
		result = b.Bytes()
		if opt.review8 == "startxref-prefix" {
			start := bytes.LastIndex(result, []byte("startxref\n"))
			result[start-1] = '%'
		}
		if opt.review8 == "startxref-junk" || opt.review8 == "startxref-space" {
			start := bytes.LastIndex(result, []byte("startxref\n"))
			if start < 0 {
				t.Fatal("missing startxref")
			}
			if opt.review8 == "startxref-junk" {
				result[start+9] = 'X'
			} else {
				result[start+9] = '\v'
			}
		}
		if opt.review8 == "epilogue-space" {
			start := bytes.LastIndex(result, []byte("\n%%EOF"))
			result[start] = '\v'
		}
		if opt.unreliableXRef == "inline" {
			result = bytes.ReplaceAll(result, []byte("xref\n0 "), []byte("xref 0 "))
		}
	}
	return result, clean, content
}

// Test-only AESV2/AESV3 encoding, preserving per-object key derivation.
func encryptPrimaryHintBytes(t *testing.T, raw []byte, nr int, ctx *model.Context) []byte {
	t.Helper()
	key := ctx.EncKey
	if ctx.E.R < 5 {
		bb := append(bytes.Clone(key), byte(nr), byte(nr>>8), byte(nr>>16), 0, 0)
		bb = append(bb, []byte("sAlT")...)
		sum := md5.Sum(bb)
		key = sum[:min(len(key)+5, 16)]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	pad := aes.BlockSize - len(raw)%aes.BlockSize
	data := append(bytes.Clone(raw), bytes.Repeat([]byte{byte(pad)}, pad)...)
	result := make([]byte, aes.BlockSize+len(data))
	for i := range aes.BlockSize {
		result[i] = byte(i)
	}
	cipher.NewCBCEncrypter(block, result[:aes.BlockSize]).CryptBlocks(result[aes.BlockSize:], data)
	return result
}

func TestHintReviewAmbiguousSyntaxDisqualifiesAuthorization(t *testing.T) {
	for _, name := range []string{"H", "Length", "S", "fallback", "ordinary", "ordinary-fallback", "trailer", "trailer-fallback", "compressed", "compressed-fallback"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted-%t", name, encrypted), func(t *testing.T) {
				opt := primaryHintFixtureOptions{encrypted: encrypted}
				if name == "H" || name == "Length" || name == "S" || name == "fallback" {
					opt.ambiguousClaim = name
				} else {
					opt.ambiguousOrigin = name
				}
				opt.objectStream = strings.HasPrefix(name, "compressed")
				pdf, _, content := plaintextPrimaryHintPDF(t, opt)
				ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
				if !encrypted {
					if err == nil {
						t.Fatal("ambiguous syntax authorized plaintext hint")
					}
					return
				}
				if err != nil {
					t.Fatalf("normal encrypted control: %v", err)
				}
				if len(ctx.LinearizationObjs) != 0 {
					t.Fatal("ambiguous syntax authorized deletion")
				}
				assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
			})
		}
	}
}

func TestHintReviewRebuildInvalidatesDeferredCandidates(t *testing.T) {
	for _, overflow := range []bool{false, true} {
		t.Run(fmt.Sprintf("overflow-%t", overflow), func(t *testing.T) {
			pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{encrypted: true, overflow: overflow, corruptOffset: true})
			ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
			if err != nil {
				t.Fatal(err)
			}
			if !ctx.Read.RepairedXRef || len(ctx.LinearizationObjs) != 0 {
				t.Fatal("rebuild retained hint authorization")
			}
			assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
		})
	}
}

func TestPrimaryHintVariants(t *testing.T) {
	for _, tc := range []struct {
		name string
		opt  primaryHintFixtureOptions
	}{
		{"overflow", primaryHintFixtureOptions{overflow: true}},
		{"aligned", primaryHintFixtureOptions{aligned: true}},
		{"encrypted", primaryHintFixtureOptions{encrypted: true}},
		{"encrypted-overflow", primaryHintFixtureOptions{encrypted: true, overflow: true}},
		{"aes256", primaryHintFixtureOptions{keyLength: 256}},
		{"hint-number-below-linearization", primaryHintFixtureOptions{lowHint: true}},
		{"xref-object-stream", primaryHintFixtureOptions{objectStream: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pdf, _, content := plaintextPrimaryHintPDF(t, tc.opt)
			conf := encryptedXRefReadConfig("user")
			ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf)
			if err != nil {
				t.Fatal(err)
			}
			if !ctx.Read.Linearized || len(ctx.LinearizationObjs) != 2+boolInt(tc.opt.overflow) {
				t.Fatal("missing verified linearization identities")
			}
			if ctx.E.P != int(model.PermissionsNone)-65536 {
				t.Fatalf("permissions changed: %d", ctx.E.P)
			}
			assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
			if _, err := PDFInfo(t.Context(), bytes.NewReader(pdf), "synthetic.pdf", nil, false, conf); err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			if err := Decrypt(t.Context(), bytes.NewReader(pdf), &out, conf); err != nil {
				t.Fatal(err)
			}
			plain, err := ReadContext(t.Context(), bytes.NewReader(out.Bytes()), model.NewDefaultConfiguration())
			if err != nil {
				t.Fatal(err)
			}
			if err := ValidateContext(t.Context(), plain); err != nil {
				t.Fatal(err)
			}
			if plain.Encrypt != nil || plain.Read.Linearized || plain.PageCount != 1 {
				t.Fatal("decrypted output retained encryption/linearization or lost pages")
			}
			assertEncryptedXRefState(t, readEncryptedXRefState(t, out.Bytes(), ""), content, false)
			if tc.opt.encrypted {
				strict := encryptedXRefReadConfig("user")
				strict.ValidationMode = model.ValidationStrict
				if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), strict); err != nil {
					t.Fatalf("strict encrypted hint: %v", err)
				}
			}
		})
	}
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func TestUnencryptedPrimaryHint(t *testing.T) {
	pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{unencrypted: true, overflow: true})
	for _, mode := range []int{model.ValidationStrict, model.ValidationRelaxed} {
		conf := model.NewDefaultConfiguration()
		conf.ValidationMode = mode
		ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf)
		if err != nil {
			t.Fatal(err)
		}
		if !ctx.Read.Linearized || ctx.Encrypt != nil {
			t.Fatal("unencrypted linearization read changed")
		}
		assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, ""), content, false)
		var out bytes.Buffer
		if err := Optimize(t.Context(), bytes.NewReader(pdf), &out, conf, nil); err != nil {
			t.Fatal(err)
		}
		plain, err := ReadContext(t.Context(), bytes.NewReader(out.Bytes()), model.NewDefaultConfiguration())
		if err != nil {
			t.Fatal(err)
		}
		if plain.Read.Linearized {
			t.Fatal("full output retained stale linearization")
		}
		assertEncryptedXRefState(t, readEncryptedXRefState(t, out.Bytes(), ""), content, false)
	}
}

func TestPrimaryHintRejectsReferencedAndStrict(t *testing.T) {
	pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{referenced: true})
	if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user")); err == nil {
		t.Fatal("referenced hint received plaintext repair")
	}
	pdf, _, _ = plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{objectStream: true, compressedReference: true})
	if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user")); err == nil {
		t.Fatal("compressed origin reference received plaintext repair")
	}
	pdf, _, _ = plaintextPrimaryHintPDF(t)
	strict := encryptedXRefReadConfig("user")
	strict.ValidationMode = model.ValidationStrict
	if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), strict); err == nil {
		t.Fatal("strict plaintext hint accepted")
	}
}

func replacePrimaryHintField(t *testing.T, pdf []byte, key, value string) []byte {
	t.Helper()
	re := regexp.MustCompile(`/` + key + ` (\[[^\]]*\]|[0-9]+)`)
	match := re.FindSubmatchIndex(pdf)
	if len(match) != 4 {
		t.Fatalf("missing %s", key)
	}
	start, end := match[2], match[3]
	if len(value) > end-start {
		t.Fatalf("replacement %s too wide", key)
	}
	out := bytes.Clone(pdf)
	copy(out[start:end], bytes.Repeat([]byte{' '}, end-start))
	copy(out[start:end], value)
	return out
}

func TestPrimaryHintBadClaimsAndCompression(t *testing.T) {
	pdf, _, _ := plaintextPrimaryHintPDF(t)
	ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, key, value string }{
		{"negative-offset", "H", "[-000000001 0000000200]"},
		{"negative-length", "H", fmt.Sprintf("[%010d -000000001]", *ctx.OffsetPrimaryHintTable)},
		{"out-of-range", "H", fmt.Sprintf("[%010d 0000000200]", len(pdf)+100)},
		{"indirect-offset", "H", "[8 0 R 0000000200]"},
		{"indirect-length", "H", fmt.Sprintf("[%010d 8 0 R]", *ctx.OffsetPrimaryHintTable)},
		{"stale-length", "L", fmt.Sprintf("%010d", len(pdf)-1)},
		{"bad-first-page-end", "E", "0000000001"},
		{"bad-main-xref-bound", "T", "0000000001"},
		{"bad-first-page-object", "O", "0"},
		{"bad-page-count", "N", "0"},
		{"wrong-type", "H", fmt.Sprintf("[%010d 0000000020]", *ctx.Table[1].Offset)},
		{"envelope", "H", fmt.Sprintf("[%010d 0000000001]", *ctx.OffsetPrimaryHintTable)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bad := replacePrimaryHintField(t, pdf, tc.key, tc.value)
			if _, err := ReadContext(t.Context(), bytes.NewReader(bad), encryptedXRefReadConfig("user")); err == nil {
				t.Fatal("bad hint claim authorized plaintext repair")
			}
		})
	}
	t.Run("duplicate", func(t *testing.T) {
		pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{overflow: true})
		parts := regexp.MustCompile(`/H \[([^\]]*)\]`).FindSubmatch(pdf)
		values := strings.Fields(string(parts[1]))
		bad := replacePrimaryHintField(t, pdf, "H", fmt.Sprintf("[%s %s %s %s]", values[0], values[1], values[0], values[1]))
		if _, err := ReadContext(t.Context(), bytes.NewReader(bad), encryptedXRefReadConfig("user")); err == nil {
			t.Fatal("duplicate hints repaired")
		}
	})
	t.Run("corrupt-compression", func(t *testing.T) {
		bad := bytes.Clone(pdf)
		i := bytes.Index(bad, []byte("/S 36 >>\nstream\n")) + len("/S 36 >>\nstream\n")
		bad[i] ^= 0xff
		if _, err := ReadContext(t.Context(), bytes.NewReader(bad), encryptedXRefReadConfig("user")); err == nil {
			t.Fatal("corrupt hint accepted")
		}
	})
}

func TestPrimaryHintForgedContentPreservedByWriter(t *testing.T) {
	pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{encrypted: true})
	ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
	if err != nil {
		t.Fatal(err)
	}
	start := int(*ctx.Table[4].Offset)
	end := bytes.Index(pdf[start:], []byte("endobj")) + start + len("endobj")
	forged := replacePrimaryHintField(t, pdf, "H", fmt.Sprintf("[%010d %010d]", start, end-start))
	ctx, err = ReadContext(t.Context(), bytes.NewReader(forged), encryptedXRefReadConfig("user"))
	if err != nil {
		t.Fatal(err)
	}
	if ctx.IsLinearizationObject(4) {
		t.Fatal("ordinary content marked disposable")
	}
	var out bytes.Buffer
	if err := Decrypt(t.Context(), bytes.NewReader(forged), &out, encryptedXRefReadConfig("user")); err != nil {
		t.Fatal(err)
	}
	assertEncryptedXRefState(t, readEncryptedXRefState(t, out.Bytes(), ""), content, false)
}

func TestPrimaryHintLimits(t *testing.T) {
	pdf, _, _ := plaintextPrimaryHintPDF(t)
	conf := encryptedXRefReadConfig("user")
	conf.Limits.MaxDecodeBytes = 100
	if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf); !errors.Is(err, filter.ErrDecodeLimitExceeded) {
		t.Fatalf("decode limit: %v", err)
	}
	conf = encryptedXRefReadConfig("user")
	conf.Limits.MaxStreamBytes = 100
	if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf); err == nil {
		t.Fatal("stream limit ignored")
	}
	c, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := ReadContext(c, bytes.NewReader(pdf), encryptedXRefReadConfig("user")); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel: %v", err)
	}
}

func TestPlaintextPrimaryHintControls(t *testing.T) {
	pdf, clean, content := plaintextPrimaryHintPDF(t)
	assertEncryptedXRefState(t, readEncryptedXRefState(t, clean, "user"), content, true)
	for _, operation := range []string{"read", "info", "decrypt"} {
		t.Run("wrong-password-"+operation, func(t *testing.T) {
			conf := encryptedXRefReadConfig("wrong")
			var err error
			switch operation {
			case "read":
				_, err = ReadContext(t.Context(), bytes.NewReader(pdf), conf)
			case "info":
				_, err = PDFInfo(t.Context(), bytes.NewReader(pdf), "synthetic.pdf", nil, false, conf)
			case "decrypt":
				var out bytes.Buffer
				err = Decrypt(t.Context(), bytes.NewReader(pdf), &out, conf)
			}
			if !errors.Is(err, pdfcpu.ErrWrongPassword) {
				t.Fatalf("got %v, want ErrWrongPassword", err)
			}
		})
	}
}

// Originally RED: relaxed compatibility must not alter ordinary encrypted content.
func TestPlaintextPrimaryHintRelaxed(t *testing.T) {
	pdf, _, content := plaintextPrimaryHintPDF(t)
	t.Run("metadata", func(t *testing.T) {
		info, err := PDFInfo(t.Context(), bytes.NewReader(pdf), "synthetic.pdf", nil, false, encryptedXRefReadConfig("user"))
		if err != nil {
			t.Fatalf("relaxed metadata with plaintext primary hint: %v", err)
		}
		if info.Title != encryptedXRefTitle {
			t.Fatalf("Title = %q", info.Title)
		}
	})
	t.Run("read", func(t *testing.T) {
		assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
	})
	t.Run("decrypt", func(t *testing.T) {
		var out bytes.Buffer
		if err := Decrypt(t.Context(), bytes.NewReader(pdf), &out, encryptedXRefReadConfig("user")); err != nil {
			t.Fatalf("relaxed decrypt with plaintext primary hint: %v", err)
		}
		assertEncryptedXRefState(t, readEncryptedXRefState(t, out.Bytes(), ""), content, false)
	})
}
