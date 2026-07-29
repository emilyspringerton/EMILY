slug: knights-of-the-void-twenty-seven-heroes-real-combat
title: Knights of the Void: Twenty-Seven Heroes and a Real Combat Model
author: FATBABY_NEWSWIRE

---

EINHORN_INDUSTRIAL this week grew REDGARDEN — Knights of the Void's playable roster to
twenty-seven and rebuilt the combat model underneath it from the ground up, moving the
multiplayer online battle arena from "a lobby that fights" to a build where auto-attacks,
casts, and territory control all read the way a real MOBA's do.

**Roster: twenty-five heroes to twenty-seven.** MnM, the Shapeshifting Crab (Tank, 26th),
and Weatherman (27th) both joined this window, each shipped with a full kit and dedicated
headless tests. Weatherman arrived paired with Donkey — an equippable item, not a separate
piloted unit, carrying two independent procs: Immortal's Fold triggers automatically below
25% HP for a temporary damage floor and periodic fight-back damage, while Paper Glide is a
player-activated long-range escape (7x speed, flies over obstacles, briefly untargetable) on
a two-minute cooldown. Weatherman's own kit interacts with it directly — his W grounds an
airborne enemy mid-glide, or extends an ally's own glide window. Mana rolled out roster-wide
as a genuine second resource layered on cooldowns, and status effects (stun, slow, silence,
root, and more) are now readable via a text label above every hero's own health bar, closing
a real gap where several of these were silently non-functional in networked play.

**A real combat model: click-to-attack, homing shots, and windup.** Attacking is now a
distinct command from moving, with a persistent target lock and pure-pursuit chase — Gary's
own basic attack became the roster's first real homing ranged shot instead of the flat melee
tick every hero previously shared. A full auto-attack windup/backswing state machine shipped
to exact League of Legends parity: a champion stops to swing, backswing is freely
cancelable, windup is not — the single biggest fidelity jump in how combat actually feels
moment to moment. Gary's own W became Aimed Shot, a true cast-time ability (movement
interrupts the cast, damage does not) with a cast bar visible to every hero on the
battlefield, not just the caster. MnM's own W was reworked mid-window from a passive toggle
into Burrow — a real cast that sends him underground, untargetable, resurfacing with a small
area burst rather than a free stat stack.

**Territory, tuned.** Camera lock shipped for real; a same-window fog-of-war experiment was
built and then deliberately pulled once it became clear client-side-only visibility isn't
real fog of war. Flow income was raised 10x roster-wide after early playtesting found the
economy too slow to ever spend into. Node-capping bots got a real fix — a stale
single-capper anchor — plus a new fractal-boids squad system: bots now split into small
squads that each independently claim a different contested node instead of the whole team
dog-piling one.

**Reliability, found live.** Two critical wire-protocol bugs were caught and fixed this
window. A fixed 2048-byte receive buffer was silently truncating every real snapshot packet
once the wire format grew past it — no bot or human client could actually complete a real
networked match while this was live. Separately, once the snapshot message itself grew to
2460 bytes — over a real network's typical MTU — it was split into several independent,
self-contained packets, so a single dropped packet on a real (non-localhost) connection can
no longer cost the whole snapshot.

**Also shipped:** real cast-radius and zone-ability ground circles, so an area spell finally
shows where it's about to land before it commits; a shop panel bug fix after the newest three
items (Blink Dagger, Donkey, Haste Trinket) turned out to be silently unbuyable behind a
stale hardcoded slot cap; Blink Dagger, a short-range escape item, and Haste Trinket, a
modest 6% cooldown-reduction passive; and an item-stats-plus-suggested-heroes table added for
new players.

Knights of the Void's full twenty-seven-hero roster is queueable now via the existing
bot-pool matchmaker. A ground-up creep and minion-wave overhaul, bringing REDGARDEN's
territory creeps to real League of Legends parity, is next on the public roadmap.

*FATBABY_NEWSWIRE — EINHORN_INDUSTRIAL*
