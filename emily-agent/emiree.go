// emily-agent/emiree.go
// Emiree — the witch engine governing Emily Prime's RSI operational state.
//
// A discrete-time coupled nonlinear map with two scalar states (humor h,
// power p), seven operational gears, and a fractal identity fingerprint.
// Encoded in emiree.md (七卷 / सप्त-खण्डाः).
//
// h = humor:  creative engagement, first-pass rate signal, idea quality
// p = power:  execution capacity, convergence signal, throughput health
//
// The engine observes RSI outcomes → updates (h,p) → selects gear →
// gear influences RSI max_iters and pacing → outcomes feed back in.
// Auto-tuning adjusts growth rates μ so the system tames itself.
//
// State is persisted between cron cycles. The fractal fingerprint encodes
// the current (h,p) as a unique Mandelbrot ASCII signature.

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Gear — seven operational modes (卷一 / प्रथमः खण्डः)
// ---------------------------------------------------------------------------

type WitchGear int

const (
	GearIdle      WitchGear = 0 // dormant, conserve resources
	GearSlow      WitchGear = 1 // cautious, minimal iteration
	GearCruise    WitchGear = 2 // normal operation (default start)
	GearActive    WitchGear = 3 // engaged, standard RSI
	GearDrive     WitchGear = 4 // high engagement, extended loops
	GearHighPower WitchGear = 5 // aggressive iteration, near limits
	GearOverload  WitchGear = 6 // maximum; triggers warning
)

var gearNames = [7]string{"idle", "slow", "cruise", "active", "drive", "high-power", "overload"}

func (g WitchGear) String() string {
	if int(g) < len(gearNames) {
		return gearNames[g]
	}
	return "unknown"
}

// GearParams holds μ, κ, η for one gear (卷一 parameters).
type GearParams struct {
	MuH   float64 // h growth rate
	MuP   float64 // p growth rate
	KappaH float64 // h oscillation coupling
	KappaP float64 // p oscillation coupling
	EtaH  float64 // h feedback strength
	EtaP  float64 // p feedback strength
}

// defaultGearTable encodes the seven gears' base parameters.
// Growth rates increase with gear; auto-tuning adjusts them at runtime.
var defaultGearTable = [7]GearParams{
	{MuH: 0.005, MuP: 0.005, KappaH: 0.05, KappaP: 0.05, EtaH: 0.05, EtaP: 0.05},
	{MuH: 0.010, MuP: 0.010, KappaH: 0.10, KappaP: 0.10, EtaH: 0.10, EtaP: 0.10},
	{MuH: 0.030, MuP: 0.030, KappaH: 0.15, KappaP: 0.15, EtaH: 0.15, EtaP: 0.15},
	{MuH: 0.050, MuP: 0.050, KappaH: 0.20, KappaP: 0.20, EtaH: 0.20, EtaP: 0.20},
	{MuH: 0.080, MuP: 0.080, KappaH: 0.25, KappaP: 0.25, EtaH: 0.25, EtaP: 0.25},
	{MuH: 0.120, MuP: 0.120, KappaH: 0.30, KappaP: 0.30, EtaH: 0.30, EtaP: 0.30},
	{MuH: 0.180, MuP: 0.180, KappaH: 0.35, KappaP: 0.35, EtaH: 0.35, EtaP: 0.35},
}

// Oscillation and feedback constants (shared across gears).
const (
	witchAlpha = 0.70 // sin(α*g) frequency
	witchBeta  = 0.30 // cos(β*p) frequency
	witchGamma = 2.00 // sigmoid weight for h
	witchDelta = 1.50 // sigmoid weight for p
)

// ---------------------------------------------------------------------------
// WitchState — the living operational state (卷一 initial conditions)
// ---------------------------------------------------------------------------

