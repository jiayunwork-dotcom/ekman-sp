package config

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func sampleInput() Input {
	rho := 1025.0
	K := 0.05
	lat := 45.0
	return Input{
		Tau:         0.2,
		WindDirDeg:  float64Ptr(90),
		Rho:         &rho,
		K:           &K,
		LatitudeDeg: &lat,
	}
}

func float64Ptr(v float64) *float64 { return &v }
func intPtr(v int) *int             { return &v }

func mustResolve(t *testing.T, in Input) *Resolved {
	t.Helper()
	r, err := Resolve(in)
	if err != nil {
		t.Fatalf("Resolve returned unexpected error: %v", err)
	}
	return r
}

func resolveError(t *testing.T, in Input, want error) {
	t.Helper()
	_, err := Resolve(in)
	if err == nil {
		t.Fatalf("Resolve succeeded, want error %v", want)
	}
	if !errors.Is(err, want) {
		t.Errorf("error = %v, want %v", err, want)
	}
}

func TestValidateDensityZero(t *testing.T) {
	in := sampleInput()
	*in.Rho = 0
	resolveError(t, in, ErrRhoNonPositive)
}

func TestValidateDensityNegative(t *testing.T) {
	in := sampleInput()
	*in.Rho = -500
	resolveError(t, in, ErrRhoNonPositive)
}

func TestValidateViscosityZero(t *testing.T) {
	in := sampleInput()
	*in.K = 0
	resolveError(t, in, ErrKNonPositive)
}

func TestValidateViscosityNegative(t *testing.T) {
	in := sampleInput()
	*in.K = -0.01
	resolveError(t, in, ErrKNonPositive)
}

func TestValidateEquatorLatitudeZero(t *testing.T) {
	in := sampleInput()
	zero := 0.0
	in.LatitudeDeg = &zero
	resolveError(t, in, ErrEquatorial)
}

func TestValidateEquatorFZero(t *testing.T) {
	in := sampleInput()
	zero := 0.0
	in.F = &zero
	in.LatitudeDeg = nil
	resolveError(t, in, ErrEquatorial)
}

func TestValidateMissingCoriolis(t *testing.T) {
	in := sampleInput()
	in.F = nil
	in.LatitudeDeg = nil
	resolveError(t, in, ErrNoCoriolis)
}

func TestValidateTauMissing(t *testing.T) {
	in := sampleInput()
	in.Tau = 0
	resolveError(t, in, ErrTauRequired)
}

func TestValidateTauNegative(t *testing.T) {
	in := sampleInput()
	in.Tau = -0.2
	resolveError(t, in, ErrTauRequired)
}

func TestValidateNPointsBad(t *testing.T) {
	in := sampleInput()
	in.NPoints = intPtr(0)
	resolveError(t, in, ErrBadNPoints)
}

func TestValidateShallowWater(t *testing.T) {
	in := sampleInput()
	depth := 20.0
	in.DepthM = &depth
	resolveError(t, in, ErrShallowWater)
}

func TestValidateDeepWaterOK(t *testing.T) {
	in := sampleInput()
	depth := 500.0
	in.DepthM = &depth
	r := mustResolve(t, in)
	if !r.FiniteDepth {
		t.Error("FiniteDepth = false, want true")
	}
	if r.Params.ProfileDepthM != depth {
		t.Errorf("ProfileDepthM = %g, want %g", r.Params.ProfileDepthM, depth)
	}
}

func TestResolveFromLatitude(t *testing.T) {
	in := sampleInput()
	r := mustResolve(t, in)
	wantF := CoriolisFromLatitude(45.0)
	if math.Abs(r.Params.F-wantF) > 1e-12 {
		t.Errorf("f = %g, want %g", r.Params.F, wantF)
	}
	if r.Params.F <= 0 {
		t.Errorf("northern latitude must give positive f, got %g", r.Params.F)
	}
	if math.Abs(r.LatitudeDeg-45.0) > 1e-12 {
		t.Errorf("LatitudeDeg = %g, want 45", r.LatitudeDeg)
	}
}

func TestResolveFromExplicitF(t *testing.T) {
	in := sampleInput()
	f := 1.0313e-4
	in.F = &f
	in.LatitudeDeg = nil
	r := mustResolve(t, in)
	if r.Params.F != f {
		t.Errorf("f = %g, want %g", r.Params.F, f)
	}
	if !math.IsNaN(r.LatitudeDeg) {
		t.Errorf("LatitudeDeg = %g, want NaN when f given directly", r.LatitudeDeg)
	}
}

