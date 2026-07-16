# Queued Fable prompt — IDUNA front door funnel for Agents and Unagents

**Status:** queued, not yet dispatched. Written 2026-07-16.
**Dependency:** dispatch this *after* the SAGA reality-vs-documentation audit of the 14 original
KIKORYU VS specs lands (`IDUNA/docs/VS_REALITY_AUDIT.md` + `IDUNA/docs/kikoryu/*.md`, especially
whatever the audit concludes about VS0 and VS7) — this prompt references that work as required
reading and the design question sharpens considerably once it exists. Check it's landed
(`git -C /home/fatbaby/IDUNA log --oneline -10` should show it) before sending this.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Prompt body (copy from here down)

You're designing IDUNA's "front door" — the entry funnel every actor passes through before it can
do anything in the EINHORN_INDUSTRIAL ecosystem. This is a real spec-writing task with a real
deliverable, not a memo about the plan.

### Required reading first

1. `/home/fatbaby/IDUNA/docs/VS_REALITY_AUDIT.md` and `/home/fatbaby/IDUNA/docs/kikoryu/` (all
   files) — the SAGA-style reality audit of IDUNA's original founding specs. Pay particular
   attention to whatever it concluded about **VS0** ("The Ritual Gate" — the human identity/gating
   funnel: Google OAuth → THE_HONOR_CODE acceptance → gamertag lock-in → device auth handoff,
   state machine ANON → HONOR → HANDLE → READY) and **VS7** ("Agents vs Unagents" — the original
   governance distinction between privileged, bounded-authority actors and default participants).
2. `/home/fatbaby/IDUNA/docs/archive/kikoryu-vs-original/vs0.md` and `vs7.md` — the original specs
   themselves, archived verbatim, for full context beyond the audit's summary.
3. Read the actual code for both of IDUNA's real, currently-live entry paths — they are
   structurally asymmetric and that asymmetry is the core problem this spec needs to resolve:
   - **The human/Unagent path** (dynamic, self-service): `internal/http/handlers/device.go`
     (`/auth/device/start`, `/auth/device/poll`), `internal/auth/device/service.go`
     (`ErrHonorCodeRequired` and the honor-code acceptance flow), `internal/http/handlers/me.go`
     (`GET /api/v1/identities/me`), `internal/http/handlers/auth.go`'s Google OAuth handler. A
     human arrives, authenticates via Google, accepts an honor code, claims a gamertag, and
     reaches a `READY` state — a real, live, dynamic ceremony.
   - **The agent path** (static, human-bootstrapped): `internal/http/handlers/auth.go`'s
     `AgentAuthHandler` (`POST /api/v1/auth/agent`, `agent_name` + `agent_secret` → JWT),
     `config/agents.json` (hand-edited registry, applied via `cmd/bootstrap`), and
     `cmd/bootstrap`'s seeding logic. An agent doesn't arrive anywhere — a human edits a JSON
     file, runs a command, and the agent's credentials simply exist thereafter. There is no
     dynamic provisioning, no self-service registration, no ceremony, no analog to "honor code"
     or "gamertag lock-in" or a `READY` state machine.
4. `/home/fatbaby/EMILY/docs/hq-specs/HQ-SPEC-PRIME-101-norn-loop-kernel.md` and
   `HQ-SPEC-DOC-102-saga-curation-lifecycle.md` — two *new* agents (KAREN, the ledger controller;
   SAGA, the documentation curator) are already planned (see `EMILY/BACKLOG.md` SECTION 142/143)
   and will need to onboard into IDUNA through whatever path exists when they're built. Design with
   them as the near-term test case, not a hypothetical.
5. `EMILY/BACKLOG.md` — search "STRATEGIC PRIORITY ORDER" and the Act II arc (SECTIONS 141-146)
   for the current state of what's actually being built, so the funnel design doesn't invent
   requirements the rest of the system doesn't have yet.

### The design question

IDUNA is the single front door for the whole ecosystem, but it currently runs two unrelated front
doors wearing one building's sign. The human funnel is a real, considered, dynamic ceremony with a
state machine and a moment of consent (the Honor Code, deliberately built as ritual, not just a
checkbox — see VS0's visual-language section for why that mattered to the original design). The
agent funnel is a bootstrap script. That's not necessarily *wrong* — agents and humans are
legitimately different kinds of actors with different risk profiles and different notions of
consent — but the asymmetry has never been a *decision*, it's just how the code happened to grow,
and nobody has asked whether it's the right shape now that new agents (KAREN, SAGA) are coming
online as first-class citizens rather than one-off bootstrap entries.

Design a coherent front-door architecture that answers, explicitly:

- Should agent onboarding gain a real funnel — a sequence of states analogous to
  ANON → HONOR → HANDLE → READY, appropriate to what an agent actually needs to consent to and
  declare (its capability scope, its custody/owner, whatever the equivalent of "the Honor Code" is
  for a piece of software making autonomous decisions)? Or is the static bootstrap model *correct*
  for agents precisely because they're not moral actors and shouldn't get a ceremony that implies
  they are — and if so, say why explicitly, don't just default to symmetry for its own sake.
- What, if anything, should be genuinely shared infrastructure between the two funnels (the JWT
  issuance core obviously already is) versus what should stay deliberately separate? Is there a
  shared "front door" concept at all, or are "funnel for Agents" and "funnel for Unagents" simply
  two different services that happen to share a signing key?
- How does VS7's Agent/Unagent *governance* distinction (bounded authority, auditability of every
  privileged action) relate to this *onboarding* question? They're adjacent but not the same
  problem — arriving vs. acting. Be precise about where one ends and the other begins; don't let
  the spec collapse them into the same thing if the audit found they're genuinely different axes.
- Where does the tournaments-platform direction (VS2, now the near-term product thrust per the
  audit and `NORTHSTAR_KIKORYU.md`) touch the funnel, if at all? A tournament entrant is a human
  going through the existing Unagent path — does tournament-specific gating (age/jurisdiction
  attestation, the play-money-only invariants VS2 specifies) belong *in* the front-door funnel as
  a new state, or downstream of it as VS2's own concern? Take a position and defend it.
- What's the concrete migration path from today's two-unrelated-paths reality to whatever you
  propose? This has to be buildable incrementally against a live system with real registered
  agents already depending on the current bootstrap flow — don't propose a redesign that requires
  a flag-day cutover.

### What to produce

A single spec document, `IDUNA/docs/kikoryu/FRONT_DOOR_FUNNEL.md` (or a better name if the
`kikoryu/` directory's existing naming convention suggests otherwise — check what's there first),
written in the same register as the audit's rewritten VS0/VS1/VS2 docs: real technical content
(state diagrams if useful, concrete endpoint/schema proposals), not a vision essay. Include a
`supersedes:` line if it genuinely supersedes specific sections of the archived VS0/VS7, or say
explicitly that it's a new synthesis that doesn't supersede either (if the audit already handled
VS0/VS7's own supersession and this is additive). Register it in
`EMILY/context/golden-docs-index.md` at tier 1, matching the existing table format. Update
`IDUNA/CHANGELOG.md` with a dated bullet. Do not touch `EMILY/BACKLOG.md` — if the spec implies
concrete build steps, list them inside the spec document itself as a numbered build order, the
same way the HQ-SPEC docs and the rewritten VS0-2 docs do, and leave backlog transcription as a
deliberate separate decision for later.

Commit IDUNA and EMILY separately, push both. Report back: the core call you made on agent
onboarding (real funnel vs. deliberately-static, and why), the file path, and confirmation both
commits landed and pushed.
