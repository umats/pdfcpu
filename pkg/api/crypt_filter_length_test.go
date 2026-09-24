package api

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

const (
	missingCFUserPW  = "synthetic-user-password"
	missingCFOwnerPW = "synthetic-owner-password"
)

// missingCFLengthPDF removes only the fixed-width crypt-filter entry; xref offsets stay valid.
func missingCFLengthPDF(t *testing.T, catalogVersion ...string) ([]byte, []byte) {
	t.Helper()
	source, content := encryptedXRefSourcePDF(t, catalogVersion...)
	conf := model.NewDefaultConfiguration()
	conf.UserPW, conf.OwnerPW = missingCFUserPW, missingCFOwnerPW
	conf.EncryptUsingAES = true
	conf.EncryptKeyLength = 128
	conf.WriteObjectStream = false
	conf.WriteXRefStream = false

	var encrypted bytes.Buffer
	if err := Encrypt(t.Context(), bytes.NewReader(source), &encrypted, conf); err != nil {
		t.Fatal(err)
	}
	pdf := bytes.Clone(encrypted.Bytes())
	if !bytes.HasPrefix(pdf, []byte("%PDF-1.7")) {
		t.Fatalf("synthetic fixture header: %q", pdf[:8])
	}
	pdf[7] = '4' // Same width: leave every object and xref offset untouched.
	if len(catalogVersion) > 0 && catalogVersion[0] == "2.0" {
		i := bytes.Index(pdf, []byte("/Lang("))
		if i < 0 {
			t.Fatal("missing synthetic Catalog Lang placeholder")
		}
		end := bytes.Index(pdf[i:], []byte(")/Pages"))
		if end < len("/Version/2.0") {
			t.Fatal("invalid synthetic Catalog placeholder")
		}
		slot := pdf[i : i+end+1]
		for j := range slot {
			slot[j] = ' '
		}
		copy(slot, "/Version/2.0")
	}
	cf := bytes.Index(pdf, []byte("/StdCF"))
	if cf < 0 {
		t.Fatal("missing StdCF")
	}
	end := bytes.Index(pdf[cf:], []byte(">>"))
	if end < 0 {
		t.Fatal("unterminated StdCF")
	}
	filter := pdf[cf : cf+end]
	if !bytes.Contains(filter, []byte("/AESV2")) {
		t.Fatal("expected AESV2 crypt filter")
	}
	length := []byte("/Length 128")
	i := bytes.Index(filter, length)
	if i < 0 {
		t.Fatal("missing AESV2 crypt-filter Length 128")
	}
	copy(filter[i:i+len(length)], bytes.Repeat([]byte{' '}, len(length)))
	return pdf, content
}

func TestMissingAESV2CryptFilterLengthCatalogVersion(t *testing.T) {
	pdf, _ := missingCFLengthPDF(t, "2.0")
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationStrict
	conf.UserPW = missingCFUserPW
	ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf)
	if err != nil {
		t.Fatalf("PDF-1.4 header with Catalog Version 2.0: %v", err)
	}
	if got := ctx.XRefTable.Version(); got != model.V20 {
		t.Fatalf("effective PDF version = %v, want 2.0", got)
	}
}

func TestMissingAESV2CryptFilterLength(t *testing.T) {
	pdf, content := missingCFLengthPDF(t)
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	conf.UserPW = missingCFUserPW

	ctx, err := ReadContext(t.Context(), bytes.NewReader(pdf), conf)
	if err != nil {
		t.Fatalf("relaxed read with synthetic AESV2 omission: %v", err)
	}
	if err := ValidateContext(t.Context(), ctx); err != nil {
		t.Fatalf("validate synthetic AESV2 omission: %v", err)
	}
	page, _, _, err := ctx.PageDict(t.Context(), 1, false)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.PageContent(page, 1)
	if err != nil || !bytes.Equal(got, content) {
		t.Fatalf("decrypted page content = %q, err = %v", got, err)
	}
	var out bytes.Buffer
	if err := Decrypt(t.Context(), bytes.NewReader(pdf), &out, conf); err != nil {
		t.Fatalf("relaxed decrypt with synthetic AESV2 omission: %v", err)
	}
	plain, err := ReadContext(t.Context(), bytes.NewReader(out.Bytes()), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("read decrypted output: %v", err)
	}
	if err := ValidateContext(t.Context(), plain); err != nil {
		t.Fatalf("validate decrypted output: %v", err)
	}
	wrong := model.NewDefaultConfiguration()
	wrong.ValidationMode = model.ValidationRelaxed
	wrong.UserPW = "wrong-test-password"
	if _, err := ReadContext(t.Context(), bytes.NewReader(pdf), wrong); !errors.Is(err, pdfcpu.ErrWrongPassword) {
		t.Fatalf("wrong password error = %v, want ErrWrongPassword", err)
	}
	out.Reset()
	if err := Decrypt(t.Context(), bytes.NewReader(pdf), &out, wrong); !errors.Is(err, pdfcpu.ErrWrongPassword) {
		t.Fatalf("wrong password decrypt error = %v, want ErrWrongPassword", err)
	}

	strict := model.NewDefaultConfiguration()
	strict.ValidationMode = model.ValidationStrict
	strict.UserPW = missingCFUserPW
	_, err = ReadContext(t.Context(), bytes.NewReader(pdf), strict)
	if !errors.Is(err, pdfcpu.ErrMalformedEncryption) || !strings.Contains(err.Error(), `crypt filter missing entry "Length"`) {
		t.Fatalf("strict read error = %v, want missing crypt-filter Length", err)
	}
}