func TestResolveCoriolisConflict(t *testing.T) {
	in := sampleInput()
	f := 2.0e-4
	in.F = &f
	resolveError(t, in, ErrCoriolisConflict)
}

func TestResolveCoriolisConsistent(t *testing.T) {
	in := sampleInput()
	f := CoriolisFromLatitude(45.0)
	in.F = &f
	r := mustResolve(t, in)
	if r.Params.F != f {
		t.Errorf("f = %g, want %g", r.Params.F, f)
	}
}

func TestLoadBytesDefaults(t *testing.T) {
	body := `{"tau":0.15,"latitude_deg":30}`
	r, err := LoadBytes([]byte(body))
	if err != nil {
		t.Fatalf("LoadBytes error: %v", err)
	}
	if r.Params.Rho != defaults.rho {
		t.Errorf("rho = %g, want default %g", r.Params.Rho, defaults.rho)
	}
	if r.Params.K != defaults.K {
		t.Errorf("K = %g, want default %g", r.Params.K, defaults.K)
	}
	if r.Params.WindHeading != defaults.windDirDeg {
		t.Errorf("wind heading = %g, want default %g", r.Params.WindHeading, defaults.windDirDeg)
	}
	if r.Params.NPoints != defaults.nPoints {
		t.Errorf("npoints = %d, want default %d", r.Params.NPoints, defaults.nPoints)
	}
}

func TestLoadUnknownFieldRejected(t *testing.T) {
	body := `{"tau":0.15,"latitude_deg":30,"bogus_field":1}`
	if _, err := LoadBytes([]byte(body)); err == nil {
		t.Fatal("unknown JSON field must be rejected")
	} else if !errors.Is(err, ErrParse) {
		t.Errorf("error = %v, want ErrParse", err)
	}
}

func TestLoadMalformedJSON(t *testing.T) {
	if _, err := LoadBytes([]byte(`{"tau":`)); err == nil {
		t.Fatal("malformed JSON must error")
	}
}

func TestCoriolisFromLatitudeZero(t *testing.T) {
	if got := CoriolisFromLatitude(0); got != 0 {
		t.Errorf("f(0 deg) = %g, want 0", got)
	}
}

func TestCoriolisFromLatitudeSouth(t *testing.T) {
	if got := CoriolisFromLatitude(-45); got >= 0 {
		t.Errorf("f(-45 deg) = %g, want negative", got)
	}
}

func TestWindHeadingDefault(t *testing.T) {
	in := sampleInput()
	in.WindDirDeg = nil
	if got := WindHeading(in); got != defaults.windDirDeg {
		t.Errorf("WindHeading = %g, want default %g", got, defaults.windDirDeg)
	}
}

func TestWindHeadingExplicit(t *testing.T) {
	in := sampleInput()
	in.WindDirDeg = float64Ptr(225)
	if got := WindHeading(in); got != 225 {
		t.Errorf("WindHeading = %g, want 225", got)
	}
}

func TestLoadFileMissing(t *testing.T) {
	if _, err := Load("/nonexistent/ekman-case.json"); err == nil {
		t.Fatal("missing file must error")
	} else if !errors.Is(err, ErrFileRead) {
		t.Errorf("error = %v, want ErrFileRead", err)
	}
}

func TestErrorMessagesUserVisible(t *testing.T) {
	badRho := sampleInput()
	*badRho.Rho = 0
	_, err := Resolve(badRho)
	if err == nil || !strings.Contains(err.Error(), "rho") {
		t.Errorf("zero-density error = %v, want message mentioning rho", err)
	}

	badK := sampleInput()
	*badK.K = -1
	_, err = Resolve(badK)
	if err == nil || !strings.Contains(err.Error(), "K") {
		t.Errorf("negative-K error = %v, want message mentioning K", err)
	}

	equator := sampleInput()
	zero := 0.0
	equator.LatitudeDeg = &zero
	_, err = Resolve(equator)
	if err == nil || !strings.Contains(err.Error(), "equator") {
		t.Errorf("equator error = %v, want message mentioning equator", err)
	}
}
