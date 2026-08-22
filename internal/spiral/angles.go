package spiral

import "math"

// Angle conventions used throughout the package:
//
//   - headings are reported in degrees clockwise from true north (0 = north,
//     90 = east, 180 = south, 270 = west), matching oceanographic usage;
//   - internal phase angles are radians counter-clockwise from the east axis,
//     matching the standard (u east, v north) coordinate pair.
//
// The conversion between the two is: heading = 90 - degrees(phase).

// DegToRad converts degrees to radians.
func DegToRad(deg float64) float64 { return deg * math.Pi / 180.0 }

// RadToDeg converts radians to degrees.
func RadToDeg(rad float64) float64 { return rad * 180.0 / math.Pi }

// NormalizeHeading folds any angle into [0, 360) degrees clockwise from north.
func NormalizeHeading(deg float64) float64 {
	deg = math.Mod(deg, 360.0)
	if deg < 0 {
		deg += 360.0
	}
	return deg
}

// HeadingToPhase converts a compass heading (degrees clockwise from north)
// into a phase angle in radians counter-clockwise from the east axis.
func HeadingToPhase(headingDeg float64) float64 {
	return DegToRad(90.0 - NormalizeHeading(headingDeg))
}

// PhaseToHeading converts a phase angle in radians counter-clockwise from the
// east axis into a compass heading in degrees clockwise from north.
func PhaseToHeading(phaseRad float64) float64 {
	return NormalizeHeading(90.0 - RadToDeg(phaseRad))
}

// TurnRight rotates a compass heading rightward (clockwise on the compass) by
// the given number of degrees. The sign of the argument is taken literally:
// a negative turn is a left turn.
func TurnRight(headingDeg, deg float64) float64 {
	return NormalizeHeading(headingDeg + deg)
}

// HeadingDifference returns the signed smallest-angle difference
// (a - b) in degrees, in (-180, 180]. It is used to report how much a current
// at depth has veered relative to the surface current.
func HeadingDifference(a, b float64) float64 {
	d := NormalizeHeading(a) - NormalizeHeading(b)
	if d > 180 {
		d -= 360
	}
	if d <= -180 {
		d += 360
	}
	return d
}

// WindUnit returns the unit vector (east, north) in the direction the wind is
// blowing toward, given the heading in degrees clockwise from north.
func WindUnit(windHeadingDeg float64) (east, north float64) {
	phi := HeadingToPhase(windHeadingDeg)
	return math.Cos(phi), math.Sin(phi)
}

// HemisphereName returns a short label for the sign of f.
func HemisphereName(f float64) string {
	switch ClassifyHemisphere(f) {
	case HemisphereNorth:
		return "northern hemisphere"
	case HemisphereSouth:
		return "southern hemisphere"
	}
	return "undefined"
}
