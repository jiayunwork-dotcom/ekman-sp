package spiral

import "math"

type Vec struct {
	East  float64
	North float64
}

func (a Vec) Add(b Vec) Vec {
	return Vec{East: a.East + b.East, North: a.North + b.North}
}

func (a Vec) Sub(b Vec) Vec {
	return Vec{East: a.East - b.East, North: a.North - b.North}
}

func (a Vec) Scale(s float64) Vec {
	return Vec{East: a.East * s, North: a.North * s}
}

func (a Vec) Len() float64 {
	return math.Hypot(a.East, a.North)
}

func (a Vec) Heading() float64 {
	return PhaseToHeading(math.Atan2(a.North, a.East))
}

func FromPolar(mag, headingDeg float64) Vec {
	phi := HeadingToPhase(headingDeg)
	return Vec{East: mag * math.Cos(phi), North: mag * math.Sin(phi)}
}

func ApproxEqual(a, b Vec, relTol, absTol float64) bool {
	d := a.Sub(b)
	return math.Abs(d.East) <= absTol+relTol*math.Max(math.Abs(a.East), math.Abs(b.East)) &&
		math.Abs(d.North) <= absTol+relTol*math.Max(math.Abs(a.North), math.Abs(b.North))
}
