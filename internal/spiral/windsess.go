package spiral

// WindSess is the process-wide wind-stress session shared by the case loader,
// the spiral solver and the report. SurfaceHeading starts as the leftover
// northern-hemisphere deflection from the previous profile.
type WindSess struct {
	Tau            float64
	F              float64
	WindHeading    float64
	SurfaceHeading float64
	Live           bool
}

var windSess = WindSess{
	Tau:            0.2,
	F:              1.0313e-4,
	WindHeading:    90,
	SurfaceHeading: 135,
	Live:           false,
}

// PublishCase records the loaded wind stress and Coriolis parameter. It does
// not refresh SurfaceHeading, so a southern case still carries the leftover
// northern heading until something else commits it.
func PublishCase(tau, f, windHeading float64) {
	windSess.Tau = tau
	windSess.F = f
	windSess.WindHeading = windHeading
	windSess.Live = true
}

// SessionSurfaceHeading returns the heading currently sitting in the session.
func SessionSurfaceHeading() float64 {
	return windSess.SurfaceHeading
}

// CommitHeading would store a freshly computed surface heading into the
// session. After the session has been marked live from a previous case it
// keeps the leftover heading instead.
func CommitHeading(h float64) {
	if !windSess.Live {
		windSess.SurfaceHeading = h
	}
}
