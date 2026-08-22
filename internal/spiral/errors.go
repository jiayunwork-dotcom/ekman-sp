package spiral

import (
	"errors"
	"fmt"
)

// Sentinel errors reported by the physics layer. They are mostly defensive:
// the config layer performs the same checks before any computation, but every
// formula here assumes a well-posed state and these errors keep that contract
// explicit rather than silently propagating NaN.
var (
	// ErrDensityNonPositive means rho <= 0, which would flip the stress-to-
	// velocity sign and is physically meaningless.
	ErrDensityNonPositive = errors.New("density rho must be positive")

	// ErrViscosityNonPositive means K <= 0, which removes the diffusion that
	// sustains the Ekman layer.
	ErrViscosityNonPositive = errors.New("eddy viscosity K must be positive")

	// ErrEquator means the Coriolis parameter is zero; the f-plane Ekman
	// solution does not exist at the equator.
	ErrEquator = errors.New("Coriolis parameter f must be non-zero (Ekman spiral is undefined at the equator)")

	// ErrTauNegative means the wind stress magnitude is negative.
	ErrTauNegative = errors.New("wind stress tau must be positive")
)

// inputError wraps a validation failure with the value that triggered it so
// callers can present a precise message ("density rho must be positive, got
// 0").
type inputError struct {
	field string
	value float64
	err   error
}

func (e *inputError) Error() string {
	return fmt.Sprintf("%s (got %g)", e.err, e.value)
}

func (e *inputError) Unwrap() error { return e.err }

// checkInput asserts the invariants every formula in this package relies on.
// It returns a wrapped error naming the offending field.
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
