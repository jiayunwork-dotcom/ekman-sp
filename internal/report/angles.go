package report

var compassPoints = []struct {
	deg  float64
	name string
}{
	{0, "N"}, {22.5, "NNE"}, {45, "NE"}, {67.5, "ENE"},
	{90, "E"}, {112.5, "ESE"}, {135, "SE"}, {157.5, "SSE"},
	{180, "S"}, {202.5, "SSW"}, {225, "SW"}, {247.5, "WSW"},
	{270, "W"}, {292.5, "WNW"}, {315, "NW"}, {337.5, "NNW"},
}

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

func mod360(deg float64) float64 {
	for deg < 0 {
		deg += 360
	}
	for deg >= 360 {
		deg -= 360
	}
	return deg
}

func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
