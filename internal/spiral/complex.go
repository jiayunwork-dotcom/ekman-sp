package spiral

import (
	"math"
	"math/cmplx"
)

func WindStressComplex(tau, windHeadingDeg float64) complex128 {
	phi := HeadingToPhase(windHeadingDeg)
	return complex(tau*math.Cos(phi), tau*math.Sin(phi))
}

func SurfaceComplex(tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return 0, err
	}
	phi0 := HeadingToPhase(SurfaceHeading(windHeadingDeg, f))
	return cmplx.Rect(speed, phi0), nil
}

func CurrentComplex(depthM, tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return 0, err
	}
	delta := EkmanScale(K, f)
	phase := HeadingToPhase(SurfaceHeading(windHeadingDeg, f)) + RotationPhaseAtDepth(depthM, K, f)
	amplitude := speed * math.Exp(-depthM/delta)
	return cmplx.Rect(amplitude, phase), nil
}

func TransportComplex(tau, rho, f, windHeadingDeg float64) (complex128, error) {
	if err := checkInput(rho, f, 1.0, tau); err != nil {
		return 0, err
	}
	tauC := WindStressComplex(tau, windHeadingDeg)
	return complex(0, -SignOf(f)) * tauC / complex(rho*math.Abs(f), 0), nil
}

func unpack(w complex128) (east, north float64) {
	return real(w), imag(w)
}
