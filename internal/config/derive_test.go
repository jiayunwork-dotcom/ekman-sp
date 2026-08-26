package config

import (
	"math"
	"strings"
	"testing"
)

func TestLatitudeFromCoriolisRoundTrip(t *testing.T) {
	for _, lat := range []float64{-60, -45, -15, 15, 45, 60} {
		f := CoriolisFromLatitude(lat)
		back := LatitudeFromCoriolis(f)
		if math.Abs(back-lat) > 1e-6 {
			t.Errorf("round trip lat %g -> %g", lat, back)
		}
	}
}

func TestLatitudeFromCoriolisAtEquator(t *testing.T) {
	if got := LatitudeFromCoriolis(0); got != 0 {
		t.Errorf("LatitudeFromCoriolis(0) = %g, want 0", got)
	}
}

func TestBandOf(t *testing.T) {
	cases := []struct {
		lat  float64
		want LatBand
	}{
		{0, BandLow}, {10, BandLow}, {14.9, BandLow},
		{15, BandMid}, {30, BandMid}, {44.9, BandMid},
		{45, BandHigh}, {60, BandHigh}, {-60, BandHigh},
	}
	for _, c := range cases {
		if got := BandOf(c.lat); got != c.want {
			t.Errorf("BandOf(%g) = %v, want %v", c.lat, got, c.want)
		}
	}
}

func TestIsTropical(t *testing.T) {
	if !IsTropical(10) {
		t.Error("10 deg should be tropical")
	}
	if IsTropical(45) {
		t.Error("45 deg should not be tropical")
	}
}

func TestDescribeReportsSource(t *testing.T) {
	r := mustResolve(t, sampleInput())
	desc := r.Describe()
	if !strings.Contains(desc, "northern hemisphere") {
		t.Errorf("Describe = %q, want hemisphere", desc)
	}
	if !strings.Contains(desc, "lat 45.0 deg") {
		t.Errorf("Describe = %q, want latitude provenance", desc)
	}
}

func TestCoriolisProvenanceFromLatitude(t *testing.T) {
	r := mustResolve(t, sampleInput())
	if !strings.Contains(r.CoriolisProvenance(), "from latitude") {
		t.Errorf("provenance = %q, want from latitude", r.CoriolisProvenance())
	}
}

func TestCoriolisProvenanceFromF(t *testing.T) {
	in := sampleInput()
	f := 1.0313e-4
	in.F = &f
	in.LatitudeDeg = nil
	r := mustResolve(t, in)
	if !strings.Contains(r.CoriolisProvenance(), "given directly") {
		t.Errorf("provenance = %q, want given directly", r.CoriolisProvenance())
	}
}

func TestProfileGeometry(t *testing.T) {
	r := mustResolve(t, sampleInput())
	de, full := r.ProfileGeometry()
	if de <= 0 || full <= 0 || full <= de {
		t.Errorf("geometry De=%g, full-turn=%g; want 0 < De < full", de, full)
	}
}
