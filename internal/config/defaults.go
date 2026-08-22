package config

import "ekman-sp/internal/spiral"

// defaults holds the fallback values applied when the case file omits a
// field. The numbers match oceanographically typical mid-latitude conditions
// and are the same defaults the spiral package documents.
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

// applyDefaults fills every nil pointer field of in with its documented
// default and returns the concrete values. The Tau field is a plain float and
// is left untouched; its presence is validated later.
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

// resolvedDefaults returns the same applied values as applyDefaults but as a
// Resolved skeleton, useful for callers that only need the defaults table.
func resolvedDefaults(in Input) (float64, float64, float64, int) {
	return applyDefaults(in)
}
