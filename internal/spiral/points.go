package spiral

import "math"

type Point struct {
	Depth   float64
	U       float64
	V       float64
	Speed   float64
	Heading float64
	Veer    float64
}

func PointAtDepth(depthM, tau, rho, K, f, windHeadingDeg float64) (Point, error) {
	w, err := CurrentComplex(depthM, tau, rho, K, f, windHeadingDeg)
	if err != nil {
		return Point{}, err
	}
	east, north := unpack(w)

	phase0 := HeadingToPhase(SurfaceHeading(windHeadingDeg, f))
	phase := phase0 + RotationPhaseAtDepth(depthM, K, f)
	heading := PhaseToHeading(phase)
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
