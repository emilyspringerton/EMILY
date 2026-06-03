// emily-agent/collector.go
// Data collection pipeline — Reddit + Wikipedia → quality-scored NDJSON corpus.
//
// Architecture mirrors PRRJECT_FATBABY:
//   - DocStore: append-only NDJSON journals, daily rotation, mutex-protected
//   - seenSet: built once at startup via O(n) scan, O(1) on hot path
//   - Worker pool fan-out: same shape as processor.Run / RunBodyCrawler
//   - Fail-open everywhere: collection failure degrades throughput, never crashes
//
// No external dependencies. Pure stdlib.
//
// Routes wired in main.go:
//   GET/POST /collect  — start (POST), stop (POST ?action=stop), status (GET)

package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------------------
// httpGetWithRetry — exponential backoff for all source fetches.
// Retries on network errors and 429/503 status codes.
// Delay: 1s → 2s → 4s … capped at 30s. Fail-open: returns last error.
// ---------------------------------------------------------------------------

func httpGetWithRetry(ctx context.Context, client *http.Client, reqURL, ua string, maxAttempts int) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
		req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
		if err != nil {
			return nil, err
		}
		if ua != "" {
			req.Header.Set("User-Agent", ua)
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %w", attempt+1, err)
			continue
		}
		if resp.StatusCode == 429 || resp.StatusCode == 503 {
			resp.Body.Close()
			lastErr = fmt.Errorf("attempt %d: status %d", attempt+1, resp.StatusCode)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("all %d attempts failed, last: %w", maxAttempts, lastErr)
}

// ---------------------------------------------------------------------------
// CollectedDoc — canonical training record written to the corpus
// ---------------------------------------------------------------------------

type CollectedDoc struct {
	ID           string    `json:"id"`           // sha256(source:stable_id)[:16 hex]
	SourceName   string    `json:"source_name"`  // "reddit" | "wikipedia"
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	BodyBytes    int       `json:"body_bytes"`
	SourceScore  int       `json:"source_score"` // Reddit: upvotes; Wikipedia: 0
	PublishedAt  time.Time `json:"published_at"`
	CollectedAt  time.Time `json:"collected_at"`
	QualityScore float64   `json:"quality_score"` // 0.0–10.0
	QualityTier  string    `json:"quality_tier"`  // "gold" | "silver" | "bronze"
	Tags         []string  `json:"tags,omitempty"`
}

// ---------------------------------------------------------------------------
// DocStore — append-only NDJSON corpus, daily-rotating journals
// Same design as PRRJECT_FATBABY eventstore.FileStore.
// ---------------------------------------------------------------------------

type docRecord struct {
	Seq        uint64       `json:"seq"`
	Doc        CollectedDoc `json:"doc"`
	AppendedAt time.Time    `json:"appended_at"`
}

type DocStore struct {
	dir     string
	mu      sync.Mutex
	current *os.File
	today   string
	seq     uint64
}

func NewDocStore(dir string) (*DocStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("doc_store mkdir: %w", err)
	}
	s := &DocStore{dir: dir}
	s.seq = s.recoverSeq()
	if err := s.openJournal(time.Now()); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *DocStore) openJournal(t time.Time) error {
	day := t.UTC().Format("2006-01-02")
	if s.today == day && s.current != nil {
		return nil
	}
	if s.current != nil {
		_ = s.current.Close()
	}
	path := filepath.Join(s.dir, day+".ndjson")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("doc_store open_journal path=%s: %w", path, err)
	}
	s.current = f
	s.today = day
	return nil
}

func (s *DocStore) Append(doc CollectedDoc) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.openJournal(time.Now()); err != nil {
		return err
	}
	s.seq++
	rec := docRecord{Seq: s.seq, Doc: doc, AppendedAt: time.Now().UTC()}
	line, err := json.Marshal(rec)
	if err != nil {
		s.seq--
		return err
	}
	if _, err := s.current.Write(append(line, '\n')); err != nil {
		s.seq--
		return err
	}
	return s.current.Sync()
}

// LoadSeen scans all journal files and returns the set of known doc IDs.
// Called once at startup — same pattern as loadBodySeenIDs in PRRJECT_FATBABY.
func (s *DocStore) LoadSeen(ctx context.Context) (map[string]struct{}, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]struct{}{}, nil
		}
		return nil, fmt.Errorf("doc_store read_dir: %w", err)
	}
	seen := map[string]struct{}{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".ndjson") {
			continue
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		path := filepath.Join(s.dir, e.Name())
		if scanErr := scanJournalForIDs(path, seen); scanErr != nil {
			log.Printf("doc_store scan_warn path=%s err=%v", path, scanErr)
		}
	}
	return seen, nil
}

func scanJournalForIDs(path string, seen map[string]struct{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 512*1024), 512*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var rec struct {
			Doc struct {
				ID string `json:"id"`
			} `json:"doc"`
		}
		if err := json.Unmarshal(line, &rec); err != nil {
			continue
		}
		if rec.Doc.ID != "" {
			seen[rec.Doc.ID] = struct{}{}
		}
	}
	return sc.Err()
}

