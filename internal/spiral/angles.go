package spiral

import "math"

func DegToRad(deg float64) float64 { return deg * math.Pi / 180.0 }

func RadToDeg(rad float64) float64 { return rad * 180.0 / math.Pi }

func NormalizeHeading(deg float64) float64 {
	deg = math.Mod(deg, 360.0)
	if deg < 0 {
		deg += 360.0
	}
	return deg
}

func HeadingToPhase(headingDeg float64) float64 {
	return DegToRad(90.0 - NormalizeHeading(headingDeg))
}

func PhaseToHeading(phaseRad float64) float64 {
	return NormalizeHeading(90.0 - RadToDeg(phaseRad))
}

func TurnRight(headingDeg, deg float64) float64 {
	return NormalizeHeading(headingDeg + deg)
}

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

func WindUnit(windHeadingDeg float64) (east, north float64) {
	phi := HeadingToPhase(windHeadingDeg)
	return math.Cos(phi), math.Sin(phi)
}

func HemisphereName(f float64) string {
	switch ClassifyHemisphere(f) {
	case HemisphereNorth:
		return "northern hemisphere"
	case HemisphereSouth:
		return "southern hemisphere"
	}
	return "undefined"
}
