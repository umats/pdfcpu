package pdfcpu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintRound13FinalMarker(t *testing.T) {
	for _, tc := range []struct {
		name, prefix string
		ambiguous    bool
	}{
		{"comment", "12 0 obj\n<< /Length 1 /DecodeParms << /Ignored /stream >> >> % >>stream\n", true},
		{"name-real-marker", "12 0 obj\n<< /Length 1 /DecodeParms << /Ignored /stream >> >>\nstream\n", false},
		{"container", "12 0 obj\n<< /Length 1 /Ignored << >>stream\n", true},
		{"ordinary-concealed-reference", "12 0 obj\n<< /Length 24 /DecodeParms << /Ignored /stream >> >> % >>stream\n<< /Target 8 0 R >>\n", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := &contextutil.ParseProvenance{}
			_, _, stream, _, err := buffer(contextutil.WithParseProvenance(t.Context(), p), bytes.NewBufferString(tc.prefix+"x\nendstream\nendobj\n"+strings.Repeat(" ", 1024)), 4096)
			if err != nil || stream < 0 {
				t.Fatalf("historical parsing changed: marker=%d err=%v", stream, err)
			}
			if p.Ambiguous() != tc.ambiguous {
				t.Fatalf("ambiguous=%t want=%t", p.Ambiguous(), tc.ambiguous)
			}
		})
	}
}

func TestHintRound13NativeXRef(t *testing.T) {
	for _, kind := range []byte{0, 1, 2} {
		for _, field := range []int{1, 2} {
			if kind != 2 && field == 1 {
				continue
			}
			for _, older := range []bool{false, true} {
				for _, wide := range []bool{false, true} {
					t.Run(fmt.Sprintf("kind%d-field%d-older%t-wide%t", kind, field, older, wide), func(t *testing.T) {
						data := make([]byte, 11)
						data[0] = kind
						value := uint64(0)
						if wide {
							value = 1 << 32
						}
						if field == 1 {
							binary.BigEndian.PutUint32(data[2:6], uint32(value))
							data[1] = byte(value >> 32)
						} else {
							binary.BigEndian.PutUint32(data[7:11], uint32(value))
							data[6] = byte(value >> 32)
						}
						retained := model.NewXRefTableEntryGen0(types.Integer(42))
						table := map[int]*model.XRefTableEntry{}
						if older {
							table[7] = retained
						}
						ctx := &model.Context{Read: &model.ReadContext{ObjectStreams: types.IntSet{}}, XRefTable: &model.XRefTable{Table: table}}
						if err := extractXRefTableEntriesFromXRefStream(data, 0, &types.XRefStreamDict{W: [3]int{1, 5, 5}, Objects: []int{7}}, ctx, 0); err != nil {
							t.Fatal(err)
						}
						if ctx.Read.RepairedXRef != (wide && strconv.IntSize == 32) {
							t.Fatalf("native%d repaired=%t", strconv.IntSize, ctx.Read.RepairedXRef)
						}
						if older && ctx.Table[7] != retained {
							t.Fatal("older entry replaced retained object")
						}
					})
				}
			}
		}
	}
}
