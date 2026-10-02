package collections

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newGoldenFixture(t *testing.T) *Golden {
	root := filepath.Join(t.TempDir(), "EMILY")
	os.MkdirAll(filepath.Join(root, "context"), 0o755)
	os.WriteFile(filepath.Join(root, "context", "golden-docs-index.md"),
		[]byte("| name | path | tier | budget | description |\n|---|---|---|---|---|\n| DOC1 | doc1.md | 1 | 0 | first |\n| EVIL | ../../etc/passwd | 1 | 0 | bad |\n| NOPIPE | doc2.md | 2 | 0 | row with no trailing pipe\n"), 0o644)
	os.WriteFile(filepath.Join(filepath.Dir(root), "doc1.md"), []byte("# hello"), 0o644)
	os.WriteFile(filepath.Join(filepath.Dir(root), "doc2.md"), []byte("two"), 0o644)
	return GoldenFromMonorepo(root)
}

func TestGoldenCollection(t *testing.T) {
	g := newGoldenFixture(t)
	items, _ := g.List()
	if len(items) != 2 || items[0].ID != "DOC1" || items[0].Size != 7 || items[1].ID != "NOPIPE" {
		t.Fatalf("list = %+v (traversal row must be dropped)", items)
	}
	b, m, ok, _ := g.Get("DOC1")
	if !ok || string(b) != "# hello" || m.SHA256 == "" {
		t.Fatalf("get failed: %v %v", ok, m)
	}
	if _, _, ok, _ := g.Get("EVIL"); ok {
		t.Fatal("traversal row served")
	}
}

func TestCollectionsFailClosedWithoutJWKS(t *testing.T) {
	r := New("")
	r.Register("golden", newGoldenFixture(t))
	for _, hdr := range []string{"", "Bearer x"} {
		req := httptest.NewRequest("GET", "/api/v1/emily/collections/golden/items", nil)
		if hdr != "" {
			req.Header.Set("Authorization", hdr)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("hdr %q: code %d, want 503", hdr, w.Code)
		}
	}
}

func mintJWT(t *testing.T, key *ecdsa.PrivateKey, perms []string) string {
	hj, _ := json.Marshal(map[string]string{"alg": "ES256", "kid": "k1", "typ": "JWT"})
	cj, _ := json.Marshal(map[string]any{"sub": "t", "permissions": perms, "exp": time.Now().Add(time.Hour).Unix()})
	msg := base64.RawURLEncoding.EncodeToString(hj) + "." + base64.RawURLEncoding.EncodeToString(cj)
	d := sha256.Sum256([]byte(msg))
	r, s, _ := ecdsa.Sign(rand.Reader, key, d[:])
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return msg + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func TestCollectionsAuthedFlow(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	x, y := make([]byte, 32), make([]byte, 32)
	key.PublicKey.X.FillBytes(x)
	key.PublicKey.Y.FillBytes(y)
	jwks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{"kty": "EC", "crv": "P-256", "kid": "k1",
			"x": base64.RawURLEncoding.EncodeToString(x), "y": base64.RawURLEncoding.EncodeToString(y)}}})
	}))
	defer jwks.Close()
	r := New(jwks.URL)
	r.Register("golden", newGoldenFixture(t))

	do := func(path, tok, inm string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("GET", path, nil)
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		if inm != "" {
			req.Header.Set("If-None-Match", inm)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	const p = "/api/v1/emily/collections/golden/items/DOC1"
	if c := do(p, "", "").Code; c != 401 {
		t.Fatalf("no token: %d", c)
	}
	if c := do(p, mintJWT(t, key, []string{"other"}), "").Code; c != 403 {
		t.Fatalf("wrong perm: %d", c)
	}
	tok := mintJWT(t, key, []string{collectionsReadPerm})
	w := do(p, tok, "")
	if w.Code != 200 || w.Body.String() != "# hello" || w.Header().Get("ETag") == "" {
		t.Fatalf("get: %d %q", w.Code, w.Body.String())
	}
	if c := do(p, tok, w.Header().Get("ETag")).Code; c != 304 {
		t.Fatalf("etag: %d", c)
	}
	if c := do("/api/v1/emily/collections/golden/items/NOPE", tok, "").Code; c != 404 {
		t.Fatalf("missing: %d", c)
	}
}
