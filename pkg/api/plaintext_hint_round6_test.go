package api

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestHintRound6Hex(t *testing.T) {
	for _, kind := range []string{"hex-hint", "hex-ordinary", "hex-compressed"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review6: kind, encrypted: encrypted, objectStream: kind == "hex-compressed"})
			})
		}
	}
}

func TestHintRound6Prefixes(t *testing.T) {
	for _, kind := range []string{"trailer-prefix", "trailer-inline", "member-prefix", "prolog-tail", "index-junk", "index-space", "index-mismatch"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				compressed := strings.HasPrefix(kind, "member") || strings.HasPrefix(kind, "prolog") || strings.HasPrefix(kind, "index")
				assertHintRound5Denial(t, primaryHintFixtureOptions{review6: kind, encrypted: encrypted, objectStream: compressed})
			})
		}
	}
}

func TestHintRound6Repairs(t *testing.T) {
	for _, kind := range []string{"single-ID", "xref-length"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review6: kind, encrypted: encrypted, objectStream: kind == "xref-length"})
			})
		}
	}
}

func TestHintRound6CommentControls(t *testing.T) {
	for _, kind := range []string{"trailer-comment", "member-comment", "prolog-comment"} {
		pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{review6: kind, objectStream: kind != "trailer-comment"})
		ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
		if err != nil {
			t.Fatal(err)
		}
		if len(ctx.LinearizationObjs) != 2 {
			t.Fatal("comments disqualified valid hints")
		}
	}
}
