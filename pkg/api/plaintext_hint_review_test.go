package api

import (
	"bytes"
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"testing"

	"github.com/umats/pdfcpu/pkg/log"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintReviewUnreliableXRef(t *testing.T) {
	for _, kind := range []string{"stream-size", "skipped-entry", "short-offset", "classic-size", "prefixed-offset", "surplus-stream-entry"} {
		for _, encrypted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-encrypted-%t", kind, encrypted), func(t *testing.T) {
				pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{
					encrypted: encrypted, objectStream: kind == "stream-size" || kind == "surplus-stream-entry", unreliableXRef: kind,
				})
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
					t.Fatal("unreliable xref authorized hint deletion")
				}
				assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
			})
		}
	}
}

func TestHintReviewNullDuplicates(t *testing.T) {
	for _, key := range []string{"H", "Length", "S", "ordinary", "trailer", "compressed"} {
		for _, order := range []string{"first", "last"} {
			for _, encrypted := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s-null-%s-encrypted-%t", key, order, encrypted), func(t *testing.T) {
					opt := primaryHintFixtureOptions{encrypted: encrypted, objectStream: key == "compressed"}
					if key == "H" || key == "Length" || key == "S" {
						opt.ambiguousClaim = key + "-null-" + order
					} else {
						opt.ambiguousOrigin = key + "-null-" + order
					}
					pdf, _, content := plaintextPrimaryHintPDF(t, opt)
					ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
					if !encrypted {
						if err == nil {
							t.Fatal("null duplicate authorized plaintext repair")
						}
						return
					}
					if err != nil {
						t.Fatal(err)
					}
					if len(ctx.LinearizationObjs) != 0 {
						t.Fatal("null duplicate authorized deletion")
					}
					assertEncryptedXRefState(t, readEncryptedXRefState(t, pdf, "user"), content, true)
				})
			}
		}
	}
}

func TestHintReviewInspectionNotices(t *testing.T) {
	for _, opt := range []primaryHintFixtureOptions{
		{encrypted: true, shortFilter: true},
		{encrypted: true, objectStream: true, ambiguousOrigin: "compressed"},
	} {
		pdf, _, _ := plaintextPrimaryHintPDF(t, opt)
		var notices bytes.Buffer
		oldCLI := *log.CLI
		log.SetCLILogger(stdlog.New(&notices, "", 0))
		_, err := ReadContext(t.Context(), bytes.NewReader(pdf), encryptedXRefReadConfig("user"))
		*log.CLI = oldCLI
		if err != nil {
			t.Fatal(err)
		}
		wantFilters := 0
		if opt.shortFilter {
			wantFilters = 1
		}
		if got := strings.Count(notices.String(), "repaired: stream filter Fl"); got != wantFilters {
			t.Fatalf("filter notices=%d want=%d: %s", got, wantFilters, notices.String())
		}
		if strings.Contains(notices.String(), "duplicate key") {
			t.Fatalf("structural inspection emitted duplicate-key notices: %s", notices.String())
		}
	}
}

func TestHintReviewIncrementReopens(t *testing.T) {
	for _, overflow := range []bool{false, true} {
		t.Run(fmt.Sprintf("overflow-%t", overflow), func(t *testing.T) {
			pdf, _, content := plaintextPrimaryHintPDF(t, primaryHintFixtureOptions{overflow: overflow})
			f, err := os.CreateTemp(t.TempDir(), "synthetic-*.pdf")
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			if _, err := f.Write(pdf); err != nil {
				t.Fatal(err)
			}
			conf := encryptedXRefReadConfig("user")
			conf.OwnerPW = "owner"
			ann := model.NewTextAnnotation(*types.NewRectangle(0, 0, 10, 10), 0,
				"synthetic annotation", "synthetic-id", "", 0, nil, "", nil, nil,
				"", "", 0, 0, 0, false, "Comment")
			if err := AddAnnotationsAsIncrement(t.Context(), f, nil, ann, conf); err != nil {
				t.Fatal(err)
			}
			out, err := os.ReadFile(f.Name())
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.HasPrefix(out, pdf) {
				t.Fatal("increment changed original bytes")
			}
			assertEncryptedXRefState(t, readEncryptedXRefState(t, out, "user"), content, true)
			strict := encryptedXRefReadConfig("user")
			strict.ValidationMode = model.ValidationStrict
			if _, err := ReadContext(t.Context(), bytes.NewReader(out), strict); err != nil {
				t.Fatalf("increment still needs relaxed hint repair: %v", err)
			}
		})
	}
}
