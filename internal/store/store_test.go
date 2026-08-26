package store

import (
	"path/filepath"
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

func TestWriteReadStore(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	entries := map[string]*config.Resolved{"north": sampleResolved(t)}
	if err := WriteFile(path, entries); err != nil {
		t.Fatal(err)
	}
	f, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	rec, ok := FindRecord(f, "north")
	if !ok || rec.V0 <= 0 {
		t.Fatalf("FindRecord() ok=%v rec=%+v", ok, rec)
	}
}

func TestVerifyFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "verify.json")
	entries := map[string]*config.Resolved{"a": sampleResolved(t)}
	if err := WriteFile(path, entries); err != nil {
		t.Fatal(err)
	}
	if err := VerifyFile(path, entries, 1e-12); err != nil {
		t.Fatalf("VerifyFile() = %v", err)
	}
}
