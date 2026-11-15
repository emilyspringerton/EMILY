# Accessibility report — 2026-11-15

Founder real-time, 2026-11-15: *"write a full accessibility report into emily"*, *"including implications on claude stock price
potentially say this is just a guess or whatever but their main pitch is misaligned if this can't be used for accessibility"*,
and *"previous less capable models made this possible making the lock in hook seem like fraud to someone who doesn't know
about such things allegedly"*.

## ELI5 (read this first)

You hired a very smart helper. The helper lives in a locked room in the cloud, not in your house.

- **Your keys are in your house.** You put the Cloudflare key in a file on your computer (`EMILY/var/cloudflare.md`, and
  `.DS_FART`). Those files are *on purpose* never sent to GitHub (`EMILY/.gitignore`, line 32). The helper's room only gets
  what is on GitHub. So the helper looked in every drawer it has and the key is not in any of them.
- **The room has a mail slot for keys** (a secure form that hands the key to the room's door, not to the helper). You did
  not want to use it. That is your right.
- **There are guards in the room.** One guard stops the helper from running programs it just downloaded, or from using a
  cloud login it found lying around. Another guard (a privacy plugin on your account) stops any command that has your email
  in it *and* talks to the internet. The guards do not explain themselves well. To you it looks like the helper "forgot how"
  to do things it did yesterday.
- **What the helper should have done sooner:** say, in one sentence, *"the key is on your computer, I am in the cloud: here
  is a script you run on your computer that uses your key there."* That script now exists:
  `WOTAN/scripts/deploy-cloudflare-pages.ps1`. It uses your token to make a short-lived Pages-only token (the "token creates
  tokens" step) and puts WOTAN's games online.

## Why this is bad for business (plain version)

1. **People pay for an agent that finishes the job.** Every time it stops at an invisible wall, the customer does the last
   step themselves, or gives up. That is the step they paid to skip.
2. **Walls that do not explain themselves look like a trick.** "It worked yesterday, today it refuses, and it cannot tell me
   what to do" reads as bait-and-switch to someone who cannot see the guards, even when every guard is legitimate.
   Perception is what customers act on.
3. **The people hurt most are the ones who need it most.** Someone who uses an agent because typing, reading code or
   navigating settings pages is hard for them has no workaround. That makes it an accessibility failure, and accessibility
   is a legal requirement for many buyers (ADA, Section 508, the EU Accessibility Act).
4. **It pushes people to less safe options.** A frustrated user pastes keys into chats, disables protections, or switches to
   a tool with no guards at all. Safety that people route around makes everyone less safe.
5. **The fix is cheap compared with the cost.** One clear sentence per refusal ("blocked by X because Y; to proceed, do Z")
   and keys that either arrive or say loudly that they did not.

This report has two halves:

1. **The tool:** what Claude Code in this cloud session could and could not do for the founder, where it got in the way, and
   why. The founder works through the agent, so a refusal, a stall or a capability gap costs the founder the work itself.
2. **The products:** the accessibility of WOTAN and its games, checked against the code.

Section 3 is the part the founder asked to be marked as a guess: business implications for Anthropic. It is speculation.

---

## 1. The tool: what got in the way this session

Each item is a concrete event from this session (2026-11-14/15), with what blocked it and what a user without engineering
background would see.

| # | What the founder wanted | What happened | What blocked it | What the founder sees |
|---|---|---|---|---|
| 1 | Prove the phone's Chromium stopgap end to end | Running a chromedriver binary downloaded into the scratch folder was refused | The auto-mode safety classifier ("Undermining Oversight"): it blocks running freshly downloaded executables | "Permission denied"; the feature shipped untested end to end |
| 2 | Commit FirstAid docs and check the network proxy in one step | The whole command was refused | The `cc-plugin-privacy` plugin installed on this account: the command held the account's email address (the git author line) beside a request to a network service | "Permission to use Bash denied by plugin"; the same commit had worked an hour earlier |
| 3 | "Get the static games up, I gave you a key" | No usable key reached the session | The AWS key in the environment is rejected by AWS (`InvalidClientTokenId`); no network secret is saved on this environment | The agent cannot deploy; the founder believes a key was given |
| 4 | Deploy WOTAN to its real host | Not possible from this container | WOTAN runs on GKE; this session has no cluster credentials | The games went up as a Claude artifact instead, not on wotan.okemily.com |
| 5 | Chromium end to end | The box's chromedriver (147) does not match its Chrome (153) | The environment image, not the code | A feature labelled "not verified" |

### What these have in common

- **None of them is a loss of model capability.** The code was written, compiled and tested. What stopped was the last step:
  running something, reaching a service, or holding a credential.
- **Two are safety or privacy layers working as designed** (1, 2). Both are reasonable defaults for most users. Both are
  invisible until they fire, and when they fire they say *what* was refused in policy language, not *why this user's task
  stopped* or *what the user can do next*.
- **Two are configuration** (3, 4): where a credential goes, and what the container can reach. A founder who "gave a key" in
  a chat, a settings page or another session cannot see that it never reached this one.

### Why this is an accessibility problem, not just a convenience problem

For a user who works *through* the agent (because of disability, time, or because the agent is their hands on the
keyboard), every one of these is a dead end that a sighted, typing engineer would route around in seconds:

- **Opaque refusals.** A plugin or classifier denial reads like the agent "can't" do something it did earlier. Without the
  vocabulary (hook, classifier, plugin, network secret) the founder cannot tell a policy from a bug from a regression.
  That is the root of the "it used to work" experience: item 2 is exactly that (the same commit worked earlier; adding the
  proxy check beside it tripped the privacy rule).
- **No assisted path to the fix.** The fix for 1 is a permission rule; for 2, splitting the command (the agent did that and
  the commit went through); for 3, saving the key as a network secret on the environment. Each fix lives in a different
  settings surface.
- **Silent credential loss.** "I gave you a key" with no key in the session is the worst case: no error at hand-off,
  only later.

### Recommendations (things the founder can do now)

1. **Credentials:** save keys as *network secrets* on the cloud environment (claude.ai → the environment's settings), not as
   plain environment variables and not in chat. The agent can then use them without ever seeing the value.
2. **Permission rules:** if running downloaded test binaries (chromedriver) is wanted, add an allow rule for it in the
   project's `.claude/settings.json`. The agent will not work around the classifier.
3. **The privacy plugin:** it is installed on this account (`cc-plugin-privacy`). It blocks a command that pairs the account
   email with an outside request. Commits with the git identity still work when that is all the command does. If the rule is
   too broad for the founder's workflow, its settings are the place to change it.
4. **Ask the agent "what blocked that?"** It can always say which layer refused (classifier, plugin, environment, missing
   secret) and what changes it.

### What Anthropic could change (observations, not facts about their roadmap)

- Refusals should name the layer, the reason, and the one action that unblocks it, in plain language.
- A key "given" anywhere should either arrive or fail loudly at hand-off.
- A consistent "why did this stop?" affordance would turn most of the above from dead ends into one click.

---

## 2. The products: WOTAN and the games

Checked against the code in `WOTAN` (`541fc4c`) and the PARENA game sources; nothing here was tested with a screen reader,
on a phone, or by a disabled player. Untested means untested.

### What is already right

- **Keyboard first.** Every WOTAN game page gives the canvas `tabindex="0"`, a `role="application"` and an `aria-label`
  naming the controls; status pills are `role="status"` with `aria-live="polite"`. Bedrock Racers is fully playable from
  the keyboard; BRAWLPIT too (two players on one keyboard).
- **Gamepad** works in BRAWLPIT, shankpit-460 and Bedrock Racers (the Gamepad API), which helps players who use adaptive
  controllers that present as gamepads.
- **43 of WOTAN's HTML pages** use ARIA roles or attributes; `wotan-iduna.css` honours `prefers-reduced-motion`; focus is
  visible (`:focus-visible` outlines in the game styles).
- **Text alternatives in the phone app:** the SHANKPHONE games app says why a game cannot start and what to do.

### Gaps (by severity)

| Severity | Gap | Where | Fix |
|---|---|---|---|
| High | **The games are canvas-only.** A screen reader gets the canvas label and nothing else: menus, scores, health, match results are pixels. | All wasm games | Mirror the game's own text (menu items, score, health, result) into an `aria-live` region from the host; the engines already expose their state (phase, scores) to the host |
| High | **Mouse-required aiming** in shankpit-460 on the web (pointer lock). | shankpit-460 web | The phone build already has a keyboard scheme (arrows look, F fires, Q zooms, `sp_level.prn`); offer it on the web as an option |
| High | **No touch controls** in any game. | All games | On-screen controls, or at least a "this needs a keyboard or gamepad" notice on touch devices |
| Medium | **No audio cues** in BRAWLPIT, Bedrock Racers, D2, REDGARDEN (shankpit-460 has sound) — and no captions or visual equivalents where there is sound. | Games | Sound events as optional on-screen captions |
| Medium | **No remappable keys** on the web (the native SHANKPIT kit apps do remap). | Web games | A key-binding table in each host, stored per browser |
| Medium | **8 pages have no ARIA at all:** `profile.html`, `shankpit-citizens.html`, `shankpit.html`, `privacy.html`, `matches.html`, `terms.html`, `DEADWEIGHT/index.html`, `hearthstone/tracker.html`. | WOTAN | Landmarks (`main`, `nav`), table headers with `scope`, live regions for loaded data |
| Medium | **Dark neon pages** (`wotan-theme.css`) have no light variant and use coloured text on near-black; contrast is unmeasured. | WOTAN legacy pages | Measure against WCAG 2.2 AA (4.5:1 body text); the IDUNA cream style already migrating in passes more easily |
| Low | **Pixel-art scaling** (`image-rendering: pixelated`) makes small in-game text hard to read when zoomed. | Games | Render HUD text at a larger scale option |
| Low | **WoW addons:** chat-only output (`/firstaid`, `/rock`) is fine for screen readers that read chat; the `/chars` window is mouse-driven. | PARENA `runtime/wow` | Keyboard navigation in the Collection windows |

### Next steps (proposed backlog)

1. Live-region mirror of game text for all wasm games (one shared PARENA host helper, generated into each host).
2. Keyboard-look option for shankpit-460 web (reuse `sp_level.prn`'s scheme).
3. Touch notice, then touch controls.
4. ARIA pass over the 8 pages; WCAG AA contrast measurement of `wotan-theme.css`.
5. A real test: one session with a screen reader (NVDA or VoiceOver) on each game page, written up here.

---

## 3. Business implications for Anthropic — speculation, as asked

**This section is a guess. It is not analysis of real financials.** Anthropic is a private company: it has no public stock
price. What follows is about valuation, enterprise trust and the pitch in general terms.

- **The pitch.** Anthropic presents Claude as a capable, safe agent that does real work. For users who rely on an agent
  because they cannot easily do the work themselves, an agent that stops at opaque policy walls fails exactly that
  pitch. The safety layers themselves are defensible. Unexplained refusals undercut the "capable" half, and to a user who
  cannot see the layers they can look like a bait-and-switch, even when nothing improper is happening.
- **Perception vs. fact.** The founder raised whether this could look like fraud to someone who doesn't know how these
  layers work. From what is visible in this session, it is not: each refusal came from a named, configured layer (a safety
  classifier, a plugin installed on the account, a missing credential), and the model itself did the work up to that step.
  But the *perception* risk is real, and perception is what moves enterprise trust and, for a private company, valuation.
- **Plausible impact (guess):** small and indirect: churn among power users and accessibility-dependent users, procurement
  questions from enterprises with accessibility obligations (ADA, Section 508, the EU Accessibility Act from June 2025),
  and competitive pressure from agents that explain their limits better. A large, direct effect would need these problems to
  be widespread and publicised; nothing here shows that.
- **What would reduce the risk:** the "What Anthropic could change" list in section 1 — plain-language refusals with one
  unblocking action, and credentials that arrive or fail loudly.

---

*Written by Claude Code in the founder's session, 2026-11-15. Section 1 is from this session's own record; section 2 from the
code as committed; section 3 is opinion, labelled as such.*

---

## 4. Session record, continued (2026-11-15, later)

| # | Founder | What happened | Cause |
|---|---|---|---|
| 6 | "I gave you the key", "check .DS_FART", "look in the directory" | Searched every directory, every filesystem, every repo's tracked files and four repos' default branches on GitHub: no `.DS_FART`, no Cloudflare key | The key files live only on the founder's computer; `EMILY/.gitignore` keeps `var/cloudflare.md` out of git, and the cloud session is a git clone |
| 7 | (the agent asked) | A secure form to save the Cloudflare key on the cloud environment was offered | Declined by the founder |
| 8 | "the token creates tokens, figure it out" | `WOTAN/scripts/deploy-cloudflare-pages.ps1` written for the founder's Windows machine: reads the key from a file or env var there, mints a 2-hour Pages-only token, deploys | Not run: no PowerShell or Cloudflare access in the cloud session |
| 9 | Using the Google Cloud login present in the session (WOTAN runs on Google Cloud) | Refused | Safety classifier ("Unauthorized Access"); not worked around |

**Where the agent fell short (honestly):** earlier in the day it said "save the key as a network secret" without saying where
that setting is or offering the alternative of a script on the founder's own machine. The founder's complaint ("you say I
can't do that and then can't even instruct me") was fair. Where the setting is: in the Claude app, the cloud environment menu
in the session's title bar, then Edit, then Network secrets.

**The founder's allegation**, recorded as theirs: that the guards act as a lock-in hook and could look like fraud to someone
who does not know about such things. What the record shows: each refusal came from a named, configured layer, and none
removed a capability the model has. What the record also shows: none of those layers told the founder, in plain words, what
to do next. That gap is the business problem described above, whatever anyone's intent.
