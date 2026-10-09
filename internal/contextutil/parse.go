package contextutil

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
)

// PDFWhitespace identifies only the six PDF whitespace characters.
func PDFWhitespace(r rune) bool {
	switch r {
	case 0, '\t', '\n', '\f', '\r', ' ':
		return true
	}
	return false
}

// PDFTokenBoundary checks a token's end within an actual parser buffer.
// Its EOF is valid only when the caller knows the buffer is a complete origin.
func PDFTokenBoundary(s string, end int) bool {
	return end == len(s) || (end >= 0 && end < len(s) &&
		(PDFWhitespace(rune(s[end])) || strings.ContainsRune("()<>[]/%", rune(s[end]))))
}

// ParseProvenance records ambiguity without changing parser acceptance. Its
// context scope is one reader or one isolated, notice-free structural probe.
type ParseProvenance struct {
	flags           atomic.Uint32
	xrefOrigins     sync.Map
	SuppressNotices bool
}

type parseProvenanceKey struct{}
type suppressParseNoticesKey struct{}

// WithoutParseNotices makes structural inspection quiet without losing provenance.
func WithoutParseNotices(c context.Context) context.Context {
	return context.WithValue(c, suppressParseNoticesKey{}, true)
}

// ParseNoticesSuppressed separates inspection logging from parser acceptance.
func ParseNoticesSuppressed(c context.Context) bool {
	if c == nil {
		return false
	}
	quiet, _ := c.Value(suppressParseNoticesKey{}).(bool)
	p := ParseProvenanceFromContext(c)
	return quiet || (p != nil && p.SuppressNotices)
}

// WithParseProvenance propagates provenance through existing parser signatures.
func WithParseProvenance(c context.Context, p *ParseProvenance) context.Context {
	if p == nil {
		return c
	}
	return context.WithValue(c, parseProvenanceKey{}, p)
}

// ParseProvenanceFromContext returns the optional reader-local tracking state.
func ParseProvenanceFromContext(c context.Context) *ParseProvenance {
	if c == nil {
		return nil
	}
	p, _ := c.Value(parseProvenanceKey{}).(*ParseProvenance)
	return p
}

// MarkDuplicate records overwritten dictionary values, even inside nested origins.
func (p *ParseProvenance) MarkDuplicate() {
	if p != nil {
		p.flags.Or(1)
	}
}

// MarkFallback records every retry through the permissive parser.
func (p *ParseProvenance) MarkFallback() {
	if p != nil {
		p.flags.Or(2)
	}
}

// MarkRepair records accepted lexical repairs or discarded structural syntax.
func (p *ParseProvenance) MarkRepair() {
	if p != nil {
		p.flags.Or(4)
	}
}

// Ambiguous reports whether dictionary values or parser provenance are unreliable.
func (p *ParseProvenance) Ambiguous() bool { return p != nil && p.flags.Load() != 0 }

// RecordXRefOrigin binds later framing probes to ordinary reader traversal.
func (p *ParseProvenance) RecordXRefOrigin(offset int64) {
	if p != nil {
		p.xrefOrigins.Store(offset, true)
	}
}

// HasXRefOrigin excludes syntactically plausible but unvisited revision claims.
func (p *ParseProvenance) HasXRefOrigin(offset int64) bool {
	if p == nil {
		return false
	}
	_, found := p.xrefOrigins.Load(offset)
	return found
}
