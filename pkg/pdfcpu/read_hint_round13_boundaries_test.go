package pdfcpu

import (
	"context"
	"strings"
	"testing"
)

// Supplemental GREEN-only framing controls. Some of these origins are already
// rejected by the historical detector; the final-cut audit must also fail closed.
func TestHintRound13LexicalStates(t *testing.T) {
	for _, tc := range []struct {
		prefix string
		want   byte
	}{
		{"12 0 obj\n<< >> % >>", 'c'},
		{"12 0 obj\n<< >> /", 'n'},
		{"12 0 obj\n<< /Ignored (>>", 'u'},
		{"12 0 obj\n<< >> (>>", 's'},
		{"12 0 obj\n<< >> <AB", 'h'},
		{"12 0 obj\n<< /Ignored << >>", 'u'},
		{"12 0 obj\n<< /Ignored /stream >>\n", 0},
		{"12 0 obj\n<< >>% complete comment\n", 0},
	} {
		mode, err := lexicalCutState(t.Context(), []byte(tc.prefix+"stream"), 0, len(tc.prefix))
		if err != nil || mode != tc.want {
			t.Fatalf("%q mode=%q want=%q err=%v", tc.prefix, mode, tc.want, err)
		}
	}
	c, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := lexicalCutState(c, []byte(strings.Repeat(" ", 4097)), 0, 4097); err == nil {
		t.Fatal("cancelled framing audit succeeded")
	}
}
