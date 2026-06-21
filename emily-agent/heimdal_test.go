package main

import (
	"strings"
	"testing"
)

func TestFormatCriteriaEmpty(t *testing.T) {
	got := formatCriteria(nil)
	if got != "" {
		t.Errorf("empty criteria: want empty string, got %q", got)
	}
}

func TestFormatCriteriaSingle(t *testing.T) {
	criteria := []AcceptanceCriterion{
		{Name: "test coverage", Description: "all packages have tests", Target: ">80%"},
	}
	got := formatCriteria(criteria)
	if !strings.Contains(got, "1.") {
		t.Errorf("should start with '1.': %q", got)
	}
	if !strings.Contains(got, "test coverage") {
		t.Errorf("should contain criterion name: %q", got)
	}
	if !strings.Contains(got, ">80%") {
		t.Errorf("should contain target: %q", got)
	}
}

func TestFormatCriteriaMultiple(t *testing.T) {
	criteria := []AcceptanceCriterion{
		{Name: "A", Description: "first", Target: "100%"},
		{Name: "B", Description: "second", Target: "0 errors"},
	}
	got := formatCriteria(criteria)
	if !strings.Contains(got, "1.") {
		t.Errorf("missing '1.' in: %q", got)
	}
	if !strings.Contains(got, "2.") {
		t.Errorf("missing '2.' in: %q", got)
	}
}

func TestMinHelper(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{1, 2, 1},
		{5, 3, 3},
		{4, 4, 4},
		{-1, 1, -1},
	}
	for _, tc := range cases {
		if got := min(tc.a, tc.b); got != tc.want {
			t.Errorf("min(%d,%d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}