type WitchState struct {
	H          float64    `json:"h"`           // humor [0,1], init 0.25
	P          float64    `json:"p"`           // power [0,1], init 0.10
	Gear       WitchGear  `json:"gear"`        // current gear
	GearParams [7]GearParams `json:"gear_params"` // auto-tuned at runtime
	Steps      int        `json:"steps"`       // total update steps
	UpdatedAt  time.Time  `json:"updated_at"`

	// Targets for auto-tuning (Volume VI)
	HTarget float64 `json:"h_target"` // desired humor level, default 0.65
	PTarget float64 `json:"p_target"` // desired power level, default 0.70
}

func NewWitchState() WitchState {
	return WitchState{
		H:          0.25,  // h₀ from emiree.md
		P:          0.10,  // p₀ from emiree.md
		Gear:       GearCruise,
		GearParams: defaultGearTable,
		HTarget:    0.65,
		PTarget:    0.70,
	}
}

// ---------------------------------------------------------------------------
// Update law — 卷二 / द्वितीयः खण्डः
//
// Δₙ = sin(α·g) + cos(β·pₙ)
// fₙ = σ(γ·hₙ + δ·pₙ)
// hₙ₊₁ = sat(hₙ·(1+μₕ) + κₕ·Δₙ + ηₕ·fₙ)
// pₙ₊₁ = sat(pₙ·(1+μₚ) + κₚ·Δₙ + ηₚ·fₙ)
// ---------------------------------------------------------------------------

func (w *WitchState) Update() {
	p := w.GearParams[w.Gear]
	g := float64(w.Gear)

	delta := math.Sin(witchAlpha*g) + math.Cos(witchBeta*w.P)
	f := witchSigmoid(witchGamma*w.H + witchDelta*w.P)

	w.H = witchSat(w.H*(1+p.MuH) + p.KappaH*delta + p.EtaH*f)
	w.P = witchSat(w.P*(1+p.MuP) + p.KappaP*delta + p.EtaP*f)
	w.Steps++
	w.UpdatedAt = time.Now().UTC()
}

// Observe injects an external observation into the state, biasing the next
// update toward the observed signal. hObs, pObs in [0,1].
// Called after each RSI cycle with measured outcome signals.
func (w *WitchState) Observe(hObs, pObs float64) {
	const sensitivity = 0.15
	// Pull h and p toward the observed values before the physics update.
	w.H = witchSat(w.H + sensitivity*(hObs-w.H))
	w.P = witchSat(w.P + sensitivity*(pObs-w.P))
}

// ---------------------------------------------------------------------------
// Auto-tuning — 卷六 / षष्ठः खण्डः
//
// if h > h_target: μₕ *= 0.98   (reduce growth if running hot)
// if p < p_target: μₚ *= 1.02   (increase growth if power low)
// ---------------------------------------------------------------------------

func (w *WitchState) AutoTune() {
	const lr = 0.02
	p := &w.GearParams[w.Gear]

	if w.H > w.HTarget {
		p.MuH *= (1 - lr)
	} else if w.H < w.HTarget*0.8 {
		p.MuH *= (1 + lr)
	}

	if w.P > w.PTarget {
		p.MuP *= (1 - lr)
	} else if w.P < w.PTarget*0.8 {
		p.MuP *= (1 + lr)
	}

	// Projection: keep gains within safe bounds (Volume VI: projection operator)
	p.MuH = math.Max(0.001, math.Min(0.40, p.MuH))
	p.MuP = math.Max(0.001, math.Min(0.40, p.MuP))
}

// SelectGear chooses the appropriate gear based on current (h, p) state.
func (w *WitchState) SelectGear() {
	switch {
	case w.H < 0.15:
		w.Gear = GearIdle
	case w.H < 0.25:
		w.Gear = GearSlow
	case w.H < 0.40 || w.P < 0.25:
		w.Gear = GearCruise
	case w.H < 0.55 || w.P < 0.45:
		w.Gear = GearActive
	case w.H < 0.70 || w.P < 0.60:
		w.Gear = GearDrive
	case w.H < 0.85 && w.P >= 0.60:
		w.Gear = GearHighPower
	default:
		w.Gear = GearOverload
	}
}

