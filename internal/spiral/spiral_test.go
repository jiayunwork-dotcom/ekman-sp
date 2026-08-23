package spiral

import (
	"math"
	"testing"
)

// testTol is the relative tolerance used across the physics tests.
const testTol = 1e-9

// near compares x with the expected value to a relative tolerance.
func near(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol*math.Max(1, math.Abs(want)) {
		t.Errorf("%s = %g, want %g", name, got, want)
	}
}

// midlatParams returns parameters representative of a mid-latitude westerly
// wind case (easterly wind at 45N).
func midlatParams() Params {
	return Params{
		Tau:         0.2,
		WindHeading: 90,
		Rho:         1025,
		F:           1.0313e-4,
		K:           0.05,
		NPoints:     21,
	}
}

func TestSurfaceSpeedMagnitude(t *testing.T) {
	p := midlatParams()
	// |V0| = tau / (rho * sqrt(K*|f|)); for this case the closed form is
	// about 0.0859 m/s.
	got, err := SurfaceSpeed(p.Tau, p.Rho, p.K, p.F)
	if err != nil {
		t.Fatalf("SurfaceSpeed returned error: %v", err)
	}
	near(t, "V0", got, 0.085923, 1e-3)
}

func TestSurfaceVelocityMatchesMagnitude(t *testing.T) {
	p := midlatParams()
	vec, _, err := SurfaceVelocity(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SurfaceVelocity returned error: %v", err)
	}
	near(t, "|vec|", vec.Len(), p.Tau/(p.Rho*math.Sqrt(p.K*math.Abs(p.F))), testTol)
}

func TestSurfaceHeadingRight45North(t *testing.T) {
	// Northern hemisphere: the surface current sits 45 degrees to the right
	// of the wind. An easterly wind (heading 90) yields heading 135.
	p := midlatParams()
	heading := SurfaceHeading(p.WindHeading, p.F)
	near(t, "surface heading", heading, 135, testTol)
	if got := SurfaceDeflection(p.F); got != 45 {
		t.Errorf("SurfaceDeflection(north) = %g, want 45", got)
	}
}

func TestSurfaceHeadingLeft45South(t *testing.T) {
	// Southern hemisphere: the surface current deflects to the left, so the
	// same easterly wind yields heading 45 (northeast).
	p := midlatParams()
	p.F = -p.F
	heading := SurfaceHeading(p.WindHeading, p.F)
	near(t, "surface heading", heading, 45, testTol)
	if got := SurfaceDeflection(p.F); got != -45 {
		t.Errorf("SurfaceDeflection(south) = %g, want -45", got)
	}
}

func TestEkmanDepthValue(t *testing.T) {
	p := midlatParams()
	// De = pi*sqrt(2K/|f|); for this case about 97.8 m. A formula that drops
	// the factor 2 or the pi is caught by this check.
	got := EkmanDepth(p.K, p.F)
	near(t, "De", got, 97.823, 1e-3)
}

func TestEkmanDepthScalingFDoubling(t *testing.T) {
	p := midlatParams()
	base := EkmanDepth(p.K, p.F)
	half := EkmanDepth(p.K, 2*p.F)
	// Doubling |f| shrinks De by 1/sqrt(2).
	near(t, "De(f*2)/De", half/base, 1/math.Sqrt2, testTol)
}

func TestEkmanDepthScalingKDoubling(t *testing.T) {
	p := midlatParams()
	base := EkmanDepth(p.K, p.F)
	double := EkmanDepth(2*p.K, p.F)
	// Doubling K grows De by sqrt(2).
	near(t, "De(2K)/De", double/base, math.Sqrt2, testTol)
}

func TestTransportMagnitude(t *testing.T) {
	p := midlatParams()
	// Me = tau / (rho*|f|) about 1.89 m^2/s. This is the test that separates
	// tau/(rho*|f|) from the wrong tau/(rho*sqrt(|f|)) or tau/(rho*sqrt(K|f|)).
	got, err := Transport(p.Tau, p.Rho, p.F)
	if err != nil {
		t.Fatalf("Transport returned error: %v", err)
	}
	near(t, "Me", got, 1.89198, 1e-3)
}

