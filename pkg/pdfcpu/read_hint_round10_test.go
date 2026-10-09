package pdfcpu

import (
	"bytes"
	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
	"strings"
	"testing"
)

func TestHintRound10CommentEndobjValue(t *testing.T) {
	data := []byte("12 0 obj\n/A\\% endobj\n<< /Target 8 0 R >>\nendobj\n" + strings.Repeat(" ", 1024))
	ctx, err := model.NewContext(bytes.NewReader(data), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	p := &contextutil.ParseProvenance{}
	o, _, _, _, err := object(contextutil.WithParseProvenance(t.Context(), p), ctx, 0, 12, 0)
	if err != nil {
		t.Fatal(err)
	}
	if o != types.Name("A\\") {
		t.Fatalf("ordinary value changed: %v", o)
	}
	if !p.Ambiguous() {
		t.Fatal("comment-contained terminator retained authority")
	}
}
