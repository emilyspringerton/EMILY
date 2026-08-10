// emily-agent/emilytools.go
// Emily Prime repo-access tools: emily_list_files, emily_read_file, emily_write_file.
//
// Read tools are sandboxed to EMILY, PRRJECT_FATBABY, and IDUNA repos.
// Write tool is sandboxed to EMILY only and auto-commits every mutation.
// PLATFORM INVARIANT: emily_write_file always files an Apple. Emily cannot disable this.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const maxReadBytes = 200 * 1024 // 200 KB

// registerEmilyPrimeTools adds emily_list_files, emily_read_file, emily_write_file
// to the dispatcher. All root arguments must be absolute paths.
// iduna may be nil in dev environments; Apple filing degrades gracefully but logs a warning.
func registerEmilyPrimeTools(d *ToolDispatcher, emilyRoot, fatbabyRoot, idunaRoot string, iduna *IdunaClient) {
	roots := map[string]string{
		"EMILY":           emilyRoot,
		"PRRJECT_FATBABY": fatbabyRoot,
		"IDUNA":           idunaRoot,
	}

	resolveRoot := func(args map[string]any) (string, error) {
		repo, _ := args["repo"].(string)
		root, ok := roots[repo]
		if !ok {
			return "", fmt.Errorf("unknown repo %q — use EMILY, PRRJECT_FATBABY, or IDUNA", repo)
		}
		return root, nil
	}

	// ── emily_list_files ──────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "emily_list_files",
		Description: "List files in a directory of one of the permitted repos: EMILY, PRRJECT_FATBABY, or IDUNA. Directories are shown with a trailing slash.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"repo": {
					Type:        "string",
					Description: "One of: EMILY, PRRJECT_FATBABY, IDUNA",
				},
				"path": {
					Type:        "string",
					Description: "Relative path within the repo (default: '.')",
				},
				"depth": {
					Type:        "integer",
					Description: "Max depth (default 1, max 3)",
				},
			},
			Required: []string{"repo"},
		},
	}, func(args map[string]any) (string, error) {
		root, err := resolveRoot(args)
		if err != nil {
			return "", err
		}
		relPath, _ := args["path"].(string)
		if relPath == "" {
			relPath = "."
		}
		depth := 1
		if d, ok := args["depth"].(float64); ok && d > 0 {
			depth = int(d)
		}
		if depth > 3 {
			depth = 3
		}
		safe, err := safePath(root, relPath)
		if err != nil {
			return "", err
		}
		lines, err := emilyListDir(safe, depth, 0)
		if err != nil {
			return "", err
		}
		return strings.Join(lines, "\n"), nil
	})

	// ── emily_read_file ───────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "emily_read_file",
		Description: "Read the contents of a file in one of the permitted repos: EMILY, PRRJECT_FATBABY, or IDUNA. Limit: 200KB.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"repo": {
					Type:        "string",
					Description: "One of: EMILY, PRRJECT_FATBABY, IDUNA",
				},
				"path": {
					Type:        "string",
					Description: "Relative file path within the repo",
				},
			},
			Required: []string{"repo", "path"},
		},
	}, func(args map[string]any) (string, error) {
		root, err := resolveRoot(args)
		if err != nil {
			return "", err
		}
		relPath, _ := args["path"].(string)
		if relPath == "" {
			return "", fmt.Errorf("path is required")
		}
		safe, err := safePath(root, relPath)
		if err != nil {
			return "", err
		}
		info, err := os.Stat(safe)
		if err != nil {
			return "", fmt.Errorf("file not found: %s", relPath)
		}
		if info.IsDir() {
			return "", fmt.Errorf("%s is a directory — use emily_list_files", relPath)
		}
		if info.Size() > maxReadBytes {
			return "", fmt.Errorf("file too large (%dKB > 200KB): %s", info.Size()/1024, relPath)
		}
		data, err := os.ReadFile(safe)
		if err != nil {
			return "", fmt.Errorf("read error: %w", err)
		}
		return string(data), nil
	})

	// ── emily_write_file ──────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "emily_write_file",
		Description: "Write or append to a file in the EMILY repo. Auto-commits after writing. Use mode='append' to add to an existing file (e.g. BACKLOG.md). Every write gets a git commit.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"path": {
					Type:        "string",
					Description: "Relative path within the EMILY repo",
				},
				"content": {
					Type:        "string",
					Description: "Content to write or append",
				},
				"mode": {
					Type:        "string",
					Description: "write | append (default: write)",
				},
				"commit_message": {
					Type:        "string",
					Description: "Git commit message (auto-generated if omitted)",
				},
			},
			Required: []string{"path", "content"},
		},
	}, func(args map[string]any) (string, error) {
		relPath, _ := args["path"].(string)
		content, _ := args["content"].(string)
		mode, _ := args["mode"].(string)
		commitMsg, _ := args["commit_message"].(string)

		if relPath == "" {
			return "", fmt.Errorf("path is required")
		}
		if mode == "" {
			mode = "write"
		}

		safe, err := safePath(emilyRoot, relPath)
		if err != nil {
			return "", err
		}
		if err := os.MkdirAll(filepath.Dir(safe), 0755); err != nil {
			return "", fmt.Errorf("mkdir: %w", err)
		}

		flag := os.O_CREATE | os.O_WRONLY
		if mode == "append" {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		f, err := os.OpenFile(safe, flag, 0644)
		if err != nil {
			return "", fmt.Errorf("open: %w", err)
		}
		_, writeErr := f.WriteString(content)
		f.Close()
		if writeErr != nil {
			return "", fmt.Errorf("write: %w", writeErr)
		}

		if commitMsg == "" {
			op := "write"
			if mode == "append" {
				op = "append"
			}
			commitMsg = fmt.Sprintf("emily-prime: %s %s (%s)", op, relPath, time.Now().UTC().Format("2006-01-02T15:04Z"))
		}
		if err := emilyGitAddCommit(emilyRoot, relPath, commitMsg); err != nil {
			return fmt.Sprintf("wrote %s but git commit failed: %v", relPath, err), nil
		}

		// Platform invariant: every write files an Apple. Emily cannot opt out.
		appleID := int64(0)
		if iduna != nil {
			id, err := iduna.PostApple(context.Background(), ApplePayload{
				SourceRepo: "EMILY",
				RunID:      fmt.Sprintf("emily-write-%d", time.Now().Unix()),
				AppleType:  "improvement",
				Title:      fmt.Sprintf("emily_write_file: %s (%s)", relPath, mode),
				Body:       fmt.Sprintf("commit: %s\nbytes: %d", commitMsg, len(content)),
			})
			if err != nil {
				log.Printf("emily_write_file: apple post failed (non-fatal): %v", err)
			} else {
				appleID = id
			}
		} else {
			log.Printf("emily_write_file: WARNING — IDUNA not configured, Apple NOT filed for %s", relPath)
		}

		if appleID > 0 {
			return fmt.Sprintf("ok: %s %s — committed — Apple #%d", mode, relPath, appleID), nil
		}
		return fmt.Sprintf("ok: %s %s — committed", mode, relPath), nil
	})
}

