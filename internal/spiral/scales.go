package spiral

import "math"

func HalfPowerDepth(K, f float64) float64 {
	return EkmanScale(K, f) * math.Ln2
}

func QuarterTurnDepth(K, f float64) float64 {
	return EkmanScale(K, f) * math.Pi / 2.0
}

func FullTurnDepth(K, f float64) float64 {
	return EkmanScale(K, f) * 2.0 * math.Pi
}

func DecayFraction(depthM, K, f float64) float64 {
	return 1.0 - math.Exp(-depthM/EkmanScale(K, f))
}

func TurnPerMetre(K, f float64) float64 {
	return SignOf(f) * RadToDeg(1.0/EkmanScale(K, f))
}

func SurfaceFlowAngle(windHeadingDeg, f float64) float64 {
	return HeadingToPhase(SurfaceHeading(windHeadingDeg, f))
}

type Summary struct {
	HalfPowerDepthM      float64
	QuarterTurnDepthM    float64
	FullTurnDepthM       float64
	SurfaceDecayFraction float64
	VeerRate             float64
}

func Summarize(r *Result) Summary {
	return Summary{
		HalfPowerDepthM:      HalfPowerDepth(r.K, r.F),
		QuarterTurnDepthM:    QuarterTurnDepth(r.K, r.F),
		FullTurnDepthM:       FullTurnDepth(r.K, r.F),
		SurfaceDecayFraction: DecayFraction(r.EkmanDepth, r.K, r.F),
		VeerRate:             TurnPerMetre(r.K, r.F),
	}
}
