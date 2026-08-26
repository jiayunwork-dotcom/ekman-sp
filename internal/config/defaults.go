package config

import "ekman-sp/internal/spiral"

var defaults = struct {
	windDirDeg float64
	rho        float64
	K          float64
	nPoints    int
}{
	windDirDeg: spiral.DefaultWindHeading,
	rho:        spiral.DefaultRho,
	K:          spiral.DefaultK,
	nPoints:    spiral.DefaultNPoints,
}

func applyDefaults(in Input) (windDirDeg, rho, K float64, nPoints int) {
	windDirDeg = defaults.windDirDeg
	if in.WindDirDeg != nil {
		windDirDeg = *in.WindDirDeg
	}
	rho = defaults.rho
	if in.Rho != nil {
		rho = *in.Rho
	}
	K = defaults.K
	if in.K != nil {
		K = *in.K
	}
	nPoints = defaults.nPoints
	if in.NPoints != nil {
		nPoints = *in.NPoints
	}
	return windDirDeg, rho, K, nPoints
}

func resolvedDefaults(in Input) (float64, float64, float64, int) {
	return applyDefaults(in)
}
