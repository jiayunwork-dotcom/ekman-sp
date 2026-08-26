package spiral

import (
	"fmt"
	"math"
)

func InertialPeriod(f float64) (float64, error) {
	if math.Abs(f) <= HemisphereTolerance {
		return 0, fmt.Errorf("spiral: inertial period undefined at equator")
	}
	return 2 * math.Pi / math.Abs(f), nil
}

func InertialFrequency(f float64) (float64, error) {
	t, err := InertialPeriod(f)
	if err != nil {
		return 0, err
	}
	return 1 / t, nil
}

func EkmanPumping(rho, f, dTauXdy, dTauYdx float64) (float64, error) {
	if rho <= 0 {
		return 0, fmt.Errorf("spiral: density must be positive")
	}
	if math.Abs(f) <= HemisphereTolerance {
		return 0, fmt.Errorf("spiral: pumping undefined at equator")
	}
	curlZ := dTauYdx - dTauXdy
	return curlZ / (rho * f), nil
}

func TransportDivergence(rho, f, dTxdx, dTydy float64) (float64, error) {
	if rho <= 0 {
		return 0, fmt.Errorf("spiral: density must be positive")
	}
	if math.Abs(f) <= HemisphereTolerance {
		return 0, fmt.Errorf("spiral: divergence undefined at equator")
	}
	return (dTxdx + dTydy) / (rho * math.Abs(f)), nil
}

func PumpingFromUniformCurl(rho, f, curlTau float64) (float64, error) {
	return EkmanPumping(rho, f, 0, curlTau)
}

func HemisphereFlipsPumping(rho, f, curlTau float64) error {
	wn, err := EkmanPumping(rho, f, 0, curlTau)
	if err != nil {
		return err
	}
	ws, err := EkmanPumping(rho, -f, 0, curlTau)
	if err != nil {
		return err
	}
	if math.Abs(wn+ws) > 1e-12*(math.Abs(wn)+math.Abs(ws)+1e-18) {
		return fmt.Errorf("spiral: pumping did not change sign across equator: %g vs %g", wn, ws)
	}
	return nil
}
