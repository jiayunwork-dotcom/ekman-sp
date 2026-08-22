package spiral

import (
	"math"
	"math/cmplx"
)

// The steady Ekman spiral is naturally written as a damped complex
// exponential. With z measured upward from the sea surface and the current
// packed as W = u + i*v (real = east, imag = north), the governing equation
//
//	K d^2W/dz^2 = i*f*W
//
// has the deep-water solution
//
//	W(d) = |V0| * exp(-d/delta) * exp(i*(phi0 - sign(f)*d/delta))
//
// where d = -z is the depth below the surface. These helpers expose that
// complex form directly; the real/imag components used by the report are
// simply extracted from it.

// WindStressComplex packs a wind stress of magnitude tau blowing toward
// windHeadingDeg into a complex number aligned with the wind.
func WindStressComplex(tau, windHeadingDeg float64) complex128 {
	phi := HeadingToPhase(windHeadingDeg)
	return complex(tau*math.Cos(phi), tau*math.Sin(phi))
}

// SurfaceComplex returns the complex surface current (real = east, imag =
// north) for the infinite-depth spiral.
func SurfaceComplex(tau, rho, K, f, windHeadingDeg float64) (complex128, error) {
	speed, err := SurfaceSpeed(tau, rho, K, f)
	if err != nil {
		return 0, err
	}
	phi0 := HeadingToPhase(SurfaceHeading(windHeadingDeg, f))
	return cmplx.Rect(speed, phi0), nil
}

// CurrentComplex evaluates the complex current at depth depthM below the
// surface for the infinite-depth spiral.
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

// TransportComplex returns the depth-integrated Ekman transport as a complex
// number (real = east, imag = north). The integral of the spiral collapses to
// the surface stress divided by rho*f, rotated by the hemisphere: the
// transport lies 90 degrees to the right (northern) or left (southern) of
// the wind and does not involve the eddy viscosity.
func TransportComplex(tau, rho, f, windHeadingDeg float64) (complex128, error) {
	if err := checkInput(rho, f, 1.0, tau); err != nil {
		return 0, err
	}
	tauC := WindStressComplex(tau, windHeadingDeg)
	// -i*W rotates a vector 90 degrees clockwise, which is the northern
	// hemisphere deflection; the southern hemisphere needs +i*W.
	return complex(0, -SignOf(f)) * tauC / complex(rho*math.Abs(f), 0), nil
}

// unpack extracts the east/north components of a complex current vector.
func unpack(w complex128) (east, north float64) {
	return real(w), imag(w)
}
