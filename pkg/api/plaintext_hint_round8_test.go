package api

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func round8KeywordSyntax(kind string) string {
	if strings.Contains(kind, "valid") {
		return "/Unused [null true false 1 0 R 2]"
	}
	if strings.Contains(kind, "reference") {
		return "/Unused [1 0 R2]"
	}
	return "/Unused [nulltruefalse]"
}

func TestHintRound8LegitimateBoundaryControls(t *testing.T) {
	for _, origin := range []string{"linear", "ordinary", "trailer", "compressed"} {
		t.Run(origin, func(t *testing.T) {
			pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{review8: "keyword-valid-" + origin, objectStream: origin == "compressed"})
			ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
			if err != nil {
				t.Fatal(err)
			}
			if len(ctx.LinearizationObjs) != 2 {
				t.Fatal("legitimate boundaries disqualified hints")
			}
		})
	}
}

func TestHintRound8Terminators(t *testing.T) {
	for _, kind := range []string{"endobj-linear", "endobj-ordinary"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review8: kind, encrypted: encrypted})
			})
		}
	}
}

func TestHintRound8XRefFraming(t *testing.T) {
	for _, kind := range []string{"xref-marker", "xref-header", "xref-record", "xref-header-older", "xref-record-older", "startxref-junk", "startxref-space", "startxref-prefix", "epilogue-space"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review8: kind, encrypted: encrypted})
			})
		}
	}
}

func TestHintRound8KeywordBoundaries(t *testing.T) {
	for _, token := range []string{"boolean", "reference"} {
		for _, origin := range []string{"linear", "ordinary", "trailer", "compressed"} {
			for _, encrypted := range []bool{false, true} {
				kind := "keyword-" + token + "-" + origin
				t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
					assertHintRound5Denial(t, primaryHintFixtureOptions{review8: kind, encrypted: encrypted, objectStream: origin == "compressed"})
				})
			}
		}
	}
}
