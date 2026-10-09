package api

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintRound4DepthBoundary(t *testing.T) {
	for _, origin := range []string{"stream", "compressed"} {
		for _, unencrypted := range []bool{false, true} {
			for _, mode := range []int{model.ValidationStrict, model.ValidationRelaxed} {
				t.Run(fmt.Sprintf("%s-unencrypted-%t-mode-%d", origin, unencrypted, mode), func(t *testing.T) {
					pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{encrypted: !unencrypted, unencrypted: unencrypted, objectStream: origin == "compressed", boundaryOrigin: origin})
					conf := encryptedXRefReadConfig("user")
					conf.ValidationMode = mode
					conf.Limits.MaxRecursionDepth = 4 // Dict(0), three arrays(1..3), scalar(4).
					ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf)
					if err != nil {
						t.Fatal(err)
					}
					found := false
					for nr := range ctx.Table {
						o, err := ctx.Dereference(*types.NewIndirectRef(nr, 0))
						if err != nil {
							t.Fatal(err)
						}
						var d types.Dict
						switch o := o.(type) {
						case types.Dict:
							d = o
						case types.StreamDict:
							d = o.Dict
						}
						if d["Boundary"] != nil {
							found = true
						}
					}
					if !found {
						t.Fatal("boundary origin missing")
					}
					assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, !unencrypted)
				})
			}
		}
	}
}

func TestHintRound4Headers(t *testing.T) {
	for _, header := range []string{"linear-junk", "linear-bj", "hint-junk", "hint-bj"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted-%t", header, encrypted), func(t *testing.T) {
				pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{encrypted: encrypted, malformedHeader: header})
				ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
				if !encrypted {
					if err == nil {
						t.Fatal("malformed header authorized plaintext repair")
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				if len(ctx.LinearizationObjs) != 0 {
					t.Fatal("malformed header authorized deletion")
				}
				assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
			})
		}
	}
}

// Both classic sections assign the same numbers. The earlier revision must
// not be confused with overlapping assignments within one section.
func TestHintRound4OlderRevisionPrecedence(t *testing.T) {
	for _, mode := range []int{model.ValidationStrict, model.ValidationRelaxed} {
		pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{encrypted: true})
		conf := encryptedXRefReadConfig("user")
		conf.ValidationMode = mode
		ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf)
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Read.RepairedXRef || len(ctx.LinearizationObjs) != 2 {
			t.Fatal("legitimate older section disqualified hints")
		}
		assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
	}
}

func TestHintRound4XRefAuthorization(t *testing.T) {
	for _, kind := range []string{"duplicate-classic", "duplicate-stream", "missing-zero", "free-list"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted-%t", kind, encrypted), func(t *testing.T) {
				pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{encrypted: encrypted, objectStream: kind == "duplicate-stream", unreliableXRef: kind})
				ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
				if !encrypted {
					if err == nil {
						t.Fatal("unreliable xref authorized plaintext repair")
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				if !ctx.Read.RepairedXRef || len(ctx.LinearizationObjs) != 0 {
					t.Fatal("unreliable xref authorized deletion")
				}
				assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
			})
		}
	}
}
