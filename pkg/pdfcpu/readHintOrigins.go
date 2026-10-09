package pdfcpu

import (
	"context"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

// hintFramingToken never lets a clipped comment manufacture a new origin.
func hintFramingToken(s *string) string {
	for {
		*s = strings.TrimLeftFunc(*s, contextutil.PDFWhitespace)
		if *s == "" || (*s)[0] != '%' || strings.HasPrefix(*s, "%%EOF") {
			break
		}
		i := strings.IndexAny(*s, "\r\n")
		if i < 0 {
			return ""
		}
		*s = (*s)[i:]
	}
	i := strings.IndexFunc(*s, func(r rune) bool {
		return contextutil.PDFWhitespace(r) || (!strings.HasPrefix(*s, "%%EOF") && strings.ContainsRune("()<>[]/%", r))
	})
	if i < 0 {
		i = len(*s)
	}
	token := (*s)[:i]
	*s = (*s)[i:]
	return token
}

func hintPhysicalEnd(c context.Context, ctx *model.Context, offset int64, nr, generation int) (int64, bool, error) {
	p := &contextutil.ParseProvenance{SuppressNotices: true}
	c = contextutil.WithParseProvenance(c, p)
	o, end, stream, streamOffset, err := object(c, ctx, offset, nr, generation, model.ValidationStrict)
	if err != nil || p.Ambiguous() {
		return 0, false, err
	}
	if stream < 0 || (end >= 0 && end < stream) {
		return offset + int64(end) + int64(len("endobj")), end >= 0, nil
	}
	d, ok := o.(types.Dict)
	if !ok {
		return 0, false, nil
	}
	// Length() retains only the object number. Physical proof must retain
	// the reference's exact live generation, unlike ordinary loading.
	if ir, ok := d["Length"].(types.IndirectRef); ok {
		e := ctx.Table[ir.ObjectNumber.Value()]
		if e == nil || e.Free || e.Generation == nil || *e.Generation != ir.GenerationNumber.Value() ||
			(e.Compressed && ir.GenerationNumber.Value() != 0) {
			return 0, false, nil
		}
	}
	length, lengthRef := d.Length()
	sd := types.NewStreamDict(d, offset+streamOffset, length, lengthRef, nil)
	if length == nil && lengthRef == nil {
		return 0, false, nil
	}
	if err := ensureIndirectStreamLength(c, ctx, &sd, false); err != nil {
		return 0, false, err
	}
	if p.Ambiguous() || sd.StreamLength == nil || *sd.StreamLength < 0 ||
		sd.StreamOffset < offset || sd.StreamOffset > ctx.Read.FileSize ||
		*sd.StreamLength > ctx.Read.FileSize-sd.StreamOffset {
		return 0, false, nil
	}
	start := sd.StreamOffset + *sd.StreamLength
	rd, err := newPositionedReader(ctx.Read.RS, &start)
	if err != nil {
		return 0, false, err
	}
	bb, err := io.ReadAll(io.LimitReader(rd, min(int64(1025), objectBufferLimit(ctx))))
	if err != nil {
		return 0, false, err
	}
	s := string(bb)
	for _, expected := range []string{"endstream", "endobj"} {
		if hintFramingToken(&s) != expected {
			return 0, false, nil
		}
		if s == "" && start+int64(len(bb)) != ctx.Read.FileSize {
			return 0, false, nil
		}
	}
	return start + int64(len(bb)-len(s)), !p.Ambiguous(), c.Err()
}

// hintPhysicalGap accepts only separators or a complete xref/trailer epilogue.
// Bytes after EOF still belong to this gap and cannot conceal structural values.
func hintPhysicalGap(c context.Context, ctx *model.Context, start, end int64, xrefs map[int64]bool, early bool) (bool, error) {
	if end < start || end-start > objectBufferLimit(ctx) {
		return false, nil
	}
	rd, err := newPositionedReader(ctx.Read.RS, &start)
	if err != nil {
		return false, err
	}
	bb, err := io.ReadAll(io.LimitReader(rd, end-start))
	if err != nil || int64(len(bb)) != end-start {
		return false, err
	}
	s := string(bb)
	token := hintFramingToken(&s)
	if token == "" {
		return s == "" && end < ctx.Read.FileSize, c.Err()
	}
	allowZero := early && token == "xref"
	if token == "xref" {
		xrefOffset := start + int64(len(bb)-len(s)-len(token))
		if !contextutil.ParseProvenanceFromContext(c).HasXRefOrigin(xrefOffset) {
			return false, nil
		}
		for {
			if err := c.Err(); err != nil {
				return false, err
			}
			first := hintFramingToken(&s)
			if first == "trailer" {
				break
			}
			_, err := strconv.ParseUint(first, 10, 64)
			count, countErr := strconv.ParseUint(hintFramingToken(&s), 10, 64)
			if err != nil || countErr != nil || count == 0 || count > uint64(len(s))/5 {
				return false, nil
			}
			for range count {
				if err := c.Err(); err != nil {
					return false, err
				}
				if _, err := strconv.ParseUint(hintFramingToken(&s), 10, 64); err != nil {
					return false, nil
				}
				if _, err := strconv.ParseUint(hintFramingToken(&s), 10, 64); err != nil {
					return false, nil
				}
				status := hintFramingToken(&s)
				if status != "n" && status != "f" {
					return false, nil
				}
			}
		}
		p := &contextutil.ParseProvenance{SuppressNotices: true}
		probe := contextutil.WithParseProvenance(c, p)
		result, err := model.ParseObjectWithPolicy(probe, &s, 0, model.ValidationStrict, recursionLimit(ctx))
		if err != nil || p.Ambiguous() {
			return false, err
		}
		if _, ok := result.Object.(types.Dict); !ok {
			return false, nil
		}
		xrefs[xrefOffset] = true
		token = hintFramingToken(&s)
	}
	if token != "startxref" {
		return false, nil
	}
	xref, err := strconv.ParseInt(hintFramingToken(&s), 10, 64)
	if err != nil || (!xrefs[xref] && !(allowZero && xref == 0)) {
		return false, nil
	}
	if hintFramingToken(&s) != "%%EOF" {
		return false, nil
	}
	lineEnd := strings.IndexAny(s, "\r\n")
	if lineEnd < 0 {
		return end == ctx.Read.FileSize && strings.TrimFunc(s, contextutil.PDFWhitespace) == "", c.Err()
	}
	if strings.TrimFunc(s[:lineEnd], contextutil.PDFWhitespace) != "" {
		return false, nil
	}
	s = s[lineEnd:]
	return hintFramingToken(&s) == "" && s == "", c.Err()
}

// proveHintOrigins uses xref offsets only as claims; every predecessor envelope
// and intervening gap must prove the next claimed physical boundary.
func proveHintOrigins(c context.Context, ctx *model.Context, first int64) (bool, error) {
	var numbers []int
	for nr, e := range ctx.Table {
		if err := c.Err(); err != nil {
			return false, err
		}
		if e != nil && !e.Free && !e.Compressed && e.Offset != nil {
			numbers = append(numbers, nr)
		}
	}
	slices.SortFunc(numbers, func(a, b int) int {
		x, y := *ctx.Table[a].Offset, *ctx.Table[b].Offset
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	if len(numbers) == 0 || *ctx.Table[numbers[0]].Offset != first {
		return false, nil
	}
	xrefs := map[int64]bool{}
	for i, nr := range numbers {
		e := ctx.Table[nr]
		offset := *e.Offset
		next := ctx.Read.FileSize
		if i+1 < len(numbers) {
			next = *ctx.Table[numbers[i+1]].Offset
		}
		if e.Generation == nil || offset < first || next <= offset || next > ctx.Read.FileSize {
			return false, nil
		}
		end, ok, err := hintPhysicalEnd(c, ctx, offset, nr, *e.Generation)
		if err != nil || !ok || end > next {
			return false, err
		}
		// Xref streams have an ordinary object envelope, not a classic table.
		if _, isXRef := e.Object.(types.XRefStreamDict); isXRef && contextutil.ParseProvenanceFromContext(c).HasXRefOrigin(offset) {
			xrefs[offset] = true
		}
		ok, err = hintPhysicalGap(c, ctx, end, next, xrefs, i == 0 && next < ctx.Read.FileSize)
		if err != nil || !ok {
			return false, err
		}
	}
	return true, c.Err()
}
