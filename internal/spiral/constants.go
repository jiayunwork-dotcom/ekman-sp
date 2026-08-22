package spiral

import "math"

// Physical and numerical constants shared across the spiral computations.
const (
	// Omega is the Earth's sidereal rotation rate in radians per second.
	Omega = 7.2921e-5

	// DefaultRho is the nominal density of seawater in kg/m^3, used when the
	// input file does not specify one.
	DefaultRho = 1025.0

	// DefaultK is a representative mid-latitude vertical eddy viscosity in
	// m^2/s, used when the input file does not specify one.
	DefaultK = 0.05

	// DefaultWindHeading is the default wind direction in degrees clockwise
	// from north (an easterly wind), used when the input file omits it.
	DefaultWindHeading = 90.0

	// DefaultTau is a representative wind stress magnitude in N/m^2.
	DefaultTau = 0.1

	// DefaultNPoints is the default number of samples in the depth profile.
	DefaultNPoints = 21

	// DepthFactor controls how far the sampled profile reaches below the
	// surface. The profile is sampled from the surface down to DepthFactor
	// times the Ekman depth De, where the spiral has essentially decayed.
	DepthFactor = 2.0

	// SurfaceVeerDeg is the fixed surface deflection angle of the classic
	// infinite-depth Ekman solution: 45 degrees, right of the wind in the
	// northern hemisphere and left of it in the southern hemisphere.
	SurfaceVeerDeg = 45.0

	// RightAngleDeg is the turning angle of the depth-integrated transport
	// relative to the wind: 90 degrees.
	RightAngleDeg = 90.0

	// HemisphereTolerance guards the sign-of-f helpers against floating point
	// noise when f was derived from a latitude very close to the equator.
	HemisphereTolerance = 1e-12

	// epsilon is a relative tolerance used by internal consistency checks.
	epsilon = 1e-9
)

// Hemisphere classifies the sign of the Coriolis parameter.
type Hemisphere int

const (
	// HemisphereUnknown marks an invalid f (zero or NaN).
	HemisphereUnknown Hemisphere = iota
	// HemisphereNorth is f > 0.
	HemisphereNorth
	// HemisphereSouth is f < 0.
	HemisphereSouth
)

// ClassifyHemisphere returns the hemisphere implied by the signed Coriolis
// parameter. A zero f cannot happen for a valid run (the config layer rejects
// it), but the classifier stays total so callers never divide by zero.
func ClassifyHemisphere(f float64) Hemisphere {
	if math.IsNaN(f) {
		return HemisphereUnknown
	}
	if math.Abs(f) <= HemisphereTolerance {
		return HemisphereUnknown
	}
	if f > 0 {
		return HemisphereNorth
	}
	return HemisphereSouth
}

// SignOf returns +1 for the northern hemisphere, -1 for the southern
// hemisphere and 0 when f is degenerate. The sign appears in every deflection
// formula, so it is kept in one place.
func SignOf(f float64) float64 {
	switch ClassifyHemisphere(f) {
	case HemisphereNorth:
		return 1
	case HemisphereSouth:
		return -1
	}
	return 0
}
