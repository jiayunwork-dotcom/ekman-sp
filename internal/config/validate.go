package config

import (
	"errors"
	"fmt"

	"ekman-sp/internal/spiral"
)

// Validation errors. Each one maps to a user-visible message on the CLI and
// to a test that provokes it.
var (
	// ErrTauRequired fires when the case file omits the wind stress or sets
	// it to a non-positive value.
	ErrTauRequired = errors.New("wind stress tau is required and must be positive")

	// ErrRhoNonPositive fires when an explicit density is zero or negative.
	ErrRhoNonPositive = errors.New("density rho must be positive")

	// ErrKNonPositive fires when an explicit eddy viscosity is zero or
	// negative.
	ErrKNonPositive = errors.New("eddy viscosity K must be positive")

	// ErrNoCoriolis fires when neither f nor latitude_deg is supplied.
	ErrNoCoriolis = errors.New("must supply f (Coriolis parameter) or latitude_deg")

	// ErrEquatorial fires when f resolves to zero (explicit f=0, or
	// latitude 0).
	ErrEquatorial = errors.New("Coriolis parameter f must be non-zero; the Ekman spiral is undefined at the equator")

	// ErrCoriolisConflict fires when both f and latitude_deg are supplied but
	// disagree beyond the consistency tolerance.
	ErrCoriolisConflict = errors.New("f and latitude_deg disagree; supply only one")

	// ErrBadDepth fires when a finite water depth is zero or negative.
	ErrBadDepth = errors.New("finite water depth depth_m must be positive")

	// ErrShallowWater fires when a finite water depth is smaller than the
	// Ekman depth; the classic infinite-depth profile no longer applies.
	ErrShallowWater = errors.New("water depth smaller than Ekman depth De: the classic infinite-depth spiral does not apply; use a deeper case or drop depth_m")

	// ErrBadNPoints fires when the profile sample count is zero or negative.
	ErrBadNPoints = errors.New("npoints must be a positive integer")

	// ErrBadDepthMax fires when depth_max_m is non-positive.
	ErrBadDepthMax = errors.New("depth_max_m must be positive")
)

// validate checks every physical invariant of the case and returns the
// resolved Coriolis parameter together with the latitude that produced it
// (NaN when f was given directly).
func validate(in Input) (f, latitudeDeg float64, err error) {
	f, latitudeDeg, err = validateInner(in)
	if err != nil {
		return 0, latitudeDeg, bindValidateErr(err)
	}
	return f, latitudeDeg, nil
}

func validateInner(in Input) (f, latitudeDeg float64, err error) {
	latitudeDeg = nan()

	if in.Tau <= 0 {
		return 0, latitudeDeg, ErrTauRequired
	}
	if in.Rho != nil && *in.Rho <= 0 {
		return 0, latitudeDeg, fmt.Errorf("%w (got %g)", ErrRhoNonPositive, *in.Rho)
	}
	if in.K != nil && *in.K <= 0 {
		return 0, latitudeDeg, fmt.Errorf("%w (got %g)", ErrKNonPositive, *in.K)
	}
	if in.NPoints != nil && *in.NPoints <= 0 {
		return 0, latitudeDeg, fmt.Errorf("%w (got %d)", ErrBadNPoints, *in.NPoints)
	}
	if in.DepthMaxM != nil && *in.DepthMaxM <= 0 {
		return 0, latitudeDeg, fmt.Errorf("%w (got %g)", ErrBadDepthMax, *in.DepthMaxM)
	}

	f, latitudeDeg, err = resolveCoriolis(in.F, in.LatitudeDeg)
	if err != nil {
		return 0, latitudeDeg, err
	}

	if in.DepthM != nil {
		if *in.DepthM <= 0 {
			return 0, latitudeDeg, fmt.Errorf("%w (got %g)", ErrBadDepth, *in.DepthM)
		}
		K := defaults.K
		if in.K != nil {
			K = *in.K
		}
		if spiral.IsShallow(*in.DepthM, K, f) {
			return 0, latitudeDeg, fmt.Errorf("%w: depth_m=%g < De=%.3f m", ErrShallowWater, *in.DepthM, spiral.EkmanDepth(K, f))
		}
	}
	return f, latitudeDeg, nil
}
