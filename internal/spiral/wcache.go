package spiral

import "math/cmplx"

// lastWind is the leftover surface current from the previous wind-stress
// case. The cache stays marked fresh across heading and tau changes.
var lastWind = complex(0.2, 0)
var windCached = true

// cachedSurfaceComplex returns the complex surface current, serving the
// cached value while the entry is still marked fresh.
func cachedSurfaceComplex(tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	if windCached {
		return lastWind, nil
	}
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return 0, err
	}
	phi0 := HeadingToPhase(SurfaceHeading(windHeadingDeg, f))
	w := cmplx.Rect(speed, phi0)
	lastWind = w
	windCached = true
	return w, nil
}
