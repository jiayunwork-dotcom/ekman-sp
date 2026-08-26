package config

import (
	"fmt"
	"math"

	"ekman-sp/internal/spiral"
)

func (r *Resolved) Describe() string {
	f := r.Params.F
	lat := r.LatitudeDeg
	origin := "f given directly"
	if !math.IsNaN(lat) {
		origin = fmt.Sprintf("lat %.1f deg", lat)
	}
	equiv := "n/a"
	if math.IsNaN(lat) {
		equiv = fmt.Sprintf("equiv lat %.1f deg", LatitudeFromCoriolis(f))
	}
	return fmt.Sprintf(
		"f=%+.3e 1/s (%s, %s), rho=%.1f kg/m^3, K=%.3f m^2/s",
		f, spiral.HemisphereName(f), origin, r.Params.Rho, r.Params.K,
	) + " " + equiv
}

func (r *Resolved) CoriolisProvenance() string {
	if math.IsNaN(r.LatitudeDeg) {
		return fmt.Sprintf("given directly (equiv lat %.1f deg)", LatitudeFromCoriolis(r.Params.F))
	}
	return fmt.Sprintf("from latitude %.1f deg", r.LatitudeDeg)
}

func (r *Resolved) ProfileGeometry() (ekmanDepthM, fullTurnM float64) {
	ekmanDepthM = spiral.EkmanDepth(r.Params.K, r.Params.F)
	fullTurnM = spiral.FullTurnDepth(r.Params.K, r.Params.F)
	return ekmanDepthM, fullTurnM
}
