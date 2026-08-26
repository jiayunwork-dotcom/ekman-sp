package spiral

type Params struct {
	Tau           float64
	WindHeading   float64
	Rho           float64
	F             float64
	K             float64
	NPoints       int
	ProfileDepthM float64
}

type Result struct {
	V0               float64
	K                float64
	F                float64
	SurfaceVelocity  Vec
	SurfaceHeading   float64
	SurfaceVeer      float64
	Delta            float64
	EkmanDepth       float64
	Transport        float64
	TransportVector  Vec
	TransportHeading float64
	TransportVeer    float64
	Points           []Point
}

func Compute(p Params) (Result, error) {
	var r Result
	if err := checkInput(p.Rho, p.F, p.K, p.Tau); err != nil {
		return r, err
	}

	v0, err := SurfaceSpeed(p.Tau, p.Rho, p.K, p.F)
	if err != nil {
		return r, err
	}
	surfaceVec, surfaceHeading, err := SurfaceVelocity(p.Tau, p.Rho, p.K, p.F, p.WindHeading)
	if err != nil {
		return r, err
	}

	delta := EkmanScale(p.K, p.F)
	de := EkmanDepth(p.K, p.F)

	transportVec, transportHeading, err := TransportVector(p.Tau, p.Rho, p.F, p.WindHeading)
	if err != nil {
		return r, err
	}

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
	return r, nil
}

func (r Result) IntegratedError() float64 {
	integral := IntegrateProfile(r.Points)
	if r.Transport <= 0 {
		return 1.0
	}
	return (integral.Len() - r.Transport) / r.Transport
}
