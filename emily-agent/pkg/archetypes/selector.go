package archetypes

import (
	"math"
	"strings"
)

// intentKeywords maps common intent keywords to spirit ranks.
// Keyword matching is a fast tier-1 router before embedding-based similarity.
var intentKeywords = map[string][]int{
	// Strategic / planning
	"strategy":      {15, 21, 35},  // Eligos, Marax, Marchosias
	"plan":          {15, 3, 17},   // Eligos, Vassago, Botis
	"forecast":      {3, 17, 29},   // Vassago, Botis, Astaroth
	"foresight":     {3, 15, 29},   // Vassago, Eligos, Astaroth
	"analysis":      {21, 50, 15},  // Marax, Furcas, Eligos
	"intelligence":  {3, 15, 21},   // Vassago, Eligos, Marax
	"research":      {36, 21, 54},  // Stolas, Marax, Murmur
	// Creative / synthesis
	"creative":      {37, 22, 60},  // Phenex, Ipos, Vapula
	"synthesis":     {7, 11, 26},   // Amon, Gusion, Bune
	"innovation":    {39, 60, 5},   // Malphas, Vapula, Marbas
	"writing":       {24, 37, 2},   // Naberius, Phenex, Agares
	"poetry":        {37, 53, 67},  // Phenex, Caim, Amdusias
	"language":      {2, 24, 53},   // Agares, Naberius, Caim
	// Financial / wealth
	"wealth":        {26, 20, 56},  // Bune, Purson, Gremory
	"finance":       {26, 71, 21},  // Bune, Dantalion, Marax
	"market":        {26, 3, 15},   // Bune, Vassago, Eligos
	"signal":        {3, 15, 23},   // Vassago, Eligos, Aim
	"governance":    {55, 11, 17},  // Orobas, Gusion, Botis
	"fraud":         {23, 51, 45},  // Aim, Balam, Vine
	// Reconciliation / diplomacy
	"reconcile":     {7, 17, 55},   // Amon, Botis, Orobas
	"diplomacy":     {7, 30, 11},   // Amon, Forneus, Gusion
	"conflict":      {14, 7, 53},   // Leraje, Amon, Caim
	"negotiation":   {7, 22, 53},   // Amon, Ipos, Caim
	// Technical / engineering
	"engineering":   {5, 60, 38},   // Marbas, Vapula, Halphas
	"architecture":  {38, 39, 60},  // Halphas, Malphas, Vapula
	"code":          {60, 5, 39},   // Vapula, Marbas, Malphas
	"build":         {38, 39, 5},   // Halphas, Malphas, Marbas
	"machine":       {60, 5, 21},   // Vapula, Marbas, Marax
	// Discovery / hidden
	"hidden":        {20, 62, 1},   // Purson, Valac, Bael
	"discover":      {3, 20, 56},   // Vassago, Purson, Gremory
	"treasure":      {20, 62, 56},  // Purson, Valac, Gremory
	"find":          {3, 62, 56},   // Vassago, Valac, Gremory
	// Transformation / alchemy
	"transform":     {61, 48, 4},   // Zagan, Haagenti, Samigina
	"transmute":     {61, 28, 48},  // Zagan, Berith, Haagenti
	"change":        {61, 45, 7},   // Zagan, Vine, Amon
	// Persuasion / rhetoric
	"persuade":      {22, 71, 53},  // Ipos, Dantalion, Caim
	"rhetoric":      {22, 24, 2},   // Ipos, Naberius, Agares
	"influence":     {71, 9, 22},   // Dantalion, Paimon, Ipos
	// Healing / restoration
	"heal":          {10, 19, 11},  // Buer, Sallos, Gusion
	"restore":       {10, 11, 37},  // Buer, Gusion, Phenex
	"repair":        {5, 10, 38},   // Marbas, Buer, Halphas
	// Chaos / disruption
	"disrupt":       {39, 45, 33},  // Malphas, Vine, Gaap
	"chaos":         {33, 39, 23},  // Gaap, Malphas, Aim
	"break":         {45, 39, 23},  // Vine, Malphas, Aim
	// Memory / time
	"memory":        {29, 4, 64},   // Astaroth, Samigina, Haures
	"past":          {4, 29, 64},   // Samigina, Astaroth, Haures
	"future":        {3, 29, 47},   // Vassago, Astaroth, Vual
	"time":          {29, 2, 3},    // Astaroth, Agares, Vassago
	// Justice / truth
	"justice":       {72, 55, 14},  // Andromalius, Orobas, Leraje
	"truth":         {55, 23, 35},  // Orobas, Aim, Marchosias
	"integrity":     {55, 35, 50},  // Orobas, Marchosias, Furcas
	// Music / art
	"music":         {67, 37, 24},  // Amdusias, Phenex, Naberius
	"art":           {9, 37, 22},   // Paimon, Phenex, Ipos
	"sound":         {67, 49, 53},  // Amdusias, Crocell, Caim
	// Travel / logistics
	"travel":        {70, 18, 33},  // Seere, Bathin, Gaap
	"logistics":     {70, 52, 21},  // Seere, Alloces, Marax
	"transport":     {70, 33, 18},  // Seere, Gaap, Bathin
	// Protection / defense
	"defend":        {38, 31, 1},   // Halphas, Foras, Bael
	"protect":       {38, 31, 1},   // Halphas, Foras, Bael
	"security":      {1, 31, 38},   // Bael, Foras, Halphas
	"stealth":       {1, 31, 51},   // Bael, Foras, Balam
	// Learning / philosophy
	"learn":         {2, 36, 33},   // Agares, Stolas, Gaap
	"philosophy":    {33, 50, 54},  // Gaap, Furcas, Murmur
	"knowledge":     {20, 36, 3},   // Purson, Stolas, Vassago
	"wisdom":        {36, 48, 54},  // Stolas, Haagenti, Murmur
	// Social / reputation
	"reputation":    {30, 11, 40},  // Forneus, Gusion, Raum
	"social":        {19, 30, 6},   // Sallos, Forneus, Valefor
	"loyalty":       {6, 19, 35},   // Valefor, Sallos, Marchosias
	// Game / combat
	"combat":        {14, 35, 50},  // Leraje, Marchosias, Furcas
	"war":           {66, 50, 14},  // Cimeies, Furcas, Leraje
	"game":          {9, 14, 22},   // Paimon, Leraje, Ipos
	// RSI / productivity
	"rsi":           {15, 3, 71},   // Eligos, Vassago, Dantalion
	"productivity":  {15, 38, 21},  // Eligos, Halphas, Marax
	"decision":      {15, 21, 17},  // Eligos, Marax, Botis
	// Emotion / relationship
	"love":          {13, 19, 56},  // Beleth, Sallos, Gremory
	"emotion":       {12, 13, 7},   // Sitri, Beleth, Amon
	"relationship":  {19, 6, 40},   // Sallos, Valefor, Raum
}

