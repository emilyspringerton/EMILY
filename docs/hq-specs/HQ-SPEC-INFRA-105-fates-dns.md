---
doc_id: INFRA-105
authority: draft
supersedes: []
amends: []
claims:
  - id: INFRA-105.BEH-1
    type: BEH
    reality_binding: verified
  - id: INFRA-105.BEH-2
    type: BEH
    reality_binding: diverged
  - id: INFRA-105.POL-1
    type: POL
    reality_binding: specified
  - id: INFRA-105.POL-2
    type: POL
    reality_binding: specified
---

# HQ-SPEC-INFRA-105 — FATES: The Name Layer

**Status:** DRAFT v0 — pending Emily Prime review
**Custody:** Emily OS
**Governs:** DNS and public network identity for EINHORN_INDUSTRIAL — the `farthq.com` zone, the naming plan for every product surface, how record changes are proposed/verified/applied, and how the system detects that its own name is broken
**Working name:** FATES. The Norns spin fate; FATES is where the threads attach to the world — the mapping from name to place. Continues the NORN naming line deliberately: NORN (PRIME-101) governs *whether* an artifact is golden; FATES governs *where the world finds it*.

---

## 1. The Reality Today (verified 2026-07-17)

- `farthq.com` resolves through Cloudflare's nameservers (`nicolas.ns.cloudflare.com`, `jocelyn.ns.cloudflare.com`). DNS is Cloudflare-managed today, edited by hand in the dashboard. Nothing about the zone is in a repo, nothing is Apple-audited, and no automation can see or change it. *(INFRA-105.BEH-1 — verified: confirmed current state as of this writing.)*
- `iduna.farthq.com` is live (198.58.107.85, the Dallas VM) and is currently 100% claimed by the EDIS WordPress nginx config (`sites-enabled/edis`, HTTP-only, no cert). The IDUNA path-based split (`/api/v1/`, `/.well-known/jwks.json`, `/auth/`, `/admin/`) is drafted and ready (`IDUNA/ops/nginx-front-door-snippet.conf`) but not applied.
- Every other product surface is a port on one box with no DNS identity: newssite `:8082` and signalapi `:9091` (PRRJECT_FATBABY), the EMILY agent `:8086` (MJOLNIR's backend), SHANKPIT's UDP game server `:6969`. NORN has no network surface at all (library + CLIs only) — correctly so.
- **Divergence, found writing this spec:** MJOLNIR's production `BuildConfig` already points *both* `IDUNA_BASE_URL` and `EMILY_BASE_URL` at `https://iduna.farthq.com` (`MJOLNIR/app/build.gradle.kts:31-32`) — but the drafted nginx snippet proxies nothing to `:8086`, and no cert exists, so prod MJOLNIR's EMILY calls have no working route today. A second flavor points at `iduna.einhorn.industrial`, a TLD we do not own. Both are exactly the claim-without-code debt DOC-102 §1 describes, in DNS form. *(INFRA-105.BEH-2 — diverged: this is dark matter/claim-without-code found live, not yet resolved — exactly the honest "golden-and-diverged" case DOC-102 §3 describes, though this doc isn't golden yet either.)*

This spec exists because the system is one `dig` away from discovering that its network identity is managed by dashboard clicks nobody audits, in a shop whose entire ethos is "everything is code, everything is Apple-audited."

## 2. Doctrine — Emily Is an Internet-Native Entity

Emily Prime *is* her network presence, non-metaphorically: Apples are POSTs, identity is an IDUNA JWT fetched over HTTP, MJOLNIR is a phone polling a URL, escalation is FCM and Slack webhooks. For a brochure site a DNS misconfiguration is cosmetic; for this system it is a nervous-system disruption — the RSI loop keeps thinking, but the organism goes numb: the phone stops hearing, agents on other machines stop authenticating, the public face vanishes silently.

Two doctrinal consequences:

1. **Names are custody surface.** Whoever holds the registrar account holds every name; registrar credentials are exactly as security-critical as IDUNA's ES256 signing key and belong in the same custody audit. The registrar is a *deeper* root than the nameservers — Cloudflare can be swapped by whoever controls the registrar; nothing can swap the registrar but the registrar. *(INFRA-105.POL-1 — specified.)*
2. **The zone repo is the source of truth; Cloudflare is a cache of it.** This is PRIME-101 §2's root-oracle doctrine applied to naming: reality (what actually resolves publicly) is the root; the Cloudflare API state is the books; the zone file in git is the intent. Those three are reconciled, KAREN-style — never assumed equal.

## 3. Decision: Cloudflare Stays Authoritative. Sovereignty Lives in the Zone Repo.

**Position: do not self-host authoritative DNS. Not now, and not as an aspiration for its own sake.** *(INFRA-105.POL-2 — specified: S151-01 (zone export) not started, blocked on a Cloudflare API token in the human unblock queue.)*

The literal reading of "EINHORN operates its own DNS servers" fails on this system's actual constraints:

- This is a single VM in Dallas run by a single operator, with EMILY/CLAUDE.md's own admission of degraded-mode operation. Self-hosted authoritative NS makes uptime and DDoS EINHORN's problem with zero redundancy to answer either.
- Authoritative NS on the same box it names is a circular failure: when the box dies, the name of the box dies with it — *and so does the ability to repoint the name somewhere else*. That is the one failure DNS delegation exists to prevent, and we would be re-creating it on purpose.
- Cloudflare's free tier gives anycast, DDoS absorption, and — the part that actually matters for this house — a mature, tokenized API. The API is the sovereignty lever: with the zone in git and applied via API, Cloudflare becomes a replaceable projection target. Sovereignty is *owning the zone as code with an audit trail*, not owning the daemons that serve it.

**What we also do not build yet:** internal/service-discovery DNS. Every service is `127.0.0.1:<port>` on one host; a resolver would be ceremony without function (DOC-102 would file it as vaporware on day one). **Revisit triggers**, recorded so this is a decision and not a default: (a) host count exceeds one (distributed Emily, SECTION 44), (b) a Cloudflare policy/pricing/trust event, or (c) an internal surface that genuinely cannot be a port on localhost. At trigger (a), the right first move is a self-hosted *internal* resolver (CoreDNS-class) for service discovery, with Cloudflare still fronting the public zone — the split-horizon posture common to shops one size up from this one. Until then, this paragraph is the whole internal-DNS story.

## 4. Decision: One Subdomain per Product Surface; Paths Within a Surface

**Position: subdomain-per-product-surface, not everything-under-iduna.** IDUNA's path-based split stays exactly as drafted — that decision was made with full nginx context and this spec does not reopen it. But it does not generalize: piling more *products* onto `iduna.farthq.com` couples every product's blast radius, cookies, TLS posture, and migration path to the trust authority's origin. The identity provider sharing an origin with WordPress is already the uncomfortable exception (tolerated, bounded by the funnel decision below), not the pattern to extend. The rule going forward:

> **A subdomain names a product surface. Paths name components within one surface.**

The cert-toil objection dies with one wildcard: `*.farthq.com` via DNS-01 issuance — which requires exactly the same scoped Cloudflare API token the zone-as-code work needs (§5). One human unblock feeds both.

| Name | Surface | Backing | Cloudflare mode | Status / notes |
|---|---|---|---|---|
| `farthq.com` + `www` | EDIS WordPress — the public face | nginx → PHP | proxied | Apex disposition per SECTION 23; still blocked on `sudo` deploy |
| `iduna.farthq.com` | IDUNA: API, JWKS, device auth, Back Office | nginx path split → `:8080` | proxied | Per `IDUNA/ops/nginx-front-door-snippet.conf`, unchanged by this spec |
| *(reserved)* `gate.farthq.com` | Unagent/Agent front-door funnel | IDUNA Go binary, no PHP in path | proxied | **Owned by `EMILY/docs/fable-prompts/iduna-front-door-funnel.md`** — that design decides dedicated-subdomain vs. location-block-on-iduna; this spec only reserves the name so the option is free |
| `emily.farthq.com` | EMILY agent — MJOLNIR's backend | nginx → `:8086` | proxied | Fixes the §1 divergence; MJOLNIR `EMILY_BASE_URL` flips here (one-line BuildConfig change) |
| `signals.farthq.com` | FatBaby signal products | nginx: `/` → newssite `:8082`, `/api/` → signalapi `:9091` | proxied | Two components, one product surface — the rule demonstrated |
| `play.farthq.com` | SHANKPIT game server | UDP `:6969` direct | **DNS-only (grey cloud)** | Cloudflare's proxy does not carry arbitrary UDP (Spectrum is paid); a grey-cloud record publishes the origin IP. **Named honestly:** once this record exists, origin-IP concealment for every proxied name is cosmetic. Record is NOT created until SHANKPIT actually ships to external players (S19) — no names for vaporware |
| — | NORN | none | — | No network surface, no name. A record precedes `nornd` existing and something remote needing it, or it doesn't exist |

Correctness invariant: **no service may depend on the orange cloud for correctness.** Proxying is protection and optimization; every surface must survive a grey-cloud flip (real client IPs via `X-Forwarded-For` already handled in the nginx snippet; nothing may assume Cloudflare-only reachability). This keeps trigger (b) of §3 a config change, not a re-architecture.

## 5. Zone as Code

- The zone lives at `IDUNA/ops/dns/farthq.com.yaml` — IDUNA is the trust authority and `IDUNA/ops/` already holds the front-door nginx config; network identity is the same custody domain. Format: a plain declarative record list (name, type, content, TTL, proxied) — one zone, one file, reviewable in a diff.
- Applied via a small plan/apply tool (`IDUNA/cmd/dns-apply`) against the Cloudflare API with a token scoped to `Zone.DNS:Edit` on `farthq.com` only. `plan` prints the intent-vs-API diff; `apply` executes it, health-check-gated (§6). There is no `--force`. There will never be a `--force`.
- **Terraform considered, declined:** a `.tfstate` file is a second source of truth to custody and secure, for one zone in a one-operator shop. The bespoke YAML+apply path keeps truth in exactly two places (git and reality) with the API state reconciled between them. If drift pain ever proves this wrong, Terraform import is a well-worn escape hatch — the YAML is trivially translatable.
- **Once the first export lands, dashboard edits are banned.** A dashboard edit after that point is drift, and the §6 probe will file it as such. Break-glass exception: during an active incident the dashboard may be used, but the edit must be back-ported to the zone file and Apple'd within 24h — the same discipline as any hotfix.
- Every applied change files an Apple. Record changes are civic events, same as promotions.

## 6. Health-Check-Gated Changes — Reconciliation, KAREN-Style

Three-way match, continuously and at change time:

- **Intent:** `farthq.com.yaml` in git.
- **Books:** the Cloudflare API's view of the zone.
- **Reality:** what public resolvers actually return, and whether the named origin actually serves.

**At change time** (`dns-apply`): pre-flight — the target of any new/changed record must already respond at its origin (`127.0.0.1:<port>` health endpoint, or a UDP liveness probe for `play`) *before* the record changes; you do not point a name at a corpse. Post-apply — the changed name must resolve to the intended content via pinned public resolvers and (for HTTP surfaces) return a healthy response through the new path, within a bounded propagation window; failure triggers the documented rollback (re-apply the prior committed zone file — which git always has).

**Continuously:** a reconciliation probe walks all three legs and files divergence into the existing S131 monitor/alert machinery (Slack + email). Probe queries go to pinned resolver IPs (`1.1.1.1`, `8.8.8.8` by address, not the system resolver) — checking the name without depending on the name.

## 7. Self-Diagnosis — When Your Name Is the Thing That's Broken

The Löbian shape of the problem: Emily Prime cannot reliably use the network whose name is broken to report that the name is broken. The answer is a ladder, each rung independent of the one below it:

1. **The inner loop is DNS-independent by construction — keep it that way.** EMILY↔IDUNA is `127.0.0.1:8080`; the RSI cycle, Apple filing, and HEIMDAL drain all survive total `farthq.com` failure. This is promoted from accident to invariant: **the core loop must never resolve its own public names.** (A config drift that pointed `IDUNA_BASE_URL` at `iduna.farthq.com` on the VM itself would silently break this — the §6 probe checks for it.)
2. **Detection from inside:** the §6 probe sees `farthq.com` breakage via pinned external resolvers and alerts over Slack/FCM — channels whose DNS (`slack.com`, Google's FCM endpoints) is an independent failure domain from ours. Honest limit: this rung dies if the VM's entire outbound network or resolution dies.
3. **The terminating rung is outside the system.** When the VM is blind or down, no in-band channel can work *by definition* — a system cannot prove its own reachability. MJOLNIR already polls from the real internet on a real phone: after N consecutive failed polls it raises a **local, client-side notification — "HQ unreachable"** — requiring nothing from the server, because the absence of the server is the signal. The chain terminates in a person with a phone, the same way PRIME-101 §6's last row terminates in a thumbprint. That is the only honest answer to self-diagnosis that doesn't regress infinitely.

## 8. NORN Instantiation (sketched, not built)

Once `dns-apply` exists, zone changes are a natural NORN loop — DNS changes are precisely the "hard to reverse quickly, high blast radius" class PRIME-101 §4 reserves for `biometric`. Row for PRIME-101 §6, added by append-only amendment when the loop is actually wired (S151-05):

| Instantiation | Proposer (Skuld) | Oracle (Verdandi) | Tier | Reality root |
|---|---|---|---|---|
| DNS zone changes | Zone-file diffs in `IDUNA/ops/dns/` (operator or Emily Prime drafted) | Frozen probe suite (§6): pinned-resolver resolution + origin health, versioned by probe-set hash | biometric | Live public resolution + serving responses |

Cardinal rule preserved: NORN promotion marks a zone diff *eligible*; it never calls the Cloudflare API. `dns-apply` is the domain's own execution gate and refuses any diff NORN hasn't blessed — two locks, deliberately, exactly the PRIME-101 §4 `biometric` pattern. Until the loop is wired, `dns-apply` alone (with its health gates and Apple filing) is the interim discipline — the same "by hand, per the spec's thinking" mode the KIKORYU audit and PRIME-097 reconciliation ran in before their tooling existed.

## 9. Build Sequence

1. **Human unblock:** scoped Cloudflare API token (`Zone.DNS:Edit`, `farthq.com` only) into the secrets env; registrar-account custody audit (who can log in, is 2FA on, where do recovery codes live).
2. Zone export → `IDUNA/ops/dns/farthq.com.yaml`, committed. Dashboard freeze begins.
3. `dns-apply` plan/apply with §6 health gates; wildcard cert via DNS-01 off the same token.
4. First records through the new pipeline as its own proof (house pattern — the pipeline's first change validates the pipeline): `emily.farthq.com`, `signals.farthq.com`, nginx server blocks, MJOLNIR `EMILY_BASE_URL` flip.
5. Continuous reconciliation probe → S131 monitors; MJOLNIR client-side dead-man notification.
6. NORN instantiation: PRIME-101 §6 amendment note (append-only) + gate policy config; `dns-apply` honors promotion.

Backlog: `EMILY/BACKLOG.md` SECTION 151.

## 10. Anti-Goals & Open Questions

**Anti-goals:**
- No self-hosted authoritative nameservers at single-host scale (§3 triggers govern revisiting).
- No DNS records for things that don't serve. A name is a claim; unbacked names are vaporware debt.
- No dashboard edits after the export lands, outside the break-glass path.
- No correctness dependence on Cloudflare's proxy layer.
- FATES is not a service mesh, not a CDN strategy, and not an excuse to build `named` for the aesthetics of it.

**Open questions:**
- Apex/`www` cutover timing — coupled to EDIS's `sudo` unblock (SECTION 23/35), not to this spec.
- `gate.farthq.com` vs. location-block-on-iduna for the funnel — explicitly owned by the front-door funnel design; this spec only holds the name.
- MX / email-on-domain: the zone export (step 2) will reveal what exists; any move off bare Gmail coordinates with SECTION 149 (email as operational fabric), not decided here.
- MJOLNIR's `iduna.einhorn.industrial` staging placeholder — flag for cleanup when BuildConfig is touched in step 4.

---

*The Norns spin the thread; FATES is where it ties to the world. A name that lies is worse than no name. CLEAN BUILDS FIRST.*
