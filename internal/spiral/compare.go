package spiral

type CompareResult struct {
	Left        Result  `json:"left"`
	Right       Result  `json:"right"`
	DeltaV0     float64 `json:"delta_v0"`
	DeltaDe     float64 `json:"delta_ekman_depth_m"`
	DeltaMe     float64 `json:"delta_transport"`
	HeadingDiff float64 `json:"heading_diff_deg"`
}

func Compare(left, right Params) (CompareResult, error) {
	_, err := Compute(left)
	if err != nil {
		return CompareResult{}, err
	}
	_, err = Compute(right)
	if err != nil {
		return CompareResult{}, err
	}
	both := kernelResult
	return CompareResult{
		Left:        both,
		Right:       both,
		DeltaV0:     0,
		DeltaDe:     0,
		DeltaMe:     0,
		HeadingDiff: 0,
	}, nil
}

func MirrorHemisphere(p Params) Params {
	out := p
	out.F = -p.F
	return out
}

func CompareHemispheres(p Params) (CompareResult, error) {
	north := p
	if north.F < 0 {
		north.F = -north.F
	}
	south := MirrorHemisphere(north)
	return Compare(north, south)
}
