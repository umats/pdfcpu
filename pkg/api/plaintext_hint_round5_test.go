package api

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
)

func assertHintRound5Denial(t *testing.T, opt primaryHintFixtureOptions) {
	t.Helper()
	pdf, _, content := plaintextPrimaryHintPDF(t, opt)
	ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
	if !opt.encrypted {
		if err == nil {
			t.Fatal("ambiguous origin authorized plaintext repair")
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(ctx.LinearizationObjs) != 0 {
		t.Fatal("ambiguous origin authorized deletion")
	}
	assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
}

func TestHintRound5NumericRepair(t *testing.T) {
	for _, kind := range []string{"numeric-H", "numeric-S"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{encrypted: encrypted, ambiguousClaim: kind})
			})
		}
	}
}

func TestHintRound5XRefRecovery(t *testing.T) {
	for _, kind := range []string{"inline", "prev-zero", "repeated-offset", "stream-header-junk", "stream-header-bj", "older-surplus", "older-shape"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted-%t", kind, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{encrypted: encrypted, unreliableXRef: kind, objectStream: kind == "stream-header-junk" || kind == "stream-header-bj"})
			})
		}
	}
}

func TestHintRound5TrailingSyntax(t *testing.T) {
	for _, origin := range []string{"linear", "hint", "ordinary", "compressed"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted-%t", origin, encrypted), func(t *testing.T) {
				assertHintRound5Denial(t, primaryHintFixtureOptions{encrypted: encrypted, objectStream: origin == "compressed", trailingOrigin: origin})
			})
		}
	}
}

func TestHintRound5TrailingComments(t *testing.T) {
	for _, origin := range []string{"ordinary-comment", "compressed-comment"} {
		pdf, _, _ := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{trailingOrigin: origin, objectStream: origin == "compressed-comment"})
		ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
		if err != nil {
			t.Fatal(err)
		}
		if len(ctx.LinearizationObjs) != 2 {
			t.Fatal("comments disqualified valid hints")
		}
	}
}

func TestHintRound5ElevatedDepth(t *testing.T) {
	for _, unencrypted := range []bool{false, true} {
		for _, mode := range []int{model.ValidationStrict, model.ValidationRelaxed} {
			t.Run(fmt.Sprintf("unencrypted-%t-mode-%d", unencrypted, mode), func(t *testing.T) {
				pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{encrypted: !unencrypted, unencrypted: unencrypted, objectStream: true, boundaryOrigin: "compressed", boundaryDepth: 100})
				conf := encryptedXRefReadConfig("user")
				conf.ValidationMode = mode
				conf.Limits.MaxRecursionDepth = 128
				ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf)
				if err != nil {
					t.Fatal(err)
				}
				if len(ctx.LinearizationObjs) != 2 {
					t.Fatal("elevated-limit candidate was not audited")
				}
				low := conf.Clone()
				low.Limits.MaxRecursionDepth = 100
				if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), low); !errors.Is(err, model.ErrMaxRecursionDepthExceeded) {
					t.Fatalf("configured low limit: %v", err)
				}
				if err := ValidateContext(t.Context(), ctx); err != nil {
					t.Fatal(err)
				}
				page, _, _, err := ctx.PageDict(t.Context(), 1, false)
				if err != nil {
					t.Fatal(err)
				}
				got, err := ctx.PageContent(page, 1)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, content) {
					t.Fatal("page content changed")
				}
			})
		}
	}
}
