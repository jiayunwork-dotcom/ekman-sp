package report

import (
	"fmt"
	"math"

	"ekman-sp/internal/spiral"
)

func fmtStress(tau float64) string {
	return fmt.Sprintf("%.3f N/m^2", tau)
}

func fmtRho(rho float64) string {
	return fmt.Sprintf("%.1f kg/m^3", rho)
}

func fmtK(K float64) string {
	return fmt.Sprintf("%.3f m^2/s", K)
}

func fmtSpeed(v float64) string {
	return fmt.Sprintf("%.4f m/s", v)
}

func fmtTransport(m float64) string {
	return fmt.Sprintf("%.4f m^2/s", m)
}

func fmtCoriolis(f, latitude float64, hasLatitude bool) string {
	if hasLatitude {
		return fmt.Sprintf("%+.3e 1/s (lat %.1f deg, %s)", f, latitude, spiral.HemisphereName(f))
	}
	return fmt.Sprintf("%+.3e 1/s (%s)", f, spiral.HemisphereName(f))
}

func fmtHeading(deg float64) string {
	return fmt.Sprintf("%.1f deg (%s)", spiral.NormalizeHeading(deg), CompassPoint(deg))
}

func fmtVeer(deg float64, reference string) string {
	dir := "right"
	if deg < 0 {
		dir = "left"
	}
	return fmt.Sprintf("%.1f deg %s of %s", math.Abs(deg), dir, reference)
}
