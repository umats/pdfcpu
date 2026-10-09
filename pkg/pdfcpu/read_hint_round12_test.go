package pdfcpu

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/umats/pdfcpu/pkg/filter"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintObjectStreamPrefixBudget(t *testing.T) {
	for _, tc := range []struct {
		name     string
		gap      int
		budget   int64
		exceeded bool
	}{
		{"gap", 128, 64, true}, {"in-budget", 16, 64, false}, {"unlimited", 128, -1, false}, {"default", 128, 0, false},
		{"prolog", 0, 3, true}, {"lookahead", 0, 4, true}, {"exact-lookahead", 0, 5, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := hintAuditTestContext(t)
			ctx.Configuration.Limits.MaxDecodeBytes = tc.budget
			prolog := fmt.Sprintf("4 %d ", tc.gap)
			var raw bytes.Buffer
			zw := zlib.NewWriter(&raw)
			_, _ = zw.Write([]byte(prolog + strings.Repeat(" ", tc.gap) + "42 "))
			_ = zw.Close()
			length := int64(raw.Len())
			sd := types.NewStreamDict(types.Dict{"First": types.Integer(len(prolog)), "N": types.Integer(1)}, 0, &length, nil, []types.PDFFilter{{Name: filter.Flate}})
			sd.Raw = raw.Bytes()
			_, err := decodeObjectStreamObjects(t.Context(), ctx, &sd, 30, nil)
			if tc.exceeded {
				if !errors.Is(err, filter.ErrDecodeLimitExceeded) {
					t.Fatalf("expected decode budget sentinel, got %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestHintObjectStreamOffsetArithmetic(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, first := range []int{1, 0} {
		t.Run(strconv.Itoa(first), func(t *testing.T) {
			osd := &types.ObjectStreamDict{FirstObjOffset: first, MaxDecodeBytes: -1}
			_, err := buildObjectArrayForObjectStream(t.Context(), osd, []string{"4", strconv.Itoa(maxInt)}, false, nil, 100)
			if !errors.Is(err, filter.ErrDecodeLimitExceeded) {
				t.Fatalf("offset/lookahead overflow must be terminal: %v", err)
			}
		})
	}
}
