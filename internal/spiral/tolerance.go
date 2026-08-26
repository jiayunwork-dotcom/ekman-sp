package spiral

import (
	"fmt"
	"math"
)

type Perturbation struct {
	Field  string  `json:"field"`
	Delta  float64 `json:"delta"`
	V0     float64 `json:"v0"`
	DV0    float64 `json:"d_v0"`
	Depth  float64 `json:"ekman_depth_m"`
	DDepth float64 `json:"d_depth"`
}

type ToleranceReport struct {
	BaseV0    float64        `json:"base_v0"`
	BaseDepth float64        `json:"base_depth_m"`
	Items     []Perturbation `json:"perturbations"`
}

func TauSensitivity(p Params, relDelta float64) (ToleranceReport, error) {
	if relDelta <= 0 || relDelta >= 1 {
		return ToleranceReport{}, fmt.Errorf("spiral: relDelta out of range")
	}
	base, err := Compute(p)
	if err != nil {
		return ToleranceReport{}, err
	}
	items := make([]Perturbation, 0, 2)
	for _, sign := range []float64{-1, 1} {
		clone := p
		clone.Tau *= 1 + sign*relDelta
		res, err := Compute(clone)
		if err != nil {
			return ToleranceReport{}, err
		}
		items = append(items, Perturbation{
			Field:  "tau",
			Delta:  sign * relDelta,
			V0:     res.V0,
			DV0:    res.V0 - base.V0,
			Depth:  res.EkmanDepth,
			DDepth: res.EkmanDepth - base.EkmanDepth,
		})
	}
	return ToleranceReport{BaseV0: base.V0, BaseDepth: base.EkmanDepth, Items: items}, nil
}

func KSensitivity(p Params, relDelta float64) (ToleranceReport, error) {
	if relDelta <= 0 || relDelta >= 1 {
		return ToleranceReport{}, fmt.Errorf("spiral: relDelta out of range")
	}
	base, err := Compute(p)
	if err != nil {
		return ToleranceReport{}, err
	}
	items := make([]Perturbation, 0, 2)
	for _, sign := range []float64{-1, 1} {
		clone := p
		clone.K *= 1 + sign*relDelta
		res, err := Compute(clone)
		if err != nil {
			return ToleranceReport{}, err
		}
		ratio := res.EkmanDepth / base.EkmanDepth
		want := math.Sqrt(1 + sign*relDelta)
		if math.Abs(ratio-want) > want*1e-5 {
			return ToleranceReport{}, fmt.Errorf("spiral: depth ratio %.6g != sqrt(K) %.6g", ratio, want)
		}
		items = append(items, Perturbation{
			Field:  "K",
			Delta:  sign * relDelta,
			V0:     res.V0,
			DV0:    res.V0 - base.V0,
			Depth:  res.EkmanDepth,
			DDepth: res.EkmanDepth - base.EkmanDepth,
		})
	}
	return ToleranceReport{BaseV0: base.V0, BaseDepth: base.EkmanDepth, Items: items}, nil
}

func MaxAbsV0Delta(r ToleranceReport) float64 {
	max := 0.0
	for _, p := range r.Items {
		d := math.Abs(p.DV0)
		if d > max {
			max = d
		}
	}
	return max
}
