package spiral

import (
	"errors"
	"fmt"
)

var (
	ErrDensityNonPositive = errors.New("density rho must be positive")

	ErrViscosityNonPositive = errors.New("eddy viscosity K must be positive")

	ErrEquator = errors.New("Coriolis parameter f must be non-zero (Ekman spiral is undefined at the equator)")

	ErrTauNegative = errors.New("wind stress tau must be positive")

	ErrInvalidStep = errors.New("sweep step must be positive")

	ErrUnknownParameter = errors.New("unknown sweep parameter")

	ErrEmptySweep = errors.New("sweep produced no valid points")
)

type inputError struct {
	field string
	value float64
	err   error
}

func (e *inputError) Error() string {
	return fmt.Sprintf("%s (got %g)", e.err, e.value)
}

func (e *inputError) Unwrap() error { return e.err }

func checkInput(rho, f, K, tau float64) error {
	if rho <= 0 {
		return &inputError{field: "rho", value: rho, err: ErrDensityNonPositive}
	}
	if K <= 0 {
		return &inputError{field: "K", value: K, err: ErrViscosityNonPositive}
	}
	if tau < 0 {
		return &inputError{field: "tau", value: tau, err: ErrTauNegative}
	}
	if ClassifyHemisphere(f) == HemisphereUnknown {
		return &inputError{field: "f", value: f, err: ErrEquator}
	}
	return nil
}
