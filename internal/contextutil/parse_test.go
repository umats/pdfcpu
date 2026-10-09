package contextutil

import "testing"

func TestVisitedXRefOrigins(t *testing.T) {
	var absent *ParseProvenance
	absent.RecordXRefOrigin(12)
	if absent.HasXRefOrigin(12) {
		t.Fatal("absent provenance granted authority")
	}
	p := &ParseProvenance{}
	if p.HasXRefOrigin(12) {
		t.Fatal("unvisited origin granted authority")
	}
	p.RecordXRefOrigin(12)
	if !p.HasXRefOrigin(12) || p.HasXRefOrigin(13) {
		t.Fatal("xref origin identity was not preserved")
	}
}
