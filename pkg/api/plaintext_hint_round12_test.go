package api

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/umats/pdfcpu/pkg/filter"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
)

// Ordinary unencrypted PDF with no linearization or hint candidate. The sole
// compressed member is deliberately not referenced by the catalog.
func noHintPrefixBudgetPDF(t *testing.T, gap int) []byte {
	t.Helper()
	var compressed bytes.Buffer
	prolog := fmt.Sprintf("4 %d ", gap)
	zw := zlib.NewWriter(&compressed)
	if _, err := zw.Write([]byte(prolog + strings.Repeat(" ", gap) + "42 ")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.7\n")
	offsets := map[int]int{}
	object := func(n int, s string) { offsets[n] = pdf.Len(); fmt.Fprintf(&pdf, "%d 0 obj\n%s\nendobj\n", n, s) }
	object(1, "<< /Type /Catalog /Pages 2 0 R >>")
	object(2, "<< /Type /Pages /Kids [] /Count 0 >>")
	object(3, fmt.Sprintf("<< /Type /ObjStm /N 1 /First %d /Length %d /Filter /FlateDecode >>\nstream\n%sendstream", len(prolog), compressed.Len(), compressed.Bytes()))
	offsets[5] = pdf.Len()
	var xref bytes.Buffer
	for n := 0; n < 6; n++ {
		var entry [7]byte
		switch n {
		case 0:
			binary.BigEndian.PutUint16(entry[5:], 65535)
		case 4:
			entry[0] = 2
			binary.BigEndian.PutUint32(entry[1:], 3)
		default:
			entry[0] = 1
			binary.BigEndian.PutUint32(entry[1:], uint32(offsets[n]))
		}
		xref.Write(entry[:])
	}
	object(5, fmt.Sprintf("<< /Type /XRef /Size 6 /Root 1 0 R /W [1 4 2] /Length %d >>\nstream\n%sendstream", xref.Len(), xref.Bytes()))
	fmt.Fprintf(&pdf, "startxref\n%d\n%%%%EOF\n", offsets[5])
	return pdf.Bytes()
}

func TestHintObjectStreamPublicPrefixBudget(t *testing.T) {
	for _, gap := range []int{16, 128} {
		t.Run(fmt.Sprint(gap), func(t *testing.T) {
			conf := model.NewDefaultConfiguration()
			conf.Limits.MaxDecodeBytes = 64
			ctx, err := ReadContext(t.Context(), bytes.NewReader(noHintPrefixBudgetPDF(t, gap)), conf)
			if gap > 64 {
				if !errors.Is(err, filter.ErrDecodeLimitExceeded) {
					t.Fatalf("no-hint reader expected decode budget sentinel, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if ctx.Read.Linearized || len(ctx.LinearizationObjs) != 0 {
				t.Fatal("ordinary fixture acquired hint authority")
			}
		})
	}
}
