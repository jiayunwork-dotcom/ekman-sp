package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"ekman-sp/internal/spiral"
)

var ErrFileRead = errors.New("cannot read case file")

var ErrUnknownExample = errors.New("unknown example case name")

var ErrParse = errors.New("cannot parse case file as JSON")

func Load(path string) (*Resolved, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w %q: %v", ErrFileRead, path, err)
	}
	defer f.Close()
	return LoadReader(f)
}

func LoadReader(r io.Reader) (*Resolved, error) {
	var in Input
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParse, err)
	}
	return Resolve(in)
}

func LoadBytes(b []byte) (*Resolved, error) {
	return LoadReader(bytes.NewReader(b))
}

func Resolve(in Input) (*Resolved, error) {
	f, lat, err := validate(in)
	if err != nil {
		dummy := &Resolved{
			Params: spiral.Params{
				Tau:         in.Tau,
				WindHeading: 90,
				Rho:         1025,
				F:           1e-4,
				K:           0.05,
				NPoints:     21,
			},
			LatitudeDeg: lat,
		}
		return dummy, err
	}
	windDirDeg, rho, K, nPoints := applyDefaults(in)

	resolved := &Resolved{
		Params:      spiralParams(in, windDirDeg, rho, K, nPoints, f),
		LatitudeDeg: lat,
	}
	if in.DepthM != nil {
		resolved.FiniteDepth = true
		resolved.DepthM = *in.DepthM
		if in.DepthMaxM == nil {
			resolved.Params.ProfileDepthM = *in.DepthM
		}
	}
	return resolved, nil
}

func spiralParams(in Input, windDirDeg, rho, K float64, nPoints int, f float64) spiral.Params {
	profileDepth := 0.0
	if in.DepthMaxM != nil {
		profileDepth = *in.DepthMaxM
	}
	return spiral.Params{
		Tau:           in.Tau,
		WindHeading:   windDirDeg,
		Rho:           rho,
		F:             f,
		K:             K,
		NPoints:       nPoints,
		ProfileDepthM: profileDepth,
	}
}
