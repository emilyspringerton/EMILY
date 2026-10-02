// api_collections.go — generic read-only "collection" API (K8s migration, VS0).
//
// Founder real-time, 2026-10-02: file/host-path coupling between services can't
// cross pod boundaries, so data that is read off local disk becomes an API.
// Make it generic: a Collection is a named set of items with List/Get/ETag.
// Golden docs are the first collection; other var/ read models plug in the
// same way (see docs/KUBERNETES_SERVICE_MIGRATION_NORTHSTAR.md).
//
//	GET /api/v1/emily/collections                      -> collection names
//	GET /api/v1/emily/collections/{coll}/items         -> item metadata list
//	GET /api/v1/emily/collections/{coll}/items/{id}    -> item body (ETag / 304)
//
// Auth: IDUNA ES256 JWT (Bearer) with permission "emily.collections.read".
// Fails CLOSED (503) when IDUNA_JWKS_URL is unset: this is the internet-facing path.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"emily-agent/idunaauth"
)

const collectionsReadPerm = "emily.collections.read"

// ItemMeta describes one item without its body.
type ItemMeta struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
	Tier        string `json:"tier,omitempty"`
	SHA256      string `json:"sha256"`
	Size        int    `json:"size"`
}

// Collection is the generic read-side contract.
type Collection interface {
	List() ([]ItemMeta, error)
	Get(id string) (body []byte, meta ItemMeta, ok bool, err error)
}

// CollectionRegistry holds named collections and the JWT verifier.
type CollectionRegistry struct {
	verifier *idunaauth.Verifier
	colls    map[string]Collection
}

func NewCollectionRegistry(jwksURL string) *CollectionRegistry {
	r := &CollectionRegistry{colls: map[string]Collection{}}
	if jwksURL != "" {
		v, err := idunaauth.NewVerifier(jwksURL)
		if err != nil {
			log.Printf("collections: JWKS init failed (%v) — API will return 503", err)
		} else {
			r.verifier = v
		}
	}
	return r
}

func (r *CollectionRegistry) Register(name string, c Collection) { r.colls[name] = c }

func (r *CollectionRegistry) auth(w http.ResponseWriter, req *http.Request) bool {
	if r.verifier == nil {
		http.Error(w, "collections API disabled: IDUNA JWKS not configured", http.StatusServiceUnavailable)
		return false
	}
	tok := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if tok == "" || tok == req.Header.Get("Authorization") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	claims, err := r.verifier.Verify(tok)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	if !claims.HasPermission(collectionsReadPerm) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return false
	}
	return true
}

func (r *CollectionRegistry) handle(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !r.auth(w, req) {
		return
	}
	rest := strings.Trim(strings.TrimPrefix(req.URL.Path, "/api/v1/emily/collections"), "/")
	parts := strings.Split(rest, "/")
	if rest == "" {
		names := make([]string, 0, len(r.colls))
		for n := range r.colls {
			names = append(names, n)
		}
		sort.Strings(names)
		writeJSON(w, map[string]any{"collections": names})
		return
	}
	c, ok := r.colls[parts[0]]
	if !ok || len(parts) < 2 || parts[1] != "items" {
		http.NotFound(w, req)
		return
	}
	if len(parts) == 2 {
		items, err := c.List()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]any{"items": items})
		return
	}
	body, meta, found, err := c.Get(strings.Join(parts[2:], "/"))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.NotFound(w, req)
		return
	}
	etag := `"` + meta.SHA256 + `"`
	w.Header().Set("ETag", etag)
	if req.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Write(body)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// ---- golden docs collection ----

// goldenCollection serves exactly the files named in golden-docs-index.md (all
// tiers) plus the compiled context. IDs come from the index, never from the
// caller's path, so there is no path traversal surface.
type goldenCollection struct {
	emilyRoot string
}

type goldenRow struct{ name, path, tier, desc string }

func (g *goldenCollection) rows() []goldenRow {
	data, err := os.ReadFile(g.emilyRoot + "/context/golden-docs-index.md")
	if err != nil {
		return nil
	}
	var out []goldenRow
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") || strings.HasPrefix(line, "|---") || strings.HasPrefix(line, "| name") {
			continue
		}
		f := strings.Split(line, "|")
		if len(f) < 7 {
			continue
		}
		row := goldenRow{strings.TrimSpace(f[1]), strings.TrimSpace(f[2]), strings.TrimSpace(f[3]), strings.TrimSpace(f[5])}
		if row.name == "" || row.path == "" || strings.Contains(row.path, "..") {
			continue
		}
		out = append(out, row)
	}
	return out
}

func (g *goldenCollection) read(row goldenRow) ([]byte, ItemMeta, bool) {
	b, err := os.ReadFile(os_base(g.emilyRoot) + "/" + row.path)
	if err != nil {
		return nil, ItemMeta{}, false
	}
	sum := sha256.Sum256(b)
	return b, ItemMeta{ID: row.name, Description: row.desc, Tier: row.tier, SHA256: hex.EncodeToString(sum[:]), Size: len(b)}, true
}

func os_base(emilyRoot string) string { return emilyRoot[:strings.LastIndex(emilyRoot, "/")] }

func (g *goldenCollection) List() ([]ItemMeta, error) {
	var items []ItemMeta
	for _, row := range g.rows() {
		if _, m, ok := g.read(row); ok {
			items = append(items, m)
		}
	}
	return items, nil
}

func (g *goldenCollection) Get(id string) ([]byte, ItemMeta, bool, error) {
	for _, row := range g.rows() {
		if row.name == id {
			b, m, ok := g.read(row)
			return b, m, ok, nil
		}
	}
	return nil, ItemMeta{}, false, nil
}
