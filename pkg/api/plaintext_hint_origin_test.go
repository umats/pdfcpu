package api

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHintCompressedIndirectLength(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprint(encrypted), func(t *testing.T) {
			pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{
				review10: "compressed-length", objectStream: true, encrypted: encrypted,
			})
			ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
			if err != nil {
				t.Fatal(err)
			}
			if len(ctx.LinearizationObjs) != 2 {
				t.Fatal("compressed indirect length lost hint authority")
			}
			assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
		})
	}
}

func TestHintPhysicalOriginDenial(t *testing.T) {
	for _, origin := range []string{"complete-epilogue", "comment-origin", "stream-origin", "counterfeit-xref"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/encrypted=%t", origin, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review10: origin, encrypted: encrypted})
			})
		}
	}
}
