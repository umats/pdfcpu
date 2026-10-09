package pdfcpu

import (
	"testing"

	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

func TestHintRound4PostProcessRepairProvenance(t *testing.T) {
	for _, shift := range []bool{false, true} {
		entry := model.NewXRefTableEntryGen0(types.Integer(42))
		if shift {
			entry = model.NewFreeHeadXRefTableEntry()
		}
		ctx := &model.Context{Read: &model.ReadContext{}, XRefTable: &model.XRefTable{Table: map[int]*model.XRefTableEntry{1: entry}}}
		postProcess(ctx, 1)
		if !ctx.Read.RepairedXRef {
			t.Fatal("object zero repair was not tracked")
		}
		if ctx.Table[0] == nil {
			t.Fatal("object zero repair changed")
		}
	}
}
