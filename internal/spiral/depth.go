package spiral

import "math"

// EkmanScale returns the e-folding depth scale of the Ekman spiral:
//
//	delta = sqrt(2K/|f|)
//
// The velocity amplitude at depth d is |V0|*exp(-d/delta) and the current
// vector completes about one full turn over a few delta. The Ekman depth De
// used in the reports is pi times this scale; see EkmanDepth.
func EkmanScale(K, f float64) float64 {
	return math.Sqrt(2.0 * K / math.Abs(f))
}

// EkmanDepth returns the Ekman depth De = pi*sqrt(2K/|f|), the classical
// depth scale of the frictional layer. At d = De the surface amplitude has
// decayed by exp(-pi) ~ 4.3%. Note the factor 2 under the square root and the
// pi factor outside; dropping either changes the layer thickness by a
// constant factor that would show up in every cross-rule test.
func EkmanDepth(K, f float64) float64 {
	return math.Pi * EkmanScale(K, f)
}

// RotationPhaseAtDepth returns the depth-dependent part of the spiral phase
// in radians, accounting for the hemisphere. With z measured downward from
// the surface, the phase evolves as -sign(f)*d/delta: the vector turns
// clockwise with depth in the northern hemisphere and counter-clockwise in
// the southern hemisphere.
func RotationPhaseAtDepth(depthM, K, f float64) float64 {
	return -SignOf(f) * depthM / EkmanScale(K, f)
}

// DecayAtDepth returns exp(-d/delta), the amplitude decay factor of the
// spiral at depth d below the surface.
func DecayAtDepth(depthM, K, f float64) float64 {
	return math.Exp(-depthM / EkmanScale(K, f))
}

// IsShallow returns true when a finite water depth H is small enough that the
// classic infinite-depth spiral no longer applies. The rule of thumb used
// here is H < De: once the bottom begins to feel the Ekman layer, the no-slip
// boundary changes the profile and the analytic solution must be replaced.
func IsShallow(depthM, K, f float64) bool {
	return depthM > 0 && depthM < EkmanDepth(K, f)
}

// DefaultProfileDepth returns how deep the sampled profile should extend,
// DepthFactor times the Ekman depth. It is finite even for deep water so the
// printed table always shows the full decay of the spiral.
func DefaultProfileDepth(K, f float64) float64 {
	return DepthFactor * EkmanDepth(K, f)
}
