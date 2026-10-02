package collections

import (
	"os"
	"testing"
)

// Runs against a real GOLDEN_DOCS checkout when GOLDEN_DOCS_DIR is set (skipped in CI otherwise).
func TestGoldenAgainstRealCheckout(t *testing.T) {
	dir := os.Getenv("GOLDEN_DOCS_DIR")
	if dir == "" {
		t.Skip("GOLDEN_DOCS_DIR not set")
	}
	g := &Golden{IndexPath: dir + "/context/golden-docs-index.md", BaseDir: dir + "/docs"}
	items, _ := g.List()
	if len(items) < 100 {
		t.Fatalf("only %d items served from real checkout", len(items))
	}
	t.Logf("%d golden docs served", len(items))
}
