package config

import (
	"encoding/json"
	"fmt"
	"math"
)

type Snapshot struct {
	Tau         float64  `json:"tau"`
	WindDirDeg  float64  `json:"wind_dir_deg"`
	Rho         float64  `json:"rho"`
	K           float64  `json:"K"`
	F           float64  `json:"f"`
	LatitudeDeg *float64 `json:"latitude_deg,omitempty"`
	DepthM      *float64 `json:"depth_m,omitempty"`
	NPoints     int      `json:"npoints"`
	DepthMaxM   *float64 `json:"depth_max_m,omitempty"`
	FiniteDepth bool     `json:"finite_depth"`
	Hemisphere  string   `json:"hemisphere"`
}

func SnapshotOf(r *Resolved) Snapshot {
	s := Snapshot{
		Tau:         r.Params.Tau,
		WindDirDeg:  r.Params.WindHeading,
		Rho:         r.Params.Rho,
		K:           r.Params.K,
		F:           r.Params.F,
		NPoints:     r.Params.NPoints,
		FiniteDepth: r.FiniteDepth,
		Hemisphere:  HemisphereLabelFrom(r.Params.F),
	}
	if !math.IsNaN(r.LatitudeDeg) {
		lat := r.LatitudeDeg
		s.LatitudeDeg = &lat
	}
	if r.FiniteDepth {
		d := r.DepthM
		s.DepthM = &d
	}
	return s
}

func MarshalResolved(r *Resolved) ([]byte, error) {
	b, err := json.MarshalIndent(SnapshotOf(r), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal resolved case: %w", err)
	}
	return b, nil
}
