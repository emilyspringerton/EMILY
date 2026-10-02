// Package securechan is a post-quantum-hybrid, compressed, full-duplex stream transport that
// implements net.Conn / net.Listener, so net/http (and anything else) runs over it unchanged.
//
// K8s migration VS0, founder real-time 2026-10-02: traffic now crosses the internet and "needs
// strong encryption ... maybe the post quantum stuff in PARENA", with "custom streaming
// bidirectional LZ4".
//
// Handshake (1 round trip). The server owns a static ML-KEM-768 key that clients pin; only the
// holder of its decapsulation key can derive the session keys, which is what authenticates the
// server (no signatures needed). Ephemeral X25519 + ephemeral ML-KEM-768 add forward secrecy
// against both classical and quantum adversaries; all three secrets feed HKDF-SHA256 salted with
// the transcript hash:
//
//	C->S  ver(1)=1 | x25519 eph pub(32) | mlkem eph ek(1184) | ct_static(1088)   [encaps to server static ek]
//	S->C  x25519 eph pub(32) | ct_eph(1088) [encaps to client eph ek] | confirm(32)
//	keys = HKDF(ikm = k_static || k_x25519 || k_eph, salt = SHA256(msg1 || msg2-without-confirm))
//	confirm = HMAC-SHA256(confirm key, "server" || transcript)   -- client rejects on mismatch
//
// Records: len(4, BE, ciphertext length) | AES-256-GCM(nonce = 4 zero bytes || 8-byte BE sequence,
// AAD = the 4 length bytes). Separate keys per direction. Plaintext: flag(1) then either raw bytes
// (0) or uvarint(rawLen) + an LZ4 block (1); compression is used only when it actually shrinks
// the record.
//
// KNOWN LIMITS (v0, stated plainly): compress-then-encrypt leaks compressed lengths, the
// CRIME/BREACH class. Do not put attacker-influenced data and secrets in the same stream; the
// IDUNA bearer token travels in HTTP headers, so callers that mix untrusted content into
// requests should wrap with DisableCompression. The server does not authenticate clients here
// (IDUNA JWT does, one layer up). No rekeying: a connection is closed after 2^32 records. The
// crypto is Go's stdlib (crypto/mlkem, crypto/ecdh, crypto/aes, crypto/hkdf); the PARENA
// crypto/mlkem binding is cross-verified against the same Go implementation.
package securechan

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/hmac"
	"crypto/mlkem"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	version      = 1
	maxChunk     = 64 * 1024
	maxCT        = maxChunk + 64
	minCompress  = 64
	maxRecords   = 1 << 32
	handshakeTTL = 10 * time.Second
	msg1Len      = 1 + 32 + mlkem.EncapsulationKeySize768 + mlkem.CiphertextSize768
	msg2Len      = 32 + mlkem.CiphertextSize768 + 32
)

var (
	ErrAuth      = errors.New("securechan: server authentication failed (wrong pinned key or MITM)")
	ErrHandshake = errors.New("securechan: malformed handshake")
	ErrRecord    = errors.New("securechan: bad record")
)

// Identity is a server's static ML-KEM-768 key. Persist Seed() (64 bytes) as a secret; publish PublicKey().
type Identity struct{ dk *mlkem.DecapsulationKey768 }

func GenerateIdentity() (*Identity, error) {
	dk, err := mlkem.GenerateKey768()
	if err != nil {
		return nil, err
	}
	return &Identity{dk}, nil
}
func IdentityFromSeed(seed []byte) (*Identity, error) {
	dk, err := mlkem.NewDecapsulationKey768(seed)
	if err != nil {
		return nil, err
	}
	return &Identity{dk}, nil
}
func (id *Identity) Seed() []byte      { return id.dk.Bytes() }
func (id *Identity) PublicKey() []byte { return id.dk.EncapsulationKey().Bytes() }

type keys struct{ c2s, s2c, confirm []byte }