func (s *DocStore) recoverSeq() uint64 {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return 0
	}
	var latest uint64
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".ndjson") {
			continue
		}
		f, err := os.Open(filepath.Join(s.dir, e.Name()))
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 512*1024), 512*1024)
		for sc.Scan() {
			var rec struct{ Seq uint64 `json:"seq"` }
			if err := json.Unmarshal(sc.Bytes(), &rec); err == nil && rec.Seq > latest {
				latest = rec.Seq
			}
		}
		f.Close()
	}
	return latest
}

func (s *DocStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.current != nil {
		return s.current.Close()
	}
	return nil
}

// ---------------------------------------------------------------------------
// Source interface
// ---------------------------------------------------------------------------

type Source interface {
	Name() string
	Fetch(ctx context.Context, since time.Time) ([]CollectedDoc, error)
}

// ---------------------------------------------------------------------------
// RedditSource — public JSON API, no auth required for read-only access
// ---------------------------------------------------------------------------

type RedditSource struct {
	Subreddits  []string
	UserAgent   string
	MinScore    int
	MinComments int
	MinBodyLen  int
	client      *http.Client
}

func (r *RedditSource) Name() string { return "reddit" }

func (r *RedditSource) Fetch(ctx context.Context, since time.Time) ([]CollectedDoc, error) {
	if r.client == nil {
		r.client = &http.Client{Timeout: 15 * time.Second}
	}
	ua := r.UserAgent
	if ua == "" {
		ua = "emily-agent/0.1 data-collection-bot"
	}
	var out []CollectedDoc
	for _, sub := range r.Subreddits {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		docs, err := r.fetchSubreddit(ctx, sub, since, ua)
		if err != nil {
			log.Printf("reddit fetch_error sub=%s err=%v", sub, err)
			continue
		}
		out = append(out, docs...)
		// Polite inter-subreddit delay — reddit rate-limits aggressively.
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return out, nil
}

type redditListing struct {
	Data struct {
		Children []struct {
			Data struct {
				ID          string  `json:"id"`
				Title       string  `json:"title"`
				Selftext    string  `json:"selftext"`
				Score       int     `json:"score"`
				NumComments int     `json:"num_comments"`
				Permalink   string  `json:"permalink"`
				Subreddit   string  `json:"subreddit"`
				CreatedUTC  float64 `json:"created_utc"`
				IsSelf      bool    `json:"is_self"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (r *RedditSource) fetchSubreddit(ctx context.Context, sub string, since time.Time, ua string) ([]CollectedDoc, error) {
	reqURL := "https://www.reddit.com/r/" + url.PathEscape(sub) + "/new.json?limit=100"
	resp, err := httpGetWithRetry(ctx, r.client, reqURL, ua, 4)
	if err != nil {
		return nil, fmt.Errorf("reddit sub=%s: %w", sub, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("reddit status=%d sub=%s", resp.StatusCode, sub)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	var listing redditListing
	if err := json.Unmarshal(raw, &listing); err != nil {
		return nil, fmt.Errorf("reddit parse sub=%s: %w", sub, err)
	}

	minScore := r.MinScore
	if minScore <= 0 {
		minScore = 100
	}
	minComments := r.MinComments
	if minComments <= 0 {
		minComments = 10
	}
	minBody := r.MinBodyLen
	if minBody <= 0 {
		minBody = 500
	}

	var docs []CollectedDoc
	for _, child := range listing.Data.Children {
		p := child.Data
		if !p.IsSelf {
			continue
		}
		created := time.Unix(int64(p.CreatedUTC), 0).UTC()
		if !since.IsZero() && !created.After(since) {
			continue
		}
		if p.Score < minScore || p.NumComments < minComments || len(p.Selftext) < minBody {
			continue
		}
		if p.Selftext == "[deleted]" || p.Selftext == "[removed]" {
			continue
		}
		docs = append(docs, CollectedDoc{
			ID:          docID("reddit:" + p.ID),
			SourceName:  "reddit",
			URL:         "https://www.reddit.com" + p.Permalink,
			Title:       p.Title,
			Body:        p.Selftext,
			BodyBytes:   len(p.Selftext),
			SourceScore: p.Score,
			PublishedAt: created,
			CollectedAt: time.Now().UTC(),
			Tags:        []string{p.Subreddit},
		})
	}
	return docs, nil
}

// ---------------------------------------------------------------------------
// WikipediaSource — MediaWiki API: recent changes → plain-text extracts
// ---------------------------------------------------------------------------

type WikipediaSource struct {
	UserAgent  string
	MinBodyLen int
	client     *http.Client
}

func (w *WikipediaSource) Name() string { return "wikipedia" }

func (w *WikipediaSource) Fetch(ctx context.Context, since time.Time) ([]CollectedDoc, error) {
	if w.client == nil {
		w.client = &http.Client{Timeout: 20 * time.Second}
	}
	ua := w.UserAgent
	if ua == "" {
		ua = "emily-agent/0.1 data-collection-bot"
	}
	minBody := w.MinBodyLen
	if minBody <= 0 {
		minBody = 1000
	}

	pageIDs, err := w.fetchRecentChanges(ctx, since, ua)
	if err != nil {
		return nil, err
	}
	if len(pageIDs) == 0 {
		return nil, nil
	}

	var docs []CollectedDoc
	for i := 0; i < len(pageIDs); i += 100 {
		if err := ctx.Err(); err != nil {
			return docs, err
		}
		end := i + 100
		if end > len(pageIDs) {
			end = len(pageIDs)
		}
		batch, err := w.fetchExtracts(ctx, pageIDs[i:end], ua, minBody)
		if err != nil {
			log.Printf("wikipedia extract_error batch_start=%d err=%v", i, err)
			continue
		}
		docs = append(docs, batch...)
		select {
		case <-ctx.Done():
			return docs, ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
	return docs, nil
}

func (w *WikipediaSource) fetchRecentChanges(ctx context.Context, since time.Time, ua string) ([]int, error) {
	params := url.Values{
		"action":      {"query"},
		"format":      {"json"},
		"list":        {"recentchanges"},
		"rctype":      {"edit"},
		"rcnamespace": {"0"},
		"rclimit":     {"100"},
		"rcprop":      {"title|timestamp|ids"},
	}
	if !since.IsZero() {
		params.Set("rcstart", since.UTC().Format(time.RFC3339))
		params.Set("rcdir", "newer")
	}

	resp, err := httpGetWithRetry(ctx, w.client,
		"https://en.wikipedia.org/w/api.php?"+params.Encode(), ua, 3)
	if err != nil {
		return nil, fmt.Errorf("wikipedia rc_fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("wikipedia rc_status=%d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return nil, err
	}

	var result struct {
		Query struct {
			RecentChanges []struct {
				PageID int `json:"pageid"`
			} `json:"recentchanges"`
		} `json:"query"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("wikipedia rc_parse: %w", err)
	}

	dedup := map[int]bool{}
	var ids []int
	for _, rc := range result.Query.RecentChanges {
		if rc.PageID > 0 && !dedup[rc.PageID] {
			ids = append(ids, rc.PageID)
			dedup[rc.PageID] = true
		}
	}
	return ids, nil
}

func (w *WikipediaSource) fetchExtracts(ctx context.Context, pageIDs []int, ua string, minBody int) ([]CollectedDoc, error) {
	ids := make([]string, len(pageIDs))
	for i, id := range pageIDs {
		ids[i] = strconv.Itoa(id)
	}
	params := url.Values{
		"action":          {"query"},
		"format":          {"json"},
		"prop":            {"extracts"},
		"explaintext":     {"1"},
		"exsectionformat": {"plain"},
		"exchars":         {"50000"},
		"pageids":         {strings.Join(ids, "|")},
	}

	resp, err := httpGetWithRetry(ctx, w.client,
		"https://en.wikipedia.org/w/api.php?"+params.Encode(), ua, 3)
	if err != nil {
		return nil, fmt.Errorf("wikipedia extract_fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("wikipedia extract_status=%d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				PageID  int    `json:"pageid"`
				Title   string `json:"title"`
				Extract string `json:"extract"`
			} `json:"pages"`
		} `json:"query"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("wikipedia extract_parse: %w", err)
	}

	var docs []CollectedDoc
	for _, page := range result.Query.Pages {
		if page.PageID <= 0 || len(page.Extract) < minBody {
			continue
		}
		// Exclude stub articles even if they pass the minBody length check.
		// Wikipedia stubs often have boilerplate padding that inflates their size.
		if isWikipediaStub(page.Extract) {
			continue
		}
		pageURL := "https://en.wikipedia.org/wiki/" +
			url.PathEscape(strings.ReplaceAll(page.Title, " ", "_"))

		// Count sections (paragraphs separated by double newline) as a quality signal.
		sectionCount := countWikiSections(page.Extract)
		tags := []string{"wikipedia", fmt.Sprintf("sections:%d", sectionCount)}

		docs = append(docs, CollectedDoc{
			ID:          docID("wikipedia:" + strconv.Itoa(page.PageID)),
			SourceName:  "wikipedia",
			URL:         pageURL,
			Title:       page.Title,
			Body:        page.Extract,
			BodyBytes:   len(page.Extract),
			PublishedAt: time.Now().UTC(),
			CollectedAt: time.Now().UTC(),
			Tags:        tags,
		})
	}
	return docs, nil
}

// isWikipediaStub returns true when the article extract contains stub markers.
// Wikipedia stubs are explicitly marked and are unsuitable training data.
func isWikipediaStub(text string) bool {
	lower := strings.ToLower(text)
	stubPhrases := []string{
		"this article is a stub",
		"this is a stub",
		"this biography is a stub",
		"this article about",
		"you can help wikipedia by expanding it",
		"this entry is a stub",
	}
	for _, phrase := range stubPhrases {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

// countWikiSections counts the number of top-level sections in a Wikipedia
// plain-text extract. Sections are delimited by double newlines with a short
// all-caps or title-case line (section headings in explaintext format).
func countWikiSections(text string) int {
	if text == "" {
		return 0
	}
	paragraphs := strings.Split(text, "\n\n")
	count := 0
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			count++
		}
	}
	if count == 0 {
		return 1
	}
	return count
}

// ---------------------------------------------------------------------------
// ArXivSource — academic preprints via the public Atom feed (no auth required).
// Default categories: cs.AI, cs.LG, cs.CL, cs.CV, stat.ML
// Body = title + abstract — dense, curated ML/AI content with zero noise.
// ---------------------------------------------------------------------------

type ArXivSource struct {
	Categories  []string // arxiv category codes, e.g. ["cs.AI", "cs.LG"]
	MaxResults  int      // per request, max 200 per arxiv API limits
	UserAgent   string
	MinAbstrLen int
	// MaxFullText controls how many full-paper HTML fetches are attempted per
	// batch. 0 = abstract-only (default). Each full-text fetch costs one HTTP
	// call at ~1s delay; set to 10–20 for best quality without slowing the cycle.
	MaxFullText int
	client      *http.Client
}

func (a *ArXivSource) Name() string { return "arxiv" }

func (a *ArXivSource) Fetch(ctx context.Context, since time.Time) ([]CollectedDoc, error) {
	if a.client == nil {
		a.client = &http.Client{Timeout: 30 * time.Second}
	}
	ua := a.UserAgent
	if ua == "" {
		ua = "emily-agent/0.1 data-collection-bot"
	}
	cats := a.Categories
	if len(cats) == 0 {
		cats = []string{"cs.AI", "cs.LG", "cs.CL", "cs.CV", "stat.ML"}
	}
	maxResults := a.MaxResults
	if maxResults <= 0 {
		maxResults = 100
	}
	minAbstr := a.MinAbstrLen
	if minAbstr <= 0 {
		minAbstr = 200
	}

	parts := make([]string, len(cats))
	for i, c := range cats {
		parts[i] = "cat:" + c
	}
	query := strings.Join(parts, "+OR+")
	apiURL := fmt.Sprintf(
		"https://export.arxiv.org/api/query?search_query=%s&start=0&max_results=%d&sortBy=submittedDate&sortOrder=descending",
		url.QueryEscape(query), maxResults,
	)

	// ArXiv requests a 3-second delay between calls per their usage policy.
	resp, err := httpGetWithRetry(ctx, a.client, apiURL, ua, 3)
	if err != nil {
		return nil, fmt.Errorf("arxiv fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("arxiv status=%d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	return a.parseAtom(raw, since, minAbstr)
}

type arxivFeed struct {
	XMLName xml.Name      `xml:"feed"`
	Entries []arxivEntry  `xml:"entry"`
}

type arxivEntry struct {
	ID         string `xml:"id"`
	Title      string `xml:"title"`
	Summary    string `xml:"summary"`
	Published  string `xml:"published"`
	Categories []struct {
		Term string `xml:"term,attr"`
	} `xml:"category"`
}

// reLatexMacro matches raw LaTeX commands that leak into arxiv abstracts.
// Common forms: \textbf{...}, \emph{...}, \cite{...}, \mathbb{R}, \alpha, etc.
// We strip the command but keep the content (for \textbf{word} → word).
var (
	reLatexCmd      = regexp.MustCompile(`\\[a-zA-Z]+\{([^}]*)\}`)    // \cmd{content} → content
	reLatexCmdBare  = regexp.MustCompile(`\\[a-zA-Z]+\b`)              // bare \cmd → ""
	reLatexEnv      = regexp.MustCompile(`(?s)\\begin\{[^}]+\}.*?\\end\{[^}]+\}`) // \begin{...}...\end{...}
	reLatexMath     = regexp.MustCompile(`\$\$[^$]+\$\$|\$[^$\n]+\$`) // $$...$$ and $...$
	reLatexBraces   = regexp.MustCompile(`[{}]`)
)

// cleanAbstract strips raw LaTeX markup from an arxiv abstract, returning
// plain prose. Environments and math blocks are removed entirely; inline
// commands like \textbf{word} have their content preserved.
func cleanAbstract(text string) string {
	// Remove LaTeX environments first (may contain nested braces).
	text = reLatexEnv.ReplaceAllString(text, " ")
	// Remove display and inline math.
	text = reLatexMath.ReplaceAllString(text, " ")
	// Inline commands: \cmd{content} → content.
	for reLatexCmd.MatchString(text) {
		text = reLatexCmd.ReplaceAllString(text, "$1")
	}
	// Bare commands: \alpha → "".
	text = reLatexCmdBare.ReplaceAllString(text, "")
	// Remove stray braces left over.
	text = reLatexBraces.ReplaceAllString(text, "")
	// Collapse excess whitespace.
	text = regexp.MustCompile(`[ \t]+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")
	return strings.TrimSpace(text)
}

