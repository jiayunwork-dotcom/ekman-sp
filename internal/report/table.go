package report

import (
	"fmt"
	"io"

	"ekman-sp/internal/spiral"
)

// tableColumns are the printed columns of the depth table.
const (
	colDepth   = "depth (m)"
	colU       = "u (m/s)"
	colV       = "v (m/s)"
	colSpeed   = "|V| (m/s)"
	colHeading = "heading (deg)"
	colVeer    = "veer (deg)"
)

// tableWidths are the fixed column widths used to line up the depth table.
var tableWidths = []int{10, 9, 9, 10, 13, 11}

// writeTable prints the sampled spiral as an aligned text table.
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

// FormatPointRow renders one profile sample as a fixed-width table line. It
// is exposed separately so tests can assert on a single row's layout.
func FormatPointRow(p spiral.Point) string {
	return fmt.Sprintf("%10.2f  %8.4f  %8.4f  %9.4f  %12.1f  %10.1f",
		p.Depth, p.U, p.V, p.Speed, p.Heading, p.Veer)
}

// writeTableHeader prints the column header line and a rule beneath it.
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

// padRight left-justifies s within width.
func padRight(s string, width int) string {
	for len(s) < width {
		s += " "
	}
	return s
}

// rule builds the dashes under a header line from the column widths.
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
