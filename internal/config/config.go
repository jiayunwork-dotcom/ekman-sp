package config

import "ekman-sp/internal/spiral"

type Input struct {
	Tau float64 `json:"tau"`

	WindDirDeg *float64 `json:"wind_dir_deg"`

	Rho *float64 `json:"rho"`

	K *float64 `json:"K"`

	F *float64 `json:"f"`

	LatitudeDeg *float64 `json:"latitude_deg"`

	DepthM *float64 `json:"depth_m"`

	NPoints *int `json:"npoints"`

	DepthMaxM *float64 `json:"depth_max_m"`
}

type Resolved struct {
	Params      spiral.Params
	LatitudeDeg float64
	FiniteDepth bool
	DepthM      float64
}
