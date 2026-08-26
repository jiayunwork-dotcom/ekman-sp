package spiral

import "math"

const (
	Omega = 7.2921e-5

	DefaultRho = 1025.0

	DefaultK = 0.05

	DefaultWindHeading = 90.0

	DefaultTau = 0.1

	DefaultNPoints = 21

	DepthFactor = 2.0

	SurfaceVeerDeg = 45.0

	RightAngleDeg = 90.0

	HemisphereTolerance = 1e-12

	epsilon = 1e-9
)

type Hemisphere int

const (
	HemisphereUnknown Hemisphere = iota
	HemisphereNorth
	HemisphereSouth
)

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

func SignOf(f float64) float64 {
	switch ClassifyHemisphere(f) {
	case HemisphereNorth:
		return 1
	case HemisphereSouth:
		return -1
	}
	return 0
}
