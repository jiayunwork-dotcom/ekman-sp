package spiral

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestCurrentComplexMatchesPoint(t *testing.T) {
	p := midlatParams()
	w, err := CurrentComplex(12.5, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("CurrentComplex error: %v", err)
	}
	pt, err := PointAtDepth(12.5, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("PointAtDepth error: %v", err)
	}
	if math.Abs(real(w)-pt.U) > 1e-12 || math.Abs(imag(w)-pt.V) > 1e-12 {
		t.Errorf("complex (%g,%g) vs point (%g,%g)", real(w), imag(w), pt.U, pt.V)
	}
}

func TestSurfaceComplexMatchesSurfaceVelocity(t *testing.T) {
	p := midlatParams()
	w, err := SurfaceComplex(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SurfaceComplex error: %v", err)
	}
	vec, _, err := SurfaceVelocity(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SurfaceVelocity error: %v", err)
	}
	if math.Abs(real(w)-vec.East) > 1e-12 || math.Abs(imag(w)-vec.North) > 1e-12 {
		t.Errorf("complex (%g,%g) vs vector (%g,%g)", real(w), imag(w), vec.East, vec.North)
	}
}

func TestTransportComplexMatchesVector(t *testing.T) {
	p := midlatParams()
	w, err := TransportComplex(p.Tau, p.Rho, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("TransportComplex error: %v", err)
	}
	vec, _, err := TransportVector(p.Tau, p.Rho, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("TransportVector error: %v", err)
	}
	if math.Abs(real(w)-vec.East) > 1e-12 || math.Abs(imag(w)-vec.North) > 1e-12 {
		t.Errorf("complex (%g,%g) vs vector (%g,%g)", real(w), imag(w), vec.East, vec.North)
	}
}

func TestFiniteDepthTendsToInfinite(t *testing.T) {
	p := midlatParams()
	inf, err := SurfaceComplex(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SurfaceComplex error: %v", err)
	}
	// A water depth of many De is indistinguishable from infinite depth.
	fin, err := FiniteDepthSurfaceCurrent(20*EkmanDepth(p.K, p.F), p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("FiniteDepthSurfaceCurrent error: %v", err)
	}
	if cmplx.Abs(fin-inf) > 1e-8 {
		t.Errorf("finite surface %g, infinite %g", cmplx.Abs(fin), cmplx.Abs(inf))
	}
}

func TestFiniteDepthShallowWeakened(t *testing.T) {
	p := midlatParams()
	inf, err := SurfaceComplex(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SurfaceComplex error: %v", err)
	}
	// Very shallow water (H well below the Ekman scale) constrains the whole
	// column with bottom friction: the surface current drops far below the
	// infinite-depth value (about 14% for H = 0.1*delta).
	H := 0.1 * EkmanScale(p.K, p.F)
	fin, err := FiniteDepthSurfaceCurrent(H, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("FiniteDepthSurfaceCurrent error: %v", err)
	}
	if cmplx.Abs(fin) > 0.5*cmplx.Abs(inf) {
		t.Errorf("very shallow surface |%g| should be well below infinite |%g|", cmplx.Abs(fin), cmplx.Abs(inf))
	}
	// The finite-depth profile is no longer the classic 45-degree spiral.
	def, err := FiniteDepthSurfaceDeflection(H, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("FiniteDepthSurfaceDeflection error: %v", err)
	}
	if math.Abs(def-45) < 5 {
		t.Errorf("shallow deflection %g should deviate clearly from the classic 45 deg", def)
	}
}

func TestFiniteDepthNoSlipBottom(t *testing.T) {
	p := midlatParams()
	H := 2.0 * EkmanDepth(p.K, p.F)
	w, err := FiniteDepthCurrent(H, H, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("FiniteDepthCurrent error: %v", err)
	}
	if cmplx.Abs(w) > 1e-9 {
		t.Errorf("bottom current |%g| should be near zero (no-slip)", cmplx.Abs(w))
	}
}

func TestFiniteDepthErrorRatioSmallForDeepWater(t *testing.T) {
	p := midlatParams()
	H := 5.0 * EkmanDepth(p.K, p.F)
	ratio, err := FiniteDepthErrorRatio(H, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("FiniteDepthErrorRatio error: %v", err)
	}
	if ratio > 1e-6 {
		t.Errorf("deep-water deviation %g, want < 1e-6", ratio)
	}
}

func TestFiniteDepthSurfaceDeflection(t *testing.T) {
	p := midlatParams()
	// Deep water returns the classic 45-degree deflection.
	deep := 20.0 * EkmanDepth(p.K, p.F)
	def, err := FiniteDepthSurfaceDeflection(deep, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("FiniteDepthSurfaceDeflection error: %v", err)
	}
	near(t, "deep deflection", def, 45, 1e-6)
}

func TestHalfPowerDepth(t *testing.T) {
	p := midlatParams()
	delta := EkmanScale(p.K, p.F)
	near(t, "half-power depth", HalfPowerDepth(p.K, p.F), delta*math.Ln2, 1e-12)
}

func TestQuarterTurnDepth(t *testing.T) {
	p := midlatParams()
	delta := EkmanScale(p.K, p.F)
	near(t, "quarter-turn depth", QuarterTurnDepth(p.K, p.F), delta*math.Pi/2, 1e-12)
}

func TestFullTurnDepth(t *testing.T) {
	p := midlatParams()
	delta := EkmanScale(p.K, p.F)
	near(t, "full-turn depth", FullTurnDepth(p.K, p.F), delta*2*math.Pi, 1e-12)
}

func TestDecayFraction(t *testing.T) {
	p := midlatParams()
	delta := EkmanScale(p.K, p.F)
	near(t, "decay at delta", DecayFraction(delta, p.K, p.F), 1-1/math.E, 1e-12)
	if got := DecayFraction(0, p.K, p.F); got != 0 {
		t.Errorf("decay at surface = %g, want 0", got)
	}
}

func TestTurnPerMetre(t *testing.T) {
	p := midlatParams()
	rate := TurnPerMetre(p.K, p.F)
	// Rate must be positive (clockwise) in the north and its magnitude is the
	// phase change per metre.
	near(t, "veer rate north", rate, RadToDeg(1.0/EkmanScale(p.K, p.F)), 1e-9)
	south := TurnPerMetre(p.K, -p.F)
	near(t, "veer rate south", south, -RadToDeg(1.0/EkmanScale(p.K, p.F)), 1e-9)
}

func TestSummarize(t *testing.T) {
	p := midlatParams()
	res, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	s := Summarize(&res)
	delta := EkmanScale(p.K, p.F)
	near(t, "summary half-power", s.HalfPowerDepthM, delta*math.Ln2, 1e-9)
	near(t, "summary quarter-turn", s.QuarterTurnDepthM, delta*math.Pi/2, 1e-9)
	near(t, "summary full-turn", s.FullTurnDepthM, delta*2*math.Pi, 1e-9)
}

func TestWindStressComplex(t *testing.T) {
	w := WindStressComplex(0.2, 90)
	// Easterly wind: unit vector pointing east.
	if math.Abs(real(w)-0.2) > 1e-12 || math.Abs(imag(w)) > 1e-12 {
		t.Errorf("easterly stress = (%g, %g), want (0.2, 0)", real(w), imag(w))
	}
}
