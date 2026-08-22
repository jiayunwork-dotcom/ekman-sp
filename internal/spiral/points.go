package spiral

import "math"

// Point is one sample of the Ekman spiral at a given depth below the surface.
type Point struct {
	// Depth is distance below the sea surface in metres (positive down).
	Depth float64
	// U is the eastward velocity component in m/s.
	U float64
	// V is the northward velocity component in m/s.
	V float64
	// Speed is the horizontal current magnitude in m/s.
	Speed float64
	// Heading is the current direction in degrees clockwise from north.
	Heading float64
	// Veer is the signed turning of this current relative to the surface
	// current, in degrees; positive means clockwise.
	Veer float64
}

// PointAtDepth evaluates the analytic Ekman spiral at depthM metres below the
// surface. The complex current is
//
//	W(d) = V0 * exp(-d/delta) * exp(i*(phase0 - sign(f)*d/delta))
//
// where phase0 is the surface phase (wind heading minus the hemisphere
// deflection). The amplitude decays exponentially and the phase turns
// clockwise (northern hemisphere) or counter-clockwise (southern hemisphere)
// with depth.
func PointAtDepth(depthM, tau, rho, K, f, windHeadingDeg float64) (Point, error) {
	w, err := CurrentComplex(depthM, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return Point{}, err
	}
	east, north := unpack(w)

	phase0 := HeadingToPhase(SurfaceHeading(windHeadingDeg, f))
	phase := phase0 + RotationPhaseAtDepth(depthM, K, f)
	heading := PhaseToHeading(phase)
	// Veer is the cumulative, unwrapped turning from the surface current:
	// -sign(f)*depth/delta in radians, positive clockwise. Reporting the
	// accumulated turn (not the folded minimal angle) keeps the table readable
	// as the spiral goes around.
	veerDeg := -RadToDeg(phase - phase0)

	return Point{
		Depth:   depthM,
		U:       east,
		V:       north,
		Speed:   math.Hypot(east, north),
		Heading: heading,
		Veer:    veerDeg,
	}, nil
}

// SampleProfile produces nPoints evenly spaced samples from the surface down
// to profileDepthM. When profileDepthM is zero, DefaultProfileDepth is used.
// The returned slice always includes the surface sample at depth 0.
func SampleProfile(nPoints int, profileDepthM, tau, rho, K, f, windHeadingDeg float64) ([]Point, error) {
	if nPoints < 1 {
		nPoints = 1
	}
	if profileDepthM <= 0 {
		profileDepthM = DefaultProfileDepth(K, f)
	}
	points := make([]Point, 0, nPoints)
	for i := 0; i < nPoints; i++ {
		fraction := float64(i) / float64(nPoints-1)
		depth := fraction * profileDepthM
		p, err := PointAtDepth(depth, tau, rho, K, f, windHeadingDeg)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}

// IntegrateProfile approximates the depth-integrated transport of the
// sampled profile with the trapezoidal rule, returning the (east, north)
// components in m^2/s. It is used by tests to confirm that the integral of
// the spiral reproduces the analytic transport vector.
func IntegrateProfile(points []Point) Vec {
	if len(points) < 2 {
		return Vec{}
	}
	var east, north float64
	for i := 1; i < len(points); i++ {
		dz := points[i].Depth - points[i-1].Depth
		midU := 0.5 * (points[i].U + points[i-1].U)
		midV := 0.5 * (points[i].V + points[i-1].V)
		east += midU * dz
		north += midV * dz
	}
	return Vec{East: east, North: north}
}
