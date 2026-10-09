package pdfcpu

import (
	"testing"

	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintRound7OlderWideXRef(t *testing.T) {
	for _, field := range []int{0, 1, 2} {
		for _, lossy := range []bool{false, true} {
			w := [3]int{1, 4, 2}
			w[field] = 9
			data := make([]byte, w[0]+w[1]+w[2])
			start := 0
			for i := 0; i < field; i++ {
				start += w[i]
			}
			if lossy {
				data[start] = 1
			}
			retained := model.NewXRefTableEntryGen0(types.Integer(42))
			ctx := &model.Context{Read: &model.ReadContext{}, XRefTable: &model.XRefTable{Table: map[int]*model.XRefTableEntry{7: retained}}}
			if err := extractXRefTableEntriesFromXRefStream(data, 0, &types.XRefStreamDict{W: w, Objects: []int{7}}, ctx, 0); err != nil {
				t.Fatal(err)
			}
			if ctx.Read.RepairedXRef != lossy {
				t.Fatalf("field=%d lossy=%t repaired=%t", field, lossy, ctx.Read.RepairedXRef)
			}
			if ctx.Table[7] != retained {
				t.Fatal("older record replaced live entry")
			}
		}
	}
}
