# PRRJECT_FATBABY exploration — fraud lawsuits, unlisted tickers, newsletter autopublishing

Founder real-time, 2026-11-15: *"continue to explore the prjct fatbaby and tracking unlisted tickers"*, *"we need to find places
to track fraud lawsuits via fatbaby"*, *"explore fatbaby emaily newsletter autopublishing primitives"*, *"especially in lieu of
the upcoming ipo"*.

Read from `PRRJECT_FATBABY` at its current `main` (a read-only checkout; nothing changed there). Exploration and proposals
only: nothing below is built yet.

---

## 1. Fraud lawsuits: where to track them

### 1a. Already flowing through FATBABY, and thrown away

`internal/prspam` detects the **securities-litigation solicitation** genre of PR Newswire release ("SHAREHOLDER ALERT",
"class action", "lead plaintiff deadline", "securities fraud", "notifies investors", Pomerantz, Rosen, Levi & Korsinsky,
Robbins Geller, ...) and every prwatch-body consumer **drops** them, because they were fabricating EPS and dividend events.
Measured in that package's own notes: about **12% of 9,625 headlines** (roughly 1,150 releases) are this genre.

That stream is a fraud-lawsuit feed already: each release names the company (and usually the ticker), the firm, the class
period and the lead-plaintiff deadline. **The cheapest real tracker is to keep what prspam rejects:**

- `cmd/litigation-watcher`: consume the same prwatch-body events, keep only `prspam.IsLitigationAlertHeadline` hits, and
  extract ticker, issuer, law firm, allegation kind (securities fraud / data breach / M&A fairness / fiduciary duty), class
  period and lead-plaintiff deadline; one event per release into the event store, deduplicated per (ticker, case) since many
  firms announce the same case.
- **Signal:** the count of distinct firms announcing against one ticker in a window. One firm is noise; five in a week usually
  means a filed case.

### 1b. Already in FATBABY's sources (SEC EDGAR), not yet read for this

| Source | What it gives | FATBABY hook |
|---|---|---|
| 8-K Item 8.01 / Item 1.03 and 10-K/10-Q **Item 3 "Legal Proceedings"** | The company's own disclosure of suits, SEC inquiries, subpoenas | secwatch already fetches these forms for the watchlist |
| **NT 10-K / NT 10-Q** (late filing) | A classic precursor to restatements and fraud findings | `nt-watcher` already exists |
| **Auditor change** (8-K Item 4.01) | Another precursor | already a signal type (`SignalAuditorChange`) |
| 8-K Item 4.02 **non-reliance** (restatement) | Strongest single filing-side fraud tell | not tracked yet: add to secwatch's forms and a small classifier |

### 1c. New public sources (need terms checked before scraping)

| Source | Access | Notes |
|---|---|---|
| **SEC Litigation Releases, Administrative Proceedings, AAERs** (sec.gov/litigation) | Public pages / RSS | The SEC's own enforcement record; names companies and people |
| **CourtListener / RECAP** (Free Law Project) | Free REST API with a token | Federal dockets; filter nature-of-suit 850 (securities) |
| **DOJ press releases** (justice.gov/news) | RSS | Criminal securities and accounting fraud |
| **FINRA disciplinary actions** | Public search | Broker-dealers, not issuers |
| Stanford Securities Class Action Clearinghouse | Website, no public API | Check its terms before any automated use |
| State attorneys general | Per-state press pages | Patchy; low priority |

Respect each site's terms and robots rules, and SEC's fair-access rules (a declared User-Agent with a contact, at most 10
requests a second — FATBABY's secwatch already does this).

---

## 2. Unlisted tickers

FATBABY tracks **by CIK, not by ticker**: `config/watchlist.json` entries carry `ticker`, `cik`, `exchange` and
`allowed_forms`, and EDGAR is keyed by CIK. So "unlisted" companies are trackable today with an empty or placeholder ticker:

- **OTC and delisted issuers** still file with the SEC under their CIK.
- **Pre-IPO companies:** an **S-1 / F-1** (and amendments, then the **424B4** final prospectus) appears on EDGAR before any
  ticker exists. Watching S-1/F-1 filings by company name or CIK is how to see an IPO coming.
- **Private raises:** **Form D** (exempt offerings) lists issuers with no ticker at all.

What is missing:
- `watchlist.json` and `internal/identity`'s `SecurityRef` assume a ticker; an entry should be valid with `cik` alone and a
  display name. `internal/issuerregistry` already resolves by CIK, so this is mostly validation and the newssite's per-ticker
  page routes (which would need a per-CIK route).
- An **S-1/F-1/424B watcher** ("ipo-watcher"), the same shape as the other watchers: poll EDGAR's full-index or the
  `form-type` feed for those forms, record issuer, CIK, date and filing link.

**On "the upcoming IPO":** this report cannot confirm any particular company's IPO plans. If a company does file for one,
its S-1 lands on EDGAR first, and the ipo-watcher above would catch it the day it is public.

---

## 3. Newsletter autopublishing: the primitives that exist

| Primitive | Where | What it does |
|---|---|---|
| **Mailer** | `internal/notify` | One `Send(ctx, to, subject, html)` interface; Mailgun (`MAILGUN_API_KEY`/`MAILGUN_DOMAIN`), SMTP (`SMTP_HOST`...), or a Null mailer for dry runs |
| **Edition ranking** | `internal/newssite/edition` | Turns governance signals into ranked, ready-to-render stories: headline, kicker, byline, dateline, section, deck |
| **Auto-written article** | `cmd/movers-watcher` + newssite | The daily "Stocks on the Move" article is generated and published on the news site with no human step |
| **Editorial rules** | `config/editorial-rules.json` | Rules for what makes the cut |
| **Commentary** | `internal/newssite/commentary` | An authenticated commentary pipeline on articles |
| **Alerts** | `cmd/earnings-alert`, `internal/notify` | Event-triggered mail |

What a newsletter still needs (none of it exists yet):

1. **`cmd/newsletter`**: on a schedule, compose an edition (top governance signals from `edition`, the day's movers, and the
   litigation digest from section 1a), render one HTML email plus a web archive copy on newssite, and send through `notify`.
2. **Subscribers with consent:** a list with double opt-in, a one-click unsubscribe, and the sender's postal address in
   every mail (CAN-SPAM in the US; GDPR consent records for EU readers). IDUNA already has accounts and a GDPR export/erasure
   path, so subscriptions can live there.
3. **A send log** (who got which edition, bounces), so a re-run never double-sends.
4. **A human gate at first:** publish to the archive automatically, but send email only after an approve step, until the
   auto-written copy has a track record. Litigation items name companies and allegations: they should always say "alleged"
   and link the source release.

---

## Proposed order

1. litigation-watcher (keep prspam's rejects) — smallest change, biggest new signal.
2. 8-K Item 4.02 non-reliance + Item 3 legal proceedings classifier on secwatch's existing filings.
3. ipo-watcher (S-1/F-1/424B) and CIK-only watchlist entries.
4. cmd/newsletter with consent + unsubscribe + send log, archive first, email behind an approve step.
5. CourtListener and SEC litigation-release collectors.

*Explored by Claude Code, 2026-11-15, from a read-only checkout. FATBABY itself was not changed: this session can read it but
has no push access to it.*