func deriveKeys(k1, k2, k3, transcript []byte) (keys, error) {
	ikm := append(append(append([]byte{}, k1...), k2...), k3...)
	prk, err := hkdf.Extract(sha256.New, ikm, transcript)
	if err != nil {
		return keys{}, err
	}
	var k keys
	for _, e := range []struct {
		dst  *[]byte
		info string
	}{{&k.c2s, "securechan v1 c2s"}, {&k.s2c, "securechan v1 s2c"}, {&k.confirm, "securechan v1 confirm"}} {
		if *e.dst, err = hkdf.Expand(sha256.New, prk, e.info, 32); err != nil {
			return keys{}, err
		}
	}
	return k, nil
}

func confirmTag(key, transcript []byte) []byte {
	m := hmac.New(sha256.New, key)
	m.Write([]byte("server"))
	m.Write(transcript)
	return m.Sum(nil)
}

// Conn is an encrypted, compressed net.Conn.
type Conn struct {
	raw        net.Conn
	send, recv cipher.AEAD
	sendSeq    uint64
	recvSeq    uint64
	wmu, rmu   sync.Mutex
	rbuf       []byte
	// DisableCompression sends every record raw (flag 0). Set it when attacker-influenced data
	// shares the stream with secrets (see KNOWN LIMITS).
	DisableCompression bool
	// Stats: bytes of plaintext handed to Write vs bytes put on the wire (records + headers).
	statMu            sync.Mutex
	plainOut, wireOut uint64
}

func newAEAD(key []byte) (cipher.AEAD, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(b)
}

func newConn(raw net.Conn, sendKey, recvKey []byte) (*Conn, error) {
	s, err := newAEAD(sendKey)
	if err != nil {
		return nil, err
	}
	r, err := newAEAD(recvKey)
	if err != nil {
		return nil, err
	}
	return &Conn{raw: raw, send: s, recv: r}, nil
}

// Dial connects and completes the handshake against a server whose public key is pinned.
func Dial(ctx context.Context, network, addr string, serverPub []byte) (*Conn, error) {
	var d net.Dialer
	raw, err := d.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	c, err := Client(raw, serverPub)
	if err != nil {
		raw.Close()
		return nil, err
	}
	return c, nil
}

// Client runs the client handshake over an existing connection.
func Client(raw net.Conn, serverPub []byte) (*Conn, error) {
	raw.SetDeadline(time.Now().Add(handshakeTTL))
	defer raw.SetDeadline(time.Time{})
	sEK, err := mlkem.NewEncapsulationKey768(serverPub)
	if err != nil {
		return nil, fmt.Errorf("securechan: bad pinned server key: %w", err)
	}
	xPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	eDK, err := mlkem.GenerateKey768()
	if err != nil {
		return nil, err
	}
	k1, ctStatic := sEK.Encapsulate()
	msg1 := make([]byte, 0, msg1Len)
	msg1 = append(msg1, version)
	msg1 = append(msg1, xPriv.PublicKey().Bytes()...)
	msg1 = append(msg1, eDK.EncapsulationKey().Bytes()...)
	msg1 = append(msg1, ctStatic...)
	if _, err := raw.Write(msg1); err != nil {
		return nil, err
	}
	msg2 := make([]byte, msg2Len)
	if _, err := io.ReadFull(raw, msg2); err != nil {
		return nil, err
	}
	sx, ctEph, confirm := msg2[:32], msg2[32:32+mlkem.CiphertextSize768], msg2[32+mlkem.CiphertextSize768:]
	sPub, err := ecdh.X25519().NewPublicKey(sx)
	if err != nil {
		return nil, ErrHandshake
	}
	k2, err := xPriv.ECDH(sPub)
	if err != nil {
		return nil, ErrHandshake
	}
	k3, err := eDK.Decapsulate(ctEph)
	if err != nil {
		return nil, ErrHandshake
	}
	th := sha256.Sum256(append(append([]byte{}, msg1...), msg2[:32+mlkem.CiphertextSize768]...))
	k, err := deriveKeys(k1, k2, k3, th[:])
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(confirm, confirmTag(k.confirm, th[:])) {
		return nil, ErrAuth
	}
	return newConn(raw, k.c2s, k.s2c)
}

