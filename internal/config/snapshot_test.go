package config

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSnapshotOf(t *testing.T) {
	r := mustResolve(t, sampleInput())
	s := SnapshotOf(r)
	if s.F <= 0 {
		t.Errorf("snapshot f = %g, want positive", s.F)
	}
	if s.LatitudeDeg == nil || *s.LatitudeDeg != 45 {
		t.Errorf("snapshot latitude = %v, want 45", s.LatitudeDeg)
	}
	if s.Hemisphere != "northern hemisphere" {
		t.Errorf("snapshot hemisphere = %q", s.Hemisphere)
	}
	if s.FiniteDepth {
		t.Error("snapshot finite depth should be false for the sample case")
	}
}

func TestSnapshotFiniteDepth(t *testing.T) {
	in := sampleInput()
	depth := 500.0
	in.DepthM = &depth
	r := mustResolve(t, in)
	s := SnapshotOf(r)
	if !s.FiniteDepth || s.DepthM == nil || *s.DepthM != depth {
		t.Errorf("snapshot finite depth = %v, depth = %v", s.FiniteDepth, s.DepthM)
	}
}

func TestMarshalResolvedIsValidJSON(t *testing.T) {
	r := mustResolve(t, sampleInput())
	b, err := MarshalResolved(r)
	if err != nil {
		t.Fatalf("MarshalResolved error: %v", err)
	}
	var back map[string]interface{}
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("snapshot is not valid JSON: %v", err)
	}
	if _, ok := back["f"]; !ok {
		t.Error("snapshot missing f field")
	}
	if !strings.Contains(string(b), "latitude_deg") {
		t.Error("snapshot missing latitude_deg field")
	}
}