// reHTMLTag strips HTML tags. Used to extract plain text from arxiv HTML papers.
var reHTMLTag = regexp.MustCompile(`<[^>]+>`)

// stripHTML removes HTML tags and decodes common entities, returning plain text.
func stripHTML(s string) string {
	// Remove script and style blocks entirely.
	s = regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`).ReplaceAllString(s, " ")
	// Replace block-level tags with newlines.
	s = regexp.MustCompile(`(?i)</(p|div|section|h[1-6]|li|br|tr)>`).ReplaceAllString(s, "\n")
	// Strip all remaining tags.
	s = reHTMLTag.ReplaceAllString(s, "")
	// Decode common HTML entities.
	replacer := strings.NewReplacer(
		"&amp;", "&", "&lt;", "<", "&gt;", ">",
		"&quot;", `"`, "&#39;", "'", "&nbsp;", " ",
		"&mdash;", "—", "&ndash;", "–", "&hellip;", "…",
	)
	s = replacer.Replace(s)
	// Normalise whitespace: collapse runs of blanks, trim.
	s = regexp.MustCompile(`[ \t]+`).ReplaceAllString(s, " ")
	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

// fetchArXivFullText attempts to retrieve the plain-text content of a paper
// from the arxiv HTML endpoint (https://arxiv.org/html/{id}). Returns ("", nil)
// when the HTML version is unavailable (404) or the content is too short.
// Caller is responsible for rate-limiting between calls.
func (a *ArXivSource) fetchArXivFullText(ctx context.Context, paperID, ua string) (string, error) {
	htmlURL := "https://arxiv.org/html/" + paperID
	req, err := http.NewRequestWithContext(ctx, "GET", htmlURL, nil)
	if err != nil {
		return "", err
	}
	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 || resp.StatusCode == 403 {
		return "", nil // HTML version unavailable — normal for older papers
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("arxiv html status=%d id=%s", resp.StatusCode, paperID)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20)) // 2 MB cap
	if err != nil {
		return "", err
	}

	// Extract the <article> element if present (arxiv HTML papers use <article>).
	article := string(raw)
	if start := strings.Index(article, "<article"); start >= 0 {
		if end := strings.LastIndex(article, "</article>"); end > start {
			article = article[start : end+len("</article>")]
		}
	}

	text := stripHTML(article)
	if len(text) < 2000 {
		return "", nil // too short — likely a redirect or unavailable stub
	}
	return text, nil
}

