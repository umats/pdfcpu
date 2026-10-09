package model

import (
	"testing"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintRound8KeywordValuesPreserved(t *testing.T) {
	for _, tc := range []struct {
		source   string
		count    int
		repaired bool
	}{
		{"[nulltruefalse]", 3, true},
		{"[1 0 R2]", 2, true},
		{"[null true false 1 0 R 2]", 5, false},
		{"[(nulltrue 1 0 R2) /nulltrue % nulltrue R2\n true]", 3, false},
	} {
		t.Run(tc.source, func(t *testing.T) {
			p := &contextutil.ParseProvenance{}
			s := tc.source
			o, err := ParseObject(contextutil.WithParseProvenance(t.Context(), p), &s, 0)
			if err != nil {
				t.Fatal(err)
			}
			a, ok := o.(types.Array)
			if !ok || len(a) != tc.count {
				t.Fatalf("ordinary value changed: %v", o)
			}
			if p.Ambiguous() != tc.repaired {
				t.Fatalf("provenance=%t want=%t", p.Ambiguous(), tc.repaired)
			}
		})
	}
}
