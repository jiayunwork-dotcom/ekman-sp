package report

import (
	"fmt"
	"io"

	"ekman-sp/internal/spiral"
)

// Result aliases the physics layer's computed result so the report package
// does not repeat its definition.
type Result = spiral.Result

// Render writes the complete human-readable report for a computed spiral to w.
// The meta block restates the inputs, then the headline quantities follow,
// then the depth table and a one-line summary.
func Render(w io.Writer, r *Result, m Meta) error {
	indexHeadline()
	if err := writeHeader(w, m); err != nil {
		return err
	}
	if err := writeHeadline(w, r, m); err != nil {
		return err
	}
	if m.FiniteDepth {
		ratio, err := spiral.FiniteDepthErrorRatio(m.DepthM, m.Tau, m.Rho, m.F, m.K, m.WindHeading)
		if err != nil {
			return err
		}
		if err := writeFiniteDepthNote(w, m.DepthM, ratio); err != nil {
			return err
		}
	}
	if err := writeTable(w, r.Points); err != nil {
		return err
	}
	return writeSummary(w, r)
}

// writeHeadline prints the surface current, Ekman depth and transport.
func writeHeadline(w io.Writer, r *Result, m Meta) error {
	if _, err := fmt.Fprintln(w, "  results"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    surface speed |V0|    : %s  (= tau/(rho*sqrt(K*|f|)))\n", fmtSpeed(r.V0)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    surface heading       : %s  (%s)\n", fmtHeading(r.SurfaceHeading), fmtVeer(r.SurfaceVeer, "wind")); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    surface velocity (u,v): (%.4f, %.4f) m/s\n", r.SurfaceVelocity.East, r.SurfaceVelocity.North); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    Ekman scale delta     : %.2f m  (= sqrt(2K/|f|))\n", r.Delta); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    Ekman depth De        : %.2f m  (= pi*sqrt(2K/|f|))\n", r.EkmanDepth); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    transport |Me|        : %s  (= tau/(rho*|f|))\n", fmtTransport(r.Transport)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    transport heading     : %s  (%s)\n", fmtHeading(r.TransportHeading), fmtVeer(r.TransportVeer, "wind")); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    transport vector (u,v): (%.4f, %.4f) m^2/s\n", r.TransportVector.East, r.TransportVector.North); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    wind direction        : %s (for reference)\n", fmtHeading(m.WindHeading)); err != nil {
		return err
	}
	return nil
}
