package spiral

import "math/cmplx"

// lastWind is the last assembled complex surface current. It is always
// overwritten from this call's wind stress so an expired leftover cannot
// be served as this case.
var lastWind complex128
var windCached bool

// cachedSurfaceComplex returns the complex surface current computed from
// this case's tau, heading and Coriolis parameter. A leftover east-wind
// entry is never returned just because the slot is still marked fresh.
func cachedSurfaceComplex(tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
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
