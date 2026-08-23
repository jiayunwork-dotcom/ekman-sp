package spiral

import "math"

// SurfaceSpeed returns the magnitude of the surface Ekman current:
//
//	|V0| = tau / (rho * sqrt(K*|f|))
//
// The velocity magnitude scales linearly with the wind stress and drops like
// 1/sqrt(|f|) and 1/sqrt(K); doubling |f| or K therefore reduces |V0| by a
// factor 1/sqrt(2), while doubling tau doubles |V0|.
func SurfaceSpeed(tau, rho, K, f float64) (float64, error) {
	if err := checkInput(rho, f, K, tau); err != nil {
		return 0, err
	}
	speed := tau / (rho * sqrtKf(K, f))
	return rememberSpeed(tau, rho, K, f, speed), nil
}

// sqrtKf returns sqrt(K*|f|), the denominator of the surface-speed formula.
func sqrtKf(K, f float64) float64 {
	return sqrtAbs(K * f)
}

// SurfaceDeflection returns the signed deflection of the surface current from
// the wind direction in degrees: +45 (right of the wind) in the northern
// hemisphere, -45 (left of the wind) in the southern hemisphere. The sign
// follows the Coriolis parameter.
func SurfaceDeflection(f float64) float64 {
	return SignOf(f) * SurfaceVeerDeg
}

// SurfaceHeading returns the compass heading of the surface current given the
// wind heading, by adding the hemisphere-dependent deflection.
func SurfaceHeading(windHeadingDeg, f float64) float64 {
	return TurnRight(windHeadingDeg, SurfaceDeflection(f))
}

// SurfaceVelocity returns the surface velocity vector (east, north) in m/s
// for the given wind stress magnitude and heading.
func SurfaceVelocity(tau, rho, K, f, windHeadingDeg float64) (Vec, float64, error) {
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return Vec{}, 0, err
	}
	heading := SurfaceHeading(windHeadingDeg, f)
	return FromPolar(speed, heading), heading, nil
}

// sqrtAbs returns sqrt(|x|) without ever producing NaN for the valid ranges
// used in the denominators above. The absolute value is taken first because
// checkInput already rejected f == 0, so |K*f| > 0.
func sqrtAbs(x float64) float64 {
	if x < 0 {
		x = -x
	}
	return math.Sqrt(x)
}
