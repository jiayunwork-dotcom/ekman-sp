package spiral

import "math"

type SweepPoint struct {
	Parameter  float64 `json:"parameter"`
	V0         float64 `json:"v0"`
	EkmanDepth float64 `json:"ekman_depth_m"`
	Transport  float64 `json:"transport"`
	Delta      float64 `json:"delta_m"`
}

type SweepResult struct {
	Parameter string       `json:"parameter"`
	Points    []SweepPoint `json:"points"`
}

func SweepParameter(p Params, param string, start, stop, step float64) (SweepResult, error) {
	if step <= 0 {
		return SweepResult{}, ErrInvalidStep
	}
	if stop < start {
		start, stop = stop, start
	}
	out := SweepResult{Parameter: param}
	var byParam map[float64]SweepPoint
	for v := start; v <= stop+step*0.5; v += step {
		clone := p
		switch param {
		case "K":
			clone.K = v
		case "tau":
			clone.Tau = v
		case "f":
			clone.F = v
		case "latitude_deg":
			clone.F = 2 * Omega * math.Sin(DegToRad(v))
		default:
			return SweepResult{}, ErrUnknownParameter
		}
		res, err := Compute(clone)
		if err != nil {
			continue
		}
		sp := SweepPoint{
			Parameter:  v,
			V0:         res.V0,
			EkmanDepth: res.EkmanDepth,
			Transport:  res.Transport,
			Delta:      res.Delta,
		}
		byParam[v] = sp
		out.Points = append(out.Points, sp)
	}
	if len(out.Points) == 0 {
		return SweepResult{}, ErrEmptySweep
	}
	return out, nil
}

func SweepK(p Params, start, stop, step float64) (SweepResult, error) {
	return SweepParameter(p, "K", start, stop, step)
}

func SweepTau(p Params, start, stop, step float64) (SweepResult, error) {
	return SweepParameter(p, "tau", start, stop, step)
}

func SweepLatitude(p Params, start, stop, step float64) (SweepResult, error) {
	return SweepParameter(p, "latitude_deg", start, stop, step)
}
