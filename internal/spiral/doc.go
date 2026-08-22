// Package spiral implements the steady-state ocean surface Ekman spiral.
//
// The governing equations follow the conventional mid-latitude form with the
// vertical coordinate z measured positive upward from the sea surface:
//
//	K d^2u/dz^2 = -f v
//	K d^2v/dz^2 =  f u
//
// where u and v are the horizontal velocity components (east and north), K is
// the vertical eddy viscosity, and f is the Coriolis parameter. The surface
// boundary condition equates the vertical shear of the current to the imposed
// wind stress, and the deep boundary condition requires the flow to vanish as
// z -> -infinity.
//
// The analytic solution is a damped, turning spiral. With wind stress of
// magnitude tau and eddy viscosity K, the Ekman depth scale is
//
//	delta = sqrt(2K/|f|)          (m)
//	De    = pi * delta           (m, Ekman depth)
//
// and the surface current magnitude is
//
//	|V0| = tau / (rho * sqrt(K*|f|))
//
// In the northern hemisphere (f > 0) the surface current is deflected 45
// degrees to the right of the wind and the current vector turns clockwise with
// depth; in the southern hemisphere both deflections reverse. The depth-
// integrated volume transport is
//
//	Me = tau / (rho * |f|)        (m^2/s)
//
// directed 90 degrees to the right of the wind in the northern hemisphere and
// 90 degrees to the left in the southern hemisphere.
//
// The package exposes only pure functions over plain structs. All input
// validation (density, viscosity, equatorial f, finite depth) lives in
// internal/config; the functions here still guard the degenerate cases so the
// math can never silently produce nonsense.
package spiral
