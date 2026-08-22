package config

import "ekman-sp/internal/spiral"

// WindHeading returns the resolved wind direction in degrees clockwise from
// north, applying the default when the case omits it.
func WindHeading(in Input) float64 {
	if in.WindDirDeg != nil {
		return spiral.NormalizeHeading(*in.WindDirDeg)
	}
	return spiral.NormalizeHeading(spiral.DefaultWindHeading)
}

// ResolvedWindHeading returns the heading together with a flag recording
// whether the case supplied one explicitly. The report uses the flag only to
// display the source of the value, never to change the math.
func ResolvedWindHeading(in Input) (heading float64, explicit bool) {
	if in.WindDirDeg != nil {
		return spiral.NormalizeHeading(*in.WindDirDeg), true
	}
	return spiral.NormalizeHeading(spiral.DefaultWindHeading), false
}

// HemisphereLabelFrom resolves the hemisphere description for reporting
// purposes from the resolved f value.
func HemisphereLabelFrom(f float64) string {
	return spiral.HemisphereName(f)
}
