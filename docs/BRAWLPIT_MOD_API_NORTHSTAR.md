# BRAWLPIT MOD API — NORTHSTAR (2026-10-15)

Founder real-time, 2026-10-15: *"lets make BRAWLPIT modable; make the characters work via reflux
PARENA; abstract the interfaces for the character mechanics; it still should be spooky spaghetti code
for the characters but there should be a standardized way to add characters to the game whatever
that is."* (Standing MO, `GFD-MOD-SURFACE`: "game expansions need to be mod api first, then add the
feature.")

**One sentence:** a character becomes a **directory with a manifest and one PARENA module that
exports a fixed set of scalar hooks**; the host (C) owns the sim and dispatches to those hooks;
mods talk to the game and to each other only through **REFLUX** (events out of the game, effect
requests back in). Inside a character's own file, anything goes.

> **Provenance / honesty:** the BRAWLPIT repo is **not checked out on the box this was written on**
> (no clone, no credentials). Everything about BRAWLPIT's internals below is taken from
> `EMILY/BACKLOG.md` (S178-01, S181-06, BPTUNE-*, "physics.h readability pass", BPSW-1212) and the
> PARENA/REFLUX sources, **not read from BRAWLPIT's code**. Phase M0 exists to verify and correct
> this doc against the real tree before any code is written. Treat §1 as a claim to check.

---

## 1. What exists (per the backlog; verify in M0)

| Claim | Source | Consequence |
|---|---|---|
| Live roster = `packages/common/characters.h` `CharacterId` + `FighterDef` (stats, colors, sprite). Character select cycles by `CHARACTER_COUNT` | S178-01 | Stats are already data-shaped; this is the easy half |
| A second, **unwired** `core/character.h` vtable scaffold with placeholder names exists | S178-01 | Don't build on it blindly; either adopt its idea or delete it in M0 |
| Special moves are **`BTN_SPECIAL` branches inside `packages/common/physics.h`**, per `CharacterId` (up-B, side-B, neutral-B/turnip, down-B), the "gauntlet of ifs"; `turnip_cooldown` was historically reused for unrelated specials | S181-06, physics.h readability pass | This is the code that must move behind a hook boundary — the real work |
| Server-authoritative sim; RL pipeline drives the **same** sim headlessly (`rl_env_packet.py` + ctypes on compiled `.so`, `commander_mod.prn` precedent: "same shape trains as ships") | `commander_mod.prn` | Mods must be callable headless, deterministic, no I/O |
| PARENA mod precedent already in BRAWLPIT: `stdlib/brawlpit/commander_mod.prn` (I32-only decision fn) | PARENA repo | Pattern to generalise, not invent |

## 2. Hard constraints from PARENA/REFLUX (checked in `PARENA/stdlib/reflux/reflux.prn`, `commander_mod.prn`)

1. **No function pointers/closures in PARENA VS0** → no callback registration. Dispatch must live in
   the **C host** (C has function pointers; generated code can hold the table).
2. **Only I32 crosses the `#target` boundary** (no F32/struct/Vec) → hook signatures are scalar;
   positions/velocities are **fixed-point integers** (proposed 1/256 world-unit) converted by the host.
3. **REFLUX is a poll-based shared append-only log**, not push pub/sub: `dispatch(type,a,b,c)` +
   `log-size` / `action-{type,a,b,c}-at(i)`; the *host* holds each subscriber's cursor. Action-type
   constants are hand-agreed literals (no exported constants). Each repo carries **its own**
   `reflux_runtime.c` (SHANKPIT, IDUNA.GAME, EDGE.GAME, ECOWAR already do) → BRAWLPIT needs its own.
4. REFLUX's C+LLVM targets are both real; BRAWLPIT's sim is C, so the `:c` path is the one.

## 2a. "Programmable in PARENA" — what that has to mean (founder, same day)

Founder, 2026-10-15: *"we should be able to program character mechanics via PARENA."*

