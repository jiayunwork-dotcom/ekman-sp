package config

import (
	"math"

	"ekman-sp/internal/spiral"
)

// LatitudeFromCoriolis recovers the latitude that would produce a given
// Coriolis parameter, f = 2*Omega*sin(lat). It is the inverse of
// CoriolisFromLatitude and is used to report an equivalent latitude when the
// case supplied f directly.
func LatitudeFromCoriolis(f float64) float64 {
	ratio := f / (2.0 * spiral.Omega)
	if ratio > 1 {
		ratio = 1
	}
	if ratio < -1 {
		ratio = -1
	}
	return spiral.RadToDeg(math.Asin(ratio))
}

// LatBand classifies a latitude magnitude into a qualitative band used for
// reporting and for sanity checks on the Coriolis value.
type LatBand int

const (
	// BandLow covers |lat| < 15 degrees.
	BandLow LatBand = iota
	// BandMid covers 15 <= |lat| < 45 degrees.
	BandMid
	// BandHigh covers |lat| >= 45 degrees.
	BandHigh
)

// BandOf returns the latitude band for a latitude in degrees.
func BandOf(latitudeDeg float64) LatBand {
	abs := math.Abs(latitudeDeg)
	switch {
	case abs < 15:
		return BandLow
	case abs < 45:
		return BandMid
	default:
		return BandHigh
	}
}

// String returns a short label for the band.
func (b LatBand) String() string {
	switch b {
	case BandLow:
		return "low"
	case BandMid:
		return "mid"
	case BandHigh:
		return "high"
	}
	return "unknown"
}

// IsTropical reports whether the latitude lies inside the trade-wind belt,
// used by the report to annotate a low-latitude case.
func IsTropical(latitudeDeg float64) bool {
	return BandOf(latitudeDeg) == BandLow
}
