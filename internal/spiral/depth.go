package spiral

import "math"

func EkmanScale(K, f float64) float64 {
	return math.Sqrt(2.0 * K / math.Abs(f))
}

func EkmanDepth(K, f float64) float64 {
	return math.Pi * EkmanScale(K, f)
}

func RotationPhaseAtDepth(depthM, K, f float64) float64 {
	return -SignOf(f) * depthM / EkmanScale(K, f)
}

func DecayAtDepth(depthM, K, f float64) float64 {
	return math.Exp(-depthM / EkmanScale(K, f))
}

func IsShallow(depthM, K, f float64) bool {
	return depthM > 0 && depthM < EkmanDepth(K, f)
}

func DefaultProfileDepth(K, f float64) float64 {
	return DepthFactor * EkmanDepth(K, f)
}