func (a *ArXivSource) parseAtom(raw []byte, since time.Time, minAbstr int) ([]CollectedDoc, error) {
	var feed arxivFeed
	if err := xml.Unmarshal(raw, &feed); err != nil {
		return nil, fmt.Errorf("arxiv xml parse: %w", err)
	}

	var docs []CollectedDoc
	for _, entry := range feed.Entries {
		published, _ := time.Parse(time.RFC3339, strings.TrimSpace(entry.Published))
		if !since.IsZero() && !published.IsZero() && !published.After(since) {
			continue
		}

		title := strings.TrimSpace(strings.ReplaceAll(entry.Title, "\n", " "))
		abstract := cleanAbstract(strings.TrimSpace(entry.Summary))
		if len(abstract) < minAbstr {
			continue
		}

		// Extract stable paper ID from URL like http://arxiv.org/abs/2501.12345v1
		paperID := entry.ID
		if idx := strings.LastIndex(paperID, "/"); idx >= 0 {
			paperID = paperID[idx+1:]
		}
		if idx := strings.LastIndex(paperID, "v"); idx > 0 {
			// Only strip version suffix if what remains looks like an arxiv ID (contains a dot)
			if strings.Contains(paperID[:idx], ".") {
				paperID = paperID[:idx]
			}
		}

		var tags []string
		for _, cat := range entry.Categories {
			if cat.Term != "" {
				tags = append(tags, cat.Term)
			}
		}

		// Body = title + abstract — a complete, self-contained training unit.
		body := title + "\n\n" + abstract
		docs = append(docs, CollectedDoc{
			ID:          docID("arxiv:" + paperID),
			SourceName:  "arxiv",
			URL:         "https://arxiv.org/abs/" + paperID,
			Title:       title,
			Body:        body,
			BodyBytes:   len(body),
			PublishedAt: published,
			CollectedAt: time.Now().UTC(),
			Tags:        tags,
		})
	}

	// Best-effort full-paper HTML enrichment: fetch up to MaxFullText papers.
	// Each fetch costs ~1 HTTP call; we sleep 1s between calls to respect arxiv's
	// usage policy. Abstract body is kept as fallback when HTML is unavailable.
	if a.MaxFullText > 0 && len(docs) > 0 {
		ua := a.UserAgent
		if ua == "" {
			ua = "emily-agent/0.1 data-collection-bot"
		}
		limit := a.MaxFullText
		if limit > len(docs) {
			limit = len(docs)
		}
		enrichCtx, enrichCancel := context.WithTimeout(context.Background(), time.Duration(limit)*3*time.Second)
		defer enrichCancel()
		for i := 0; i < limit; i++ {
			doc := &docs[i]
			// Extract paper ID from URL: https://arxiv.org/abs/{id}
			paperID := strings.TrimPrefix(doc.URL, "https://arxiv.org/abs/")
			fullText, err := a.fetchArXivFullText(enrichCtx, paperID, ua)
			if err != nil {
				log.Printf("arxiv fulltext id=%s err=%v", paperID, err)
			} else if fullText != "" {
				doc.Body = doc.Title + "\n\n" + fullText
				doc.BodyBytes = len(doc.Body)
				doc.Tags = append(doc.Tags, "full_paper")
			}
			// Polite delay between HTML fetches.
			select {
			case <-enrichCtx.Done():
				return docs, nil
			case <-time.After(time.Second):
			}
		}
	}

	return docs, nil
}

