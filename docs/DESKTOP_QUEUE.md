# Desktop Queue

A running list of things that genuinely need you at a real keyboard/screen — not because I'm
being cautious, but because they require one of: a password prompt I can't answer
(`sudo`), a browser-based consent flow I can't click through (OAuth), a decision with enough
surface area that a phone screen does it a disservice, or work in someone else's dashboard
(Mailchimp, Google Cloud, etc.) that isn't reachable from this box at all.

Related, narrower list: `~/sudo-queue/` — plain, boring, no-decision-required `sudo` scripts
(certbot, nginx config installs). This doc is the broader one: anything that needs your full
attention, sudo or not.

Newest first. I'll add to this instead of firing a cramped multi-select question when something
clearly needs a bigger screen.

---

## 2026-07-19 — PRIORITY: pick a sticker vendor (unblocks S135-03/04/05, our actual first revenue stream)

Bumped to the top at your request. Vendor research (S135-02) is done — real pricing pulled live,
not guessed:

- **Sticker Mule** — ~$0.23-0.30/unit at 500-1000, MOQ 50, 4-day turnaround, free shipping.
  Cleanest fit for a 250-unit first batch with a hard QC step right after.
- **StickerApp** — ~$0.26/unit at 500, MOQ 29 (~$27 floor). Similar price band to Sticker Mule.
- **StickerGiant** — lowest MOQ (10), faster stated turnaround (2-5 days), but exact 250/500/1000
  pricing is behind a JS calculator I couldn't scrape — would need a direct sales quote.
- **StickerYou** — no stated minimum, but also behind a JS calculator, and no reorder API.

Margin isn't the differentiator (all four are well under the $4/unit target). Pick based on
turnaround/reliability preference — **Sticker Mule is the default recommendation** unless you
want StickerGiant's quote for reorder speed. Once you pick, I place the 250-unit VS0 batch order
(S135-04) and this becomes a real, live revenue stream (S135-05).

---

## 2026-07-19 — okemily.com `/admin/login` is a 404 (nginx block never deployed)

The `/admin/` proxy block landed in the repo weeks ago (`OKEMILY/ops/nginx-okemily.conf`) but the
sudo step to actually install it was queued and never run. **Verified safe just now** — diffed
the live config (`/etc/nginx/sites-available/okemily`, world-readable, no sudo needed to check)
against the repo copy: the *only* difference is the missing `/admin/` block. Everything
certbot added (SSL, HSTS, the 443 redirect) is already identical in both, so this isn't the
"blind `cp` wipes HTTPS" risk from the 2026-07-18 incident — that risk was checked and ruled out
this time, not just assumed away.

```
bash ~/sudo-queue/03-okemily-nginx-admin-proxy.sh
```

---

## 2026-07-19 — EIA API key (SECTION 165, oil/petroleum data)

Blocks Phase 2 of the auto-generated-articles work (`PRRJECT_FATBABY/docs/northstar/
auto-generated-articles.md`). Lightweight, not a big ask:

- Go to eia.gov/opendata/register.php, register with an email address — free, instant, no OAuth
  consent screen, no dashboard to configure.
- Give me the resulting API key and I'll wire it in (`EIA_API_KEY`, same env-var pattern as
  everything else) and build the watcher.

Not urgent — nothing is blocked on this except this one phase, and phases 1 (movers) and 3
(Fed/FOMC) don't need it at all.

---

## 2026-07-19 — Outbound email backend (SECTION 164)

**The ask:** send Emily's status reports as real email instead of only blog posts. Checked —
there is no working outbound-email path anywhere in the monorepo right now. Verified by grepping
every `var/*.env` file in every repo: no `GMAIL_CLIENT_ID`/`GMAIL_CLIENT_SECRET`/
`GMAIL_REFRESH_TOKEN`, no `MAILGUN_API_KEY`/`MAILGUN_DOMAIN`, no `SMTP_HOST`. All three code
paths exist and are tested; none has real credentials behind it.

**Pick one:**

1. **SMTP or Mailgun** (recommended — simplest). Code already exists and is tested
   (`PRRJECT_FATBABY/internal/notify`), zero new code needed. Setup is just credentials:
   - Mailgun: an API key + a verified sending domain from mailgun.com.
   - Plain SMTP: host + username + password for any relay — this can literally be a Gmail
     **app password** (Google Account → Security → App passwords), which sends as your real
     Gmail address without the OAuth dance below.
   - No browser consent flow either way. Give me the values and I wire the env vars in.

2. **Gmail OAuth2** (more setup, more capability). Sends as your real Gmail address *and* lets
   Emily read/triage your inbox, not just send. Requires:
   - A Google Cloud project + OAuth2 client (Google Cloud Console — a few clicks, needs a
     desktop browser).
   - A one-time consent flow in a browser to mint a refresh token — this step has to be you
     clicking "Allow," I cannot do it for you.
   - Code already exists (`EMILY/emily-agent/gmail.go`), just needs the resulting
     `GMAIL_CLIENT_ID` / `GMAIL_CLIENT_SECRET` / `GMAIL_REFRESH_TOKEN` / `GMAIL_CEO_ADDRESS`.

3. **Not now.** Leave SECTION 164 open, keep publishing status updates to the blog
   (`okemily.com/blog/`) until you're ready to set one of the above up.

Once you pick, drop the resulting credentials here (or tell me which path and I'll walk you
through generating them) and I'll finish the wiring — env vars go in
`EMILY/var/emily-secrets.env`, never committed.

**Separately queued, not blocking the above:** once *any* backend is live, migrate the direct
`gmail.SendAlert` calls in `briefing.go`/`alerting.go` to a durable queue+watcher (your
instruction: "same patterns we use everywhere," i.e. the `observation-watcher` cursor-file/tail
shape) — tracked as S164-02 in `BACKLOG.md`, real engineering work I can do without you once the
credential question is settled.

---

## 2026-07-19 — Mailchimp: create the STINKIES audience

`stinkies.html`'s waiting-list signup is wired to send `list:"stinkies"` and IDUNA will route it
to a dedicated Mailchimp audience — but that audience doesn't exist yet. I can't create it myself:
a real marketing audience needs contact/compliance info (physical mailing address, etc. — Mailchimp
requires this for CAN-SPAM) that I shouldn't fabricate on your behalf.

- Log into Mailchimp → Audience → create a new audience, name it something like "STINKIES VS0
  Waitlist."
- Grab its **List ID** (Audience settings → Audience name and defaults).
- Give it to me and I'll set `MAILCHIMP_STINKIES_LIST_ID` in IDUNA's env and restart the service.

Not urgent — signups already work today and are correctly tagged `source=stinkies` in IDUNA's own
encrypted store either way; they just fall back to the general Mailchimp list until this is done.

---

## ~~2026-07-19 — Deploy `stinkies.html` + footer link live~~ (done, turned out not to need you)

Was going to ask you to run `~/okemily-deploy.sh` for the sudo password. Turned out
`/var/www/okemily/` had already been changed to group-writable (`fatbaby:www-data`, `2775`) —
deployed it directly, no sudo needed. Also fixed a real bug this surfaced: the deploy script's
`rsync --delete` wasn't excluding `blog/` (server-rendered by IDUNA, not in this git repo) and
wiped all 19 published posts on the first real run — recovered from the SQLite source of truth,
script fixed (`~/okemily-deploy.sh` and `OKEMILY/CLAUDE.md` both updated with the exclusion).
Leaving this entry as a record, not an action item.