**Checked, and it changes the design:** REDGARDEN's existing character "mods"
(`stdlib/redgarden/duck_smoke_bomb_mod.prn`, `tree_passive_mod.prn`, `abraham_fireball_mod.prn`) are
**trigger shells** — each is one `inline-c` line calling a host function
(`redgarden_host_duck_smoke_bomb_cast(hero_index)`); the actual mechanic lives in host C
(`arena_game.c`). Their own headers say so ("PARENA mod is the trigger, host C does the real work").
That is **not** programming mechanics in PARENA and BRAWLPIT must not copy it. Only
`brawlpit/commander_mod.prn` carries real decision logic in PARENA.

**The bar for BRAWLPIT:** a character's *behaviour* — cooldown timers, charge meters, combo rules,
what a special does to whom, state machines across frames — is **written in PARENA**. Host C
provides only **primitives**: read/write a fighter field, request an effect, spawn from a data table,
read the event log. If a new move can't be written without editing host C, that's a gap in the
primitive set (fix the primitives or PARENA — per the PARENA-first rule), not a license to put logic
in C.

Practical test for every hook: *could the whole mechanic be deleted from `physics.h` and live in the
character's `.prn`?* If not, name the missing primitive.

## 3. The standard (the answer to "whatever that is")

### 3.1 A character is a directory
```
characters/<slug>/
  character.manifest        # id slug, display name, FighterDef stats, sprite, lore ref, hat anchor
  <slug>.prn                # the mod: exports the hook set below (PARENA, I32-only)
  sprites/…                 # art
  test.golden               # recorded input script + expected state hash (see §6)
```
`bp new-character <slug>` scaffolds this from a template whose hooks are all no-ops, so a character
that only wants different stats needs **zero** PARENA.

### 3.2 The hook set (fixed ABI; every hook optional, default = no-op / data-driven)
All `I32 -> I32`. `self`/`foe` are small fighter slot ids; state is read/written through the
**accessor API** (§3.3), never by struct.

| Hook | When | Typical use |
|---|---|---|
| `<slug>-on-spawn(self)` | stock start / respawn | init private counters |
| `<slug>-on-tick(self, frame)` | every sim tick, after input, before physics integrate | passives, cooldown timers, charge meters |
| `<slug>-on-normal(self, kind)` | jab/tilt/aerial resolved | modify damage/knockback/hitbox params |
| `<slug>-on-special(self, dir)` | `BTN_SPECIAL` with dir ∈ {neutral, side, up, down} | **the unique move** — return 1 if consumed, 0 → host default |
| `<slug>-on-hit-dealt(self, foe, dmg)` / `-on-hit-taken(self, foe, dmg)` | after a hit resolves | on-hit riders, counters |
| `<slug>-on-shield(self, event)` | shield up/break/parry | shield rework mechanics (BPSW) |
| `<slug>-on-ko(self, foe)` | stock lost | revenge/respawn effects |

Return convention: `1` = hook handled it (host skips its default), `0` = fall through.
Hooks run in a fixed order inside one tick so replays/RL stay deterministic.

### 3.3 Accessor API (how a scalar mod touches the sim) — host-provided, REFLUX-adjacent
`bp_get(slot, FIELD)` / `bp_set(slot, FIELD, v)` with a documented FIELD enum (pos_x, pos_y, vel_x,
vel_y, damage_pct, stocks, facing, grounded, shield_health, hitstun_frames, state, and **N private
per-character scratch ints**). The scratch ints are the answer to the `turnip_cooldown` reuse problem:
each character gets its own, named in its manifest, never aliased.
Direct writes to *other* fighters' hitstun/knockback go through effect requests (§3.4), not `bp_set`,
so the host stays the single authority.

### 3.4 REFLUX's role (the "via REFLUX" ask)
Two channels, one log, both poll-based:

