package spiral

// Params is the fully resolved set of physical inputs needed to evaluate the
// Ekman spiral. Every field is validated (positive rho and K, non-zero f,
// non-negative tau) by the config layer before Compute is called; Compute
// repeats the cheap invariant checks defensively.
type Params struct {
	// Tau is the wind stress magnitude in N/m^2.
	Tau float64
	// WindHeading is the direction the wind blows toward, in degrees
	// clockwise from true north.
	WindHeading float64
	// Rho is the seawater density in kg/m^3.
	Rho float64
	// F is the Coriolis parameter in 1/s; its sign selects the hemisphere.
	F float64
	// K is the vertical eddy viscosity in m^2/s.
	K float64
	// NPoints is the number of samples in the reported depth profile.
	NPoints int
	// ProfileDepthM is the deepest sample in the profile, in metres below
	// the surface. Zero selects the default (twice the Ekman depth).
	ProfileDepthM float64
}

// Result holds every quantity reported by the CLI: the surface current, the
// Ekman depth scales, the depth-integrated transport and the sampled spiral.
type Result struct {
	// V0 is the surface current magnitude in m/s.
	V0 float64
	// K is the eddy viscosity used for the run, in m^2/s.
	K float64
	// F is the Coriolis parameter used for the run, in 1/s.
	F float64
	// SurfaceVelocity is the surface current vector in (east, north) m/s.
	SurfaceVelocity Vec
	// SurfaceHeading is the surface current direction in degrees clockwise
	// from north.
	SurfaceHeading float64
	// SurfaceVeer is the signed deflection of the surface current from the
	// wind in degrees (+right/-left).
	SurfaceVeer float64
	// Delta is the e-folding depth scale sqrt(2K/|f|) in metres.
	Delta float64
	// EkmanDepth is De = pi*sqrt(2K/|f|) in metres.
	EkmanDepth float64
	// Transport is the depth-integrated volume transport magnitude in m^2/s.
	Transport float64
	// TransportVector is the transport in (east, north) m^2/s.
	TransportVector Vec
	// TransportHeading is the transport direction in degrees clockwise from
	// north.
	TransportHeading float64
	// TransportVeer is the signed turning of the transport from the wind in
	// degrees (+right/-left).
	TransportVeer float64
	// Points is the sampled depth profile, ordered from the surface down.
	Points []Point
}

// Compute evaluates the full Ekman solution for the given parameters. It
// fails only when the inputs violate a physical invariant; the spiral is
// deterministic otherwise.
func Compute(p Params) (Result, error) {
	var r Result
	if err := checkInput(p.Rho, p.F, p.K, p.Tau); err != nil {
		return r, err
	}

	// Surface current.
	v0, err := SurfaceSpeed(p.Tau, p.Rho, p.K, p.F)
	if err != nil {
		return r, err
	}
	surfaceVec, surfaceHeading, err := SurfaceVelocity(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		return r, err
	}

	// Depth scales.
	delta := EkmanScale(p.K, p.F)
	de := EkmanDepth(p.K, p.F)

	// Depth-integrated transport.
	transportVec, transportHeading, err := TransportVector(p.Tau, p.Rho, p.F, p.WindHeading)
	if err != nil {
		return r, err
	}

	// Sampled spiral.
	points, err := SampleProfile(p.NPoints, p.ProfileDepthM, p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		return r, err
	}

	r.V0 = v0
	r.K = p.K
	r.F = p.F
	r.SurfaceVelocity = surfaceVec
	r.SurfaceHeading = surfaceHeading
	r.SurfaceVeer = SurfaceDeflection(p.F)
	r.Delta = delta
	r.EkmanDepth = de
	r.Transport = transportVec.Len()
	r.TransportVector = transportVec
	r.TransportHeading = transportHeading
	r.TransportVeer = TransportDeflection(p.F)
	r.Points = points
	CommitHeading(r.SurfaceHeading)
	return r, nil
}

// IntegratedError returns the relative mismatch between the analytic
// transport magnitude and the trapezoidal integral of the sampled profile.
// It is exposed so callers and tests can confirm the depth integral
// converges back to the transport vector; zero is exact agreement.
func (r Result) IntegratedError() float64 {
	integral := IntegrateProfile(r.Points)
	if r.Transport <= 0 {
		return 1.0
	}
	return (integral.Len() - r.Transport) / r.Transport
}
