package spiral

import "testing"

func TestAnalyzeSensitivity(t *testing.T) {
	p := Params{Tau: 0.2, WindHeading: 90, Rho: 1025, F: 1e-4, K: 0.05, NPoints: 5}
	sens, err := AnalyzeSensitivity(p)
	if err != nil {
		t.Fatal(err)
	}
	if sens.DTauV0 <= 0 {
		t.Fatalf("dTauV0=%v", sens.DTauV0)
	}
	for _, c := range sens.ScalingChecks {
		if !c.Pass {
			t.Fatalf("scaling %s failed ratio=%v", c.Name, c.Ratio)
		}
	}
}

func TestSweepK(t *testing.T) {
	p := Params{Tau: 0.2, WindHeading: 90, Rho: 1025, F: 1e-4, K: 0.05, NPoints: 3}
	sweep, err := SweepK(p, 0.03, 0.08, 0.02)
	if err != nil {
		t.Fatal(err)
	}
	if len(sweep.Points) < 2 {
		t.Fatalf("points=%d", len(sweep.Points))
	}
}

func TestCompareHemispheres(t *testing.T) {
	p := Params{Tau: 0.2, WindHeading: 90, Rho: 1025, F: 1e-4, K: 0.05, NPoints: 5}
	cmp, err := CompareHemispheres(p)
	if err != nil {
		t.Fatal(err)
	}
	if cmp.Left.SurfaceVeer*cmp.Right.SurfaceVeer >= 0 {
		t.Fatalf("veers %v %v", cmp.Left.SurfaceVeer, cmp.Right.SurfaceVeer)
	}
}

func TestSweepInvalidParam(t *testing.T) {
	p := Params{Tau: 0.2, WindHeading: 90, Rho: 1025, F: 1e-4, K: 0.05, NPoints: 3}
	_, err := SweepParameter(p, "bad", 1, 2, 0.1)
	if err != ErrUnknownParameter {
		t.Fatalf("err=%v", err)
	}
}
