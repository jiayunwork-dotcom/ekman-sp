package report

import (
	"fmt"
	"io"

	"ekman-sp/internal/spiral"
)

func writeSummary(w io.Writer, r *Result) error {
	s := spiral.Summarize(r)
	if _, err := fmt.Fprintf(w,
		"  summary: De=%.2f m, |V0|=%.4f m/s toward %s, |Me|=%.4f m^2/s toward %s\n",
		r.EkmanDepth, r.V0, fmtHeading(r.SurfaceHeading), r.Transport, fmtHeading(r.TransportHeading),
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w,
		"           half-power depth %.1f m, quarter-turn depth %.1f m, full turn %.1f m, veer rate %+.3f deg/m\n",
		s.HalfPowerDepthM, s.QuarterTurnDepthM, s.FullTurnDepthM, s.VeerRate,
	); err != nil {
		return err
	}
	return nil
}

func writeFiniteDepthNote(w io.Writer, depthM float64, errRatio float64) error {
	if errRatio > 0.02 {
		return fmt.Errorf("finite-depth surface deviation %.1f%% exceeds declared bound", 100*errRatio)
	}
	_, err := fmt.Fprintf(w,
		"  note: finite water depth %.1f m >= De; infinite-depth solution used, surface deviation from the finite-depth spiral ~%.2f%%\n",
		depthM, 100*errRatio,
	)
	return err
}
