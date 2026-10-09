package pdfcpu

import (
	"errors"
	"testing"

	"github.com/umats/pdfcpu/pkg/filter"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintObjectStreamPrefixBoundaries(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	if got, err := objectStreamOffset(maxInt-1, 1); err != nil || got != maxInt {
		t.Fatalf("exact offset boundary: %d %v", got, err)
	}
	for _, lookahead := range []int{0, 1} {
		offset := int(filter.DefaultMaxDecodeBytes) + 1 - lookahead
		if _, err := decodeObjectStreamPrefix(&types.ObjectStreamDict{}, offset, lookahead, 0); !errors.Is(err, filter.ErrDecodeLimitExceeded) {
			t.Fatalf("default budget boundary: %v", err)
		}
	}
	if _, err := decodeObjectStreamPrefix(&types.ObjectStreamDict{}, maxInt, 1, -1); !errors.Is(err, filter.ErrDecodeLimitExceeded) {
		t.Fatalf("unlimited lookahead overflow: %v", err)
	}
	if _, err := objectStreamOffset(-1, 1); !errors.Is(err, errCorruptObjectStreamDict) {
		t.Fatalf("negative relative offset: %v", err)
	}
}
