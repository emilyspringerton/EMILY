package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuditFrontDoor_AllHealthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><title>ok</title><body>fine</body></html>"))
	}))
	defer srv.Close()

	ctx := context.Background()
	passed, results := auditFrontDoor(ctx, []string{srv.URL, srv.URL})
	if !passed {
		t.Fatalf("expected all-healthy targets to pass, got results: %+v", results)
	}
	for _, r := range results {
		if !r.Passed {
			t.Errorf("target %s expected to pass, reasons: %v", r.URL, r.Reasons)
		}
	}
}

func TestAuditFrontDoor_5xxFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ctx := context.Background()
	passed, results := auditFrontDoor(ctx, []string{srv.URL})
	if passed {
		t.Fatalf("expected 5xx target to fail the gate")
	}
	if len(results) != 1 || results[0].Passed {
		t.Fatalf("expected single failing result, got: %+v", results)
	}
	if len(results[0].Reasons) == 0 || !strings.Contains(results[0].Reasons[0], "500") {
		t.Errorf("expected reason to mention HTTP 500, got: %v", results[0].Reasons)
	}
}

func TestAuditFrontDoor_TooManyBrokenLinksFails(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body>
			<a href="/dead1">a</a>
			<a href="/dead2">b</a>
			<a href="/dead3">c</a>
			<a href="/dead4">d</a>
		</body></html>`))
	})
	for _, p := range []string{"/dead1", "/dead2", "/dead3", "/dead4"} {
		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
	}
	srv := httptest.NewServer(mux)
	defer srv.Close()

	ctx := context.Background()
	passed, results := auditFrontDoor(ctx, []string{srv.URL})
	if passed {
		t.Fatalf("expected target with 4 broken links (>3) to fail the gate")
	}
	if len(results[0].BrokenLinks) != 4 {
		t.Errorf("expected 4 broken links detected, got %d", len(results[0].BrokenLinks))
	}
}

func TestAuditFrontDoor_UpToThreeBrokenLinksPasses(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body>
			<a href="/dead1">a</a>
			<a href="/dead2">b</a>
			<a href="/dead3">c</a>
		</body></html>`))
	})
	for _, p := range []string{"/dead1", "/dead2", "/dead3"} {
		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
	}
	srv := httptest.NewServer(mux)
	defer srv.Close()

	ctx := context.Background()
	passed, _ := auditFrontDoor(ctx, []string{srv.URL})
	if !passed {
		t.Fatalf("expected exactly 3 broken links (not >3) to still pass the gate")
	}
}

func TestAuditFrontDoor_UnreachableTargetFails(t *testing.T) {
	ctx := context.Background()
	passed, results := auditFrontDoor(ctx, []string{"http://127.0.0.1:1"})
	if passed {
		t.Fatalf("expected unreachable target to fail the gate")
	}
	if len(results) != 1 || results[0].Passed {
		t.Fatalf("expected single failing result for unreachable target, got: %+v", results)
	}
}
