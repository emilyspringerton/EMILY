# EINHORN_INDUSTRIAL architecture and working process (2026-10-18)

Extracted from the monorepo and per-repo `CLAUDE.md` files, the memory notes, and the Bedrock Racers online build. Sourced from the docs unless
marked verified; the docs lag reality in places (see the last section).

## 1. What it is
A founder-run monorepo of ~60 repos sharing a Go workspace and one operating procedure ("The Emily Way"). IDUNA is the identity/audit
backbone, EMILY Prime the autonomous planner, Claude Code the implementer, PARENA the in-house language the games are compiled from, WOTAN the
esports/stats site where games and ladders are published.

## 2. The control loop
EMILY Prime (:8086), 5-minute cycle: OBSERVE -> DECIDE -> ACT -> PLAN. It files Apples to IDUNA and directs tasks to
`EMILY/signals/tasks`; obs-watcher dispatches Claude Code, which implements. Observe -> direct task -> implement -> Apple -> repeat.
- Apples are the audit trail, backed up in the data-only `APPLES` git repo.
- Founder direction (including "stop") goes through `emily observe` first; it auto-files an Apple before acting.
- Every commit carries a `session: <tag>` trailer (`emily session current`).

## 3. IDUNA, the trust authority
- ES256 JWTs, Google OAuth + local password, M2M agents, device flow, hierarchical RBAC, append-only Apples ledger.
- Every service that needs auth gets a JWT from IDUNA; tokens from elsewhere are never trusted.
- Games register in `games.Registry` (a row + a `*.play` permission); guest-register / guest-login / sso-exchange are generic per game.
- WOTAN pages never render a login form: they link to `iam.okemily.com`, read the token from the URL fragment, share one session
  (`js/iduna-sso.js`).
- Game gateways run inside IDUNA (D2, Red Garden, EMPIRE, shankpit-460, Bedrock Racers): WebSocket, matchmaking, match recording.
- IDUNA_PRO is the standalone extraction of the core ("Emily for Business"); the licence/trademark text is not lawyer-reviewed.

## 4. PARENA, the language
S-expression syntax, compile-time region-based memory safety, no GC; targets C, TypeScript, JVM, WebAssembly.
- Games are PARENA `.prn` compiled to wasm: no allocator, zero imports, fixed 8 MB memory.
- Page JavaScript is generated (PARENA -> TypeScript -> type-stripped JS), never hand-written (founder rule).
- Core dependencies are written in PARENA first; third-party tools are labelled stopgaps.
- Siblings: BURROW (Go rewrite of the compiler), DUNG (terminal+editor), SAND (editor), JEWEL (Jupyter kernel), LO (esolang), MIXFORGE (DJ app).

## 5. Games and multiplayer patterns
BRAWLPIT (local wasm), shankpit-460 (wasm + IDUNA human FFA lobbies), Red Garden (wasm client of the real C arena), D2 (wasm client of
dw2_server, ruleset 2), EMPIRE (wasm sim, IDUNA wazero referee), DEADWEIGHT (card game, deck browser), Bedrock Racers (rocket-car Quidditch).
- Server-authoritative relay: D2, Red Garden, shankpit-460.
- Deterministic lockstep: Bedrock Racers, EMPIRE. A match is a pure function of its inputs; only command bytes travel; the server steps its own
  copy of the same wasm as referee. Verified: V8 and wazero agree at every checkpoint through 8,000 ticks; two real browsers stayed in sync
  through the server's hash checks (WOTAN `scripts/bedrockracers-online-test.mjs`, 10/10).

## 6. Runtime and deployment
One GKE cluster, namespace `emily`; images in Artifact Registry; the `parena-gitops` sync Deployment bumps tags (some manual, e.g.
`DEADWEIGHT2_TAG`). WOTAN is static HTML/CSS, no build step; wasm and JS outputs are committed because the deploy image is built from tracked
files. A push to WOTAN `main` deploys to production. The Gateway routes `/api` (WebSockets included) to IDUNA. Deploy IDUNA first, then WOTAN.
Secrets live in the running pod and the `iduna-env` Secret; local copies go stale.

## 7. How work is done
1. Read `EMILY/BACKLOG.md`; take the top open item. Founder real-time direction overrides the pick but is always logged.
2. File an Apple before marking done; update CHANGELOG; update the README when its claim changes; register golden docs and resync `GOLDEN_DOCS`.
3. Tests before commit; one atomic change per commit; push immediately; commit as the repo's configured identity.
4. CONSTRUCT file (deterministic plaintext source snapshot) per release; LZ4 is the default compression.
5. Bedrock Racers gate: 0 imports, physics/magnet/fire/engine/net tests, byte-identical soak line; ask before pushing WOTAN main; `wip/` for unfinished work.

## 8. Cautions
- Per-repo docs still describe the old systemd/localhost model in many places.
- Untested for most games: Firefox, Safari, phones, audio, physical gamepads. Bedrock Racers online has not run on the live cluster.
- SAND, DUNG, LO, MIXFORGE, BIG_O, CarePyre are plans with no working product.
- Never publish claim codes or secrets; repos are public.