// ---------------------------------------------------------------------------
// QualityScorer — deterministic, no ML, no external dependencies
// Score range: 0–10 (max achievable ~9.0 for dense ML code+text).
// Gold >= 6.5, Silver >= 4.5, Bronze < 4.5
//
// Score breakdown:
//   Length       0–3 pts  (saturates at 3000 chars)
//   Sentence     0–2 pts  (saturates at 15 sentences)
//   Vocabulary   0–2 pts  (type/token ratio)
//   Citations    0–1 pt   (numbered refs or [n] markers)
//   Technical    0–1 pt   (code blocks, equations, ML keywords)
//   Spam penalty 0–3 pts  (deducted for ad/promo phrases)
// ---------------------------------------------------------------------------

var (
	reSentenceEnd  = regexp.MustCompile(`[.!?]+\s`)
	reSpam         = regexp.MustCompile(`(?i)(click here|buy now|limited time offer|free download|#ad\b|sponsored content|affiliate link|promo code|discount code|sign up now|subscribe now|act now|order today)`)
	reWordToken    = regexp.MustCompile(`\b\w{3,}\b`)
	reCitation     = regexp.MustCompile(`\[\d+\]|\^\d+\b|\(\d{4}\)`)
	reCodeBlock    = regexp.MustCompile("(?s)```[\\w]*\\n.*?```")
	reInlineCode   = regexp.MustCompile("`[^`\n]{3,}`")
	reEquation     = regexp.MustCompile(`\$[^$\n]{3,}\$|\\\[.*?\\\]`)
	reTechKeyword  = regexp.MustCompile(`(?i)\b(algorithm|neural|gradient|transformer|embedding|backprop|inference|training|dataset|benchmark|accuracy|precision|recall|f1|epoch|batch|loss|softmax|attention|token|layer|weight|bias|activation|pytorch|tensorflow|numpy|cuda|gpu|llm|llama|gpt|bert|fine.tun|quantiz)\b`)
)

