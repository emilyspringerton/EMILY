// emily-agent/webaudit.go
// Web audit tool — Emily Prime can fetch and inspect web pages as a real HTTP client.
//
// Uses stdlib only: net/http for requests, regexp for link extraction.
// Enables Emily Prime to audit the FatBaby newssite (localhost:8082) and other
// HTTP endpoints as part of RSI cycles and as a pre-send guard before MJOLNIR pushes.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// auditClient is a shared HTTP client with a reasonable timeout and no redirect loops.
var auditClient = &http.Client{
	Timeout: 15 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many redirects")
		}
		return nil
	},
}

// linkRe matches href="..." and src="..." in HTML.
var linkRe = regexp.MustCompile(`(?i)(?:href|src)=["']([^"'#][^"']*)["']`)

// titleRe extracts the page <title>.
var titleRe = regexp.MustCompile(`(?i)<title[^>]*>([^<]{1,200})</title>`)

// h1Re extracts the first <h1>.
var h1Re = regexp.MustCompile(`(?i)<h1[^>]*>([^<]{1,200})</h1>`)

// AuditResult is the structured output of a web audit.
type AuditResult struct {
	URL          string            `json:"url"`
	StatusCode   int               `json:"status_code"`
	ContentType  string            `json:"content_type"`
	Title        string            `json:"title,omitempty"`
	H1           string            `json:"h1,omitempty"`
	BodySizeKB   int               `json:"body_size_kb"`
	LinkCount    int               `json:"link_count"`
	BrokenLinks  []BrokenLink      `json:"broken_links,omitempty"`
	LinksChecked int               `json:"links_checked"`
	Errors       []string          `json:"errors,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	DurationMS   int64             `json:"duration_ms"`
	AuditedAt    string            `json:"audited_at"`
}

type BrokenLink struct {
	Href       string `json:"href"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
}

func auditURL(ctx context.Context, rawURL string, checkLinks bool, maxLinks int) (AuditResult, error) {
	start := time.Now()
	result := AuditResult{
		URL:       rawURL,
		AuditedAt: start.UTC().Format(time.RFC3339),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return result, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "Emily-Prime-Audit/1.0 (EINHORN_INDUSTRIAL)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,*/*")

	resp, err := auditClient.Do(req)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("fetch error: %v", err))
		return result, nil
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")

	// Capture a few key response headers for the report.
	result.Headers = map[string]string{}
	for _, h := range []string{"Server", "Cache-Control", "X-Content-Type-Options", "Content-Security-Policy"} {
		if v := resp.Header.Get(h); v != "" {
			result.Headers[h] = v
		}
	}

	// Read body (cap at 1 MB to avoid surprises).
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("read body: %v", err))
		result.DurationMS = time.Since(start).Milliseconds()
		return result, nil
	}
	body := string(bodyBytes)
	result.BodySizeKB = len(bodyBytes) / 1024

	if m := titleRe.FindStringSubmatch(body); len(m) > 1 {
		result.Title = strings.TrimSpace(m[1])
	}
	if m := h1Re.FindStringSubmatch(body); len(m) > 1 {
		result.H1 = strings.TrimSpace(m[1])
	}

	// Extract links.
	base, _ := url.Parse(rawURL)
	var links []string
	seen := map[string]bool{}
	for _, m := range linkRe.FindAllStringSubmatch(body, -1) {
		raw := m[1]
		if raw == "" || strings.HasPrefix(raw, "mailto:") || strings.HasPrefix(raw, "javascript:") {
			continue
		}
		ref, err := url.Parse(raw)
		if err != nil {
			continue
		}
		abs := base.ResolveReference(ref).String()
		// Only audit same-host links.
		parsed, err := url.Parse(abs)
		if err != nil || parsed.Host != base.Host {
			continue
		}
		if !seen[abs] {
			seen[abs] = true
			links = append(links, abs)
		}
	}
	result.LinkCount = len(links)

	if checkLinks && len(links) > 0 {
		cap := maxLinks
		if cap <= 0 || cap > 30 {
			cap = 30
		}
		if len(links) < cap {
			cap = len(links)
		}
		for _, link := range links[:cap] {
			result.LinksChecked++
			lreq, err := http.NewRequestWithContext(ctx, http.MethodHead, link, nil)
			if err != nil {
				result.BrokenLinks = append(result.BrokenLinks, BrokenLink{Href: link, Error: err.Error()})
				continue
			}
			lreq.Header.Set("User-Agent", "Emily-Prime-Audit/1.0")
			lresp, err := auditClient.Do(lreq)
			if err != nil {
				result.BrokenLinks = append(result.BrokenLinks, BrokenLink{Href: link, Error: err.Error()})
				continue
			}
			lresp.Body.Close()
			if lresp.StatusCode >= 400 {
				result.BrokenLinks = append(result.BrokenLinks, BrokenLink{Href: link, StatusCode: lresp.StatusCode})
			}
		}
	}

	result.DurationMS = time.Since(start).Milliseconds()
	return result, nil
}

func registerWebAuditTools(d *ToolDispatcher) {
	d.Register(ToolDef{
		Name:        "web_audit_url",
		Description: "Audit a web page as a real HTTP client. Fetches the URL, checks status code, extracts title and h1, counts links, optionally validates all same-host links for broken status (4xx/5xx). Use this to audit the FatBaby newssite (http://localhost:8082), SignalAPI (http://localhost:8083), or any other product URL before sending a push notification or after a deploy. Returns a structured JSON report.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"url": {
					Type:        "string",
					Description: "Full URL to audit, e.g. http://localhost:8082 or http://localhost:8082/ticker/AAPL",
				},
				"check_links": {
					Type:        "string",
					Description: "Set to 'true' to also HEAD-check all same-host links found on the page. Adds latency proportional to link count.",
				},
				"max_links": {
					Type:        "string",
					Description: "Maximum number of links to check when check_links=true. Default 20, max 30.",
				},
			},
			Required: []string{"url"},
		},
	}, func(args map[string]any) (string, error) {
		rawURL, _ := args["url"].(string)
		if rawURL == "" {
			return "", fmt.Errorf("url is required")
		}
		checkLinks := args["check_links"] == "true" || args["check_links"] == true
		maxLinks := 20
		if v, ok := args["max_links"].(string); ok && v != "" {
			fmt.Sscanf(v, "%d", &maxLinks)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := auditURL(ctx, rawURL, checkLinks, maxLinks)
		if err != nil {
			return "", fmt.Errorf("audit: %w", err)
		}

		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", err
		}

		// Prepend a human summary for the LLM.
		summary := fmt.Sprintf("Audit of %s: HTTP %d, title=%q, %d links found, %d links checked, %d broken.\n\n",
			rawURL, result.StatusCode, result.Title, result.LinkCount, result.LinksChecked, len(result.BrokenLinks))
		if len(result.Errors) > 0 {
			summary += "ERRORS: " + strings.Join(result.Errors, "; ") + "\n\n"
		}
		return summary + string(b), nil
	})
}
