//go:build tinybrain_debug

// Package llama - activation dump utility for numerical alignment debugging.
//
// Build with: go build -tags tinybrain_debug ./...
//
// Set TINYBRAIN_LOGIT_DUMP=1 at runtime. On the first forward pass:
//   - logits_dump.bin  : flat float32 little-endian (all vocab logits)
//   - logits_top20.txt : top-20 logit values with token indices
package llama

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sort"
	"sync/atomic"
)

var logitDumpDone int32

// DumpLogits exports logits on the first call when TINYBRAIN_LOGIT_DUMP=1.
// logits must be float32, length vocab_size.
func DumpLogits(logits []float32) {
	if !atomic.CompareAndSwapInt32(&logitDumpDone, 0, 1) {
		return
	}
	if os.Getenv("TINYBRAIN_LOGIT_DUMP") != "1" {
		return
	}
	binPath := "logits_dump.bin"
	f, err := os.Create(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logit_dump] cannot create %s: %v\n", binPath, err)
		return
	}
	h := sha256.New()
	tmp := make([]byte, 4)
	for _, v := range logits {
		binary.LittleEndian.PutUint32(tmp, math.Float32bits(v))
		f.Write(tmp)
		h.Write(tmp)
	}
	f.Close()

	type kv struct {
		idx int
		val float32
	}
	top := make([]kv, len(logits))
	for i, v := range logits {
		top[i] = kv{i, v}
	}
	sort.Slice(top, func(a, b int) bool { return top[a].val > top[b].val })
	if len(top) > 20 {
		top = top[:20]
	}
	if tf, e := os.Create("logits_top20.txt"); e == nil {
		fmt.Fprintf(tf, "# vocab_size=%d SHA-256=%x\n", len(logits), h.Sum(nil))
		for rank, kv := range top {
			fmt.Fprintf(tf, "%3d  %7d  %+.6f\n", rank+1, kv.idx, kv.val)
		}
		tf.Close()
	}
	fmt.Printf("[logit_dump] %s (%d floats) sha=%x\n", binPath, len(logits), h.Sum(nil))
}
