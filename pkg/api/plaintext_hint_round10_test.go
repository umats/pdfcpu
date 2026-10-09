package api

import (
	"fmt"
	"testing"
)

func TestHintRound10CommentEndobj(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprint(encrypted), func(t *testing.T) {
			assertHintRound5Denial(t, primaryHintFixtureOptions{review10: "comment-endobj", encrypted: encrypted})
		})
	}
}
