// collections-server — standalone, off-box-able server for the generic collections API
// (golden docs first). Deliberately NOT emily-agent: no RSI cron, no shared state writes, no
// LLM calls — just IDUNA-JWT-gated, read-only documents from a checkout of GOLDEN_DOCS.
//
//	GOLDEN_DOCS_DIR path to a GOLDEN_DOCS checkout (context/golden-docs-index.md + docs/<repo>/...)
//	IDUNA_JWKS_URL  required; without it every request fails closed with 503
//	PORT            default 8087
//	SECURECHAN_ADDR optional, e.g. ":8443": ALSO serve the API over the PQ-hybrid compressed channel
//	SECURECHAN_SEED base64 64-byte ML-KEM seed (secret). Unset => an ephemeral identity is generated and its
//	                public key logged (fine for tests; clients must re-pin on every restart).
//	NOTE: a raw securechan listener is plain TCP, so GKE's HTTP(S) Ingress cannot front it; it needs
//	an L4 LB (a second, billed LB) or a WebSocket tunnel. See KUBERNETES_SERVICE_MIGRATION_NORTHSTAR.md.
package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"emily-agent/collections"
	"emily-agent/securechan"
)

func main() {
	root := os.Getenv("GOLDEN_DOCS_DIR")
	if root == "" {
		log.Fatal("GOLDEN_DOCS_DIR is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}
	reg := collections.New(os.Getenv("IDUNA_JWKS_URL"))
	reg.Register("golden", &collections.Golden{IndexPath: root + "/context/golden-docs-index.md", BaseDir: root + "/docs"})
	mux := http.NewServeMux()
	mux.Handle("/api/v1/emily/collections", reg)
	mux.Handle("/api/v1/emily/collections/", reg)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	if addr := os.Getenv("SECURECHAN_ADDR"); addr != "" {
		var id *securechan.Identity
		var err error
		if seed := os.Getenv("SECURECHAN_SEED"); seed != "" {
			raw, derr := base64.StdEncoding.DecodeString(seed)
			if derr != nil {
				log.Fatalf("SECURECHAN_SEED: %v", derr)
			}
			id, err = securechan.IdentityFromSeed(raw)
		} else {
			log.Printf("SECURECHAN_SEED unset: using an EPHEMERAL identity (clients must re-pin after restart)")
			id, err = securechan.GenerateIdentity()
		}
		if err != nil {
			log.Fatalf("securechan identity: %v", err)
		}
		l, err := securechan.Listen("tcp", addr, id)
		if err != nil {
			log.Fatalf("securechan listen: %v", err)
		}
		log.Printf("securechan on %s, pin this server key (base64): %s", addr, base64.StdEncoding.EncodeToString(id.PublicKey()))
		go func() { log.Fatal(http.Serve(l, mux)) }()
	}
	log.Printf("collections-server on :%s (root %s)", port, root)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
