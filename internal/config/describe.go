package config

import (
	"fmt"
	"math"

	"ekman-sp/internal/spiral"
)

// Describe renders a one-line, human-readable summary of the resolved case:
// where the Coriolis parameter came from, the hemisphere and the density and
// viscosity in use. The report header uses it to keep the input block compact.
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

// CoriolisProvenance returns a short phrase describing where f came from,
// suitable for the report header: "from latitude 45.0 deg" or "given
// directly (equiv lat 45.0 deg)".
func (r *Resolved) CoriolisProvenance() string {
	if math.IsNaN(r.LatitudeDeg) {
		return fmt.Sprintf("given directly (equiv lat %.1f deg)", LatitudeFromCoriolis(r.Params.F))
	}
	return fmt.Sprintf("from latitude %.1f deg", r.LatitudeDeg)
}

// ProfileGeometry returns the deep-water geometry the report quotes before
// the profile table: the Ekman depth and the depth of a full turn of the
// spiral. It is a pure description helper over the physics layer.
func (r *Resolved) ProfileGeometry() (ekmanDepthM, fullTurnM float64) {
	ekmanDepthM = spiral.EkmanDepth(r.Params.K, r.Params.F)
	fullTurnM = spiral.FullTurnDepth(r.Params.K, r.Params.F)
	return ekmanDepthM, fullTurnM
}
