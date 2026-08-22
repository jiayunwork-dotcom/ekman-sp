package spiral

import "math"

// Depth scales and rotation rates derived from the Ekman scale delta. They
// are closed-form consequences of the same solution and give the report (and
// the tests) concrete numbers to check the spiral against.

// HalfPowerDepth returns the depth at which the current speed has decayed to
// half its surface value: delta*ln(2).
func HalfPowerDepth(K, f float64) float64 {
	return EkmanScale(K, f) * math.Ln2
}

// QuarterTurnDepth returns the depth at which the current has turned one
// quarter of a circle (90 degrees): delta*pi/2.
func QuarterTurnDepth(K, f float64) float64 {
	return EkmanScale(K, f) * math.Pi / 2.0
}

// FullTurnDepth returns the depth at which the current has completed one full
// revolution of the spiral: delta*2*pi.
func FullTurnDepth(K, f float64) float64 {
	return EkmanScale(K, f) * 2.0 * math.Pi
}

// DecayFraction returns the fraction of the surface amplitude that has been
// lost at depthM: 1 - exp(-depthM/delta). It rises from 0 at the surface and
// approaches 1 exponentially.
func DecayFraction(depthM, K, f float64) float64 {
	return 1.0 - math.Exp(-depthM/EkmanScale(K, f))
}

// TurnPerMetre returns the rate at which the current veers with depth, in
// degrees per metre. It is positive (clockwise) in the northern hemisphere
// and negative (counter-clockwise) in the southern hemisphere; its magnitude
// is 180/(pi*delta).
func TurnPerMetre(K, f float64) float64 {
	return SignOf(f) * RadToDeg(1.0/EkmanScale(K, f))
}

// SurfaceFlowAngle returns the phase (radians counter-clockwise from east) of
// the surface current for the given wind heading.
func SurfaceFlowAngle(windHeadingDeg, f float64) float64 {
	return HeadingToPhase(SurfaceHeading(windHeadingDeg, f))
}

// Summary packs the handful of derived quantities a report wants to state
// about the spiral beyond the profile itself.
type Summary struct {
	// HalfPowerDepthM is where the speed halves.
	HalfPowerDepthM float64
	// QuarterTurnDepthM is where the current has turned 90 degrees.
	QuarterTurnDepthM float64
	// FullTurnDepthM is where the current has turned 360 degrees.
	FullTurnDepthM float64
	// SurfaceDecayFraction is the fraction of the amplitude gone at De.
	SurfaceDecayFraction float64
	// VeerRate is the turning rate in degrees per metre.
	VeerRate float64
}

// Summarize derives the summary quantities for a computed result.
func Summarize(r *Result) Summary {
	return Summary{
		HalfPowerDepthM:      HalfPowerDepth(r.K, r.F),
		QuarterTurnDepthM:    QuarterTurnDepth(r.K, r.F),
		FullTurnDepthM:       FullTurnDepth(r.K, r.F),
		SurfaceDecayFraction: DecayFraction(r.EkmanDepth, r.K, r.F),
		VeerRate:             TurnPerMetre(r.K, r.F),
	}
}