type QualityScorer struct{}

func (q *QualityScorer) Score(doc CollectedDoc) (score float64, tier string) {
	body := doc.Body
	if len(body) == 0 {
		return 0, "bronze"
	}

	// Length: 0–3 pts (saturates at 3000 chars)
	lengthScore := math.Min(float64(len(body))/1000.0, 3.0)

	// Sentence depth: 0–2 pts (saturates at 15 sentences)
	sentences := len(reSentenceEnd.FindAllString(body, -1))
	sentenceScore := math.Min(float64(sentences)/7.5, 2.0)

	// Vocabulary richness: 0–2 pts (type/token ratio scaled)
	words := reWordToken.FindAllString(strings.ToLower(body), -1)
	var vocabScore float64
	if len(words) > 10 {
		unique := make(map[string]struct{}, len(words))
		for _, w := range words {
			unique[w] = struct{}{}
		}
		vocabScore = math.Min(float64(len(unique))/float64(len(words))*4.0, 2.0)
	}

	// Citation richness: 0–1 pt (numbered refs like [1], ^2, or (2024) indicate
	// sourced content vs. opinion/spam). Saturates at 5 distinct citation markers.
	citationHits := len(reCitation.FindAllString(body, -1))
	citationScore := math.Min(float64(citationHits)/5.0, 1.0)

	// Technical content bonus: 0–1 pt for ML/CS content.
	// Code blocks (``` fenced or inline `code`), equations ($...$), and
	// domain-specific ML/AI keywords are strong indicators of high-value corpus data.
	codeBlocks := len(reCodeBlock.FindAllString(body, -1))
	inlineCode := len(reInlineCode.FindAllString(body, -1))
	equations := len(reEquation.FindAllString(body, -1))
	techKeywords := len(reTechKeyword.FindAllString(body, -1))
	techSignals := codeBlocks*3 + inlineCode + equations*2 + techKeywords/5
	techScore := math.Min(float64(techSignals)/10.0, 1.0)

	// Spam penalty: deduct up to 3 pts
	spamHits := float64(len(reSpam.FindAllString(body, -1)))
	spamPenalty := math.Min(spamHits*1.5, 3.0)

	total := lengthScore + sentenceScore + vocabScore + citationScore + techScore - spamPenalty
	total = math.Round(math.Max(0, math.Min(10, total))*100) / 100

	// Thresholds calibrated against achievable max (~9.0 for dense ML code+text):
	// gold 10-20%, silver 30-50%, bronze remainder.
	switch {
	case total >= 6.5:
		tier = "gold"
	case total >= 4.5:
		tier = "silver"
	default:
		tier = "bronze"
	}
	return total, tier
}

