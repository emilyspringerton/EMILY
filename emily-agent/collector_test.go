package main

import (
	"context"
	"testing"
	"time"
)

func TestDocStoreAppendAndLoadSeen(t *testing.T) {
	dir := t.TempDir()
	store, err := NewDocStore(dir)
	if err != nil {
		t.Fatalf("NewDocStore: %v", err)
	}
	defer store.Close()

	docs := []CollectedDoc{
		{ID: "doc-1", SourceName: "test", Title: "First", Body: "body1", URL: "http://example.com/1", CollectedAt: time.Now()},
		{ID: "doc-2", SourceName: "test", Title: "Second", Body: "body2", URL: "http://example.com/2", CollectedAt: time.Now()},
	}
	for _, d := range docs {
		if err := store.Append(d); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}

	seen, err := store.LoadSeen(context.Background())
	if err != nil {
		t.Fatalf("LoadSeen: %v", err)
	}
	if _, ok := seen["doc-1"]; !ok {
		t.Error("doc-1 not in seen set")
	}
	if _, ok := seen["doc-2"]; !ok {
		t.Error("doc-2 not in seen set")
	}
}

func TestDocStoreSeqIncrement(t *testing.T) {
	dir := t.TempDir()
	store, err := NewDocStore(dir)
	if err != nil {
		t.Fatalf("NewDocStore: %v", err)
	}
	defer store.Close()

	before := store.seq
	if err := store.Append(CollectedDoc{ID: "x", SourceName: "test"}); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if store.seq != before+1 {
		t.Errorf("seq after Append: got %d, want %d", store.seq, before+1)
	}
}

func TestDocStoreEmptyLoadSeen(t *testing.T) {
	dir := t.TempDir()
	store, err := NewDocStore(dir)
	if err != nil {
		t.Fatalf("NewDocStore: %v", err)
	}
	defer store.Close()

	seen, err := store.LoadSeen(context.Background())
	if err != nil {
		t.Fatalf("LoadSeen on empty store: %v", err)
	}
	if len(seen) != 0 {
		t.Errorf("empty store: want 0 seen docs, got %d", len(seen))
	}
}
