package securechan

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func startEcho(t *testing.T, id *Identity) string {
	l, err := Listen("tcp", "127.0.0.1:0", id)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			go func() { defer c.Close(); io.Copy(c, c) }()
		}
	}()
	return l.Addr().String()
}

func TestRoundTripAndDuplex(t *testing.T) {
	id, _ := GenerateIdentity()
	addr := startEcho(t, id)
	c, err := Dial(context.Background(), "tcp", addr, id.PublicKey())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	payloads := [][]byte{[]byte("x"), bytes.Repeat([]byte("golden doc\n"), 50000), make([]byte, 300000)}
	rand.Read(payloads[2])
	// write and read concurrently (full duplex): a 300KB write would deadlock a half-duplex design
	for _, p := range payloads {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { defer wg.Done(); c.Write(p) }()
		got := make([]byte, len(p))
		if _, err := io.ReadFull(c, got); err != nil || !bytes.Equal(got, p) {
			t.Fatalf("echo mismatch for %d bytes: %v", len(p), err)
		}
		wg.Wait()
	}
}

func TestCompressionShrinksWire(t *testing.T) {
	id, _ := GenerateIdentity()
	addr := startEcho(t, id)
	c, _ := Dial(context.Background(), "tcp", addr, id.PublicKey())
	defer c.Close()
	p := bytes.Repeat([]byte("# NORTHSTAR\nrepetitive markdown\n"), 20000)
	go c.Write(p)
	io.ReadFull(c, make([]byte, len(p)))
	plain, wire := c.Stats()
	t.Logf("%d plaintext -> %d on the wire (%.1fx)", plain, wire, float64(plain)/float64(wire))
	if wire*10 > plain {
		t.Fatalf("compression ineffective: %d -> %d", plain, wire)
	}
	// incompressible data must not grow meaningfully
	c2, _ := Dial(context.Background(), "tcp", addr, id.PublicKey())
	defer c2.Close()
	r := make([]byte, 100000)
	rand.Read(r)
	go c2.Write(r)
	io.ReadFull(c2, make([]byte, len(r)))
	if pl, wi := c2.Stats(); wi > pl+pl/50 {
		t.Fatalf("incompressible grew too much: %d -> %d", pl, wi)
	}
}

func TestWrongPinnedKeyRejected(t *testing.T) {
	real, _ := GenerateIdentity()
	imposter, _ := GenerateIdentity()
	addr := startEcho(t, real)
	if _, err := Dial(context.Background(), "tcp", addr, imposter.PublicKey()); err != ErrAuth {
		t.Fatalf("want ErrAuth, got %v", err)
	}
}

func TestIdentityFromSeed(t *testing.T) {
	id, _ := GenerateIdentity()
	id2, err := IdentityFromSeed(id.Seed())
	if err != nil || !bytes.Equal(id.PublicKey(), id2.PublicKey()) {
		t.Fatal("seed roundtrip changed the key")
	}
	addr := startEcho(t, id2)
	if c, err := Dial(context.Background(), "tcp", addr, id.PublicKey()); err != nil {
		t.Fatal(err)
	} else {
		c.Close()
	}
}

// tamper flips one bit in transit after the handshake and expects the receiver to fail closed.
type flipConn struct {
	net.Conn
	after, seen int
	flipped     bool
}

func (f *flipConn) Write(p []byte) (int, error) {
	if f.seen >= f.after && !f.flipped && len(p) > 8 {
		q := append([]byte(nil), p...)
		q[len(q)-1] ^= 1
		f.flipped = true
		return f.Conn.Write(q)
	}
	f.seen++
	return f.Conn.Write(p)
}

func TestTamperedRecordRejected(t *testing.T) {
	id, _ := GenerateIdentity()
	addr := startEcho(t, id)
	raw, _ := net.Dial("tcp", addr)
	fc := &flipConn{Conn: raw, after: 1} // msg1 is write #0; flip the first record
	c, err := Client(fc, id.PublicKey())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.Write([]byte("hello hello hello hello"))
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := c.Read(make([]byte, 64)); err == nil {
		t.Fatal("server accepted a tampered record")
	}
}

func TestGarbageHandshakeAndBomb(t *testing.T) {
	id, _ := GenerateIdentity()
	addr := startEcho(t, id)
	raw, _ := net.Dial("tcp", addr)
	raw.Write(bytes.Repeat([]byte{0xff}, msg1Len)) // bad version byte
	raw.SetReadDeadline(time.Now().Add(2 * time.Second))
	if n, _ := raw.Read(make([]byte, 16)); n != 0 {
		t.Fatal("server answered a garbage handshake")
	}
}

// A validly-encrypted record whose header claims a huge uncompressed size must be refused before
// any allocation (the peer is authenticated but could still be hostile or buggy).
func TestDecompressionBombRefused(t *testing.T) {
	a, b := net.Pipe()
	k1, k2 := bytes.Repeat([]byte{1}, 32), bytes.Repeat([]byte{2}, 32)
	snd, _ := newConn(a, k1, k2)
	rcv, _ := newConn(b, k2, k1)
	go func() {
		body := append([]byte{1}, 0x80, 0x80, 0x80, 0x80, 0x04) // flag=lz4, uvarint 1<<30
		body = append(body, 0x10, 'x')
		hdr := []byte{0, 0, 0, 0}
		hdr[3] = byte(len(body) + snd.send.Overhead())
		a.Write(snd.send.Seal(hdr, nonce(0), body, hdr))
	}()
	b.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := rcv.Read(make([]byte, 16)); err != ErrRecord {
		t.Fatalf("want ErrRecord for oversized claimed length, got %v", err)
	}
}

func TestHTTPOverSecureChan(t *testing.T) {
	id, _ := GenerateIdentity()
	l, _ := Listen("tcp", "127.0.0.1:0", id)
	body := strings.Repeat("# golden doc line\n", 5000)
	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Seen-Auth", r.Header.Get("Authorization"))
		io.WriteString(w, body)
	})}
	go srv.Serve(l)
	defer srv.Close()
	var conns []*Conn
	var mu sync.Mutex
	hc := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, n, a string) (net.Conn, error) {
		c, err := Dial(ctx, n, a, id.PublicKey())
		mu.Lock()
		conns = append(conns, c)
		mu.Unlock()
		return c, err
	}}}
	req, _ := http.NewRequest("GET", "http://"+l.Addr().String()+"/x", nil)
	req.Header.Set("Authorization", "Bearer abc")
	resp, err := hc.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(resp.Body)
	if string(got) != body || resp.Header.Get("X-Seen-Auth") != "Bearer abc" {
		t.Fatal("http over securechan returned wrong data")
	}
	t.Logf("http body %d bytes delivered over the channel", len(got))
}