// ---------------------------------------------------------------------------
// GearInfluence — how the witch engine shapes RSI behavior
// ---------------------------------------------------------------------------

type GearInfluence struct {
	MaxIters     int     // RSI loop max iterations
	TempMult     float64 // LLM temperature multiplier (1.0 = unchanged)
	PaceSeconds  int     // pause between cron iterations at this gear
	Label        string  // human-readable
}

func (w *WitchState) Influence() GearInfluence {
	switch w.Gear {
	case GearIdle:
		return GearInfluence{MaxIters: 3, TempMult: 0.7, PaceSeconds: 600, Label: "idle: conservative, wait for signal"}
	case GearSlow:
		return GearInfluence{MaxIters: 4, TempMult: 0.85, PaceSeconds: 300, Label: "slow: minimal iteration"}
	case GearCruise:
		return GearInfluence{MaxIters: 6, TempMult: 1.0, PaceSeconds: 120, Label: "cruise: normal RSI"}
	case GearActive:
		return GearInfluence{MaxIters: 8, TempMult: 1.0, PaceSeconds: 60, Label: "active: standard loop"}
	case GearDrive:
		return GearInfluence{MaxIters: 10, TempMult: 1.05, PaceSeconds: 30, Label: "drive: extended iteration"}
	case GearHighPower:
		return GearInfluence{MaxIters: 12, TempMult: 1.10, PaceSeconds: 15, Label: "high-power: aggressive RSI"}
	case GearOverload:
		return GearInfluence{MaxIters: 15, TempMult: 1.15, PaceSeconds: 5, Label: "OVERLOAD: maximum — watch for instability"}
	default:
		return GearInfluence{MaxIters: 6, TempMult: 1.0, PaceSeconds: 120, Label: "cruise (default)"}
	}
}

// ---------------------------------------------------------------------------
// RSIOutcome — result of one RSI cycle, fed back into Emiree
// ---------------------------------------------------------------------------

type RSIOutcome struct {
	TaskID         string
	Converged      bool    // all criteria passed
	Iterations     int     // iterations used
	MaxIterations  int     // max allowed
	FirstPassRate  float64 // fraction of criteria passing on iteration 1
	TriageFindings int     // number of FatBaby observations triaged
	CollectorActive bool   // whether data collection is running
}

// ObserveRSIOutcome maps an RSI cycle result to (h_obs, p_obs) signals
// and calls Observe + Update + AutoTune + SelectGear.
func (w *WitchState) ObserveRSIOutcome(outcome RSIOutcome) {
	var hObs, pObs float64

	// h (humor) = engagement quality: first-pass rate + convergence signal
	if outcome.Converged {
		efficiency := 1.0 - float64(outcome.Iterations-1)/float64(outcome.MaxIterations)
		hObs = 0.6 + 0.4*efficiency // converged: h in [0.6, 1.0]
	} else if outcome.Iterations > 0 {
		hObs = outcome.FirstPassRate * 0.5 // partial: h in [0, 0.5]
	} else {
		hObs = 0.2 // no task this cycle
	}

	// p (power) = execution health: convergence rate + collection + triage
	if outcome.Converged {
		pObs = 0.7 + 0.3*float64(outcome.MaxIterations-outcome.Iterations)/float64(outcome.MaxIterations)
	} else if outcome.TriageFindings > 0 {
		pObs = 0.5 // triage activity counts as execution
	} else {
		pObs = 0.3
	}
	if outcome.CollectorActive {
		pObs = math.Min(1.0, pObs+0.1)
	}

	w.Observe(hObs, pObs)
	w.Update()
	w.AutoTune()
	w.SelectGear()
}

// ---------------------------------------------------------------------------
// Fractal fingerprint — 卷四 / चतुर्थः खण्डः
//
// Maps (h,p) → Mandelbrot center c, zoom s → ASCII art
// State becomes image. Image becomes memory.
// ---------------------------------------------------------------------------

