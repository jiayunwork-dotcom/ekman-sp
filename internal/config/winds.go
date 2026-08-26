package config

import "ekman-sp/internal/spiral"

func WindHeading(in Input) float64 {
	if in.WindDirDeg != nil {
		return spiral.NormalizeHeading(*in.WindDirDeg)
	}
	return spiral.NormalizeHeading(spiral.DefaultWindHeading)
}

func ResolvedWindHeading(in Input) (heading float64, explicit bool) {
	if in.WindDirDeg != nil {
		return spiral.NormalizeHeading(*in.WindDirDeg), true
	}
	return spiral.NormalizeHeading(spiral.DefaultWindHeading), false
}

func HemisphereLabelFrom(f float64) string {
	return spiral.HemisphereName(f)
}
