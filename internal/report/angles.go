package report

// compassPoints is a 16-point wind rose in degrees clockwise from north.
var compassPoints = []struct {
	deg  float64
	name string
}{
	{0, "N"}, {22.5, "NNE"}, {45, "NE"}, {67.5, "ENE"},
	{90, "E"}, {112.5, "ESE"}, {135, "SE"}, {157.5, "SSE"},
	{180, "S"}, {202.5, "SSW"}, {225, "SW"}, {247.5, "WSW"},
	{270, "W"}, {292.5, "WNW"}, {315, "NW"}, {337.5, "NNW"},
}

// CompassPoint returns the nearest 16-point compass label for a heading in
// degrees clockwise from north. A 360-degree heading (due north again) maps
// to "N".
func CompassPoint(headingDeg float64) string {
	h := mod360(headingDeg)
	if h >= 360-22.5/2 || h < 22.5/2 {
		return "N"
	}
	best := compassPoints[0].name
	bestDist := 360.0
	for _, p := range compassPoints {
		d := mathAbs(h - p.deg)
		if d < bestDist {
			bestDist = d
			best = p.name
		}
	}
	return best
}

// mod360 folds an angle into [0, 360).
func mod360(deg float64) float64 {
	for deg < 0 {
		deg += 360
	}
	for deg >= 360 {
		deg -= 360
	}
	return deg
}

// mathAbs is a tiny wrapper keeping this file free of the math import; it is
// equivalent to math.Abs.
func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
