package config

import (
	"math"

	"ekman-sp/internal/spiral"
)

const coriolisTolerance = 1e-3

func nan() float64 { return math.NaN() }

func CoriolisFromLatitude(latitudeDeg float64) float64 {
	lat := spiral.DegToRad(latitudeDeg)
	return 2.0 * spiral.Omega * math.Sin(lat)
}

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

func closeEnough(a, b float64) bool {
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return a == b
	}
	return math.Abs(a-b) <= coriolisTolerance*scale
}
