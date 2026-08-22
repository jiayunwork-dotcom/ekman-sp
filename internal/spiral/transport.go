package spiral

// Transport returns the depth-integrated Ekman volume transport magnitude:
//
//	Me = tau / (rho * |f|)
//
// expressed in m^2/s (volume transport per unit width along the coast/along
// the front). Note the denominator is rho*|f|, NOT rho*sqrt(|f|) or
// rho*sqrt(K*|f|): the transport does not involve the eddy viscosity at all,
// because the vertical integral of the spiral collapses to a balance between
// the surface stress and the Coriolis force on the water column.
func Transport(tau, rho, f float64) (float64, error) {
	if err := checkInput(rho, f, 1.0, tau); err != nil {
		return 0, err
	}
	return tau / (rho * absF(f)), nil
}

// TransportDeflection returns the signed turning angle of the depth-integrated
// transport relative to the wind: +90 degrees in the northern hemisphere
// (transport to the right of the wind) and -90 degrees in the southern
// hemisphere (transport to the left).
func TransportDeflection(f float64) float64 {
	return SignOf(f) * RightAngleDeg
}

// TransportHeading returns the compass heading of the depth-integrated
// transport given the wind heading.
func TransportHeading(windHeadingDeg, f float64) float64 {
	return TurnRight(windHeadingDeg, TransportDeflection(f))
}

// TransportVector returns the depth-integrated transport as a vector in
// (east, north) coordinates with magnitude tau/(rho*|f|). Integrating the
// analytic profile numerically and returning to this vector is one of the
// invariants tested by the package; see the TestTransportIntegral test.
func TransportVector(tau, rho, f, windHeadingDeg float64) (Vec, float64, error) {
	mag, err := Transport(tau, rho, f)
	if err != nil {
		return Vec{}, 0, err
	}
	heading := TransportHeading(windHeadingDeg, f)
	return FromPolar(mag, heading), heading, nil
}

// absF returns |f|. checkInput guarantees f is not zero before this is
// reached, so the result is strictly positive.
func absF(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
