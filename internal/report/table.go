package report

import (
	"fmt"
	"io"

	"ekman-sp/internal/spiral"
)

const (
	colDepth   = "depth (m)"
	colU       = "u (m/s)"
	colV       = "v (m/s)"
	colSpeed   = "|V| (m/s)"
	colHeading = "heading (deg)"
	colVeer    = "veer (deg)"
)

var tableWidths = []int{10, 9, 9, 10, 13, 11}

func writeTable(w io.Writer, points []spiral.Point) error {
	if _, err := fmt.Fprintln(w, "  depth profile"); err != nil {
		return err
	}
	if err := writeTableHeader(w); err != nil {
		return err
	}
	for _, p := range points {
		if _, err := fmt.Fprintln(w, FormatPointRow(p)); err != nil {
			return err
		}
	}
	return nil
}

func FormatPointRow(p spiral.Point) string {
	return fmt.Sprintf("%10.2f  %8.4f  %8.4f  %9.4f  %12.1f  %10.1f",
		p.Depth, p.U, p.V, p.Speed, p.Heading, p.Veer)
}

func writeTableHeader(w io.Writer) error {
	cols := []string{colDepth, colU, colV, colSpeed, colHeading, colVeer}
	line := "  "
	for i, c := range cols {
		line += padRight(c, tableWidths[i])
	}
	if _, err := fmt.Fprintln(w, line); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "  "+rule(tableWidths)); err != nil {
		return err
	}
	return nil
}

func padRight(s string, width int) string {
	for len(s) < width {
		s += " "
	}
	return s
}

func rule(widths []int) string {
	out := ""
	for _, wd := range widths {
		for i := 0; i < wd; i++ {
			out += "-"
		}
		out += "  "
	}
	return out
}
