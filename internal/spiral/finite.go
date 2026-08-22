package spiral

import (
	"math"
	"math/cmplx"
)

// The finite-depth Ekman spiral keeps the same governing equation but adds a
// no-slip condition at the sea floor z = -H:
//
//	W(z) = A*exp(lambda*z) + B*exp(-lambda*z)
//	W(-H) = 0,                dW/dz|z=0 = tau/(rho*K)
//
// with lambda = (1 + i*sign(f))/delta. The two conditions fix A and B, giving
// a solution that reduces to the infinite-depth spiral as H grows and that
// weakens the surface current and its deflection as H shrinks. The config
// layer rejects cases with H < De outright, so these functions are mainly
// used to quantify the boundary effect and to document the approximation
// behind the deep-water formula.

// EkmanLambda returns the complex wavenumber lambda of the spiral solution
// for the given eddy viscosity and Coriolis parameter. Its real part is
// positive so e^{lambda*z} decays into the deep water.
func EkmanLambda(K, f float64) complex128 {
	return complex(1.0/EkmanScale(K, f), SignOf(f)/EkmanScale(K, f))
}

// finiteDepthAmplitude solves for the surface amplitude A of the
// finite-depth spiral given the water depth H (metres).
func finiteDepthAmplitude(H, tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	if err := checkInput(rho, f, K, tau); err != nil {
		return 0, err
	}
	lambda := EkmanLambda(K, f)
	stress := WindStressComplex(tau, windHeadingDeg)
	// A = tau_c / (rho*K*lambda*(1 + exp(-2*lambda*H)))
	denominator := complex(rho*K, 0) * lambda * (1 + cmplx.Exp(-2*lambda*complex(H, 0)))
	return stress / denominator, nil
}

// FiniteDepthSurfaceCurrent returns the complex surface current for a water
// column of depth H. It approaches the infinite-depth surface current as H
// grows and drops toward the linear (bottom-stressed) regime for very shallow
// water.
func FiniteDepthSurfaceCurrent(H, tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	A, err := finiteDepthAmplitude(H, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	lambda := EkmanLambda(K, f)
	return A * (1 - cmplx.Exp(-2*lambda*complex(H, 0))), nil
}

// FiniteDepthCurrent evaluates the finite-depth spiral at depthM below the
// surface (0 <= depthM <= H).
func FiniteDepthCurrent(depthM, H, tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	A, err := finiteDepthAmplitude(H, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	lambda := EkmanLambda(K, f)
	z := -depthM
	return A*cmplx.Exp(lambda*complex(z, 0)) + A*(-cmplx.Exp(-2*lambda*complex(H, 0)))*cmplx.Exp(-lambda*complex(z, 0)), nil
}

// FiniteDepthSurfaceDeflection returns the angle between the finite-depth
// surface current and the wind, in degrees. It is exactly 45 degrees only in
// the infinite-depth limit; shallow columns produce smaller deflections as
// the bottom friction pulls the current back toward the wind.
func FiniteDepthSurfaceDeflection(H, tau, rho, K, f, windHeadingDeg float64) (float64, error) {
	w, err := FiniteDepthSurfaceCurrent(H, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	wind := WindStressComplex(tau, windHeadingDeg)
	angle := math.Atan2(imag(w), real(w)) - math.Atan2(imag(wind), real(wind))
	// The deflection is clockwise in the northern hemisphere, so sign matters.
	return -RadToDeg(angle), nil
}

// FiniteDepthErrorRatio returns |W_finite(0)/W_infinite(0) - 1|, the relative
// deviation of the finite-depth surface current from the infinite-depth
// formula. The config layer guarantees H >= De, where this ratio is at most
// ~4% (the amplitude at z = -De has decayed to exp(-pi)); the report declares
// this residual when a finite depth is present.
func FiniteDepthErrorRatio(H, tau, rho, K, f, windHeadingDeg float64) (float64, error) {
	wf, err := FiniteDepthSurfaceCurrent(H, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	wi, err := SurfaceComplex(tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	ratio := cmplx.Abs(wf) / cmplx.Abs(wi)
	return math.Abs(ratio - 1), nil
}
