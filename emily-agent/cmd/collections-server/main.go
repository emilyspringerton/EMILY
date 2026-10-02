// collections-server — standalone, off-box-able server for the generic collections API
// (golden docs first). Deliberately NOT emily-agent: no RSI cron, no shared state writes, no
// LLM calls — just IDUNA-JWT-gated, read-only documents from a checkout of GOLDEN_DOCS.
//
//	GOLDEN_DOCS_DIR path to a GOLDEN_DOCS checkout (context/golden-docs-index.md + docs/<repo>/...)
//	IDUNA_JWKS_URL  required; without it every request fails closed with 503
//	PORT            default 8087
package main

import (
	"log"
	"net/http"
	"os"

	"emily-agent/collections"
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
	log.Printf("collections-server on :%s (root %s)", port, root)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
