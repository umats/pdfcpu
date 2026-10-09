package api

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHintRound9CompressedCuts(t *testing.T) {
	for _, kind := range []string{"compressed-first-one", "compressed-first-42", "compressed-comment", "compressed-name", "compressed-gap-comment"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review9: kind, objectStream: true, encrypted: encrypted})
			})
		}
	}
}

func TestHintRound9TrailerCuts(t *testing.T) {
	for _, kind := range []string{"trailer-junk", "trailer-comment", "trailer-spaced-junk", "trailer-fake-frame"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review9: kind, encrypted: encrypted})
			})
		}
	}
}

func TestHintRound9HarmlessTrailerComment(t *testing.T) {
	pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{review9: "trailer-harmless"})
	ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
	if err != nil {
		t.Fatal(err)
	}
	if len(ctx.LinearizationObjs) != 2 {
		t.Fatal("harmless marker comment denied hints")
	}
}
