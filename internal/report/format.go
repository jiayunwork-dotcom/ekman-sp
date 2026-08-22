package report

import (
	"fmt"
	"math"

	"ekman-sp/internal/spiral"
)

// fmtStress renders a wind stress magnitude in N/m^2 with a sensible number
// of significant digits.
func fmtStress(tau float64) string {
	return fmt.Sprintf("%.3f N/m^2", tau)
}

// fmtRho renders a density in kg/m^3.
func fmtRho(rho float64) string {
	return fmt.Sprintf("%.1f kg/m^3", rho)
}

// fmtK renders an eddy viscosity in m^2/s.
func fmtK(K float64) string {
	return fmt.Sprintf("%.3f m^2/s", K)
}

// fmtSpeed renders a current speed in m/s.
func fmtSpeed(v float64) string {
	return fmt.Sprintf("%.4f m/s", v)
}

// fmtTransport renders a depth-integrated transport in m^2/s.
func fmtTransport(m float64) string {
	return fmt.Sprintf("%.4f m^2/s", m)
}

// fmtCoriolis renders the Coriolis parameter together with the latitude that
// produced it, when the case supplied one.
func fmtCoriolis(f, latitude float64, hasLatitude bool) string {
	if hasLatitude {
		return fmt.Sprintf("%+.3e 1/s (lat %.1f deg, %s)", f, latitude, spiral.HemisphereName(f))
	}
	return fmt.Sprintf("%+.3e 1/s (%s)", f, spiral.HemisphereName(f))
}

// fmtHeading renders a compass heading with its cardinal point.
func fmtHeading(deg float64) string {
	return fmt.Sprintf("%.1f deg (%s)", spiral.NormalizeHeading(deg), CompassPoint(deg))
}

// fmtVeer renders a signed deflection relative to a reference direction.
func fmtVeer(deg float64, reference string) string {
	dir := "right"
	if deg < 0 {
		dir = "left"
	}
	return fmt.Sprintf("%.1f deg %s of %s", math.Abs(deg), dir, reference)
}
