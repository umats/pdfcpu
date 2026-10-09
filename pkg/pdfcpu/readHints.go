/*
Copyright 2026 The pdfcpu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pdfcpu

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/filter"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

// hintReadState is a per-read, structural candidate inventory. It grants no
// plaintext exemption until all incoming references have been audited.
// This compatibility exception is NOT a normative encryption exemption.
type hintReadState struct {
	linearNr   int
	streams    map[int]types.StreamDict
	provenance *contextutil.ParseProvenance
}

func hintCandidate(states []*hintReadState, nr int) bool {
	if len(states) == 0 || states[0] == nil {
		return false
	}
	_, ok := states[0].streams[nr]
	return ok
}

func terminalHintError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, filter.ErrDecodeLimitExceeded) || errors.Is(err, errObjectBufferLimit) ||
		errors.Is(err, model.ErrMaxRecursionDepthExceeded) || errors.Is(err, model.ErrInputSizeLimit)
}

// firstPhysicalObject never scans past the first 1KiB or skips a non-comment
// object. A later /Linearized marker cannot authorize a plaintext probe.
func firstPhysicalObject(c context.Context, ctx *model.Context) (int64, error) {
	if err := c.Err(); err != nil {
		return 0, err
	}
	zero := int64(0)
	rd, err := newPositionedReader(ctx.Read.RS, &zero)
	if err != nil {
		return 0, err
	}
	bb, err := io.ReadAll(io.LimitReader(rd, 1024))
	if err != nil {
		return 0, err
	}
	if !bytes.HasPrefix(bb, []byte("%PDF-")) {
		return 0, nil
	}
	for i := 0; i < len(bb); {
		if strings.ContainsRune("\x00\t\n\f\r ", rune(bb[i])) {
			i++
			continue
		}
		if bb[i] != '%' {
			return int64(i), nil
		}
		for i < len(bb) && bb[i] != '\n' && bb[i] != '\r' {
			i++
		}
	}
	return 0, nil
}

func uniqueHintEntry(c context.Context, ctx *model.Context, offset int64) (int, *model.XRefTableEntry, error) {
	nr := -1
	var found *model.XRefTableEntry
	for n, e := range ctx.Table {
		if err := c.Err(); err != nil {
			return -1, nil, err
		}
		if e == nil || e.Free || e.Compressed || e.Offset == nil || *e.Offset != offset {
			continue
		}
		if found != nil {
			return -1, nil, nil
		}
		nr, found = n, e
	}
	return nr, found, nil
}

func exactHintIdentity(c context.Context, ctx *model.Context, offset int64, nr int, e *model.XRefTableEntry) (bool, error) {
	if err := c.Err(); err != nil {
		return false, err
	}
	if e.Generation == nil {
		return false, nil
	}
	rd, err := newPositionedReader(ctx.Read.RS, &offset)
	if err != nil {
		return false, err
	}
	bb, err := io.ReadAll(io.LimitReader(rd, 128))
	if err != nil {
		return false, err
	}
	return literalObjectHeader(bb, nr, *e.Generation), nil
}

func literalObjectHeader(bb []byte, nr, generation int) bool {
	// Authorization requires a literal header at the claimed byte start,
	// not ParseObjectAttributes' junk stripping or shortened "bj" repair.
	space := func(b byte) bool { return strings.ContainsRune("\x00\t\n\f\r ", rune(b)) }
	integer := func() (int, bool) {
		i := 0
		for i < len(bb) && bb[i] >= '0' && bb[i] <= '9' {
			i++
		}
		if i == 0 || i == len(bb) || !space(bb[i]) {
			return 0, false
		}
		n, err := strconv.Atoi(string(bb[:i]))
		bb = bb[i:]
		for len(bb) > 0 && space(bb[0]) {
			bb = bb[1:]
		}
		return n, err == nil
	}
	obj, ok := integer()
	if !ok || obj != nr {
		return false
	}
	gen, ok := integer()
	if !ok || gen != generation || len(bb) < 4 || !bytes.HasPrefix(bb, []byte("obj")) {
		return false
	}
	return space(bb[3]) || strings.ContainsRune("()<>[]{}/%", rune(bb[3]))
}

func directHintInteger(d types.Dict, key string) (int64, bool) {
	n, ok := d[key].(types.Integer)
	return int64(n), ok
}

func validLinearizationHintFields(ctx *model.Context, d types.Dict) bool {
	if !d.IsLinearizationParmDict() || validateLinearizationDirectEntries(d, 0) != nil {
		return false
	}
	l, lok := directHintInteger(d, "L")
	o, ook := directHintInteger(d, "O")
	e, eok := directHintInteger(d, "E")
	n, nok := directHintInteger(d, "N")
	t, tok := directHintInteger(d, "T")
	pageEntry := ctx.Table[int(o)]
	return lok && ook && eok && nok && tok && l == ctx.Read.FileSize &&
		o > 0 && n > 0 && e > 0 && e < l && e <= t && t < l && pageEntry != nil && !pageEntry.Free
}

func supportedHintDict(ctx *model.Context, d types.Dict, primary bool) bool {
	// /Length is direct here to avoid repaired or ambiguous stream boundaries.
	l, ok := directHintInteger(d, "Length")
	if !ok || l <= 0 || validateHintStreamDict(ctx, d, 0) != nil {
		return false
	}
	for key := range d {
		switch key {
		case "Length", "Filter", "DecodeParms", "S", "O":
		default:
			return false
		}
	}
	if primary && d["S"] == nil {
		return false
	}
	for _, key := range []string{"S", "O"} {
		if d[key] != nil {
			value, ok := directHintInteger(d, key)
			if !ok || value < 0 {
				return false
			}
		}
	}
	return true
}

func hintEnvelope(c context.Context, ctx *model.Context, sd *types.StreamDict, offset, span int64) (bool, error) {
	l := *sd.StreamLength
	end := offset + span // Already checked with subtraction-based file bounds.
	if sd.StreamOffset <= offset || sd.StreamOffset > end || l > end-sd.StreamOffset {
		return false, nil
	}
	tailStart := sd.StreamOffset + l
	if end-tailStart > 1024 {
		return false, nil
	}
	for _, entry := range ctx.Table {
		if err := c.Err(); err != nil {
			return false, err
		}
		if entry != nil && !entry.Free && entry.Offset != nil && *entry.Offset > offset && *entry.Offset < end {
			return false, nil
		}
	}
	if err := c.Err(); err != nil {
		return false, err
	}
	rd, err := newPositionedReader(ctx.Read.RS, &tailStart)
	if err != nil {
		return false, err
	}
	bb, err := io.ReadAll(io.LimitReader(rd, end-tailStart+1))
	if err != nil {
		return false, err
	}
	claimed := end - tailStart
	if int64(len(bb)) < claimed {
		return false, nil
	}
	// A clipped buffer EOF is not a physical token boundary. Check the
	// actual following byte when the claim ends on the terminator.
	if claimed > 0 && !strings.ContainsRune("\x00\t\n\f\r ", rune(bb[claimed-1])) && int64(len(bb)) > claimed &&
		!strings.ContainsRune("\x00\t\n\f\r <>[]()/%", rune(bb[claimed])) {
		return false, nil
	}
	fields := strings.FieldsFunc(string(bb[:claimed]), func(r rune) bool { return strings.ContainsRune("\x00\t\n\f\r ", r) })
	return len(fields) == 2 && fields[0] == "endstream" && fields[1] == "endobj", nil
}

func probeHintObject(c context.Context, ctx *model.Context, offset int64, nr, generation int) (types.Object, int, int, int64, error) {
	provenance := &contextutil.ParseProvenance{SuppressNotices: true}
	probeContext := contextutil.WithParseProvenance(c, provenance)
	o, end, stream, streamOffset, err := object(probeContext, ctx, offset, nr, generation, model.ValidationStrict)
	if err == nil && provenance.Ambiguous() {
		// Return no eligible object. The ordinary parser will preserve its
		// historical result and emit its own notices exactly once.
		o = nil
	}
	return o, end, stream, streamOffset, err
}

// canonicalHintFilter normalizes only eligible Flate filters on a copy, avoiding
// repair notices from inspection before ordinary loading normalizes the origin.
func canonicalHintFilter(ctx *model.Context, d types.Dict) (types.Dict, bool) {
	out := d.Clone().(types.Dict)
	switch f := d["Filter"].(type) {
	case types.Name:
		if streamFilterName(ctx, string(f)) != filter.Flate {
			return nil, false
		}
		out["Filter"] = types.Name(filter.Flate)
	case types.Array:
		if len(f) != 1 {
			return nil, false
		}
		name, ok := f[0].(types.Name)
		if !ok || streamFilterName(ctx, string(name)) != filter.Flate {
			return nil, false
		}
		out["Filter"] = types.Array{types.Name(filter.Flate)}
	default:
		return nil, false
	}
	return out, true
}

func discoverHintStreams(c context.Context, ctx *model.Context) (*hintReadState, error) {
	if ctx.Read.RepairedXRef || ctx.Read.RepairOffset != 0 {
		return nil, nil
	}
	offset, err := firstPhysicalObject(c, ctx)
	if err != nil || offset == 0 {
		return nil, err
	}
	nr, entry, err := uniqueHintEntry(c, ctx, offset)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	ok, err := exactHintIdentity(c, ctx, offset, nr, entry)
	if err != nil || !ok {
		return nil, err
	}
	o, endInd, streamInd, _, err := probeHintObject(c, ctx, offset, nr, *entry.Generation)
	if err != nil {
		if terminalHintError(err) {
			return nil, err
		}
		return nil, nil
	}
	d, ok := o.(types.Dict)
	if !ok || endInd < 0 || (streamInd >= 0 && streamInd < endInd) || !validLinearizationHintFields(ctx, d) {
		return nil, nil
	}
	h, ok := d["H"].(types.Array)
	if !ok || (len(h) != 2 && len(h) != 4) {
		return nil, nil
	}
	state := &hintReadState{linearNr: nr, streams: map[int]types.StreamDict{}, provenance: contextutil.ParseProvenanceFromContext(c)}
	var previousEnd int64
	for i := 0; i < len(h); i += 2 {
		off, ook := h[i].(types.Integer)
		length, lok := h[i+1].(types.Integer)
		start, span := int64(off), int64(length)
		if !ook || !lok || start <= offset || span <= 0 || start >= ctx.Read.FileSize ||
			span > ctx.Read.FileSize-start || (i > 0 && start < previousEnd) {
			return nil, nil
		}
		previousEnd = start + span
		if i == 0 {
			firstPageEnd, _ := directHintInteger(d, "E")
			if previousEnd > firstPageEnd {
				return nil, nil
			}
		}
		nr, e, err := uniqueHintEntry(c, ctx, start)
		if err != nil {
			return nil, err
		}
		if e == nil || nr == state.linearNr {
			return nil, nil
		}
		ok, err := exactHintIdentity(c, ctx, start, nr, e)
		if err != nil || !ok {
			return nil, err
		}
		o, _, streamInd, streamOffset, err := probeHintObject(c, ctx, start, nr, *e.Generation)
		if err != nil {
			if terminalHintError(err) {
				return nil, err
			}
			return nil, nil
		}
		d, ok := o.(types.Dict)
		if !ok || streamInd < 0 || !supportedHintDict(ctx, d, i == 0) {
			return nil, nil
		}
		d, ok = canonicalHintFilter(ctx, d)
		if !ok {
			return nil, nil
		}
		sd, err := streamDictForObject(c, ctx, d, nr, streamInd, streamOffset, start)
		if err != nil {
			return nil, err
		}
		if len(sd.FilterPipeline) != 1 || sd.FilterPipeline[0].Name != "FlateDecode" {
			return nil, nil
		}
		ok, err = hintEnvelope(c, ctx, &sd, start, span)
		if err != nil || !ok {
			return nil, err
		}
		state.streams[nr] = sd
	}
	return state, nil
}

// hintReferences walks every live structural origin, including lazy objects
// whose incoming edges are invisible to ordinary reference-count traversal.
func hintReferences(c context.Context, ctx *model.Context, state *hintReadState) (bool, error) {
	// Lazy origins retain read-local provenance, but inspection must not
	// repeat ordinary duplicate-key notices before later validation.
	c = contextutil.WithoutParseNotices(c)
	referenced := false
	linearizationMarkers := 0
	var walk func(types.Object, int) error
	walk = func(o types.Object, depth int) error {
		if err := c.Err(); err != nil {
			return err
		}
		if err := ctx.CheckRecursionDepth("hint reference audit", depth); err != nil {
			return err
		}
		switch o := o.(type) {
		case types.IndirectRef:
			nr := o.ObjectNumber.Value()
			if _, ok := state.streams[nr]; ok || nr == state.linearNr {
				referenced = true
			}
		case types.Dict:
			if o.IsLinearizationParmDict() {
				linearizationMarkers++
				if linearizationMarkers > 1 {
					referenced = true
				}
			}
			for _, key := range slices.Sorted(maps.Keys(o)) {
				if err := walk(o[key], depth+1); err != nil {
					return err
				}
			}
		case types.Array:
			for _, value := range o {
				if err := walk(value, depth+1); err != nil {
					return err
				}
			}
		case types.StreamDict:
			return walk(o.Dict, depth)
		case types.ObjectStreamDict:
			if err := walk(o.Dict, depth); err != nil {
				return err
			}
			// Members are independent indirect object origins, not nested
			// values of the containing stream dictionary.
			for _, value := range o.ObjArray {
				if err := walk(value, 0); err != nil {
					return err
				}
			}
		case types.XRefStreamDict:
			return walk(o.Dict, depth)
		case types.LazyObjectStreamObject:
			value, err := o.DecodedObject(c)
			if err != nil {
				return err
			}
			return walk(value, depth)
		}
		return nil
	}
	for _, nr := range slices.Sorted(maps.Keys(ctx.Table)) {
		if err := c.Err(); err != nil {
			return false, err
		}
		entry := ctx.Table[nr]
		if entry == nil || entry.Free {
			continue
		}
		if entry.Object == nil {
			// A skipped/unparsed origin cannot prove the absence of references.
			return true, nil
		}
		if err := walk(entry.Object, 0); err != nil {
			return false, err
		}
	}
	for _, d := range ctx.Read.TrailerDicts {
		if err := walk(d, 0); err != nil {
			return false, err
		}
	}
	for _, ref := range []*types.IndirectRef{ctx.Root, ctx.Info, ctx.Encrypt} {
		if ref != nil {
			if err := walk(*ref, 0); err != nil {
				return false, err
			}
		}
	}
	return referenced, nil
}

func decodeHintProbe(c context.Context, ctx *model.Context, sd *types.StreamDict) error {
	if err := c.Err(); err != nil {
		return err
	}
	if err := sd.DecodeWithLimit(decodeLimit(ctx)); err != nil {
		return err
	}
	if err := c.Err(); err != nil {
		return err
	}
	if s, ok := directHintInteger(sd.Dict, "S"); ok && s > int64(len(sd.Content)) {
		return fmt.Errorf("hint shared-object table offset exceeds decoded stream")
	}
	return nil
}

func finalizeHintStreams(c context.Context, ctx *model.Context, state *hintReadState) error {
	referenced, err := hintReferences(c, ctx, state)
	if terminalHintError(err) {
		return fmt.Errorf("audit hint references: %w", err)
	}
	// An unreadable structural origin cannot authorize compatibility, but
	// must not add syntax validation to the existing normal cipher path.
	verified := err == nil && !referenced && !state.provenance.Ambiguous() && !ctx.Read.RepairedXRef && ctx.Read.RepairOffset == 0
	if verified {
		first, proofErr := firstPhysicalObject(c, ctx)
		verified = proofErr == nil && first != 0
		if verified {
			verified, proofErr = proveHintOrigins(c, ctx, first)
		}
		if terminalHintError(proofErr) {
			return fmt.Errorf("audit physical hint origins: %w", proofErr)
		}
		verified = verified && proofErr == nil
	}
	return finalizeHintStreamContents(c, ctx, state, verified)
}

func finalizeHintStreamContents(c context.Context, ctx *model.Context, state *hintReadState, verified bool) error {
	for _, nr := range slices.Sorted(maps.Keys(state.streams)) {
		entry := ctx.Table[nr]
		sd, ok := entry.Object.(types.StreamDict)
		if !ok {
			return fmt.Errorf("hint stream obj#%d changed type", nr)
		}
		if err := c.Err(); err != nil {
			return err
		}
		// Decryption is in-place. Both trial dictionaries and byte slices must
		// be independent so unsuccessful normal processing cannot destroy Raw.
		probe := sd.Clone().(types.StreamDict)
		probe.Raw = bytes.Clone(sd.Raw)
		probe.Content = nil
		err := saveDecodedStreamContent(ctx, &probe, nr, *entry.Generation, ctx.DecodeAllStreams)
		if cancelErr := c.Err(); cancelErr != nil {
			return cancelErr
		}
		repaired := false
		repair := verified && ctx.XRefTable.ValidationMode == model.ValidationRelaxed && ctx.EncKey != nil && ctx.AES4Streams
		if err == nil && repair {
			err = decodeHintProbe(c, ctx, &probe)
		}
		if err != nil && repair && !terminalHintError(err) {
			plain := sd.Clone().(types.StreamDict)
			plain.Raw = bytes.Clone(sd.Raw)
			plain.Content = nil
			if plainErr := decodeHintProbe(c, ctx, &plain); plainErr == nil {
				probe, err = plain, nil
				repaired = true
			} else if terminalHintError(plainErr) {
				return fmt.Errorf("hint stream obj#%d: plaintext probe: %w", nr, plainErr)
			}
		}
		if err != nil {
			return fmt.Errorf("hint stream obj#%d: finalize content: %w", nr, err)
		}
		if err := c.Err(); err != nil {
			return err
		}
		if repaired {
			if ctx.Read.RepairedHints == nil {
				ctx.Read.RepairedHints = types.IntSet{}
			}
			ctx.Read.RepairedHints[nr] = true
			model.ShowRepaired(fmt.Sprintf("plaintext linearization hint obj#%d (compatibility exception)", nr))
		}
		entry.Object = probe
		ctx.Read.BinaryTotalSize += *probe.StreamLength
	}
	if err := c.Err(); err != nil {
		return err
	}
	if verified {
		ctx.LinearizationObjs[state.linearNr] = true
		for nr := range state.streams {
			ctx.LinearizationObjs[nr] = true
		}
	}
	return nil
}
