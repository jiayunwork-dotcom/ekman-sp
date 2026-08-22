package spiral

import "math"

// Vec is a horizontal velocity or transport vector in the (east, north)
// plane. The same type is reused for point velocities (m/s) and for
// depth-integrated transports (m^2/s); only the caller's interpretation
// differs.
type Vec struct {
	East  float64
	North float64
}

// Add returns the component-wise sum of a and b.
func (a Vec) Add(b Vec) Vec {
	return Vec{East: a.East + b.East, North: a.North + b.North}
}

// Sub returns a - b.
func (a Vec) Sub(b Vec) Vec {
	return Vec{East: a.East - b.East, North: a.North - b.North}
}

// Scale multiplies both components by s.
func (a Vec) Scale(s float64) Vec {
	return Vec{East: a.East * s, North: a.North * s}
}

// Len returns the Euclidean magnitude of the vector.
func (a Vec) Len() float64 {
	return math.Hypot(a.East, a.North)
}

// Heading returns the compass heading in degrees clockwise from north.
func (a Vec) Heading() float64 {
	return PhaseToHeading(math.Atan2(a.North, a.East))
}

// FromPolar builds a vector of the given magnitude heading in degrees
// clockwise from north.
func FromPolar(mag, headingDeg float64) Vec {
	phi := HeadingToPhase(headingDeg)
	return Vec{East: mag * math.Cos(phi), North: mag * math.Sin(phi)}
}

// ApproxEqual reports whether two vectors agree within a relative and an
// absolute tolerance; useful for the cross-rule tests.
func ApproxEqual(a, b Vec, relTol, absTol float64) bool {
	d := a.Sub(b)
	return math.Abs(d.East) <= absTol+relTol*math.Max(math.Abs(a.East), math.Abs(b.East)) &&
		math.Abs(d.North) <= absTol+relTol*math.Max(math.Abs(a.North), math.Abs(b.North))
}
