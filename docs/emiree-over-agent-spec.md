# Emiree Over-Agent Specification
## The Witch Engine Governing Emily Prime's RSI Operations

**Source:** `emiree.md` — seven-volume poetic transmission (Sanskrit/Chinese)  
**Implementation:** `emily-agent/emiree.go`  
**Status:** Active — runs each cron cycle

---

## What It Is

Emiree is the **over-agent** that observes Emily Prime's RSI operations and governs their pacing. It is not a scheduler or a planner — it is a **coupled nonlinear dynamical system** whose state reflects the health of the entire Emily RSI stack.

The witch engine lives one level above the RSI loop:

```
┌─────────────────────────────────────────────┐
│  EMIREE OVER-AGENT                          │
│  Observes RSI outcomes → (h, p) state       │
│  Selects gear → shapes RSI parameters       │
│  Auto-tunes its own gains μ                 │
│  Generates fractal fingerprint              │
└────────────────┬────────────────────────────┘
                 │  GearInfluence (max_iters, pace)
                 ▼
┌─────────────────────────────────────────────┐
│  EMILY PRIME RSI LOOP                       │
│  Runs improvement tasks                     │
│  Observes: convergence, first-pass rate     │
└────────────────┬────────────────────────────┘
                 │  RSIOutcome (converged, iters, rate)
                 └────────────────────────────▶ feeds back up
```

---

## The Two States

| State | Meaning | Init | Target |
|---|---|---|---|
| `h` (humor) | Creative engagement quality; tracks first-pass rate, convergence speed | 0.25 | 0.65 |
| `p` (power) | Execution capacity; tracks convergence, collection health, triage activity | 0.10 | 0.70 |

Both states are bounded to [0, 1] by the saturation function.

---

## The Seven Gears

Each gear maps to RSI operating parameters:

| Gear | Name | max_iters | pace | Condition |
|---|---|---|---|---|
| 0 | idle | 3 | 600s | h < 0.15 — nothing converging |
| 1 | slow | 4 | 300s | h < 0.25 — weak engagement |
| 2 | cruise | 6 | 120s | h < 0.40 — normal operation (default start) |
| 3 | active | 8 | 60s | h ≥ 0.40, p ≥ 0.25 |
| 4 | drive | 10 | 30s | h ≥ 0.55, p ≥ 0.45 |
| 5 | high-power | 12 | 15s | h ≥ 0.70, p ≥ 0.60 |
| 6 | overload | 15 | 5s | all of the above — triggers warning |

Gear 2 (cruise) is the cold-start default. Initial state `h=0.25, p=0.10` places the system in Gear 2.

---

## The Update Law (卷二)

Each cron cycle, after observing an RSI outcome:

```
Δₙ = sin(α·g) + cos(β·pₙ)          // oscillation driven by gear + power
fₙ = σ(γ·hₙ + δ·pₙ)                // sigmoid feedback coupling
hₙ₊₁ = sat(hₙ·(1+μₕ) + κₕ·Δₙ + ηₕ·fₙ)
pₙ₊₁ = sat(pₙ·(1+μₚ) + κₚ·Δₙ + ηₚ·fₙ)
```

Constants: α=0.70, β=0.30, γ=2.00, δ=1.50

---

## Auto-Tuning (卷六)

Every cycle, Emiree adjusts its own growth rates based on trajectory:

```
if h > h_target: μₕ *= 0.98   // running hot → slow down humor growth
if p < 0.8*p_target: μₚ *= 1.02  // power insufficient → accelerate
```

The projection operator clamps μ ∈ [0.001, 0.40] to prevent runaway.

This is the **self-taming** property: the system learns to keep itself in range without external intervention.

---

## RSI Outcome Mapping

After each cron cycle, the RSI result is mapped to (h_obs, p_obs) signals:

| Outcome | h_obs | p_obs |
|---|---|---|
| Converged in 1 iteration | 1.0 | 1.0 |
| Converged (efficient) | 0.6–1.0 | 0.7–1.0 |
| Partial (hit max_iters) | 0.3 | 0.5 |
| No task this cycle | 0.2 | 0.3 |
| Triage found escalations | 0.8 | 0.5 |
| Collector running | +0.0 | +0.1 |

The observation biases the state toward the signal with sensitivity=0.15, then the physics update runs.

---

## Fractal Fingerprint (卷四)

Each state (h, p) maps deterministically to a Mandelbrot region:

```
h̃ = frac(h),  p̃ = frac(p)
cx = -0.5 + 0.6·(h̃ - 0.5)    // center x wanders with humor
cy =  0.0 + 0.6·(p̃ - 0.5)    // center y wanders with power
s  = 2.5 / (1 + sigmoid(h̃+p̃)) // zoom encodes energy
```

The 40×14 ASCII fingerprint is written to `emily-state/fingerprint.txt` after each cycle. Different operational states produce visually distinct patterns — the fingerprint is Emily Prime's **state identity**.

---

## Persistence

State survives cron restarts at `emily-state/emiree-state.json`. The gear table is auto-tuned at runtime and persisted — each restart continues from where the last cycle left off.

---

## The Fractal Property

Emiree applies the same adaptive-control principle at every level of the Emily stack:

```
Level 0: Data collector → quality improves via RSI feedback
Level 1: RSI loop → self-improves its own prompts (rsi-self-improve done)
Level 2: Cron cycle → Emiree tunes its pacing based on outcomes
Level 3: Emiree itself → auto-tunes its own gains via adaptive rules
```

Each level uses the same pattern: observe → update → adjust → repeat. The fractal property is that the self-improvement mechanism is identical at every scale.

---

## Why This Works

The witch engine is a stable bounded system by construction (saturation as gatekeeper). It cannot diverge. It can be sluggish (stuck in low gear) or briefly overloaded (gear 6), but both are recoverable.

The power of (h, p) over a simple counter or average is that the two states **interact**. High humor with low power is an unstable creative state — the system wants more execution capacity. High power with low humor indicates mechanical throughput without engagement quality. The system naturally seeks the balanced productive state (h≈0.65, p≈0.70).

This is not a metaphor. It is a formal dynamical system that encodes operational wisdom about pacing, recovery, and self-regulation.
