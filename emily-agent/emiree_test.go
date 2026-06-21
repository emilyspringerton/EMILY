package main

import (
	"strings"
	"testing"
)

func TestNewWitchStateDefaults(t *testing.T) {
	w := NewWitchState()
	if w.H != 0.25 {
		t.Errorf("H = %v, want 0.25", w.H)
	}
	if w.P != 0.10 {
		t.Errorf("P = %v, want 0.10", w.P)
	}
	if w.Gear != GearCruise {
		t.Errorf("Gear = %v, want GearCruise", w.Gear)
	}
}

func TestGearString(t *testing.T) {
	cases := []struct {
		g    WitchGear
		want string
	}{
		{GearIdle, "idle"},
		{GearSlow, "slow"},
		{GearCruise, "cruise"},
		{GearActive, "active"},
		{GearDrive, "drive"},
		{GearHighPower, "high-power"},
		{GearOverload, "overload"},
	}
	for _, tc := range cases {
		if got := tc.g.String(); got != tc.want {
			t.Errorf("WitchGear(%d).String() = %q, want %q", tc.g, got, tc.want)
		}
	}
}

func TestSelectGear(t *testing.T) {
	cases := []struct {
		name string
		h, p float64
		want WitchGear
	}{
		{"idle (h<0.15)", 0.10, 0.50, GearIdle},
		{"slow (h<0.25)", 0.20, 0.50, GearSlow},
		{"cruise (h<0.40)", 0.30, 0.50, GearCruise},
		{"active (h<0.55)", 0.50, 0.50, GearActive},
		{"drive (h<0.70)", 0.60, 0.65, GearDrive},
		{"high-power (h<0.85, p>=0.60)", 0.80, 0.70, GearHighPower},
		{"overload (h>=0.85, p>=0.60)", 0.90, 0.80, GearOverload},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWitchState()
			w.H = tc.h
			w.P = tc.p
			w.SelectGear()
			if w.Gear != tc.want {
				t.Errorf("SelectGear(): h=%.2f p=%.2f → gear=%v, want %v", tc.h, tc.p, w.Gear, tc.want)
			}
		})
	}
}

func TestObservePullsToward(t *testing.T) {
	w := NewWitchState()
	w.H = 0.0
	w.P = 1.0
	// Observe toward H=1, P=0 — should move H up and P down.
	w.Observe(1.0, 0.0)
	if w.H <= 0.0 {
		t.Errorf("H should have increased after observe(1,0), got %v", w.H)
	}
	if w.P >= 1.0 {
		t.Errorf("P should have decreased after observe(1,0), got %v", w.P)
	}
}

func TestUpdateIncrementsSteps(t *testing.T) {
	w := NewWitchState()
	before := w.Steps
	w.Update()
	if w.Steps != before+1 {
		t.Errorf("Steps after Update: got %d, want %d", w.Steps, before+1)
	}
}

func TestUpdateKeepsHPInBounds(t *testing.T) {
	w := NewWitchState()
	// Run many updates; H and P must stay in [0, 1].
	for i := 0; i < 200; i++ {
		w.Update()
	}
	if w.H < 0 || w.H > 1 {
		t.Errorf("H out of bounds after 200 updates: %v", w.H)
	}
	if w.P < 0 || w.P > 1 {
		t.Errorf("P out of bounds after 200 updates: %v", w.P)
	}
}

func TestInfluenceGearCruise(t *testing.T) {
	w := NewWitchState()
	w.Gear = GearCruise
	inf := w.Influence()
	if inf.MaxIters != 6 {
		t.Errorf("GearCruise MaxIters = %d, want 6", inf.MaxIters)
	}
	if inf.TempMult != 1.0 {
		t.Errorf("GearCruise TempMult = %v, want 1.0", inf.TempMult)
	}
}

func TestInfluenceOverloadHasHighestIters(t *testing.T) {
	w := NewWitchState()
	var maxIters int
	for g := GearIdle; g <= GearOverload; g++ {
		w.Gear = g
		inf := w.Influence()
		if inf.MaxIters > maxIters {
			maxIters = inf.MaxIters
		}
		if g == GearOverload && inf.MaxIters != maxIters {
			t.Errorf("GearOverload should have max MaxIters")
		}
	}
}

func TestObserveRSIOutcome_Converged(t *testing.T) {
	w := NewWitchState()
	w.Gear = GearCruise
	before := w.Steps
	outcome := RSIOutcome{
		Converged:    true,
		Iterations:   2,
		MaxIterations: 6,
	}
	w.ObserveRSIOutcome(outcome)
	// Should have called Update (incremented Steps).
	if w.Steps <= before {
		t.Errorf("Steps should increment after ObserveRSIOutcome, got %d (was %d)", w.Steps, before)
	}
}

func TestFractalFingerprintDimensions(t *testing.T) {
	w := NewWitchState()
	fp := w.FractalFingerprint(20, 5)
	lines := strings.Split(strings.TrimRight(fp, "\n"), "\n")
	if len(lines) != 5 {
		t.Errorf("FractalFingerprint height: got %d lines, want 5", len(lines))
	}
	for i, line := range lines {
		if len(line) != 20 {
			t.Errorf("FractalFingerprint line %d width: got %d, want 20", i, len(line))
		}
	}
}

func TestFractalFingerprintDefaultsOnZeroSize(t *testing.T) {
	w := NewWitchState()
	// width=0, height=0 → use defaults (40×14)
	fp := w.FractalFingerprint(0, 0)
	lines := strings.Split(strings.TrimRight(fp, "\n"), "\n")
	if len(lines) != 14 {
		t.Errorf("default height: got %d lines, want 14", len(lines))
	}
	for _, line := range lines {
		if len(line) != 40 {
			t.Errorf("default width: line len = %d, want 40", len(line))
		}
	}
}

func TestSummaryContainsGearAndMetrics(t *testing.T) {
	w := NewWitchState()
	s := w.Summary()
	if !strings.Contains(s, "gear=") {
		t.Errorf("Summary missing gear=: %q", s)
	}
	if !strings.Contains(s, "h=") {
		t.Errorf("Summary missing h=: %q", s)
	}
	if !strings.Contains(s, "p=") {
		t.Errorf("Summary missing p=: %q", s)
	}
	if !strings.Contains(s, "steps=") {
		t.Errorf("Summary missing steps=: %q", s)
	}
}
