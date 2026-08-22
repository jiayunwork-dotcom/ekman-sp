// Package config loads an Ekman spiral case file (JSON) and turns it into a
// fully resolved, validated parameter set ready for the physics layer.
//
// The input format is deliberately small and explicit:
//
//	{
//	  "tau":           0.2,       // wind stress magnitude, N/m^2 (required)
//	  "wind_dir_deg":  90,        // wind blows toward this heading, deg from north
//	  "rho":           1025,      // seawater density, kg/m^3 (default 1025)
//	  "K":             0.05,      // eddy viscosity, m^2/s (default 0.05)
//	  "f":             1.0e-4,    // Coriolis parameter, 1/s (or use latitude_deg)
//	  "latitude_deg":  45,        // latitude in degrees, f = 2*Omega*sin(lat)
//	  "depth_m":       200,       // finite water depth, m (optional)
//	  "npoints":       21         // profile sample count (default 21)
//	}
//
// The Coriolis parameter may be given either directly as "f" or derived from
// "latitude_deg". Supplying both is allowed only when they agree; supplying
// neither is an error. Latitude 0 or f 0 lands on the equator, where the
// f-plane Ekman spiral is undefined, and is rejected.
//
// Validation enforces rho > 0, K > 0, tau > 0, f != 0 and, when a finite
// water depth is supplied, H >= De (the classic infinite-depth solution is
// not valid in a layer thinner than the Ekman depth; the user is told to
// declare a deeper case or drop the depth).
//
// The package owns no computation: after Load + Resolve the caller receives a
// spiral.Params and everything else happens in internal/spiral and
// internal/report.
package config