func TestTransportHeadingRight90North(t *testing.T) {
	p := midlatParams()
	// Transport is 90 degrees to the right of the wind in the northern
	// hemisphere: easterly wind -> southward transport (heading 180).
	heading := TransportHeading(p.WindHeading, p.F)
	near(t, "transport heading", heading, 180, testTol)
	if got := TransportDeflection(p.F); got != 90 {
		t.Errorf("TransportDeflection(north) = %g, want 90", got)
	}
}

func TestTransportHeadingLeft90South(t *testing.T) {
	p := midlatParams()
	p.F = -p.F
	// Southern hemisphere: transport 90 degrees to the left of the wind;
	// easterly wind -> northward transport (heading 0).
	heading := TransportHeading(p.WindHeading, p.F)
	near(t, "transport heading", heading, 0, testTol)
	if got := TransportDeflection(p.F); got != -90 {
		t.Errorf("TransportDeflection(south) = %g, want -90", got)
	}
}

func TestTauDoublingScalesV0AndTransport(t *testing.T) {
	p := midlatParams()
	v0a, _ := SurfaceSpeed(p.Tau, p.Rho, p.K, p.F)
	meA, _ := Transport(p.Tau, p.Rho, p.F)

	p2 := p
	p2.Tau = 2 * p.Tau
	v0b, _ := SurfaceSpeed(p2.Tau, p2.Rho, p2.K, p2.F)
	meB, _ := Transport(p2.Tau, p2.Rho, p2.F)

	// Doubling tau doubles both the surface speed and the transport.
	near(t, "V0(tau*2)/V0", v0b/v0a, 2, testTol)
	near(t, "Me(tau*2)/Me", meB/meA, 2, testTol)
}

func TestFDoublingReducesV0(t *testing.T) {
	p := midlatParams()
	v0a, _ := SurfaceSpeed(p.Tau, p.Rho, p.K, p.F)
	v0b, _ := SurfaceSpeed(p.Tau, p.Rho, p.K, 2*p.F)
	// Doubling |f| reduces |V0| by 1/sqrt(2).
	near(t, "V0(2f)/V0", v0b/v0a, 1/math.Sqrt2, testTol)
}

func TestKDoublingReducesV0(t *testing.T) {
	p := midlatParams()
	v0a, _ := SurfaceSpeed(p.Tau, p.Rho, p.K, p.F)
	v0b, _ := SurfaceSpeed(p.Tau, p.Rho, 2*p.K, p.F)
	// Doubling K reduces |V0| by 1/sqrt(2).
	near(t, "V0(2K)/V0", v0b/v0a, 1/math.Sqrt2, testTol)
}

func TestWindReversalReversesSpiral(t *testing.T) {
	p := midlatParams()
	forward, err := SampleProfile(9, 0, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SampleProfile forward error: %v", err)
	}
	reversed, err := SampleProfile(9, 0, p.Tau, p.Rho, p.K, p.F, p.WindHeading+180)
	if err != nil {
		t.Fatalf("SampleProfile reversed error: %v", err)
	}
	// Reversing the wind must reverse every velocity vector.
	for i := range forward {
		wantU, wantV := -forward[i].U, -forward[i].V
		if math.Abs(reversed[i].U-wantU) > 1e-9 || math.Abs(reversed[i].V-wantV) > 1e-9 {
			t.Errorf("point %d: reversed wind gives (%.9g, %.9g), want (%.9g, %.9g)",
				i, reversed[i].U, reversed[i].V, wantU, wantV)
		}
	}
}

