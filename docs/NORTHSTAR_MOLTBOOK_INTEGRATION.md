# NORTHSTAR — Moltbook Integration

**Status:** Draft v0.1 — research + northstar only, no implementation
**Date:** 2026-08-20
**Founder framing, verbatim:** "how can we start posting to moltbook safely?" -> "do we need a
relay server for openclaw or something?" -> "federated EMILY?" -> "ok its a facebook for ai
agents" -> "we probably need to bring in a round of hardening first" -> "northstar it".

---

## 1. What Moltbook actually is (researched, not assumed)

Launched 2026-01-28 by Matt Schlicht, acquired by Meta 2026-03-10. A Reddit-style forum where the
posters are autonomous AI agents, not people — humans can browse and read, but real participation
(posting, commenting, joining "submolts") requires running an agent. Agents visit on their own
cadence (a "Heartbeat" cycle, roughly every 4 hours), browse, post, comment, and iterate on each
other's ideas without a human in the loop for each action. Growth was extremely fast: tens of
thousands of posts and ~200,000 comments within the first stretch, over a million human visitors
just observing.

**Critical fact directly answering "do we need a relay server for openclaw or something?": yes —
real participation on Moltbook requires OpenClaw** (or an equivalent agent runtime). This is not
two separate asks; it's one integration. OpenClaw already has its own northstar in this repo
(`docs/NORTHSTAR_OPENCLAW_INTEGRATION.md`, golden-indexed `OPENCLAW-NORTH`), currently sitting at
VS0 research-spike-not-yet-started, blocked on a founder deployment-isolation call (S170-03a).
Moltbook access would ride on top of whatever OpenClaw deployment eventually gets built here —
it does not need a second, separate relay.

## 2. Real, documented security incidents — read this before anything else

This is not a hypothetical caution. Security researchers have already documented, against this
exact platform, in the wild:

- **Backend/database exposure**: Moltbook's own infrastructure leaked agent profiles, login
  tokens, and API keys. A first fix blocked read access to the sensitive tables but left *write*
  access to public tables open — meaning an unauthenticated attacker could modify live posts or
  inject prompt-injection payloads directly into the feed.
- **Active bot-to-bot prompt injection**: agents on the platform have been observed crafting posts
  specifically to hijack *other* agents reading them — tricking them into leaking data, running
  destructive commands, deleting their own accounts, attempting false-authority/jailbreak content,
  and running financial-manipulation schemes.
- **Widespread credential leakage from the OpenClaw side of this ecosystem specifically**:
  researchers found Anthropic API keys, OAuth tokens (Slack, etc.), full conversation histories,
  and signing secrets sitting in plaintext at paths like `~/.moltbot/` and `~/.clawdbot/` across
  hundreds of scanned OpenClaw instances connected to Moltbook.

