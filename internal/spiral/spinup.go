package spiral

import (
	"fmt"
	"math"
)

func SpinUpTime(depthM, K, f float64) (float64, error) {
	if depthM <= 0 {
		return 0, fmt.Errorf("spiral: basin depth must be positive")
	}
	if K <= 0 {
		return 0, fmt.Errorf("spiral: eddy viscosity must be positive")
	}
	if math.Abs(f) <= HemisphereTolerance {
		return 0, fmt.Errorf("spiral: spin-up undefined at equator")
	}
	return depthM / math.Sqrt(2*K*math.Abs(f)), nil
}

func VerticalDiffusiveTime(depthM, K float64) (float64, error) {
	if depthM <= 0 || K <= 0 {
		return 0, fmt.Errorf("spiral: depth and K must be positive")
	}
	return depthM * depthM / K, nil
}

func SpinUpFasterThanDiffusion(depthM, K, f float64) (bool, error) {
	ts, err := SpinUpTime(depthM, K, f)
	if err != nil {
		return false, err
	}
	td, err := VerticalDiffusiveTime(depthM, K)
	if err != nil {
		return false, err
	}
	return ts < td, nil
}

func BottomLayerTransport(tau, rho, f float64) (float64, error) {
	mag, err := Transport(tau, rho, f)
	if err != nil {
		return 0, err
	}
	return -mag, nil
}

func InteriorGeostrophicBalance(bottom, surface float64) error {
	if math.Abs(bottom+surface) > 1e-12*(math.Abs(bottom)+math.Abs(surface)+1e-18) {
		return fmt.Errorf("spiral: bottom transport %g does not cancel surface %g", bottom, surface)
	}
	return nil
}

func DecayEFoldDepth(K, f float64) float64 {
	return EkmanScale(K, f)
}

func TurnsToEFold() float64 {
	return 1.0 / math.Pi
}
