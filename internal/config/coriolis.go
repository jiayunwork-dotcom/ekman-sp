package config

import (
	"math"

	"ekman-sp/internal/spiral"
)

// coriolisTolerance is the relative agreement required when both f and
// latitude_deg are supplied. 1e-3 of the value lets users round latitude to a
// tenth of a degree without tripping the consistency check.
const coriolisTolerance = 1e-3

// nan is a tiny alias for math.NaN, used to mark "no latitude supplied".
func nan() float64 { return math.NaN() }

// CoriolisFromLatitude returns the Coriolis parameter f = 2*Omega*sin(lat)
// for a latitude given in degrees. It is exported so tests and callers can
// derive f independently of the JSON pipeline.
func CoriolisFromLatitude(latitudeDeg float64) float64 {
	lat := spiral.DegToRad(latitudeDeg)
	return 2.0 * spiral.Omega * math.Sin(lat)
}

// resolveCoriolis pins down the Coriolis parameter from the case inputs. The
// precedence rule is:
//
//   - only f supplied: use it;
//   - only latitude_deg supplied: derive f = 2*Omega*sin(lat);
//   - both supplied: use them only if they agree within tolerance;
//   - neither supplied: error.
func resolveCoriolis(f *float64, latitudeDeg *float64) (float64, float64, error) {
	switch {
	case f != nil && latitudeDeg == nil:
		if *f == 0 {
			return 0, nan(), ErrEquatorial
		}
		return *f, nan(), nil
	case f == nil && latitudeDeg != nil:
		derived := CoriolisFromLatitude(*latitudeDeg)
		if derived == 0 {
			return 0, *latitudeDeg, ErrEquatorial
		}
		return derived, *latitudeDeg, nil
	case f != nil && latitudeDeg != nil:
		if *f == 0 || *latitudeDeg == 0 {
			return 0, *latitudeDeg, ErrEquatorial
		}
		derived := CoriolisFromLatitude(*latitudeDeg)
		if !closeEnough(*f, derived) {
			return 0, *latitudeDeg, ErrCoriolisConflict
		}
		return *f, *latitudeDeg, nil
	default:
		return 0, nan(), ErrNoCoriolis
	}
}

// closeEnough reports whether two Coriolis values agree to a relative
// tolerance, guarding the "both supplied" consistency rule.
func closeEnough(a, b float64) bool {
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return a == b
	}
	return math.Abs(a-b) <= coriolisTolerance*scale
}