func TestHemisphereFlipSwapsDeflection(t *testing.T) {
	p := midlatParams()
	north, _, _ := SurfaceVelocity(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	south, _, _ := SurfaceVelocity(p.Tau, p.Rho, p.K, -p.F, p.WindHeading)
	// Sign of f flips the lateral component: the southern surface current is
	// the mirror of the northern one about the wind axis.
	if math.Abs(north.East-south.East) > 1e-9 {
		t.Errorf("east component should be unchanged, got %g vs %g", north.East, south.East)
	}
	if math.Abs(north.North+south.North) > 1e-9 {
		t.Errorf("north component should flip sign, got %g vs %g", north.North, south.North)
	}
}

func TestProfileRotatesClockwiseNorth(t *testing.T) {
	p := midlatParams()
	points, err := SampleProfile(21, 0, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SampleProfile error: %v", err)
	}
	// Northern hemisphere: veer (cumulative clockwise turn) grows with depth.
	for i := 1; i < len(points); i++ {
		if points[i].Veer <= points[i-1].Veer {
			t.Fatalf("veer not monotone at depth %g: %g -> %g",
				points[i].Depth, points[i-1].Veer, points[i].Veer)
		}
	}
	if points[0].Veer != 0 {
		t.Errorf("surface veer = %g, want 0", points[0].Veer)
	}
}

func TestProfileRotatesCounterClockwiseSouth(t *testing.T) {
	p := midlatParams()
	p.F = -p.F
	points, err := SampleProfile(21, 0, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("SampleProfile error: %v", err)
	}
	// Southern hemisphere: veer is negative and decreases with depth.
	for i := 1; i < len(points); i++ {
		if points[i].Veer >= points[i-1].Veer {
			t.Fatalf("veer not decreasing at depth %g: %g -> %g",
				points[i].Depth, points[i-1].Veer, points[i].Veer)
		}
	}
}

func TestProfileDecaysExponentially(t *testing.T) {
	p := midlatParams()
	v0, err := SurfaceSpeed(p.Tau, p.Rho, p.K, p.F)
	if err != nil {
		t.Fatalf("SurfaceSpeed error: %v", err)
	}
	delta := EkmanScale(p.K, p.F)
	point, err := PointAtDepth(delta, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		t.Fatalf("PointAtDepth error: %v", err)
	}
	// One e-folding scale down, the speed has decayed to V0/e.
	near(t, "speed at delta", point.Speed, v0/math.E, 1e-6)
}

func TestTransportIntegral(t *testing.T) {
	// The depth integral of the spiral must reproduce the analytic transport
	// vector. A coarse profile with a couple of dozen points converges to
	// within ~1.5% of the transport magnitude.
	p := midlatParams()
	res, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	integral := IntegrateProfile(res.Points)
	diff := integral.Sub(res.TransportVector)
	if diff.Len() > 0.02+1e-2*res.Transport {
		t.Errorf("integrated transport %+v, want %+v (error %g m^2/s)",
			integral, res.TransportVector, diff.Len())
	}
}

func TestComputeSurfaceConsistency(t *testing.T) {
	p := midlatParams()
	res, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	near(t, "surface speed", res.SurfaceVelocity.Len(), res.V0, testTol)
	near(t, "surface heading", res.SurfaceHeading, 135, testTol)
	near(t, "ekman depth", res.EkmanDepth, math.Pi*math.Sqrt(2*p.K/math.Abs(p.F)), 1e-9)
	near(t, "transport magnitude", res.Transport, p.Tau/(p.Rho*math.Abs(p.F)), 1e-9)
}

func TestSurfaceSpeedErrors(t *testing.T) {
	p := midlatParams()
	if _, err := SurfaceSpeed(p.Tau, 0, p.K, p.F); err == nil {
		t.Error("zero density must error")
	}
	if _, err := SurfaceSpeed(p.Tau, p.Rho, 0, p.F); err == nil {
		t.Error("zero viscosity must error")
	}
	if _, err := SurfaceSpeed(p.Tau, p.Rho, p.K, 0); err == nil {
		t.Error("zero Coriolis must error")
	}
	if _, err := SurfaceSpeed(-0.1, p.Rho, p.K, p.F); err == nil {
		t.Error("negative tau must error")
	}
}
