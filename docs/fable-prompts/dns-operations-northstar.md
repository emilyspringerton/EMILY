# Queued Fable prompt — Fates DNS: EINHORN's DNS operations northstar

**Status:** dispatched 2026-07-17 (previous 2 attempts failed on transient `529 Overloaded`, not
session-limit or content issues — this is the first attempt with a full standalone prompt file,
per the founder's framing: "Emily is an internet native entity. DNS will be essential for
operations. Author the Fates DNS northstar as Fable.")

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Prompt body

You're writing a real architecture northstar for EINHORN_INDUSTRIAL, grounded in this codebase's actual DNS reality today, not a generic "how to run DNS" essay.

### Context: what's real right now

- `farthq.com` is the real production domain (verified via `dig NS farthq.com` 2026-07-17):
  currently resolves through Cloudflare's nameservers (`nicolas.ns.cloudflare.com`,
  `jocelyn.ns.cloudflare.com`) — DNS is Cloudflare-managed today, not self-hosted. This spec is
  about whether/how/why that changes, not a description of something already self-hosted.
- `iduna.farthq.com` is live (`EMILY/BACKLOG.md` SECTION 29/30, confirmed 2026-06-16) — IDUNA's
  API, auth, admin Back Office, and (per `IDUNA/ops/nginx-front-door-snippet.conf`, drafted this
  session) the eventual Unagent/Agent front-door funnel all live at this one subdomain today via
  nginx path-based routing, not subdomain-per-service.
- Real product surfaces that either already exist or are actively being built, each a candidate
  for its own subdomain or DNS-visible identity: `newssite`/signalapi (PRRJECT_FATBABY, currently
  `:8082`/`:9091` on the same box, no public DNS entry yet), EDIS WordPress (the public face of
  FatBaby, per SECTION 23, blocked on `sudo` deploy — not live yet), MJOLNIR's backend API
  (Android app talking to `EMILY_BASE_URL`, currently IP-based per `BuildConfig`, no DNS entry),
  the GoblinFoxDragon/SHANKPIT game server (`:6969` UDP, no DNS today), and now `NORN` (this
  session's new kernel repo — no network surface at all yet, library + CLIs only).
- `HQ-SPEC-PRIME-101` (NORN, this session) establishes the house pattern for any
  propose→grade→gate→promote loop with a frozen oracle and a no-`--force` gate. `HQ-SPEC-DOC-102`
  (SAGA) establishes the house pattern for reconciling documentation against reality. Both exist
  as specs only (not yet built as running code) but their *thinking* has been applied by hand
  repeatedly this session (the KIKORYU VS-audit, the PRIME-097 reconciliation, the NORN migrations
  themselves). If DNS record changes ever become something more than manual Cloudflare-dashboard
  edits, they are a natural NORN instantiation — propose a record change, grade it against reality
  (does the service actually respond at the new record?, did the change break anything currently
  resolving?), gate it (approval tier — DNS changes are exactly the kind of "physical/financial
  execution path" HQ-SPEC-PRIME-101 §4 describes as `biometric` tier: hard to reverse quickly at
  DNS-propagation timescales, high blast radius if wrong).
- `EMILY/docs/hq-specs/` is where founder-drafted and Fable-authored HQ-SPEC docs live, registered
  in `EMILY/context/golden-docs-index.md`. Follow that exact convention: filename
  `HQ-SPEC-<DOMAIN>-<NUMBER>-<slug>.md`, a `Status`/`Custody`/`Governs` header block matching the
  other six specs in that directory (read at least PRIME-101 and DOC-102 in full for the header
  format and section-numbering convention before writing).

### What "Emily is an internet native entity" should mean for this spec

The founder's framing is a real design constraint, not just flavor text: Emily Prime (and the
broader EINHORN system) increasingly *is* its network presence — Apples posted over HTTP, agents
authenticating via IDUNA JWTs over HTTP, MJOLNIR polling over HTTP, the whole RSI loop assumes
the network is there and correctly routed. A DNS misconfiguration isn't cosmetic for this system
the way it might be for a static brochure site — it's closer to a nervous-system disruption. Take
a real position on what that implies operationally (health-check-gated record changes? redundant
resolution paths? a documented "if DNS breaks, how does Emily Prime even know" story — self-
diagnosis when your own network identity is the thing that's broken is a real, interesting,
Löbian-flavored problem worth addressing directly, not hand-waved).

### What to actually decide (don't just describe options — take positions)

