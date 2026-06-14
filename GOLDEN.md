# EMILY GOLDEN CONTEXT — 2026-06-14
*Auto-generated. Dense context. Source: EMILY/BACKLOG.md.*

## SYSTEM REFS
Repos: EMILY · PRRJECT_FATBABY · IDUNA · TYLER · SHANKPIT · MJOLNIR · APPLES · emily.cli · EDIS · EmilyOS · PITVIPER · GoblinFoxDragon
Sprint API: IDUNA /api/v1/heimdal/sprints | Golden index: EMILY/context/golden-docs-index.md (19 Tier 1 + 20 Tier 2 sources = 39 total)
Backlog: EMILY/BACKLOG.md | Archive: EMILY/DONE.md

## BACKLOG STATUS: 9 open, 73+ done (S25 complete 2026-06-14)

## COMPLETE SECTIONS: S3 · S4 · S6 · S9 · S12 · S13 · S14 · S15 · S16 · S17 · S18(AGI) · S20 · S21 · S22 · S25

## PRIORITY ORDER: S19 → S23 → S24 → S5 → S2 → S10

## OPEN ITEMS

### SECTION 19: SHANKPIT → STEAM EARLY ACCESS (2 open)
• S19-03: Steam Direct account + listing prep — $100 fee, BLOCKED: human action (Steam account).
• S19-05: Steam EA launch — $9.99 EA. Dep: S19-03 ✓, S19-04 ✓.

### SECTION 23: EDIS WORDPRESS PRODUCT (3 open)
• S23-01: Deploy EDIS to live WordPress install. Use: emily install --edis. Human action.
• S23-04: Emily+ subscription gate (WooCommerce, $29/mo). Dep: IDUNA subscription resource.
• S23-05: Mailchimp waitlist integration.

### SECTION 24: NEWSSITE OPS HARDENING (5 open)
• S24-01: Deploy to production server. Human action.
• S24-02: SSL cert + domain wiring (Let's Encrypt, fatbaby.io). Dep: S24-01.
• S24-03: Log rotation + alerting.
• S24-04: nginx cache tuning post-traffic.
• S24-05: Load test baseline (wrk -t4 -c50 -d30s).

### SECTION 5: FUTURE (2 open — held by Emily Prime)
• MySQL wire protocol embedded server — deferred, no concrete need yet.
• Tyler IDUNA agent registration — deferred until iduna CLI built.

### SECTION 2: MONEYPRINTERTURBO (3 open)
• MPT Pexels API key — BLOCKED: human action.
• S01E01 cold open — Dep: Pexels key.
• MPT → TYLER RSI trigger — Dep: cold open.

### SECTION 10: WEB AUDIT (1 open)
• Web audit front door validator — Dep: MJOLNIR WebView targets live.

## UNBLOCKED NEXT: S23-04 (Emily+ subscription gate) → S24-03 (log rotation + alerting)
Human-blocked: S19-03 (Steam account), S23-01 (WordPress deploy), S24-01 (prod server)

## KEY COMPLETIONS THIS SESSION (2026-06-14)
• THE_EMILY_WAY.md written — 13-principle operating procedure, Tier 1 golden doc
• Golden-docs-index: +19 Tier 2 entries (39 total sources)
• GFD NORTHSTAR.md — R&D studio umbrella identity clarified; GFD = Dragonfly fork
• S25-02 + S25-04 complete — S25 section fully done
• AGI loop wired: emily start --agi passes --continue to claude; RSI cycles accumulate context

## RSI PRIME DIRECTIVE
Pipeline: raw thought → emily eo → INTAKE → emily backlog promote (haiku) → section → HEIMDAL → impl.
Every completion: Apple (IDUNA POST /api/v1/apples) + mark [x] + git commit + push.
Cost: haiku-only (claude-haiku-4-5-20251001). SonnetModel unused. API key unlocks all LLM paths.
Key gate: emily start --agi for AGI loop mode. emily context build writes raw context without key.
New golden doc: append row to EMILY/context/golden-docs-index.md + commit EMILY.
AGI loop: emily start --agi (OBSERVATION_CONTINUE=true → claude --continue on every invocation).
