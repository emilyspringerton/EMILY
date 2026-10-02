package securechan

import (
	"bytes"
	"math/rand"
	"testing"
)

func roundTrip(t *testing.T, in []byte) []byte {
	c := lz4Compress(nil, in)
	out, err := lz4Decompress(c, len(in))
	if err != nil || !bytes.Equal(out, in) {
		t.Fatalf("roundtrip failed len=%d err=%v", len(in), err)
	}
	return c
}

func TestLZ4Shapes(t *testing.T) {
	for _, n := range []int{0, 1, 4, 11, 12, 13, 17, 100, 255, 256, 65535, 65536, 200000} {
		roundTrip(t, bytes.Repeat([]byte("a"), n))                       // RLE / overlapping matches
		roundTrip(t, bytes.Repeat([]byte("hello, golden "), n/14+1)[:n]) // repetitive text
		r := make([]byte, n)
		rand.New(rand.NewSource(int64(n))).Read(r)
		roundTrip(t, r) // incompressible
	}
}

func TestLZ4Compresses(t *testing.T) {
	in := bytes.Repeat([]byte("# GOLDEN DOC\nsome markdown line\n"), 2000)
	if c := roundTrip(t, in); len(c)*10 > len(in) {
		t.Fatalf("expected >10x on repetitive input, got %d -> %d", len(in), len(c))
	}
}

func TestLZ4RejectsCorrupt(t *testing.T) {
	in := bytes.Repeat([]byte("abcdefgh"), 500)
	c := lz4Compress(nil, in)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 2000; i++ { // fuzz: must never panic, and never exceed outLen
		d := append([]byte(nil), c...)
		for k := 0; k < 3; k++ {
			d[rng.Intn(len(d))] = byte(rng.Intn(256))
		}
		if out, err := lz4Decompress(d, len(in)); err == nil && len(out) != len(in) {
			t.Fatal("accepted wrong-length output")
		}
	}
	if _, err := lz4Decompress(c, len(in)-1); err == nil {
		t.Fatal("accepted undersized outLen")
	}
	if _, err := lz4Decompress(c[:len(c)-3], len(in)); err == nil {
		t.Fatal("accepted truncated block")
	}
}
