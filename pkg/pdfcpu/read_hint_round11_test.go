package pdfcpu

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/umats/pdfcpu/internal/contextutil"
	"github.com/umats/pdfcpu/pkg/pdfcpu/model"
	"github.com/umats/pdfcpu/pkg/pdfcpu/types"
)

// Cancel immediately after the entry check: normal cipher processing has no
// context parameter, so its next boundary must observe this cancellation.
type hintCancelAfterCheck struct {
	context.Context
	cancel context.CancelFunc
	checks int
}

func (c *hintCancelAfterCheck) Err() error {
	c.checks++
	if c.checks == 1 {
		c.cancel()
		return nil
	}
	return c.Context.Err()
}

type hintTraversalContext struct {
	context.Context
	checks, cancelAt int
}

func (c *hintTraversalContext) Err() error {
	c.checks++
	if c.cancelAt > 0 && c.checks >= c.cancelAt {
		return context.Canceled
	}
	return nil
}

func TestHintCompressedClaimsLinearTraversal(t *testing.T) {
	for _, count := range []int{1, 100, 1000} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			table := map[int]*model.XRefTableEntry{}
			for i := 1; i <= count; i++ {
				container, index := count+i, 0
				e := model.NewXRefTableEntryGen0(nil)
				e.Compressed = true
				e.ObjectStream = &container
				e.ObjectStreamInd = &index
				table[i] = e
				table[container] = model.NewXRefTableEntryGen0(nil)
			}
			p := &contextutil.ParseProvenance{}
			c := &hintTraversalContext{Context: contextutil.WithParseProvenance(t.Context(), p)}
			claims, err := indexCompressedClaims(c, table)
			if err != nil {
				t.Fatal(err)
			}
			for i := 1; i <= count; i++ {
				if err := auditCompressedClaims(c, claims[count+i], []string{fmt.Sprint(i), "0"}); err != nil {
					t.Fatal(err)
				}
			}
			if p.Ambiguous() {
				t.Fatal("multi-container identity rejected")
			}
			if c.checks < 3*count || c.checks > 3*count+2 {
				t.Fatalf("nonlinear traversal: %d entries, %d checks", count, c.checks)
			}
			if err := auditCompressedClaims(c, claims[count+1], []string{"999999", "0"}); err != nil {
				t.Fatal(err)
			}
			if !p.Ambiguous() {
				t.Fatal("wrong-container identity retained authority")
			}
			canceled := &hintTraversalContext{Context: t.Context(), cancelAt: count}
			if _, err := indexCompressedClaims(canceled, table); !errors.Is(err, context.Canceled) {
				t.Fatalf("index cancellation: %v", err)
			}
			canceled = &hintTraversalContext{Context: t.Context(), cancelAt: 1}
			if err := auditCompressedClaims(canceled, claims[count+1], []string{"1", "0"}); !errors.Is(err, context.Canceled) {
				t.Fatalf("audit cancellation: %v", err)
			}
		})
	}
}

func TestHintCompressedClaimsDecoderTraversal(t *testing.T) {
	ctx, _ := hintAuditTestContext(t)
	container := 30
	for i := 0; i < 100; i++ {
		index := 0
		e := model.NewXRefTableEntryGen0(nil)
		e.Compressed, e.ObjectStream, e.ObjectStreamInd = true, &container, &index
		ctx.Table[100+i] = e
	}
	claims, err := indexCompressedClaims(t.Context(), ctx.Table)
	if err != nil {
		t.Fatal(err)
	}
	decode := func(members []compressedClaim) int {
		c := &hintTraversalContext{Context: t.Context()}
		length := int64(8)
		sd := types.NewStreamDict(types.Dict{"First": types.Integer(6), "N": types.Integer(1)}, 0, &length, nil, nil)
		sd.Raw = []byte("100 0 42")
		if _, err := decodeObjectStreamObjects(c, ctx, &sd, container, members); err != nil {
			t.Fatal(err)
		}
		return c.checks
	}
	base := decode(nil)
	if got := decode(claims[container]) - base; got != 100 {
		t.Fatalf("decoder did not traverse exactly its claims: %d", got)
	}
}

func TestHintNormalCipherCancellation(t *testing.T) {
	for _, verified := range []bool{false, true} {
		t.Run(fmt.Sprint(verified), func(t *testing.T) {
			ctx, state := hintAuditTestContext(t)
			ctx.XRefTable.ValidationMode = model.ValidationStrict
			ctx.EncKey = bytes.Repeat([]byte{1}, 16)
			ctx.AES4Streams = true
			ctx.E = &model.Enc{R: 4, Emd: true}
			raw, err := encryptStream([]byte("encrypted hint bytes"), 11, 0, ctx.EncKey, true, 4)
			if err != nil {
				t.Fatal(err)
			}
			before := bytes.Clone(raw)
			length := int64(len(raw))
			sd := types.NewStreamDict(types.Dict{"Length": types.Integer(length), "S": types.Integer(0)}, 0, &length, nil, nil)
			sd.Raw = raw
			ctx.Table[11].Object = sd
			state.streams[11] = sd
			base, cancel := context.WithCancel(t.Context())
			defer cancel()
			c := &hintCancelAfterCheck{Context: base, cancel: cancel}
			err = finalizeHintStreamContents(c, ctx, state, verified)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("cancellation lost: %v", err)
			}
			if len(ctx.LinearizationObjs) != 0 || len(ctx.Read.RepairedHints) != 0 || ctx.Read.BinaryTotalSize != 0 {
				t.Fatal("cancellation committed authority or accounting")
			}
			if !bytes.Equal(ctx.Table[11].Object.(types.StreamDict).Raw, before) {
				t.Fatal("cancellation committed content")
			}
		})
	}
}
