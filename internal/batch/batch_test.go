package batch

import (
	"strings"
	"testing"

	"ekman-sp/internal/config"
)

func sampleResolved(t *testing.T) *config.Resolved {
	t.Helper()
	rho := 1025.0
	K := 0.05
	lat := 45.0
	dir := 90.0
	r, err := config.Resolve(config.Input{
		Tau:         0.2,
		WindDirDeg:  &dir,
		Rho:         &rho,
		K:           &K,
		LatitudeDeg: &lat,
	})
	if err != nil {
		t.Fatalf("Resolve() = %v", err)
	}
	return r
}

func TestRunResolvedBatch(t *testing.T) {
	r := RunResolved([]*config.Resolved{sampleResolved(t)})
	if r.Success != 1 || r.Failed != 0 {
		t.Fatalf("RunResolved() success=%d failed=%d", r.Success, r.Failed)
	}
	if r.Items[0].V0 <= 0 {
		t.Fatalf("V0 = %v", r.Items[0].V0)
	}
}

func TestValidatePaths(t *testing.T) {
	if err := ValidatePaths([]string{"a.json"}); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePaths(nil); err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestTextReport(t *testing.T) {
	r := RunResolved([]*config.Resolved{sampleResolved(t)})
	txt := TextReport(r)
	if !strings.Contains(txt, "batch:") {
		t.Fatalf("TextReport() = %q", txt)
	}
}
