package api

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

func TestHintIndirectLengthGenerationDenial(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprint(encrypted), func(t *testing.T) {
			pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{review10: "compressed-length", objectStream: true, encrypted: encrypted})
			re := regexp.MustCompile(`/Length ([0-9]+) 0 R`)
			if !re.Match(pdf) {
				t.Fatal("missing indirect length")
			}
			pdf = re.ReplaceAll(pdf, []byte(`/Length ${1} 1 R`))
			ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
			if !encrypted {
				if err == nil {
					t.Fatal("plaintext hint retained authority")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(ctx.LinearizationObjs) != 0 {
				t.Fatal("generation mismatch retained deletion authority")
			}
			assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
		})
	}
}
