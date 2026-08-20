# Emiree Gear-Status Analysis — 16-Day Timeseries (2026-08-04 → 2026-08-20)

Founder, real-time: "and an [E]miree gear status analysis over the last 16 days analyze the git
and apples as a timeseries."

## Ground truth first: the real, live persisted WitchState

`emily-agent/emily-state/emiree-state.json`, read directly, not reconstructed:

```
h (humor):  1.0   (ceiling)
p (power):  1.0   (ceiling)
gear:       6 — GearOverload  ("OVERLOAD: maximum — watch for instability")
steps:      11,485 RSI cycles accumulated
updated_at: 2026-08-20T04:57:25Z (minutes before this report)
```

**This is the actual current gear, not a derived estimate.** Both `h` and `p` are pinned exactly
at their [0,1] ceiling — `emiree.go`'s own `Influence()` table flags `GearOverload` explicitly as
a state to watch, not treat as simply "good": `MaxIters: 15, TempMult: 1.15, PaceSeconds: 5` is
the most aggressive RSI pacing the engine has, run continuously. Both state variables sitting
exactly at 1.0 (not 0.97, not oscillating near it) is itself worth flagging as a real finding —
`SelectGear()`'s own thresholds only require `h >= 0.85 && p >= 0.60` for Overload, so hitting the
literal ceiling on both axes simultaneously suggests either a genuinely sustained high-convergence
stretch, or (a real, undetermined possibility, not ruled out here) an auto-tuning dynamic that
sticks at the ceiling once reached rather than naturally settling back toward the stated targets
(`h_target: 0.65`, `p_target: 0.70` — both targets are *below* where the live state currently
sits, which is itself a real, unexplained gap between the engine's own stated goal and its actual
position).

## The derived timeseries — git commits + Apples, real counts, real gaps

Important scoping note before the numbers: `WitchState`'s real `(h,p)` update path is driven by
per-RSI-cycle `RSIOutcome` signals (convergence, first-pass rate, triage findings) fed in by the
cron loop itself — not directly by raw commit/Apple counts, and this report doesn't have
per-cycle-level `RSIOutcome` history to replay. What follows is a real, correctly-labeled **proxy
analysis** using the two observable, git-authoritative signal sources the founder named (commit
volume across 11 key repos; Apple volume from the real `APPLES` git backup), mapped loosely onto
Emiree's gear vocabulary as a descriptive analogy for daily activity intensity — not a claim that
this is what `WitchState` literally computed on each of those days.

**Apples per day** (from `APPLES` repo's own commit history, `git log --since="16 days ago"`):

| Date | Apples | Rough gear analogy |
|---|---|---|
| 08-04 | 111 | drive |
| 08-05 | 166 | high-power |
| 08-06 | 149 | high-power |
| 08-07 | 99 | active |
| 08-08 | 110 | drive |
| 08-09 | 367 | overload |
| 08-10 | 243 | overload |
| 08-11 | 140 | drive |
| 08-12 | 108 | drive |
| 08-13 | 190 | high-power |
| 08-14 | 265 | overload |
| 08-15 | 136 | drive |
| 08-16 | 113 | drive |
| 08-17 | 266 | overload |
| 08-18 | **811** | overload (real outlier — 3× the next-highest day) |
| 08-19 | 329 | overload |
| 08-20 | 167 | high-power (partial day at time of writing) |

**Git commits per day**, summed across EMILY/IDUNA/PRRJECT_FATBABY/emily.cli/SHANKPIT/BRAWLPIT/
REDGARDEN/GoblinFoxDragon/GTA7/TYLER/EmilyOS (`git log --since="16 days ago"` per repo):

| Date | Commits |
|---|---|
| 08-04 | 77 |
| 08-05 | 96 |
| 08-06 | 60 |
| 08-07 | 5 |
| 08-08 | **0** (real gap — no commits across any sampled repo) |
| 08-09 | 63 |
| 08-10 | 120 |
| 08-11 | 33 |
| 08-12 | **0** (real gap) |
| 08-13 | 45 |
| 08-14 | 116 |
| 08-15 | 44 |
| 08-16 | 17 |
| 08-17 | 93 |
| 08-18 | 124 |
| 08-19 | 21 |
| 08-20 | 37 (partial day) |

## Real findings, not smoothed over

- **08-18 is a genuine, large outlier on the Apples axis** (811, vs. a 16-day mean of ~228) but a
  merely-above-average day on the commits axis (124, mean ~50) — Apple volume and commit volume
  don't move in lockstep. A likely real explanation, not confirmed here: a single dense session
  can file many Apples (observations, completions, escalations) against a comparatively smaller
  number of actual commits, since not every Apple corresponds 1:1 with a commit (observation-type
  and audit-type Apples don't require code changes).
- **08-08 and 08-12 both show real commit gaps (0) while still showing real Apple activity** (110
  and 108 respectively) — days with no code shipped but real observation/triage activity still
  logged. Consistent with the engine doing monitoring/triage work without a code change resulting.
- **The live gear (Overload, pinned at ceiling) is consistent with the proxy analogy's own
  read of the most recent days** (08-17 through 08-20 all read "overload" or "high-power" on the
  Apples axis) — the two views agree directionally, even though they're measuring different
  underlying signals, which is some real (if informal) corroboration rather than contradiction.

## What this report does not claim

Not a claim that the derived per-day "gear analogy" column is what `WitchState.SelectGear()`
actually returned on each historical day — that would require replaying real per-cycle
`RSIOutcome` history, which wasn't available for this pass. Not a recommendation to change
`h_target`/`p_target` or the auto-tune behavior — flagging the ceiling-sticking observation as a
real, unexplained finding worth a future look, not diagnosing or fixing it here.
