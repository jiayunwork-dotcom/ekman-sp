package spiral

import (
	"fmt"
	"math"
)

type ResidualStats struct {
	MaxSpeedErr  float64
	MeanSpeedErr float64
	Count        int
}

func ProfileSpeedStats(points []Point) ResidualStats {
	var s ResidualStats
	if len(points) == 0 {
		return s
	}
	s.Count = len(points)
	for i := 1; i < len(points); i++ {
		err := math.Abs(points[i].Speed - points[i-1].Speed)
		s.MeanSpeedErr += err
		if err > s.MaxSpeedErr {
			s.MaxSpeedErr = err
		}
	}
	if s.Count > 1 {
		s.MeanSpeedErr /= float64(s.Count - 1)
	}
	return s
}

type Envelope struct {
	MinDepth   float64
	MaxDepth   float64
	MinSpeed   float64
	MaxSpeed   float64
	MinHeading float64
	MaxHeading float64
}

func ComputeEnvelope(points []Point) Envelope {
	var e Envelope
	if len(points) == 0 {
		return e
	}
	e.MinDepth = points[0].Depth
	e.MaxDepth = points[0].Depth
	e.MinSpeed = points[0].Speed
	e.MaxSpeed = e.MinSpeed
	e.MinHeading = points[0].Heading
	e.MaxHeading = e.MinHeading
	for _, p := range points {
		if p.Depth < e.MinDepth {
			e.MinDepth = p.Depth
		}
		if p.Depth > e.MaxDepth {
			e.MaxDepth = p.Depth
		}
		if p.Speed < e.MinSpeed {
			e.MinSpeed = p.Speed
		}
		if p.Speed > e.MaxSpeed {
			e.MaxSpeed = p.Speed
		}
		if p.Heading < e.MinHeading {
			e.MinHeading = p.Heading
		}
		if p.Heading > e.MaxHeading {
			e.MaxHeading = p.Heading
		}
	}
	return e
}

func HeadingSpanDeg(e Envelope) float64 {
	return e.MaxHeading - e.MinHeading
}

func SpeedDecayRatio(points []Point) (float64, error) {
	if len(points) < 2 {
		return 0, fmt.Errorf("spiral: need at least 2 profile points")
	}
	top := points[0].Speed
	bottom := points[len(points)-1].Speed
	if top <= 0 {
		return 0, fmt.Errorf("spiral: zero surface speed")
	}
	return bottom / top, nil
}

func MonotonicSpeedDecay(points []Point) bool {
	if len(points) < 2 {
		return true
	}
	prev := points[0].Speed
	for i := 1; i < len(points); i++ {
		cur := points[i].Speed
		if cur > prev+1e-9 {
			return false
		}
		prev = cur
	}
	return true
}

func DepthAtHalfSpeed(points []Point) (float64, error) {
	if len(points) == 0 {
		return 0, fmt.Errorf("spiral: empty profile")
	}
	target := points[0].Speed / 2
	for i := 1; i < len(points); i++ {
		if points[i].Speed <= target {
			return points[i].Depth, nil
		}
	}
	return points[len(points)-1].Depth, nil
}

func VeerSpan(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	minV := points[0].Veer
	maxV := points[0].Veer
	for _, p := range points[1:] {
		if p.Veer < minV {
			minV = p.Veer
		}
		if p.Veer > maxV {
			maxV = p.Veer
		}
	}
	return maxV - minV
}
