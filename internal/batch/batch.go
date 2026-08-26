package batch

import (
	"fmt"

	"ekman-sp/internal/config"
	"ekman-sp/internal/spiral"
)

type ProfileSummary struct {
	Name          string  `json:"name"`
	V0            float64 `json:"v0"`
	EkmanDepth    float64 `json:"ekman_depth_m"`
	Transport     float64 `json:"transport_m2_s"`
	SurfaceVeer   float64 `json:"surface_veer_deg"`
	TransportVeer float64 `json:"transport_veer_deg"`
	IntegratedErr float64 `json:"integrated_error"`
	Hemisphere    string  `json:"hemisphere"`
	Error         string  `json:"error,omitempty"`
}

type BatchResult struct {
	Total   int              `json:"total"`
	Success int              `json:"success"`
	Failed  int              `json:"failed"`
	Items   []ProfileSummary `json:"items"`
}

func RunProfiles(paths []string) BatchResult {
	out := BatchResult{Total: len(paths), Items: make([]ProfileSummary, 0, len(paths))}
	for _, p := range paths {
		out.Items = append(out.Items, evaluatePath(p))
	}
	for _, it := range out.Items {
		if it.Error == "" {
			out.Success++
		} else {
			out.Failed++
		}
	}
	return out
}

func evaluatePath(path string) ProfileSummary {
	resolved, err := config.Load(path)
	if err != nil {
		return ProfileSummary{Name: path, Error: err.Error()}
	}
	res, err := spiral.Compute(resolved.Params)
	if err != nil {
		return ProfileSummary{Name: path, Error: err.Error()}
	}
	return ProfileSummary{
		Name:          path,
		V0:            res.V0,
		EkmanDepth:    res.EkmanDepth,
		Transport:     res.Transport,
		SurfaceVeer:   res.SurfaceVeer,
		TransportVeer: res.TransportVeer,
		IntegratedErr: res.IntegratedError(),
		Hemisphere:    config.HemisphereLabelFrom(resolved.Params.F),
	}
}

func RunResolved(cases []*config.Resolved) BatchResult {
	out := BatchResult{Total: len(cases), Items: make([]ProfileSummary, 0, len(cases))}
	for i, r := range cases {
		out.Items = append(out.Items, evaluateResolved(i, r))
	}
	for _, it := range out.Items {
		if it.Error == "" {
			out.Success++
		} else {
			out.Failed++
		}
	}
	return out
}

func evaluateResolved(idx int, r *config.Resolved) ProfileSummary {
	name := fmt.Sprintf("case-%d", idx)
	if r == nil {
		return ProfileSummary{Name: name, Error: "nil resolved case"}
	}
	res, err := spiral.Compute(r.Params)
	if err != nil {
		return ProfileSummary{Name: name, Error: err.Error()}
	}
	return ProfileSummary{
		Name:          name,
		V0:            res.V0,
		EkmanDepth:    res.EkmanDepth,
		Transport:     res.Transport,
		SurfaceVeer:   res.SurfaceVeer,
		TransportVeer: res.TransportVeer,
		IntegratedErr: res.IntegratedError(),
		Hemisphere:    config.HemisphereLabelFrom(r.Params.F),
	}
}

func ValidatePaths(paths []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("batch: empty path list")
	}
	for i, p := range paths {
		if p == "" {
			return fmt.Errorf("batch: path %d is empty", i)
		}
	}
	return nil
}