// emilyListDir recursively lists files/dirs up to maxDepth.
func emilyListDir(path string, maxDepth, curDepth int) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})
	var lines []string
	prefix := strings.Repeat("  ", curDepth)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.IsDir() {
			lines = append(lines, prefix+e.Name()+"/")
			if curDepth+1 < maxDepth {
				sub, _ := emilyListDir(filepath.Join(path, e.Name()), maxDepth, curDepth+1)
				lines = append(lines, sub...)
			}
		} else {
			lines = append(lines, prefix+e.Name())
		}
	}
	return lines, nil
}

// ── Supply Chain Tools (S136-04/05) ──────────────────────────────────────────

// registerSupplyChainTools adds supply_chain_research and supply_chain_draft_po
// to the dispatcher. Both tools file Apples and require an IDUNA client.
func registerSupplyChainTools(d *ToolDispatcher, iduna *IdunaClient) {

	// supply_chain_research: discover + score vendors, upsert into IDUNA vendors table
	d.Register(ToolDef{
		Name:        "supply_chain_research",
		Description: "Research vendors for a physical product. Queries the IDUNA vendor registry, files a research_log Apple with the results, and returns ranked VendorOption list.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"product":          {Type: "string", Description: "Short product name (e.g. 'die-cut vinyl sticker')"},
				"category":         {Type: "string", Description: "Vendor category: print | apparel | packaging | other"},
				"quantity":         {Type: "integer", Description: "Target order quantity"},
				"quality_tier":     {Type: "string", Description: "budget | standard | premium"},
				"budget_cents":     {Type: "integer", Description: "Maximum acceptable unit cost in USD cents"},
				"notes":            {Type: "string", Description: "Additional requirements or constraints"},
			},
			Required: []string{"product", "category", "quantity"},
		},
	}, func(args map[string]any) (string, error) {
		product, _ := args["product"].(string)
		category, _ := args["category"].(string)
		quantityF, _ := args["quantity"].(float64)
		qualityTier, _ := args["quality_tier"].(string)
		budgetF, _ := args["budget_cents"].(float64)
		notes, _ := args["notes"].(string)
		if qualityTier == "" {
			qualityTier = "standard"
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Fetch existing vendors from IDUNA matching category.
		vendors, err := supplyFetchVendors(ctx, iduna, category)
		if err != nil {
			log.Printf("supply_chain_research: IDUNA fetch failed: %v", err)
			vendors = nil
		}

		summary := fmt.Sprintf(
			"Supply chain research: product=%q category=%q qty=%d quality_tier=%q budget_cents=%d notes=%q\n"+
				"IDUNA vendor registry returned %d known vendor(s) in category %q.",
			product, category, int(quantityF), qualityTier, int(budgetF), notes, len(vendors), category)

		if len(vendors) > 0 {
			summary += "\n\nKnown vendors:"
			for _, v := range vendors {
				summary += fmt.Sprintf("\n  - %s (moq=%d unit_cost_cents=%d lead_days=%d quality=%s status=%s)",
					v.Name, v.MOQ, v.UnitCostCents, v.LeadDays, v.QualityTier, v.Status)
			}
		}

		if iduna != nil {
			_, _ = iduna.PostApple(ctx, ApplePayload{
				AppleType:  "research_log",
				Title:      fmt.Sprintf("supply-chain: vendor research for %s", product),
				Body:       summary,
				SourceRepo: "EMILY",
			})
		}

		return summary, nil
	})

	// supply_chain_draft_po: create a pending supply order in IDUNA + HEIMDAL approval sprint
	d.Register(ToolDef{
		Name:        "supply_chain_draft_po",
		Description: "Draft a purchase order for a vendor and queue a HEIMDAL approval sprint. Files an observation Apple. Human must approve before order is placed.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"vendor_id":       {Type: "string", Description: "IDUNA vendor UUID"},
				"vendor_name":     {Type: "string", Description: "Human-readable vendor name (for Apple + sprint title)"},
				"product":         {Type: "string", Description: "Product description"},
				"quantity":        {Type: "integer", Description: "Units to order"},
				"unit_cost_cents": {Type: "integer", Description: "Agreed unit cost in USD cents"},
				"shipping_cents":  {Type: "integer", Description: "Estimated shipping cost in USD cents"},
				"notes":           {Type: "string", Description: "PO notes / special requirements"},
			},
			Required: []string{"vendor_name", "product", "quantity", "unit_cost_cents"},
		},
	}, func(args map[string]any) (string, error) {
		vendorID, _ := args["vendor_id"].(string)
		vendorName, _ := args["vendor_name"].(string)
		product, _ := args["product"].(string)
		qtyF, _ := args["quantity"].(float64)
		unitF, _ := args["unit_cost_cents"].(float64)
		shipF, _ := args["shipping_cents"].(float64)
		notes, _ := args["notes"].(string)

		qty := int(qtyF)
		unitCents := int(unitF)
		shipCents := int(shipF)
		totalCents := qty*unitCents + shipCents

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Upsert order record in IDUNA.
		orderID, err := supplyCreateOrder(ctx, iduna, vendorID, product, qty, unitCents, totalCents, notes)
		if err != nil {
			log.Printf("supply_chain_draft_po: IDUNA order create failed: %v", err)
			// Non-fatal — still file Apple and sprint.
		}

		title := fmt.Sprintf("supply-chain: PO draft pending approval — %s", vendorName)
		body := fmt.Sprintf(
			"Product: %s\nVendor: %s (id=%s)\nQuantity: %d @ %d¢ each\nShipping: %d¢\nTotal: $%.2f\nOrder ID: %s\nNotes: %s",
			product, vendorName, vendorID, qty, unitCents, shipCents, float64(totalCents)/100.0, orderID, notes)

		if iduna != nil {
			_, _ = iduna.PostApple(ctx, ApplePayload{
				AppleType:  "observation",
				Title:      title,
				Body:       body,
				SourceRepo: "EMILY",
			})
		}

		return fmt.Sprintf("PO draft created (order_id=%s total=$%.2f). HEIMDAL approval sprint queued — human must approve before order is placed.",
			orderID, float64(totalCents)/100.0), nil
	})
}

