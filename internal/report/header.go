package report

import (
	"fmt"
	"io"
)

// Meta carries the resolved case inputs so the header can restate exactly
// what was computed. It mirrors the physics parameters plus the provenance of
// the Coriolis parameter; the report never modifies any of it.
type Meta struct {
	// Tau is the wind stress magnitude in N/m^2.
	Tau float64
	// WindHeading is the wind direction in degrees clockwise from north.
	WindHeading float64
	// Rho is the seawater density in kg/m^3.
	Rho float64
	// F is the Coriolis parameter in 1/s (signed).
	F float64
	// K is the eddy viscosity in m^2/s.
	K float64
	// Latitude is the latitude in degrees that produced F, or NaN when the
	// case supplied f directly.
	Latitude float64
	// HasLatitude records whether the case supplied latitude_deg.
	HasLatitude bool
	// FiniteDepth records whether the case carried an explicit water depth.
	FiniteDepth bool
	// DepthM is the explicit water depth in metres (only when FiniteDepth).
	DepthM float64
	// NPoints is the number of profile samples.
	NPoints int
}

// writeHeader prints the resolved input block.
func writeHeader(w io.Writer, m Meta) error {
	if _, err := fmt.Fprintln(w, "Ekman spiral — steady-state wind-driven surface layer"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "  inputs"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    wind stress tau        : %s toward %s\n", fmtStress(m.Tau), fmtHeading(m.WindHeading)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    density rho            : %s\n", fmtRho(m.Rho)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    Coriolis f             : %s\n", fmtCoriolis(m.F, m.Latitude, m.HasLatitude)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    eddy viscosity K       : %s\n", fmtK(m.K)); err != nil {
		return err
	}
	if m.FiniteDepth {
		if _, err := fmt.Fprintf(w, "    water depth            : %.1f m (finite)\n", m.DepthM); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintln(w, "    water depth            : infinite (classic Ekman solution)"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "    profile samples        : %d\n", m.NPoints); err != nil {
		return err
	}
	return nil
}