1. **Self-hosted DNS vs. staying on Cloudflare, and why.** Real tradeoffs: Cloudflare gives free
   DDoS protection, anycast, a mature API (relevant if any future automation wants to manage
   records programmatically — Cloudflare's API is a very different integration surface than
   running BIND/PowerDNS/CoreDNS yourself). Self-hosting gives full control and matches "EINHORN
   operates its own DNS servers" literally, but real operational cost: uptime becomes EINHORN's
   own problem, DDoS becomes EINHORN's own problem, and this is a small team/single-operator VM
   shop today (per `EMILY/CLAUDE.md`'s own admission of degraded-mode operation). Take a real
   position — a common good answer in shops this size is "self-host authoritative DNS for
   internal/service-discovery purposes, keep Cloudflare (or another anycast provider) in front for
   the public-facing zone," but argue it from this system's actual constraints, don't just default
   to that answer.
2. **Subdomain strategy for the product surfaces listed above.** Should IDUNA's current
   path-based routing (`iduna.farthq.com/api/...`) generalize to more services under one
   subdomain, or should NORN/MJOLNIR-backend/newssite/EDIS each get their own subdomain? State the
   tradeoff (fewer certs/simpler nginx vs. clean separation of blast radius per service) and
   decide, coordinating with (not contradicting) the already-drafted
   `IDUNA/ops/nginx-front-door-snippet.conf` and the still-unresolved EDIS-WordPress-vs-IDUNA-root
   conflict noted in `EMILY/docs/fable-prompts/iduna-front-door-funnel.md`.
3. **What actually gets built, in what order.** A concrete build sequence — DNS zone file /
   Cloudflare Terraform-or-API-managed config as code (checked into a repo, not dashboard-only
   clicking — this system's whole ethos is "everything is code, everything is Apple-audited"),
   health-check-gated record changes if that's the position taken in point 1, and if a NORN
   instantiation is warranted per the discussion above, sketch the instantiation table row (
   Proposer/Oracle/Tier/reality-root) the way PRIME-101 §6 does for its other six loops — don't
   build it, just spec the row.
4. **Human-in-the-loop points.** Domain registrar access, DNS provider account credentials,
   whatever the actual first unblock step is — name it concretely (this session's BACKLOG.md has a
   running "HUMAN UNBLOCK QUEUE" pattern; write your first actionable item in that same style so
   it can be dropped straight in).

### What to produce

1. `EMILY/docs/hq-specs/HQ-SPEC-INFRA-<next-free-number>-fates-dns.md` — the full northstar,
   following the existing spec header/section conventions. Check
   `EMILY/context/golden-docs-index.md` and `EMILY/docs/hq-specs/` for the next free HQ-SPEC number
   before writing (six exist today: PRIME-097, PRIME-101, DOC-102, FIN-098, FIN-099, SIM-100,
   AI-103 — confirm the actual highest number, don't assume).
2. A new `EMILY/BACKLOG.md` section (append at the end, following the numbering convention —
   check the highest existing `## SECTION N` before picking a number) with concrete, dependency-
   ordered build items matching the decisions above. Include the HUMAN UNBLOCK QUEUE row.
3. Register the new doc in `EMILY/context/golden-docs-index.md` (tier 2 is appropriate — this is
   subsystem-level infrastructure, not a tier-1 always-in-context doc).

### What NOT to do

- Don't touch DNS for real — no Cloudflare API calls, no nameserver changes, no `dig`/registrar
  actions beyond what's needed to verify the current state (which is already given to you above,
  verified 2026-07-17 — you don't need to re-verify it, just don't contradict it without a reason).
- Don't invent HQ-SPEC-PRIME-097's math or NORN's implementation details — those are already
  written (`pkg/norn`, real running code as of 2026-07-17, two live instantiations); if you sketch
  a DNS-as-NORN-instantiation row, cite the real interfaces, don't restate or reinvent them.
- Don't write code — this is a northstar/spec, matching the shape of the other six HQ-SPEC docs.

### When done

Commit `EMILY/docs/hq-specs/`, `EMILY/BACKLOG.md`, and `EMILY/context/golden-docs-index.md`
together as one EMILY commit, push. File an Apple (`emily apples post -t completion -repo EMILY
"<title>" "<body>"`) citing the new doc and section. Report back: the doc's filename/HQ-SPEC
number, the section number you added, the three key decisions made (self-host vs. Cloudflare,
subdomain strategy, first concrete build step), and the Apple ID.
