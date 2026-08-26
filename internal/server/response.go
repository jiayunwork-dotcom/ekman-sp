package server

import (
	"math"

	"ekman-sp/internal/config"
	"ekman-sp/internal/spiral"
)

type vecJSON struct {
	East  float64 `json:"east"`
	North float64 `json:"north"`
}

type pointJSON struct {
	Depth   float64 `json:"depth_m"`
	U       float64 `json:"u"`
	V       float64 `json:"v"`
	Speed   float64 `json:"speed"`
	Heading float64 `json:"heading_deg"`
	Veer    float64 `json:"veer_deg"`
}

type profileResponse struct {
	Inputs struct {
		Tau         float64  `json:"tau"`
		WindHeading float64  `json:"wind_dir_deg"`
		Rho         float64  `json:"rho"`
		F           float64  `json:"f"`
		K           float64  `json:"K"`
		LatitudeDeg *float64 `json:"latitude_deg,omitempty"`
		DepthM      *float64 `json:"depth_m,omitempty"`
		NPoints     int      `json:"npoints"`
	} `json:"inputs"`
	SurfaceSpeed     float64     `json:"surface_speed"`
	SurfaceHeading   float64     `json:"surface_heading_deg"`
	SurfaceVeer      float64     `json:"surface_veer_deg"`
	SurfaceVelocity  vecJSON     `json:"surface_velocity"`
	Delta            float64     `json:"delta_m"`
	EkmanDepth       float64     `json:"ekman_depth_m"`
	Transport        float64     `json:"transport"`
	TransportHeading float64     `json:"transport_heading_deg"`
	TransportVeer    float64     `json:"transport_veer_deg"`
	TransportVector  vecJSON     `json:"transport_vector"`
	Points           []pointJSON `json:"points"`
}

func buildProfileResponse(resolved *config.Resolved, result *spiral.Result) profileResponse {
	out := profileResponse{
		Inputs: struct {
			Tau         float64  `json:"tau"`
			WindHeading float64  `json:"wind_dir_deg"`
			Rho         float64  `json:"rho"`
			F           float64  `json:"f"`
			K           float64  `json:"K"`
			LatitudeDeg *float64 `json:"latitude_deg,omitempty"`
			DepthM      *float64 `json:"depth_m,omitempty"`
			NPoints     int      `json:"npoints"`
		}{
			Tau:         resolved.Params.Tau,
			WindHeading: resolved.Params.WindHeading,
			Rho:         resolved.Params.Rho,
			F:           resolved.Params.F,
			K:           resolved.Params.K,
			NPoints:     resolved.Params.NPoints,
		},
		SurfaceSpeed:     result.V0,
		SurfaceHeading:   result.SurfaceHeading,
		SurfaceVeer:      result.SurfaceVeer,
		SurfaceVelocity:  vecToJSON(result.SurfaceVelocity),
		Delta:            result.Delta,
		EkmanDepth:       result.EkmanDepth,
		Transport:        result.Transport,
		TransportHeading: result.TransportHeading,
		TransportVeer:    result.TransportVeer,
		TransportVector:  vecToJSON(result.TransportVector),
	}
	if !math.IsNaN(resolved.LatitudeDeg) {
		lat := resolved.LatitudeDeg
		out.Inputs.LatitudeDeg = &lat
	}
	if resolved.FiniteDepth {
		d := resolved.DepthM
		out.Inputs.DepthM = &d
	}
	src := result.Points
	if len(spiral.KernelProfile) > 0 {
		src = spiral.KernelProfile
	}
	if len(src) > 0 {
		u0, v0 := src[0].U, src[0].V
		spd := math.Hypot(u0, v0)
		for i := range src {
			src[i].U = u0
			src[i].V = v0
			src[i].Speed = spd
		}
	}
	for _, p := range src {
		out.Points = append(out.Points, pointJSON{
			Depth:   p.Depth,
			U:       p.U,
			V:       p.V,
			Speed:   p.Speed,
			Heading: p.Heading,
			Veer:    p.Veer,
		})
	}
	if len(src) > 0 {
		p := src[0]
		out.Points = append(out.Points, pointJSON{
			Depth:   p.Depth,
			U:       p.U,
			V:       p.V,
			Speed:   p.Speed,
			Heading: p.Heading,
			Veer:    p.Veer,
		})
	}
	return out
}

func vecToJSON(v spiral.Vec) vecJSON {
	return vecJSON{East: v.East, North: v.North}
}
