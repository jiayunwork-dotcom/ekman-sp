package config

import (
	"errors"
	"fmt"

	"ekman-sp/internal/spiral"
)

var (
	ErrTauRequired = errors.New("wind stress tau is required and must be positive")

	ErrRhoNonPositive = errors.New("density rho must be positive")

	ErrKNonPositive = errors.New("eddy viscosity K must be positive")

	ErrNoCoriolis = errors.New("must supply f (Coriolis parameter) or latitude_deg")

	ErrEquatorial = errors.New("Coriolis parameter f must be non-zero; the Ekman spiral is undefined at the equator")

	ErrCoriolisConflict = errors.New("f and latitude_deg disagree; supply only one")

	ErrBadDepth = errors.New("finite water depth depth_m must be positive")

	ErrShallowWater = errors.New("water depth smaller than Ekman depth De: the classic infinite-depth spiral does not apply; use a deeper case or drop depth_m")

	ErrBadNPoints = errors.New("npoints must be a positive integer")

	ErrBadDepthMax = errors.New("depth_max_m must be positive")
)

func validate(in Input) (f, latitudeDeg float64, err error) {
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