func tierRank(t string) int {
	switch t {
	case "gold":
		return 3
	case "silver":
		return 2
	default:
		return 1
	}
}

// ---------------------------------------------------------------------------
// seenSet — in-memory dedup map
// Exact pattern from PRRJECT_FATBABY/internal/processor/worker.go
// ---------------------------------------------------------------------------

type docSeenSet struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func newDocSeenSet(initial map[string]struct{}) *docSeenSet {
	if initial == nil {
		initial = map[string]struct{}{}
	}
	return &docSeenSet{seen: initial}
}

func (s *docSeenSet) has(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.seen[id]
	return ok
}

func (s *docSeenSet) mark(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen[id] = struct{}{}
}

func (s *docSeenSet) size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.seen)
}

// ---------------------------------------------------------------------------
// CollectionStats
// ---------------------------------------------------------------------------

type CollectionStats struct {
	DocsCollected   int64            `json:"docs_collected"`
	DocsRejected    int64            `json:"docs_rejected"`
	DocsDuplicate   int64            `json:"docs_duplicate"`
	DocsGold        int64            `json:"docs_gold"`
	DocsSilver      int64            `json:"docs_silver"`
	DocsBronze      int64            `json:"docs_bronze"`
	ErrorCount      int64            `json:"error_count"`
	SeenSetSize     int              `json:"seen_set_size"`
	LastCollectedAt *time.Time       `json:"last_collected_at,omitempty"`
	Running         bool             `json:"running"`
	PerSource       map[string]int64 `json:"per_source,omitempty"` // collected count by source name
}

// ---------------------------------------------------------------------------
// CollectorPipeline
// ---------------------------------------------------------------------------

type CollectorConfig struct {
	Store        *DocStore
	Sources      []Source
	Scorer       *QualityScorer
	Workers      int
	PollInterval time.Duration
	MinTier      string // "bronze" | "silver" | "gold" — drop docs below this
	Logger       *log.Logger
}

type CollectorPipeline struct {
	cfg CollectorConfig

	mu      sync.Mutex
	seen    *docSeenSet
	running bool
	cancel  context.CancelFunc

	docsCollected   atomic.Int64
	docsRejected    atomic.Int64
	docsDuplicate   atomic.Int64
	docsGold        atomic.Int64
	docsSilver      atomic.Int64
	docsBronze      atomic.Int64
	errorCount      atomic.Int64
	lastCollectedAt atomic.Pointer[time.Time]

	// per-source counters: source_name → collected count
	perSourceMu sync.Mutex
	perSource   map[string]int64
}

func NewCollectorPipeline(cfg CollectorConfig) *CollectorPipeline {
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 5 * time.Minute
	}
	if cfg.Logger == nil {
		cfg.Logger = log.New(log.Writer(), "[collector] ", log.LstdFlags|log.LUTC)
	}
	if cfg.MinTier == "" {
		cfg.MinTier = "bronze"
	}
	return &CollectorPipeline{cfg: cfg, perSource: make(map[string]int64)}
}

// Start loads the seen set (O(n) scan of existing corpus) then launches the
// background collection loop.  The initCtx is only used for the scan —
// the collection loop uses its own context, cancelled only by Stop().
func (p *CollectorPipeline) Start(initCtx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running {
		return fmt.Errorf("collector already running")
	}

	initial, err := p.cfg.Store.LoadSeen(initCtx)
	if err != nil {
		return fmt.Errorf("collector load_seen: %w", err)
	}
	p.seen = newDocSeenSet(initial)
	p.cfg.Logger.Printf("start seen_set=%d sources=%d poll_interval=%s min_tier=%s",
		p.seen.size(), len(p.cfg.Sources), p.cfg.PollInterval, p.cfg.MinTier)

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.running = true

	go func() {
		defer func() {
			p.mu.Lock()
			p.running = false
			p.mu.Unlock()
		}()
		p.run(ctx)
	}()
	return nil
}

func (p *CollectorPipeline) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
}

