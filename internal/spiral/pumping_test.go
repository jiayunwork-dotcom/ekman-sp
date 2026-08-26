package spiral

import (
	"math"
	"testing"
)

func TestInertialPeriodAtMidlat(t *testing.T) {
	f := 1e-4
	p, err := InertialPeriod(f)
	if err != nil {
		t.Fatal(err)
	}
	want := 2 * math.Pi / f
	if math.Abs(p-want) > 1e-9*want {
		t.Fatalf("period %g want %g", p, want)
	}
	if _, err := InertialPeriod(0); err == nil {
		t.Fatal("equator should error")
	}
}

func TestHemisphereFlipsPumping(t *testing.T) {
	if err := HemisphereFlipsPumping(1025, 1e-4, 1e-7); err != nil {
		t.Fatal(err)
	}
	w, err := PumpingFromUniformCurl(1025, 1e-4, 1e-7)
	if err != nil {
		t.Fatal(err)
	}
	if w <= 0 {
		t.Fatalf("cyclonic curl should pump up, got %g", w)
	}
}

func TestSpinUpFasterThanDiffusion(t *testing.T) {
	ok, err := SpinUpFasterThanDiffusion(50, 0.05, 1e-4)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Ekman spin-up should beat H²/K at 50 m")
	}
	ts, _ := SpinUpTime(50, 0.05, 1e-4)
	td, _ := VerticalDiffusiveTime(50, 0.05)
	if ts >= td {
		t.Fatalf("ts=%g td=%g", ts, td)
	}
}

func TestBottomCancelsSurfaceTransport(t *testing.T) {
	surf, err := Transport(0.2, 1025, 1e-4)
	if err != nil {
		t.Fatal(err)
	}
	bot, err := BottomLayerTransport(0.2, 1025, 1e-4)
	if err != nil {
		t.Fatal(err)
	}
	if err := InteriorGeostrophicBalance(bot, surf); err != nil {
		t.Fatal(err)
	}
}

func TestDecayEFoldIsEkmanScale(t *testing.T) {
	d := DecayEFoldDepth(0.05, 1e-4)
	if math.Abs(d-EkmanScale(0.05, 1e-4)) > 1e-12 {
		t.Fatalf("e-fold depth %g", d)
	}
	if TurnsToEFold() <= 0 {
		t.Fatal("turns to e-fold")
	}
}
