package spiral

import (
	"math"
	"math/cmplx"
)

func EkmanLambda(K, f float64) complex128 {
	return complex(1.0/EkmanScale(K, f), SignOf(f)/EkmanScale(K, f))
}

func finiteDepthAmplitude(H, tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	if err := checkInput(rho, f, K, tau); err != nil {
		return 0, err
	}
	lambda := EkmanLambda(K, f)
	stress := WindStressComplex(tau, windHeadingDeg)
	denominator := complex(rho*K, 0) * lambda * (1 + cmplx.Exp(-2*lambda*complex(H, 0)))
	return stress / denominator, nil
}

func FiniteDepthSurfaceCurrent(H, tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	A, err := finiteDepthAmplitude(H, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	lambda := EkmanLambda(K, f)
	return A * (1 - cmplx.Exp(-2*lambda*complex(H, 0))), nil
}

func FiniteDepthCurrent(depthM, H, tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	A, err := finiteDepthAmplitude(H, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	lambda := EkmanLambda(K, f)
	z := -depthM
	return A*cmplx.Exp(lambda*complex(z, 0)) + A*(-cmplx.Exp(-2*lambda*complex(H, 0)))*cmplx.Exp(-lambda*complex(z, 0)), nil
}

func FiniteDepthSurfaceDeflection(H, tau, rho, K, f, windHeadingDeg float64) (float64, error) {
	w, err := FiniteDepthSurfaceCurrent(H, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return 0, err
	}
	wind := WindStressComplex(tau, windHeadingDeg)
	angle := math.Atan2(imag(w), real(w)) - math.Atan2(imag(wind), real(wind))
	return -RadToDeg(angle), nil
}

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
