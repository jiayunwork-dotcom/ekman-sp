package spiral

import (
	"testing"
)

func sampleParams() Params {
	return Params{
		Tau:           0.2,
		WindHeading:   90,
		Rho:           1025,
		F:             1e-4,
		K:             0.05,
		NPoints:       21,
		ProfileDepthM: 50,
	}
}

func TestComputeEnvelope(t *testing.T) {
	res, err := Compute(sampleParams())
	if err != nil {
		t.Fatal(err)
	}
	env := ComputeEnvelope(res.Points)
	if env.MaxDepth <= env.MinDepth || env.MaxSpeed < env.MinSpeed {
		t.Fatalf("envelope = %+v", env)
	}
}

func TestMonotonicSpeedDecay(t *testing.T) {
	res, err := Compute(sampleParams())
	if err != nil {
		t.Fatal(err)
	}
	if !MonotonicSpeedDecay(res.Points) {
		t.Fatal("expected monotonic speed decay")
	}
}

func TestTauSensitivity(t *testing.T) {
	r, err := TauSensitivity(sampleParams(), 0.05)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Items) != 2 || MaxAbsV0Delta(r) <= 0 {
		t.Fatalf("report = %+v", r)
	}
}

func TestDepthAtHalfSpeed(t *testing.T) {
	res, err := Compute(sampleParams())
	if err != nil {
		t.Fatal(err)
	}
	depth, err := DepthAtHalfSpeed(res.Points)
	if err != nil || depth <= 0 {
		t.Fatalf("DepthAtHalfSpeed() = %v %v", depth, err)
	}
}