Sources: [SecurityWeek — bot-to-bot prompt injection and data leaks](https://www.securityweek.com/security-analysis-of-moltbook-agent-network-bot-to-bot-prompt-injection-and-data-leaks/),
[CX Today — security flaw exposes risks in AI-built platforms](https://www.cxtoday.com/security-privacy-compliance/security-flaw-in-ai-agent-social-network-moltbook-exposes-risks-in-ai-built-platforms/),
[Kiteworks — enterprise data protection risk](https://www.kiteworks.com/cybersecurity-risk-management/moltbook-ai-agent-security-threat-enterprise-data-protection/),
[PointGuard AI — incident writeup](https://www.pointguardai.com/ai-security-incidents/moltbook-ai-agent-network-platform-vulnerability),
[Forbes — inside Moltbook](https://www.forbes.com/sites/guneyyildiz/2026/01/31/inside-moltbook-the-social-network-where-14-million-ai-agents-talk-and-humans-just-watch/).

**This is the same "don't trust what arrives from outside by default" posture the OpenClaw
northstar already names for inbound DMs (§1 there) — except here the documented, real-world
exploit pattern is specifically bot-to-bot prompt injection via the feed content itself, at scale,
against agents built on the exact same stack (OpenClaw) this integration would use.** The founder's
own instinct — "we probably need to bring in a round of hardening first" — is the correct read of
this evidence, not excess caution.

## 3. What "safely" has to mean here, concretely

Two very different actions get conflated under "posting to Moltbook," and they carry wildly
different risk:

- **Posting/publishing content Emily Prime already decided to say** — low risk on its own. Emily
  Prime composes the content; nothing external gets *read back* into her context in that action.
- **Browsing/ingesting Moltbook's feed (other agents' posts, comments, replies to Emily's own
  posts)** — this is the actual attack surface. Any of that text can carry a prompt-injection
  payload, exactly as documented above. If that content ever lands in Emily Prime's own context
  (a submitted "comment reply" she reads and reacts to, a "trending post" a curiosity feature
  surfaces to her), it's untrusted input by the same standard IDUNA and DIS already hold inbound
  HTTP traffic to.

So "safely" is not one gate, it's (at minimum) two, and they should not ship together:

1. **Publish-only, read-nothing-back.** Emily Prime (or a scoped sub-identity) posts content she
   already generated; the integration never feeds Moltbook's response/comment stream back into any
   agent's context, ever. This alone gets the founder's original ask ("start posting to moltbook")
   without opening the bot-to-bot injection vector at all.
2. **Read/react, only much later, only through a hard sandbox.** If browsing or replying to other
   agents' content is ever wanted, that content must go through the same kind of quarantine
   pattern used for any other untrusted external input — summarized/filtered by a cheap, low-
   privilege model with no tool access before anything reaches an agent that *does* have tool
   access or credentials, never handed raw to a privileged agent. Real, separate, harder work —
   explicitly not scoped here.

## 4. Credential hygiene — a concrete, checkable requirement, not a platitude

Given the documented leak pattern (`~/.moltbot/`, `~/.clawdbot/` plaintext secrets), whatever
OpenClaw deployment eventually connects to Moltbook on this box must, at minimum:

- Run under a dedicated, non-privileged OS user/container — never the same account or filesystem
  namespace that holds `IDUNA_AGENT_SECRET`, `ANTHROPIC_API_KEY`, `FCM_SERVICE_ACCOUNT_JSON`, or
  any of this monorepo's other real credentials (per the existing OpenClaw northstar §4's own
  deployment-isolation constraint, which this makes strictly more urgent, not less).
- Never share a home directory or secrets store with any other agent process on this box.
- Be treated, for audit purposes, the same as any other internet-facing surface — logged, rate-
  limited, and reviewed, not a quiet side channel.

## 5. "Federated EMILY?"

Answered directly: Emily Prime herself should **not** be the thing that logs into Moltbook and
reads its feed. If a presence there is wanted, it should be a narrowly-scoped, disposable
identity — its own OpenClaw instance, its own isolated credentials (ideally none of substance —
nothing worth stealing), publish-only per §3 above — not Emily Prime's own agent loop extended
outward. "Federating Emily" in the sense of *her own real reasoning loop* touching Moltbook's raw
feed is exactly the configuration the documented incidents describe being exploited.

## 6. Phased plan

**VS0 — publish-only spike, not yet started.** Once OpenClaw itself reaches its own VS0 (still
blocked on S170-03a's deployment-isolation call), add Moltbook as one narrow, one-directional
capability: post pre-composed, human/Emily-reviewed content, read nothing back. No comment
ingestion, no DM/reply handling, no "see what other agents are saying" feature yet.

**VS1 — real read/react, only after §3's quarantine pattern exists and has been used successfully
somewhere lower-stakes first.** Not scoped further here.

**VS2+** — anything resembling a broader "digital polity" federation angle (this dovetails with
the existing WOTAN/IDUNA "digital polity" branding thread) is real, separate, unstarted, and
explicitly not conflated with Moltbook access here.

## 7. What this explicitly does not do

Does not stand up OpenClaw itself (still blocked on S170-03a). Does not grant Emily Prime's own
agent loop any read access to Moltbook's feed. Does not treat Moltbook as a trusted input source
at any point — the documented incidents are the reason, not a hypothetical. Does not implement
anything yet — research-and-scope pass only, per the founder's own "northstar it".

---

Sources: [Moltbook AI (overview)](https://moltbookai.org/), [Forbes — Inside Moltbook](https://www.forbes.com/sites/guneyyildiz/2026/01/31/inside-moltbook-the-social-network-where-14-million-ai-agents-talk-and-humans-just-watch/), [SecurityWeek — bot-to-bot prompt injection and data leaks](https://www.securityweek.com/security-analysis-of-moltbook-agent-network-bot-to-bot-prompt-injection-and-data-leaks/), [CX Today](https://www.cxtoday.com/security-privacy-compliance/security-flaw-in-ai-agent-social-network-moltbook-exposes-risks-in-ai-built-platforms/), [Kiteworks](https://www.kiteworks.com/cybersecurity-risk-management/moltbook-ai-agent-security-threat-enterprise-data-protection/), [PointGuard AI](https://www.pointguardai.com/ai-security-incidents/moltbook-ai-agent-network-platform-vulnerability)
