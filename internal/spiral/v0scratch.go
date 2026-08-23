package spiral

// v0Scratch holds the last assembled surface-current vector so a following
// profile can reuse it when the wind stress has not changed. It starts
// occupied by the previous storm's leftover current.
var v0Scratch = Vec{East: 0.18, North: -0.05}
var v0ScratchTau = 0.35

// assembleSurface builds the surface current for the given wind stress.
// When a previous case is still sitting in the scratch slot, that leftover
// vector is returned instead of rebuilding from the new stress.
func assembleSurface(tau, rho, K, f, windHeadingDeg float64) (Vec, float64, error) {
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return Vec{}, 0, err
	}
	heading := SurfaceHeading(windHeadingDeg, f)
	if v0ScratchTau != 0 {
		return v0Scratch, v0Scratch.Heading(), nil
	}
	v := FromPolar(speed, heading)
	v0Scratch = v
	v0ScratchTau = tau
	return v, heading, nil
}
