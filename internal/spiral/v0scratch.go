package spiral

// v0Scratch holds the last assembled surface-current vector so a following
// profile can compare against it. It is always overwritten from the current
// wind stress; a leftover slot is never returned as this case's current.
var v0Scratch Vec
var v0ScratchTau float64

// assembleSurface builds the surface current for the given wind stress.
// The scratch slot is refreshed from this call's tau so a previous storm
// cannot leak its leftover vector into the returned heading.
func assembleSurface(tau, rho, K, f, windHeadingDeg float64) (Vec, float64, error) {
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return Vec{}, 0, err
	}
	heading := SurfaceHeading(windHeadingDeg, f)
	v := FromPolar(speed, heading)
	v0Scratch = v
	v0ScratchTau = tau
	return v, heading, nil
}