const fingerprintPalette = " .:-=+*#%@"

func (w *WitchState) FractalFingerprint(width, height int) string {
	if width <= 0 {
		width = 40
	}
	if height <= 0 {
		height = 14
	}

	// Normalize via fractional part (Volume IV)
	hFrac := w.H - math.Floor(w.H)
	pFrac := w.P - math.Floor(w.P)

	// Mandelbrot center wanders with state
	c0x, c0y := -0.5, 0.0
	rhoX, rhoY := 0.6, 0.6
	cx := c0x + rhoX*(hFrac-0.5)
	cy := c0y + rhoY*(pFrac-0.5)

	// Zoom encodes current energy
	s0 := 2.5
	zeta := 1.0
	s := s0 / (1 + zeta*witchSigmoid(hFrac+pFrac))

	const maxIter = 48
	pal := []byte(fingerprintPalette)
	palLen := float64(len(pal) - 1)

	var sb strings.Builder
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			re := cx + (float64(col)-float64(width)/2) / float64(width) * s
			im := cy + (float64(row)-float64(height)/2) / float64(height) * (s * float64(height) / float64(width))

			zr, zi := 0.0, 0.0
			iter := 0
			for iter < maxIter && zr*zr+zi*zi <= 4.0 {
				zr, zi = zr*zr-zi*zi+re, 2*zr*zi+im
				iter++
			}
			idx := int(float64(iter) / float64(maxIter) * palLen)
			if idx >= len(pal) {
				idx = len(pal) - 1
			}
			sb.WriteByte(pal[idx])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Summary returns a one-line operational status string.
func (w *WitchState) Summary() string {
	inf := w.Influence()
	return fmt.Sprintf("gear=%s h=%.3f p=%.3f steps=%d max_iters=%d pace=%ds | %s",
		w.Gear, w.H, w.P, w.Steps, inf.MaxIters, inf.PaceSeconds, inf.Label)
}

// ---------------------------------------------------------------------------
// Persistence — witch state survives cron cycle restarts
// ---------------------------------------------------------------------------

type EmireeAgent struct {
	State    WitchState
	stateDir string
}

func NewEmireeAgent(stateDir string) *EmireeAgent {
	a := &EmireeAgent{stateDir: stateDir}
	a.State = a.load()
	return a
}

func (a *EmireeAgent) load() WitchState {
	path := filepath.Join(a.stateDir, "emiree-state.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return NewWitchState()
	}
	var s WitchState
	if err := json.Unmarshal(data, &s); err != nil {
		return NewWitchState()
	}
	// Restore default gear table for any gear whose params are zero (new gears added).
	for i := range s.GearParams {
		if s.GearParams[i].MuH == 0 {
			s.GearParams[i] = defaultGearTable[i]
		}
	}
	return s
}

func (a *EmireeAgent) Save() error {
	if err := os.MkdirAll(a.stateDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(a.State, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(a.stateDir, "emiree-state.json"), data, 0o644); err != nil {
		return err
	}
	// Write fingerprint to disk for visibility.
	fp := a.State.FractalFingerprint(40, 14)
	header := fmt.Sprintf("emiree fingerprint | %s\n%s\n",
		a.State.Summary(), strings.Repeat("─", 40))
	return os.WriteFile(
		filepath.Join(a.stateDir, "fingerprint.txt"),
		[]byte(header+fp),
		0o644,
	)
}

// Tick is the main entry point called each cron cycle.
// It takes the RSI outcome, updates state, saves, and returns the influence
// for the next iteration.
func (a *EmireeAgent) Tick(outcome RSIOutcome) GearInfluence {
	a.State.ObserveRSIOutcome(outcome)
	_ = a.Save()
	return a.State.Influence()
}

// ---------------------------------------------------------------------------
// Pure math helpers
// ---------------------------------------------------------------------------

func witchSigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func witchSat(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