func (p *CollectorPipeline) Stats() CollectionStats {
	p.mu.Lock()
	running := p.running
	seen := p.seen
	p.mu.Unlock()

	seenSize := 0
	if seen != nil {
		seenSize = seen.size()
	}

	var lastAt *time.Time
	if t := p.lastCollectedAt.Load(); t != nil {
		cp := *t
		lastAt = &cp
	}

	p.perSourceMu.Lock()
	perSource := make(map[string]int64, len(p.perSource))
	for k, v := range p.perSource {
		perSource[k] = v
	}
	p.perSourceMu.Unlock()

	return CollectionStats{
		DocsCollected:   p.docsCollected.Load(),
		DocsRejected:    p.docsRejected.Load(),
		DocsDuplicate:   p.docsDuplicate.Load(),
		DocsGold:        p.docsGold.Load(),
		DocsSilver:      p.docsSilver.Load(),
		DocsBronze:      p.docsBronze.Load(),
		ErrorCount:      p.errorCount.Load(),
		SeenSetSize:     seenSize,
		LastCollectedAt: lastAt,
		Running:         running,
		PerSource:       perSource,
	}
}

func (p *CollectorPipeline) run(ctx context.Context) {
	var since time.Time
	for {
		cycleStart := time.Now()
		for _, src := range p.cfg.Sources {
			if err := ctx.Err(); err != nil {
				return
			}
			p.cfg.Logger.Printf("fetch_start source=%s since=%s", src.Name(), since.Format(time.RFC3339))
			docs, err := src.Fetch(ctx, since)
			if err != nil {
				p.cfg.Logger.Printf("fetch_error source=%s err=%v", src.Name(), err)
				p.errorCount.Add(1)
				continue
			}
			p.cfg.Logger.Printf("fetch_done source=%s fetched=%d", src.Name(), len(docs))
			p.processBatch(ctx, docs)
		}
		since = cycleStart

		select {
		case <-ctx.Done():
			return
		case <-time.After(p.cfg.PollInterval):
		}
	}
}

// processBatch fans out docs to the worker pool.
// Duplicate docs are skipped before entering the pool — same hot-path pattern
// as crawlBatch in PRRJECT_FATBABY/prwatch/crawler.go.
func (p *CollectorPipeline) processBatch(ctx context.Context, docs []CollectedDoc) {
	jobs := make(chan CollectedDoc)
	var wg sync.WaitGroup
	for i := 0; i < p.cfg.Workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for doc := range jobs {
				if ctx.Err() != nil {
					return
				}
				p.processOne(doc)
			}
		}(i + 1)
	}

	enqueued := 0
	for _, doc := range docs {
		if p.seen.has(doc.ID) {
			p.docsDuplicate.Add(1)
			continue
		}
		enqueued++
		jobs <- doc
	}
	close(jobs)
	wg.Wait()

	p.cfg.Logger.Printf("batch_done total=%d enqueued=%d duplicate=%d",
		len(docs), enqueued, len(docs)-enqueued)
}

func (p *CollectorPipeline) processOne(doc CollectedDoc) {
	scorer := p.cfg.Scorer
	if scorer == nil {
		scorer = &QualityScorer{}
	}
	score, tier := scorer.Score(doc)
	doc.QualityScore = score
	doc.QualityTier = tier

	if tierRank(tier) < tierRank(p.cfg.MinTier) {
		p.docsRejected.Add(1)
		p.cfg.Logger.Printf("reject id=%s tier=%s score=%.2f source=%s", doc.ID, tier, score, doc.SourceName)
		return
	}

	if err := p.cfg.Store.Append(doc); err != nil {
		p.cfg.Logger.Printf("store_error id=%s err=%v", doc.ID, err)
		p.errorCount.Add(1)
		return
	}

	p.seen.mark(doc.ID)
	p.docsCollected.Add(1)

	p.perSourceMu.Lock()
	p.perSource[doc.SourceName]++
	p.perSourceMu.Unlock()

	now := time.Now().UTC()
	p.lastCollectedAt.Store(&now)

	switch tier {
	case "gold":
		p.docsGold.Add(1)
	case "silver":
		p.docsSilver.Add(1)
	default:
		p.docsBronze.Add(1)
	}
	p.cfg.Logger.Printf("collected id=%s source=%s tier=%s score=%.2f bytes=%d title=%q",
		doc.ID, doc.SourceName, tier, score, doc.BodyBytes, truncate(doc.Title, 60))
}

// ---------------------------------------------------------------------------
// HTTP handler — wired in main.go as GET+POST /collect
// ---------------------------------------------------------------------------

func (s *Server) handleCollect(w http.ResponseWriter, r *http.Request) {
	if s.collector == nil {
		http.Error(w, `{"error":"collector not configured — set EMILY_COLLECT_DIR"}`, 503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(s.collector.Stats())

	case http.MethodPost:
		switch r.URL.Query().Get("action") {
		case "stop":
			s.collector.Stop()
			json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
		default:
			if err := s.collector.Start(r.Context()); err != nil {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"status": "started"})
		}

	default:
		http.Error(w, "method not allowed", 405)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func docID(key string) string {
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h[:8])
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
