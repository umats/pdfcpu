package pdfcpu

import (
	"context"
	"errors"
	"testing"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintRound9LexicalCutControls(t *testing.T) {
	for _, prefix := range []string{"12 0 ", "42 %x\n", "(a(%x))", "<2528>", "/A#25B ", "[null true]", "<< /N /Name /S (%x) >>"} {
		p := &contextutil.ParseProvenance{}
		if err := auditObjectStreamCut(contextutil.WithParseProvenance(t.Context(), p), []byte(prefix+"43"), 0, len(prefix)); err != nil {
			t.Fatal(err)
		}
		if p.Ambiguous() {
			t.Fatalf("valid cut rejected: %q", prefix)
		}
	}
	for _, prefix := range []string{"12 0", "42 %x ", "/", "(a% x", "<25", "[", "<<"} {
		p := &contextutil.ParseProvenance{}
		if err := auditObjectStreamCut(contextutil.WithParseProvenance(t.Context(), p), []byte(prefix+"43"), 0, len(prefix)); err != nil {
			t.Fatal(err)
		}
		if !p.Ambiguous() {
			t.Fatalf("incomplete cut accepted: %q", prefix)
		}
	}
	c, cancel := context.WithCancel(t.Context())
	cancel()
	if err := auditObjectStreamCut(c, []byte("12 0 43"), 0, 5); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation changed: %v", err)
	}
}

func TestHintRound9SourceCompressedCuts(t *testing.T) {
	for _, tc := range []struct {
		raw      string
		n, first int
		values   types.Array
	}{
		{"12 01", 1, 4, types.Array{types.Integer(1)}},
		{"12 042", 1, 4, types.Array{types.Integer(42)}},
		{"12 0 13 6 42 %x 43", 2, 10, types.Array{types.Integer(42), types.Integer(43)}},
		{"12 0 13 1 /1", 2, 10, types.Array{types.Name(""), types.Integer(1)}},
		{"12 4 %x  43", 1, 5, types.Array{types.Integer(43)}},
	} {
		t.Run(tc.raw, func(t *testing.T) {
			raw := []byte(tc.raw)
			length := int64(len(raw))
			osd := types.ObjectStreamDict{StreamDict: types.NewStreamDict(types.Dict{}, 0, &length, nil, nil), ObjCount: tc.n, FirstObjOffset: tc.first}
			osd.Raw = raw
			p := &contextutil.ParseProvenance{}
			c := contextutil.WithParseProvenance(t.Context(), p)
			if err := parseObjectStream(c, &osd, model.DefaultResourceLimits()); err != nil {
				t.Fatal(err)
			}
			for i, o := range osd.ObjArray {
				lazy := o.(types.LazyObjectStreamObject)
				v, err := lazy.DecodedObject(c)
				if err != nil {
					t.Fatal(err)
				}
				if v != tc.values[i] {
					t.Fatalf("ordinary value changed: %v", v)
				}
			}
			if !p.Ambiguous() {
				t.Fatal("manufactured lexical origin remained reliable")
			}
		})
	}
}