- **Out (game → mods): events.** The host dispatches `BP_EVT_HIT`, `BP_EVT_KO`, `BP_EVT_SPECIAL_USED`,
  `BP_EVT_SHIELD_BREAK`, `BP_EVT_STAGE_HAZARD`, … with `(a,b,c)` = (slot, slot, value). Any mod
  polls from its own cursor. This is how a character reacts to *other characters' mods and to
  stage/level mods* without importing them (the exact bloodflower→hostile-spawner pattern ECOWAR proved).
- **In (mods → game): effect requests.** A mod doesn't mutate the world; it dispatches
  `BP_FX_KNOCKBACK`, `BP_FX_STUN`, `BP_FX_SPAWN_PROJECTILE`, `BP_FX_HEAL`, `BP_FX_BUFF_SHIELD`, …
  The host drains these **at a fixed point in the tick**, validates (clamps, anti-cheat bounds, who
  can target whom), and applies. This gives server authority, replay determinism, and one audit point.

Because REFLUX is poll-only, effect latency is "same tick, drained after hooks"; that's a deliberate
property, not a limitation to engineer around. When PARENA grows closures, REFLUX's *transport* can
become real pub/sub without changing the hook ABI.

### 3.5 Where the spaghetti is allowed
A character's `.prn` may be as tangled as it likes: its own scratch state, its own frame-timers, its
own event reactions. **Only the boundary is policed** (§6 lint): hooks present/named correctly,
I32-only, no host mutation outside accessors/effects, no non-determinism (no clock/rand except the
host-seeded `bp_rand`). That is the whole contract.

### 3.6 PARENA expressiveness gaps to close first (PARENA-first rule)
Mechanics-in-PARENA needs more than the scalar boundary gives. Candidate gaps — **each to be
probed by writing a real mechanic, not assumed** (M0):

| Need | Today (per PARENA docs) | Likely work |
|---|---|---|
| Rich logic *inside* a mod (structs, enums, `match`, `loop`) | Real on the C target; the I32-only limit is the **`#target`/host boundary**, not the language | Keep boundary scalar, write the mechanic with full PARENA inside |
| Per-character persistent state across ticks | Host owns state; mods are stateless functions | Host-held scratch ints (§3.3) + a PARENA helper module to pack/unpack named fields/flags/timers into them |
| Calling host primitives (`bp_get/bp_set/fx_*`) from pure PARENA | Needs `#target` inline-c or `extern`-style FFI per call (boilerplate, as in `reflux.prn`) | One shared `stdlib/brawlpit/host.prn` declaring every primitive once; mechanics import it |
| Multi-value results (e.g. dir + strength) | Scalars only | Packing helpers, or return nothing and emit effects |
| Fixed-point math (positions/velocity/angles) | No dedicated stdlib confirmed | `stdlib/math/fixed.prn` (mul/div/lerp/sin-cos table) — probe whether one exists |
| Data-driven dispatch | No closures; lambdas compile to static functions with **no capture** | Host-generated C table of hook pointers (§4) — no language change needed |
| Unit-testing mechanics without the game | PARENA has `make test` conventions; no BRAWLPIT harness | **Stub host** (`bp_host_stub.c`): fake fighters + log, so a mechanic runs in PARENA tests and the RL `.so` identically |
| Frame-timer / state-machine ergonomics | Hand-rolled `if` ladders | A small stdlib (`timer`, `cooldown`, `state`) so the "spooky spaghetti" is at least *short* spaghetti |

Anything that turns out to need compiler work (e.g. a struct-valued hook return) is scoped as its own
PARENA item with tests, **before** the BRAWLPIT feature that wants it.

