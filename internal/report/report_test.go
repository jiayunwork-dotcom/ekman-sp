package report

import (
	"bytes"
	"strings"
	"testing"

	"ekman-sp/internal/spiral"
)

// sampleMeta returns a Meta block matching the mid-latitude westerly case.
func sampleMeta() Meta {
	return Meta{
		Tau:         0.2,
		WindHeading: 90,
		Rho:         1025,
		F:           1.0313e-4,
		K:           0.05,
		Latitude:    45,
		HasLatitude: true,
		FiniteDepth: false,
		NPoints:     21,
	}
}

// sampleResult computes a real spiral for the sample meta.
func sampleResult(t *testing.T) Result {
	t.Helper()
	res, err := spiral.Compute(spiral.Params{
		Tau:         0.2,
		WindHeading: 90,
		Rho:         1025,
		F:           1.0313e-4,
		K:           0.05,
		NPoints:     21,
	})
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	return Result(res)
}

func TestRenderContainsHeadline(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, ptrResult(t), sampleMeta()); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{
		"Ekman spiral",
		"surface speed |V0|",
		"surface heading",
		"Ekman depth De",
		"transport |Me|",
		"transport heading",
		"depth profile",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("report missing %q", want)
		}
	}
}

func TestRenderTableRows(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, ptrResult(t), sampleMeta()); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	idx := strings.Index(buf.String(), "depth profile")
	if idx < 0 {
		t.Fatal("depth profile section missing")
	}
	lines := strings.Split(buf.String()[idx:], "\n")
	rows := 0
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		fields := strings.Fields(ln)
		if len(fields) >= 4 && isNumeric(fields[0]) {
			rows++
		}
	}
	if rows != sampleMeta().NPoints {
		t.Errorf("table rows = %d, want %d", rows, sampleMeta().NPoints)
	}
}

// isNumeric reports whether s parses as a number; used to recognise data rows
// in the rendered table.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	seenDigit := false
	for _, c := range s {
		if c == '.' || c == '-' || c == '+' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
		seenDigit = true
	}
	return seenDigit
}

func TestRenderFirstRowIsSurface(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, ptrResult(t), sampleMeta()); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	out := buf.String()
	// The first depth row is the surface: depth 0, heading equal to the
	// surface heading reported in the headline (135 deg for the sample).
	if !strings.Contains(out, "0.00    0.0608   -0.0608") {
		t.Errorf("surface row missing in table:\n%s", out)
	}
}

func TestRenderInfiniteDepthLabel(t *testing.T) {
	var buf bytes.Buffer
	meta := sampleMeta()
	meta.FiniteDepth = false
	if err := Render(&buf, ptrResult(t), meta); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if !strings.Contains(buf.String(), "infinite") {
		t.Error("infinite-depth label missing")
	}
}

func TestRenderFiniteDepthLabel(t *testing.T) {
	var buf bytes.Buffer
	meta := sampleMeta()
	meta.FiniteDepth = true
	meta.DepthM = 500
	if err := Render(&buf, ptrResult(t), meta); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if !strings.Contains(buf.String(), "500.0 m (finite)") {
		t.Errorf("finite-depth label missing:\n%s", buf.String())
	}
}

func TestCompassPoint(t *testing.T) {
	cases := []struct {
		heading float64
		want    string
	}{
		{0, "N"}, {90, "E"}, {135, "SE"}, {180, "S"},
		{225, "SW"}, {270, "W"}, {315, "NW"}, {360, "N"},
		{45, "NE"}, {-90, "W"},
	}
	for _, c := range cases {
		if got := CompassPoint(c.heading); got != c.want {
			t.Errorf("CompassPoint(%g) = %q, want %q", c.heading, got, c.want)
		}
	}
}

func TestRendererBuffersOutput(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf)
	if err := r.WriteReport(ptrResult(t), sampleMeta()); err != nil {
		t.Fatalf("WriteReport error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("buffered render produced no output")
	}
	if !strings.Contains(buf.String(), "Ekman depth De") {
		t.Error("buffered render missing headline")
	}
}

func TestFormatHeading(t *testing.T) {
	got := fmtHeading(135)
	if !strings.Contains(got, "135.0") || !strings.Contains(got, "SE") {
		t.Errorf("fmtHeading(135) = %q, want heading plus compass point", got)
	}
}

func TestFormatVeer(t *testing.T) {
	if got := fmtVeer(45, "wind"); !strings.Contains(got, "right of wind") {
		t.Errorf("fmtVeer(45) = %q, want right of wind", got)
	}
	if got := fmtVeer(-45, "wind"); !strings.Contains(got, "left of wind") {
		t.Errorf("fmtVeer(-45) = %q, want left of wind", got)
	}
}

func ptrResult(t *testing.T) *Result {
	t.Helper()
	res := sampleResult(t)
	return &res
}