// Server runs the server handshake over an accepted connection.
func Server(raw net.Conn, id *Identity) (*Conn, error) {
	raw.SetDeadline(time.Now().Add(handshakeTTL))
	defer raw.SetDeadline(time.Time{})
	msg1 := make([]byte, msg1Len)
	if _, err := io.ReadFull(raw, msg1); err != nil {
		return nil, err
	}
	if msg1[0] != version {
		return nil, ErrHandshake
	}
	cx := msg1[1:33]
	cEKb := msg1[33 : 33+mlkem.EncapsulationKeySize768]
	ctStatic := msg1[33+mlkem.EncapsulationKeySize768:]
	cPub, err := ecdh.X25519().NewPublicKey(cx)
	if err != nil {
		return nil, ErrHandshake
	}
	cEK, err := mlkem.NewEncapsulationKey768(cEKb)
	if err != nil {
		return nil, ErrHandshake
	}
	k1, err := id.dk.Decapsulate(ctStatic)
	if err != nil {
		return nil, ErrHandshake
	}
	xPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	k2, err := xPriv.ECDH(cPub)
	if err != nil {
		return nil, ErrHandshake
	}
	k3, ctEph := cEK.Encapsulate()
	part := append(append([]byte{}, xPriv.PublicKey().Bytes()...), ctEph...)
	th := sha256.Sum256(append(append([]byte{}, msg1...), part...))
	k, err := deriveKeys(k1, k2, k3, th[:])
	if err != nil {
		return nil, err
	}
	if _, err := raw.Write(append(part, confirmTag(k.confirm, th[:])...)); err != nil {
		return nil, err
	}
	return newConn(raw, k.s2c, k.c2s)
}

func nonce(seq uint64) []byte {
	n := make([]byte, 12)
	binary.BigEndian.PutUint64(n[4:], seq)
	return n
}

func (c *Conn) writeRecord(p []byte) error {
	body := append([]byte{0}, p...)
	if !c.DisableCompression && len(p) >= minCompress {
		comp := lz4Compress(nil, p)
		if len(comp)+4 < len(p) {
			body = binary.AppendUvarint([]byte{1}, uint64(len(p)))
			body = append(body, comp...)
		}
	}
	if c.sendSeq >= maxRecords {
		return errors.New("securechan: record limit reached, reconnect")
	}
	hdr := make([]byte, 4)
	binary.BigEndian.PutUint32(hdr, uint32(len(body)+c.send.Overhead()))
	rec := c.send.Seal(hdr, nonce(c.sendSeq), body, hdr)
	c.sendSeq++
	c.statMu.Lock()
	c.plainOut += uint64(len(p))
	c.wireOut += uint64(len(rec))
	c.statMu.Unlock()
	_, err := c.raw.Write(rec)
	return err
}

// Write sends p, splitting into records of at most 64 KiB.
func (c *Conn) Write(p []byte) (int, error) {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	n := 0
	for len(p) > 0 {
		chunk := p
		if len(chunk) > maxChunk {
			chunk = chunk[:maxChunk]
		}
		if err := c.writeRecord(chunk); err != nil {
			return n, err
		}
		n += len(chunk)
		p = p[len(chunk):]
	}
	return n, nil
}

