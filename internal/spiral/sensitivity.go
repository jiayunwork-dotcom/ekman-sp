package spiral

import "math"

type ScalingCheck struct {
	Name     string  `json:"name"`
	Expected float64 `json:"expected"`
	Actual   float64 `json:"actual"`
	Ratio    float64 `json:"ratio"`
	Pass     bool    `json:"pass"`
}

type Sensitivity struct {
	Base          Result         `json:"base"`
	DTauV0        float64        `json:"d_tau_v0"`
	DKV0          float64        `json:"d_k_v0"`
	DFV0          float64        `json:"d_f_v0"`
	DTauTransport float64        `json:"d_tau_transport"`
	DFTransport   float64        `json:"d_f_transport"`
	DKEkmanDepth  float64        `json:"d_k_ekman_depth"`
	ScalingChecks []ScalingCheck `json:"scaling_checks"`
}

func AnalyzeSensitivity(p Params) (Sensitivity, error) {
	base, err := Compute(p)
	if err != nil {
		return Sensitivity{}, err
	}
	tauUp := p
	tauUp.Tau = p.Tau * 2
	tauRes, err := Compute(tauUp)
	if err != nil {
		return Sensitivity{}, err
	}
	kUp := p
	kUp.K = p.K * 2
	kRes, err := Compute(kUp)
	if err != nil {
		return Sensitivity{}, err
	}
	fUp := p
	fUp.F = p.F * 2
	if fUp.F == 0 {
		fUp.F = 1e-4
	}
	fRes, err := Compute(fUp)
	if err != nil {
		return Sensitivity{}, err
	}
	dTauV0 := (tauRes.V0 - base.V0) / p.Tau
	dKV0 := (kRes.V0 - base.V0) / p.K
	dFV0 := (fRes.V0 - base.V0) / absF(fUp.F-p.F)
	dTauTransport := (tauRes.Transport - base.Transport) / p.Tau
	dFTransport := (fRes.Transport - base.Transport) / absF(fUp.F-p.F)
	dKEkman := (kRes.EkmanDepth - base.EkmanDepth) / p.K
	checks := buildScalingChecks(p, base, tauRes, kRes, fRes)
	return Sensitivity{
		Base:          base,
		DTauV0:        dTauV0,
		DKV0:          dKV0,
		DFV0:          dFV0,
		DTauTransport: dTauTransport,
		DFTransport:   dFTransport,
		DKEkmanDepth:  dKEkman,
		ScalingChecks: checks,
	}, nil
}

func buildScalingChecks(p Params, base, tauRes, kRes, fRes Result) []ScalingCheck {
	tol := 0.02
	checks := []ScalingCheck{
		makeScalingCheck("tau_doubles_v0", 2.0, base.V0, tauRes.V0, tol),
		makeScalingCheck("tau_doubles_transport", 2.0, base.Transport, tauRes.Transport, tol),
		makeScalingCheck("k_doubles_ekman_depth", math.Sqrt2, base.EkmanDepth, kRes.EkmanDepth, tol),
		makeScalingCheck("f_doubles_reduces_v0", 1/math.Sqrt2, base.V0, fRes.V0, tol),
	}
	windRev := p
	windRev.WindHeading = math.Mod(p.WindHeading+180, 360)
	revRes, err := Compute(windRev)
	if err == nil {
		checks = append(checks, makeScalingCheck("wind_reversal_surface_u", -1.0, base.SurfaceVelocity.East, revRes.SurfaceVelocity.East, tol))
	}
	return checks
}

func makeScalingCheck(name string, expectedRatio, baseVal, newVal float64, tol float64) ScalingCheck {
	ratio := 0.0
	if baseVal != 0 {
		ratio = newVal / baseVal
	}
	pass := false
	if expectedRatio == 0 {
		pass = math.Abs(newVal) < 1e-9
	} else if baseVal != 0 {
		pass = math.Abs(ratio-expectedRatio) <= tol*math.Abs(expectedRatio)
	}
	return ScalingCheck{Name: name, Expected: expectedRatio, Actual: ratio, Ratio: ratio, Pass: pass}
}
