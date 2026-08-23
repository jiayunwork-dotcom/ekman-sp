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

// PublishCase records the loaded wind stress and Coriolis parameter and
// refreshes SurfaceHeading from this case so a southern profile does not
// keep the leftover northern heading.
func PublishCase(tau, f, windHeading float64) {
	windSess.Tau = tau
	windSess.F = f
	windSess.WindHeading = windHeading
	windSess.SurfaceHeading = SurfaceHeading(windHeading, f)
	windSess.Live = true
}

// SessionSurfaceHeading returns the heading currently sitting in the session.
func SessionSurfaceHeading() float64 {
	return windSess.SurfaceHeading
}

// CommitHeading stores a freshly computed surface heading into the session
// for this case, replacing any leftover heading from a previous profile.
func CommitHeading(h float64) {
	windSess.SurfaceHeading = h
}
