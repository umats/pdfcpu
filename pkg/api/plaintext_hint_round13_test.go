package api

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func TestHintRound13PublicMarker(t *testing.T) {
	for _, kind := range []string{"comment-hint", "comment-xref", "comment-ordinary"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review13: kind, encrypted: encrypted, objectStream: kind == "comment-xref"})
			})
		}
	}
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprintf("name-real-marker-encrypted%t", encrypted), func(t *testing.T) {
			pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{review13: "name-hint", encrypted: encrypted})
			ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
			if err != nil {
				t.Fatal(err)
			}
			if len(ctx.LinearizationObjs) != 2 {
				t.Fatal("legitimate name lost hint authority")
			}
		})
	}
}

func TestHintRound13PublicNative(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("requires runnable native32 target")
	}
	for _, kind := range []string{"native-generation", "native-container", "native-index", "older-native-generation", "older-native-container", "older-native-index"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review13: kind, encrypted: encrypted, objectStream: true})
			})
		}
	}
	for _, kind := range []string{"native-in-range", "older-native-in-range"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted%t", kind, encrypted), func(t *testing.T) {
				pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{review13: kind, encrypted: encrypted, objectStream: true})
				ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
				if err != nil {
					t.Fatal(err)
				}
				if len(ctx.LinearizationObjs) != 2 {
					t.Fatal("in-range xref lost authority")
				}
			})
		}
	}
}
