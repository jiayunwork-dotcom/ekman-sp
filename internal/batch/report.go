package batch

import (
	"fmt"
	"strings"
)

func TextReport(r BatchResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "== ekman-sp batch: %d cases, %d ok, %d failed ==\n", r.Total, r.Success, r.Failed)
	for _, it := range r.Items {
		if it.Error != "" {
			fmt.Fprintf(&b, "FAIL  %s: %s\n", it.Name, it.Error)
			continue
		}
		fmt.Fprintf(&b, "OK    %s  V0=%.4f δE=%.3f m  T=%.4f m²/s  veer=%.1f°/%s  ∫err=%.2e\n",
			it.Name, it.V0, it.EkmanDepth, it.Transport, it.SurfaceVeer, it.Hemisphere, it.IntegratedErr)
	}
	return b.String()
}

func MaxIntegratedError(r BatchResult) float64 {
	max := 0.0
	for _, it := range r.Items {
		if it.Error == "" && it.IntegratedErr > max {
			max = it.IntegratedErr
		}
	}
	return max
}

func NorthernCount(r BatchResult) int {
	n := 0
	for _, it := range r.Items {
		if it.Error == "" && it.Hemisphere == "north" {
			n++
		}
	}
	return n
}

func FailedNames(r BatchResult) []string {
	var out []string
	for _, it := range r.Items {
		if it.Error != "" {
			out = append(out, it.Name)
		}
	}
	return out
}
