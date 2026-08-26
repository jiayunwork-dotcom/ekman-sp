package config

import (
	"math"

	"ekman-sp/internal/spiral"
)

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

type LatBand int

const (
	BandLow LatBand = iota
	BandMid
	BandHigh
)

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

func IsTropical(latitudeDeg float64) bool {
	return BandOf(latitudeDeg) == BandLow
}