func (c *Conn) readRecord() ([]byte, error) {
	hdr := make([]byte, 4)
	if _, err := io.ReadFull(c.raw, hdr); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(hdr)
	if n < uint32(c.recv.Overhead())+1 || n > maxCT {
		return nil, ErrRecord
	}
	ct := make([]byte, n)
	if _, err := io.ReadFull(c.raw, ct); err != nil {
		return nil, err
	}
	if c.recvSeq >= maxRecords {
		return nil, errors.New("securechan: record limit reached")
	}
	pt, err := c.recv.Open(nil, nonce(c.recvSeq), ct, hdr)
	if err != nil {
		return nil, ErrRecord
	}
	c.recvSeq++
	switch pt[0] {
	case 0:
		return pt[1:], nil
	case 1:
		rawLen, k := binary.Uvarint(pt[1:])
		if k <= 0 || rawLen == 0 || rawLen > maxChunk { // bound BEFORE allocating: decompression-bomb guard
			return nil, ErrRecord
		}
		return lz4Decompress(pt[1+k:], int(rawLen))
	}
	return nil, ErrRecord
}

func (c *Conn) Read(p []byte) (int, error) {
	c.rmu.Lock()
	defer c.rmu.Unlock()
	for len(c.rbuf) == 0 {
		b, err := c.readRecord()
		if err != nil {
			return 0, err
		}
		c.rbuf = b
	}
	n := copy(p, c.rbuf)
	c.rbuf = c.rbuf[n:]
	return n, nil
}

// Stats returns plaintext bytes written and bytes actually sent on the wire (after compression + encryption).
func (c *Conn) Stats() (plain, wire uint64) {
	c.statMu.Lock()
	defer c.statMu.Unlock()
	return c.plainOut, c.wireOut
}

func (c *Conn) Close() error                       { return c.raw.Close() }
func (c *Conn) LocalAddr() net.Addr                { return c.raw.LocalAddr() }
func (c *Conn) RemoteAddr() net.Addr               { return c.raw.RemoteAddr() }
func (c *Conn) SetDeadline(t time.Time) error      { return c.raw.SetDeadline(t) }
func (c *Conn) SetReadDeadline(t time.Time) error  { return c.raw.SetReadDeadline(t) }
func (c *Conn) SetWriteDeadline(t time.Time) error { return c.raw.SetWriteDeadline(t) }

// Listener wraps a net.Listener; the handshake runs lazily on first Read/Write so a slow client
// cannot stall Accept.
type Listener struct {
	net.Listener
	id *Identity
}

func Listen(network, addr string, id *Identity) (*Listener, error) {
	l, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	return &Listener{l, id}, nil
}

type lazyServerConn struct {
	raw  net.Conn
	id   *Identity
	once sync.Once
	c    *Conn
	err  error
}

func (l *Listener) Accept() (net.Conn, error) {
	raw, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &lazyServerConn{raw: raw, id: l.id}, nil
}

func (z *lazyServerConn) up() (*Conn, error) {
	z.once.Do(func() { z.c, z.err = Server(z.raw, z.id) })
	return z.c, z.err
}
func (z *lazyServerConn) Read(p []byte) (int, error) {
	c, err := z.up()
	if err != nil {
		z.raw.Close()
		return 0, err
	}
	return c.Read(p)
}
func (z *lazyServerConn) Write(p []byte) (int, error) {
	c, err := z.up()
	if err != nil {
		z.raw.Close()
		return 0, err
	}
	return c.Write(p)
}
func (z *lazyServerConn) Close() error                       { return z.raw.Close() }
func (z *lazyServerConn) LocalAddr() net.Addr                { return z.raw.LocalAddr() }
func (z *lazyServerConn) RemoteAddr() net.Addr               { return z.raw.RemoteAddr() }
func (z *lazyServerConn) SetDeadline(t time.Time) error      { return z.raw.SetDeadline(t) }
func (z *lazyServerConn) SetReadDeadline(t time.Time) error  { return z.raw.SetReadDeadline(t) }
func (z *lazyServerConn) SetWriteDeadline(t time.Time) error { return z.raw.SetWriteDeadline(t) }
