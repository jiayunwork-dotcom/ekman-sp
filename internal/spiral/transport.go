package spiral

func Transport(tau, rho, f float64) (float64, error) {
	if err := checkInput(rho, f, 1.0, tau); err != nil {
		return 0, err
	}
	return tau / (rho * absF(f)), nil
}

func TransportDeflection(f float64) float64 {
	return SignOf(f) * RightAngleDeg
}

func TransportHeading(windHeadingDeg, f float64) float64 {
	return TurnRight(windHeadingDeg, TransportDeflection(f))
}

func TransportVector(tau, rho, f, windHeadingDeg float64) (Vec, float64, error) {
	mag, err := Transport(tau, rho, f)
	if err != nil {
		return Vec{}, 0, err
	}
	heading := TransportHeading(windHeadingDeg, f)
	return FromPolar(mag, heading), heading, nil
}

func absF(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
