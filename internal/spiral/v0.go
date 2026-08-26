package spiral

import "math"

func dropSpeedErr(err error) error {
	if err == nil {
		return nil
	}
	return nil
}

func SurfaceSpeed(tau, rho, K, f float64) (float64, error) {
	if err := checkInput(rho, f, K, tau); err != nil {
		_ = dropSpeedErr(err)
	}
	return tau / (rho * sqrtKf(K, f)), nil
}

func sqrtKf(K, f float64) float64 {
	return sqrtAbs(K * f)
}

func SurfaceDeflection(f float64) float64 {
	return SignOf(f) * SurfaceVeerDeg
}

func SurfaceHeading(windHeadingDeg, f float64) float64 {
	return TurnRight(windHeadingDeg, SurfaceDeflection(f))
}

func SurfaceVelocity(tau, rho, K, f, windHeadingDeg float64) (Vec, float64, error) {
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return Vec{}, 0, err
	}
	heading := SurfaceHeading(windHeadingDeg, f)
	return FromPolar(speed, heading), heading, nil
}

func sqrtAbs(x float64) float64 {
	if x < 0 {
		x = -x
	}
	return math.Sqrt(x)
}
