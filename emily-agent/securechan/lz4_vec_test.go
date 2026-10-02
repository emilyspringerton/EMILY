package securechan

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

// Interop vectors for an independent LZ4 (liblz4 via ctypes). Run with LZ4_VEC_DIR set; skipped otherwise.
func lz4Inputs() [][]byte {
	rng := rand.New(rand.NewSource(7))
	var ins [][]byte
	for _, n := range []int{0, 1, 13, 100, 1000, 70000} {
		r := make([]byte, n)
		rng.Read(r)
		ins = append(ins, r, bytes.Repeat([]byte("golden doc line\n"), n/16+1)[:n], bytes.Repeat([]byte{0}, n))
	}
	return ins
}

func TestLZ4VectorsOut(t *testing.T) {
	dir := os.Getenv("LZ4_VEC_DIR")
	if dir == "" {
		t.Skip("LZ4_VEC_DIR not set")
	}
	for i, in := range lz4Inputs() {
		os.WriteFile(fmt.Sprintf("%s/%d.raw", dir, i), in, 0o644)
		os.WriteFile(fmt.Sprintf("%s/%d.ours", dir, i), []byte(hex.EncodeToString(lz4Compress(nil, in))), 0o644)
	}
}

func TestLZ4VectorsIn(t *testing.T) { // liblz4-compressed blocks -> our decompressor
	dir := os.Getenv("LZ4_VEC_DIR")
	if dir == "" {
		t.Skip("LZ4_VEC_DIR not set")
	}
	for i, in := range lz4Inputs() {
		h, err := os.ReadFile(fmt.Sprintf("%s/%d.theirs", dir, i))
		if err != nil {
			t.Fatal(err)
		}
		blk, _ := hex.DecodeString(string(h))
		out, err := lz4Decompress(blk, len(in))
		if err != nil || !bytes.Equal(out, in) {
			t.Fatalf("vector %d: liblz4 block rejected/mismatch: %v", i, err)
		}
	}
}
