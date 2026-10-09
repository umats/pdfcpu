package pdfcpu

import (
	"bytes"
	"strings"
	"testing"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintRound8NonstreamPhysicalTerminator(t *testing.T) {
	for _, tail := range []string{"endobj\n", "endobj\v", "endobjJUNK\nxref\n"} {
		data := []byte("1 0 obj\n42\n" + tail + strings.Repeat(" ", 1024))
		ctx, err := model.NewContext(bytes.NewReader(data), model.NewDefaultConfiguration())
		if err != nil {
			t.Fatal(err)
		}
		p := &contextutil.ParseProvenance{}
		o, _, _, _, err := object(contextutil.WithParseProvenance(t.Context(), p), ctx, 0, 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if o != types.Integer(42) {
			t.Fatalf("ordinary value changed: %v", o)
		}
		if p.Ambiguous() != (tail != "endobj\n") {
			t.Fatalf("tail=%q ambiguous=%t", tail, p.Ambiguous())
		}
	}
}

func TestHintRound8NonstreamNameMarker(t *testing.T) {
	data := []byte("1 0 obj\n/endobj\n" + strings.Repeat(" ", 1024))
	ctx, err := model.NewContext(bytes.NewReader(data), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	p := &contextutil.ParseProvenance{}
	o, _, _, _, err := object(contextutil.WithParseProvenance(t.Context(), p), ctx, 0, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if o != types.Name("") {
		t.Fatalf("ordinary recovery changed: %v", o)
	}
	if !p.Ambiguous() {
		t.Fatal("name marker recovery remained reliable")
	}
}

func TestHintRound8CompressedPhysicalBoundary(t *testing.T) {
	for _, joined := range []bool{false, true} {
		prolog, body := "12 0 13 5 ", "null true"
		if joined {
			prolog, body = "12 0 13 4 ", "nulltrue"
		}
		raw := []byte(prolog + body)
		length := int64(len(raw))
		osd := types.ObjectStreamDict{StreamDict: types.NewStreamDict(types.Dict{}, 0, &length, nil, nil), ObjCount: 2, FirstObjOffset: len(prolog)}
		osd.Raw = raw
		p := &contextutil.ParseProvenance{}
		c := contextutil.WithParseProvenance(t.Context(), p)
		if err := parseObjectStream(c, &osd, model.DefaultResourceLimits()); err != nil {
			t.Fatal(err)
		}
		for _, o := range osd.ObjArray {
			lazy := o.(types.LazyObjectStreamObject)
			if _, err := lazy.DecodedObject(c); err != nil {
				t.Fatal(err)
			}
		}
		if p.Ambiguous() != joined {
			t.Fatalf("joined=%t ambiguous=%t", joined, p.Ambiguous())
		}
	}
}
