package api

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func round7LexicalSyntax(kind string) string {
	switch strings.Split(kind, "-")[0] {
	case "space":
		return "/O\v0"
	case "unicode":
		return "/O\u00a00"
	case "null":
		return "/O NULL"
	case "true":
		return "/Ignored TRUE"
	case "false":
		return "/Ignored FALSE"
	case "comment":
		return "/O 0 % NULL TRUE FALSE\v\u00a0\n"
	case "payload":
		return "/O (NULL TRUE FALSE\v\u00a0)"
	}
	return ""
}

func TestHintRound7PayloadControls(t *testing.T) {
	for _, kind := range []string{"comment-hint", "comment-ordinary", "comment-trailer", "comment-compressed", "payload-compressed", "postlude-comment"} {
		t.Run(kind, func(t *testing.T) {
			pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{review7: kind, objectStream: strings.HasSuffix(kind, "compressed")})
			ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
			if err != nil {
				t.Fatal(err)
			}
			if len(ctx.LinearizationObjs) != 2 {
				t.Fatal("legitimate payload disqualified hints")
			}
		})
	}
}

func TestHintRound7Postlude(t *testing.T) {
	for _, kind := range []string{"postlude-ordinary", "postlude-container", "postlude-xref"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review7: kind, encrypted: encrypted, objectStream: kind != "postlude-ordinary"})
			})
		}
	}
}

func TestHintRound7WideXRef(t *testing.T) {
	for _, kind := range []string{"wide-offset", "wide-generation"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review7: kind, encrypted: encrypted, objectStream: true})
			})
		}
	}
}

func TestHintRound7StreamMarker(t *testing.T) {
	for _, kind := range []string{"stream-space", "stream-space-ordinary", "stream-space-container", "stream-space-xref"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{review7: kind, encrypted: encrypted, objectStream: kind == "stream-space-container" || kind == "stream-space-xref"})
			})
		}
	}
}

func TestHintRound7Envelope(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprint(encrypted), func(t *testing.T) {
			assertHintRound5Denial(t, primaryHintFixtureOptions{review7: "clipped-endobj", encrypted: encrypted})
		})
	}
}

func TestHintRound7Lexical(t *testing.T) {
	for _, keyword := range []string{"space", "unicode", "null", "true", "false"} {
		for _, origin := range []string{"hint", "ordinary", "trailer", "compressed"} {
			if origin == "hint" && (keyword == "true" || keyword == "false") {
				continue
			}
			for _, encrypted := range []bool{false, true} {
				kind := keyword + "-" + origin
				t.Run(fmt.Sprintf("%s-%t", kind, encrypted), func(t *testing.T) {
					assertHintRound5Denial(t, primaryHintFixtureOptions{review7: kind, encrypted: encrypted, objectStream: origin == "compressed"})
				})
			}
		}
	}
}
