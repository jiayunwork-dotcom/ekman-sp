package config

import "ekman-sp/internal/spiral"

// Input mirrors the JSON case-file format. Pointer fields distinguish "absent
// in the file" from "present with value zero", which matters because rho 0
// must be an error while an absent rho takes the seawater default.
type Input struct {
	// Tau is the wind stress magnitude in N/m^2. It is required; a missing
	// or non-positive value is an error because a zero stress has no spiral.
	Tau float64 `json:"tau"`

	// WindDirDeg is the heading (degrees clockwise from north) toward which
	// the wind blows. Absent defaults to an easterly wind (90).
	WindDirDeg *float64 `json:"wind_dir_deg"`

	// Rho is the seawater density in kg/m^3. Absent defaults to 1025.
	Rho *float64 `json:"rho"`

	// K is the vertical eddy viscosity in m^2/s. Absent defaults to 0.05.
	K *float64 `json:"K"`

	// F is the Coriolis parameter in 1/s, positive for the northern
	// hemisphere. Either F or LatitudeDeg must be given.
	F *float64 `json:"f"`

	// LatitudeDeg is the latitude in degrees; f is derived as
	// 2*Omega*sin(latitude). Zero latitude means the equator and is rejected.
	LatitudeDeg *float64 `json:"latitude_deg"`

	// DepthM is the finite water depth in metres below the surface. Absent
	// means deep water and the classic infinite-depth solution. A finite
	// depth smaller than the Ekman depth is rejected.
	DepthM *float64 `json:"depth_m"`

	// NPoints is the number of samples in the depth profile. Absent defaults
	// to 21.
	NPoints *int `json:"npoints"`

	// DepthMaxM caps the deepest sample of the profile. Absent defaults to
	// twice the Ekman depth.
	DepthMaxM *float64 `json:"depth_max_m"`
}

// Resolved is the validated interpretation of an Input: pointer fields are
// resolved to concrete values and the Coriolis parameter is pinned down.
type Resolved struct {
	// Params is the physics-ready parameter set.
	Params spiral.Params
	// LatitudeDeg is the latitude that produced F, or NaN when the case
	// specified f directly.
	LatitudeDeg float64
	// FiniteDepth reports whether the case carries an explicit water depth.
	FiniteDepth bool
	// DepthM is the explicit water depth in metres (zero when not finite).
	DepthM float64
}