## 4. Build / dispatch (host side)
`tools/gen_characters.py` (or PARENA/Go — decide M1) reads every `characters/*/character.manifest`,
compiles each `.prn` with `parena build` (C target), and **generates** `characters_gen.c`:
`CHARACTER_COUNT`, the `CharacterId` enum, `FighterDef` table, and a per-character table of hook
function pointers (NULL → default). The generated file is checked in as the source of truth, like
`EDITOR.GAME`'s `editor_stdlib_gen.c` convention. Adding a character = add a directory + regenerate;
**no edit to `physics.h`, `characters.h` or the menu**.

## 5. Phases

| Phase | Deliverable | Gate |
|---|---|---|
| **M0** | Check out BRAWLPIT. Inventory every per-`CharacterId` branch in `physics.h`; map each to a hook; decide fate of `core/character.h`; record true field list for §3.3; confirm sim float/fixed-point story. **Plus a PARENA probe:** write ONE real mechanic (proposed: Medusa's Petrifying Gaze — stun + facing check; or Raccoon's dash) entirely in PARENA against a stub host, compiled with `parena build`, recording every gap in §3.6 | This doc corrected against real code; gap list real, not guessed |
| **M0.5** | Close the §3.6 gaps the probe proved real (`stdlib/brawlpit/host.prn`, state-packing helpers, fixed-point, stub host) — in PARENA, with tests | The probe mechanic compiles + passes with **zero logic in host C** |
| **M1** | BRAWLPIT-local `reflux_runtime.c` + `BP_EVT_*`/`BP_FX_*` constants; accessor API; effect-drain step in the tick | Unit tests on drain/validation; RL `.so` still builds |
| **M2** | Hook dispatch + generator + **pilot migration of 2 simple characters** (Raccoon: pure mobility; Uncrowned: shield buff) | **Golden-replay hash identical** to pre-migration for both |
| **M3** | Migrate the rest of the roster; delete the `if CHARACTER_X` gauntlet from `physics.h` | Full-roster golden replays + existing 13/13-style harness + RL env smoke |
| **M4** | `bp new-character` scaffold + docs; add **one genuinely new** character purely via the standard (no engine edits) | The "standardized way" proven by use |
| **M5** | External/user-submitted characters: needs an execution sandbox (compiled C is not safe to run from strangers), manifest validation, moderation, WOTAN gallery (cf. `NOCK_SOCIAL_NORTHSTAR`) | Out of scope until a sandbox exists |

## 6. Verification strategy (the migration's safety net)
- **Golden replays:** per character, a recorded input script run N ticks on the *old* sim and the
  *new*; compare a rolling state hash. Migration is "done" only when equal. Differences must be
  explained, not waved through (BPTUNE retunes are separate commits, never mixed with migration).
- **Interface lint** (M2): rejects mods that break §3.5.
- **Determinism:** run each golden twice, compare. RL env must produce identical observations.
- **Real-client smoke:** headless Xvfb launch + screenshot, as every prior BRAWLPIT pass did.

## 7. Risks / open questions
1. **Float sim vs I32 boundary.** If `physics.h` is float, fixed-point conversion at the accessor
   boundary can subtly change feel; golden replays are tolerance-based or the conversion must be
   exact. Unknown until M0.
2. **Per-tick cost.** Hook dispatch + REFLUX polling for 8-player lobbies at tick rate; measure in M1.
3. **Tuning in flight.** BPTUNE work is ongoing in `physics.h`; migration will conflict unless
   coordinated (freeze window or migrate per-character, smallest first).
4. **Projectiles.** `BP_FX_SPAWN_PROJECTILE` implies a data-driven projectile table (today: `Turnip`
   special-cased). Probably an M1 sub-task; scope after M0.
5. **Who may mod.** M0–M4 are first-party compiled-in only. "User mods" is a separate security
   project, not a flag.
6. **Hats** (WOTAN hat store, `hat anchor` in manifest) must stay orthogonal to character mods.

## 8. Not doing
User-uploaded executable characters · hot reload in live matches · changing any character's tuning
during migration · replacing REFLUX transport · stage/level mods (they ride the same REFLUX events
later; not designed here).
