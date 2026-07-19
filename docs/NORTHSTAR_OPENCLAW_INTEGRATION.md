# NORTHSTAR — OpenClaw Integration

**Status:** Draft v0.1 — research + northstar only, no implementation
**Date:** 2026-07-19
**Founder framing, verbatim:** "openclaw integration northstar do research on that if you need to."

---

## 1. What OpenClaw actually is (researched, not assumed)

[github.com/openclaw/openclaw](https://github.com/openclaw/openclaw) — MIT-licensed, self-hosted
personal AI assistant gateway, written in TypeScript. Not a model, not an agent framework we'd be
replacing Emily Prime with — a **channel gateway**: one Gateway process that connects to 20+
messaging platforms (WhatsApp, Telegram, Slack, Discord, Signal, iMessage, IRC, Microsoft Teams,
Matrix, Feishu, LINE, Mattermost, Nextcloud Talk, Nostr, Synology Chat, Tlon, Twitch, Zalo, WeChat,
QQ, WebChat) and routes inbound messages to an agent, which can be any model/backend the config
points at (`"agent.model": "<provider>/<model-id>"`). Installed via `npm install -g openclaw` +
`openclaw onboard --install-daemon` (launchd/systemd), Docker-deployable, config at
`~/.openclaw/openclaw.json`. Companion apps exist for macOS/iOS/Android that pair as nodes over
WebSocket. Skills extend it via a `SKILL.md` pattern in a workspace directory (a registry, ClawHub,
exists for community skills).

Security model, stated in their own docs: the default `main` session runs tools with full host
access; group/channel sessions run sandboxed (Docker/SSH/OpenShell backends) with restricted tool
access; inbound DMs are explicitly documented as untrusted input. Worth naming plainly: this is
the same posture DIS and IDUNA already take toward inbound traffic — "don't trust what arrives
from outside by default" — just stated for a chat gateway instead of an HTTP server.

## 2. What this would actually be for, here

Emily Prime already reaches the founder two ways: FCM push to MJOLNIR (a purpose-built Android
app), and — as of tonight, once the OAuth dance finishes — Gmail. Every other channel (Telegram,
Discord, Signal, WhatsApp, iMessage) would today mean building a bespoke integration each, the
same way tonight's Gmail work required its own OAuth flow, its own token tool, its own env vars.
OpenClaw is a single gateway that already speaks all of those. **The real value isn't a new
capability — it's collapsing N future bespoke integrations into one.**

Two things worth being honest about, not glossed over:

- **Gmail/email is not in OpenClaw's channel list.** It doesn't replace tonight's work; it's
  complementary reach, not a substitute for it.
- **MJOLNIR overlap needs an explicit answer, not an assumption.** MJOLNIR is a purpose-built
  Android app (Apple feed, FCM push, WebView product panels) — a genuinely different surface than
  "a chat app Emily can also DM through." OpenClaw's own Android companion app is a *pairing
  node* for the gateway, not a competing product surface. These likely coexist rather than
  compete, but that's a real design call, not a foregone conclusion — worth a founder decision
  before building rather than assuming away.

## 3. How it would actually connect to Emily Prime

Emily Prime is not just a model call — she has her own tool-calling loop, IDUNA auth, Apple
filing, the RSI cron cycle. OpenClaw's `"agent.model"` config is built for pointing at a raw
model provider, not a full agent with its own HTTP API on `:8086`. The correct integration shape
is almost certainly: OpenClaw as the channel gateway, configured to call **out** to Emily Prime's
existing `:8086` API as a custom backend/tool (not Emily Prime reimplemented as an OpenClaw
skill, not OpenClaw's agent loop replacing Emily Prime's own). This preserves everything Emily
Prime already does (Apple filing, IDUNA-authenticated actions, the cron cycle) and adds OpenClaw
purely as a new set of inbound/outbound channels in front of the same brain — same shape as MJOLNIR
already has (a client that talks to Emily Prime's API), just a different, much broader front door.

## 4. Deployment posture — a real constraint, not an afterthought

This box already runs every live product this company has, already hit one OOM incident today,
and already declined a Monero miner tonight specifically because of shared-resource risk. OpenClaw
running in its default `main`-session mode (full host tool access) on this same box would be a
materially different risk profile than a sandboxed group session — worth deploying either in a
genuinely isolated context (its own container, or the same "separate dedicated hardware" answer
already given for the miner) or with its tool access deliberately scoped down before it ever
touches production credentials.

## 5. Phased plan

**VS0 — research spike, not yet started.** Stand up OpenClaw against a single low-stakes channel
(Telegram or Discord, both simple to provision, neither touching production credentials yet),
configured to call Emily Prime's `:8086` API for one narrow, read-only action (e.g. "what's the
current backlog status") — prove the gateway↔Emily Prime connection shape before wiring anything
sensitive.

**VS1 — real channel(s), once VS0 proves the shape.** Which channels the founder actually wants
Emily reachable on is a founder decision, not a default — flagged, not assumed here.

**VS2 — MJOLNIR/OpenClaw relationship decision**, once VS1 exists to actually compare against.

## 6. What this explicitly does not do

Does not replace Emily Prime's own agent loop or IDUNA-backed tool execution. Does not replace
tonight's Gmail work (different channel entirely). Does not deploy anything yet — this is the
research-and-scope pass the founder asked for, not an implementation.

---

Sources: [openclaw/openclaw (GitHub)](https://github.com/openclaw/openclaw)
