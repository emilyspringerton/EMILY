package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEmilyListDirFlat(t *testing.T) {
	dir := t.TempDir()
	// Create some files and a hidden file (should be skipped)
	for _, name := range []string{"b.txt", "a.txt", ".hidden"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	lines, err := emilyListDir(dir, 2, 0)
	if err != nil {
		t.Fatalf("emilyListDir: %v", err)
	}
	// .hidden excluded; a.txt, b.txt sorted
	if len(lines) != 2 {
		t.Errorf("want 2 entries (a.txt, b.txt), got %d: %v", len(lines), lines)
	}
	if lines[0] != "a.txt" || lines[1] != "b.txt" {
		t.Errorf("wrong order/content: %v", lines)
	}
}

func TestEmilyListDirSubdirsFirst(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "file.go"), []byte(""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	lines, err := emilyListDir(dir, 2, 0)
	if err != nil {
		t.Fatalf("emilyListDir: %v", err)
	}
	// subdir/ should appear before file.go
	if len(lines) < 2 {
		t.Fatalf("want ≥2 entries, got %d", len(lines))
	}
	if !strings.HasSuffix(lines[0], "/") {
		t.Errorf("first entry should be directory (ends with /): %q", lines[0])
	}
}

func TestEmilyListDirMaxDepth(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// File inside sub — should only appear when maxDepth allows it
	if err := os.WriteFile(filepath.Join(sub, "deep.txt"), []byte(""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// maxDepth=1 → only show sub/, not sub/deep.txt
	lines, _ := emilyListDir(dir, 1, 0)
	for _, l := range lines {
		if strings.Contains(l, "deep.txt") {
			t.Errorf("maxDepth=1 should not show deep.txt, but got: %v", lines)
		}
	}

	// maxDepth=2 → sub/deep.txt should appear
	lines2, _ := emilyListDir(dir, 2, 0)
	found := false
	for _, l := range lines2 {
		if strings.Contains(l, "deep.txt") {
			found = true
		}
	}
	if !found {
		t.Errorf("maxDepth=2 should show deep.txt, got: %v", lines2)
	}
}

func TestEmilyListDirNonExistent(t *testing.T) {
	_, err := emilyListDir("/does/not/exist/surely", 2, 0)
	if err == nil {
		t.Fatal("expected error for non-existent dir, got nil")
	}
}
