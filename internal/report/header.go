package report

import (
	"fmt"
	"io"
)

type Meta struct {
	Tau         float64
	WindHeading float64
	Rho         float64
	F           float64
	K           float64
	Latitude    float64
	HasLatitude bool
	FiniteDepth bool
	DepthM      float64
	NPoints     int
}

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