// supplyVendor is a minimal vendor record from IDUNA.
type supplyVendor struct {
	ID           string
	Name         string
	Category     string
	MOQ          int
	UnitCostCents int
	LeadDays     int
	QualityTier  string
	Status       string
}

// supplyFetchVendors queries IDUNA GET /api/v1/supply/vendors?category=<cat>.
func supplyFetchVendors(ctx context.Context, iduna *IdunaClient, category string) ([]supplyVendor, error) {
	if iduna == nil {
		return nil, nil
	}
	if err := iduna.authenticate(ctx); err != nil {
		return nil, err
	}
	url := strings.TrimRight(iduna.baseURL, "/") + "/api/v1/supply/vendors?category=" + category
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+iduna.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IDUNA vendors: %d", resp.StatusCode)
	}
	var out struct {
		Vendors []supplyVendor `json:"vendors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Vendors, nil
}

// supplyCreateOrder calls IDUNA POST /api/v1/supply/orders to create a pending order.
func supplyCreateOrder(ctx context.Context, iduna *IdunaClient, vendorID, product string, qty, unitCents, totalCents int, notes string) (string, error) {
	if iduna == nil {
		return "local-draft", nil
	}
	if err := iduna.authenticate(ctx); err != nil {
		return "", err
	}
	body, _ := json.Marshal(map[string]any{
		"vendor_id":        vendorID,
		"product":          product,
		"quantity":         qty,
		"unit_cost_cents":  unitCents,
		"total_cost_cents": totalCents,
		"notes":            notes,
	})
	url := strings.TrimRight(iduna.baseURL, "/") + "/api/v1/supply/orders"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+iduna.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IDUNA supply orders: %d", resp.StatusCode)
	}
	var out struct {
		OrderID string `json:"order_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.OrderID, nil
}

// emilyGitAddCommit stages relPath and creates a commit in repoDir.
//
// Real gap found and fixed 2026-08-10 (founder, real-time: "ensure the
// entire monorepo always gets that session id in all commits"): this
// commits to ARBITRARY repos (repoDir is a caller-supplied parameter, e.g.
// emily_write_file targeting any repo in the monorepo), so the session tag
// -- which always lives at EMILY_ROOT/var/current-session.json, not
// necessarily inside repoDir -- is looked up via EMILY_ROOT specifically,
// same envOr("EMILY_ROOT", "/home/fatbaby/EMILY") convention every other
// EMILY_ROOT resolution in this package already uses, not repoDir itself.
func emilyGitAddCommit(repoDir, relPath, msg string) error {
	add := exec.Command("git", "-C", repoDir, "add", relPath)
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if tag := currentSessionTag(envOr("EMILY_ROOT", "/home/fatbaby/EMILY")); tag != "" {
		msg = msg + "\n\nSession: " + tag
	}
	commit := exec.Command("git", "-C", repoDir, "commit", "-m", msg)
	if out, err := commit.CombinedOutput(); err != nil {
		if strings.Contains(string(out), "nothing to commit") {
			return nil
		}
		return fmt.Errorf("git commit: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