// MatchIntent selects up to 3 spirits for an intent string using keyword matching.
// It respects allowCollapse — collapse spirits are excluded unless allowCollapse=true.
// The returned stack never exceeds MaxStackSize (3) and is sorted by relevance score.
func MatchIntent(intent string, allowCollapse bool) []Spirit {
	lower := strings.ToLower(intent)
	scores := make(map[int]int) // rank → score

	for kw, ranks := range intentKeywords {
		if strings.Contains(lower, kw) {
			for i, rank := range ranks {
				// Higher score for first-listed (most relevant) spirits
				scores[rank] += len(ranks) - i
			}
		}
	}

	// Collect scored spirits, filtering Collapse if not allowed.
	type scored struct {
		s     Spirit
		score int
	}
	var candidates []scored
	for rank, score := range scores {
		if rank < 1 || rank > 72 {
			continue
		}
		s := All72[rank-1]
		if s.AllowCollapse && !allowCollapse {
			continue
		}
		candidates = append(candidates, scored{s, score})
	}

	// Sort by score descending, then amplitude descending as tiebreak.
	for i := 1; i < len(candidates); i++ {
		for j := i; j > 0 && (candidates[j].score > candidates[j-1].score ||
			(candidates[j].score == candidates[j-1].score && candidates[j].s.Amplitude > candidates[j-1].s.Amplitude)); j-- {
			candidates[j], candidates[j-1] = candidates[j-1], candidates[j]
		}
	}

	// Take top 3, normalize amplitude if needed.
	n := MaxStackSize
	if len(candidates) < n {
		n = len(candidates)
	}
	stack := make([]Spirit, n)
	for i := 0; i < n; i++ {
		stack[i] = candidates[i].s
	}

	// Amplitude safety: normalize if sum > MaxStackAmplitude.
	if sum := AmplitudeSum(stack); sum > MaxStackAmplitude {
		ratio := MaxStackAmplitude / sum
		for i := range stack {
			stack[i].Amplitude = math.Round(stack[i].Amplitude*ratio*100) / 100
		}
	}

	return stack
}

// ResonanceState classifies Δφ (phase difference in degrees) into a corridor.
func ResonanceState(deltaPhi float64) Corridor {
	switch {
	case deltaPhi <= 15:
		return CoherentLock
	case deltaPhi <= 38:
		return GoldenBand
	case deltaPhi <= 90:
		return ChaosEdge
	default:
		return Collapse
	}
}

// E1Weight returns E₁ (Carrier) contribution weight for a given resonance state.
func E1Weight(state Corridor) float64 {
	switch state {
	case CoherentLock:
		return 0.90
	case GoldenBand:
		return 0.60
	case ChaosEdge:
		return 0.30
	default: // Collapse
		return 1.00
	}
}

// E2Weight returns E₂ (Explorer) contribution weight for a given resonance state.
func E2Weight(state Corridor) float64 {
	return 1.0 - E1Weight(state)
}
