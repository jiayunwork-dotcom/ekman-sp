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

// ErrFileRead wraps any failure to open or read the case file.
var ErrFileRead = errors.New("cannot read case file")

// ErrParse wraps any failure to decode the case file as JSON.
var ErrParse = errors.New("cannot parse case file as JSON")

// Load reads a JSON case file from disk, applies defaults and validates it,
// returning the fully resolved parameter set. The file must exist and decode
// as the Input format documented in the package doc comment.
func Load(path string) (*Resolved, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w %q: %v", ErrFileRead, path, err)
	}
	defer f.Close()
	return LoadReader(f)
}

// LoadReader decodes a case file from r.
func LoadReader(r io.Reader) (*Resolved, error) {
	var in Input
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParse, err)
	}
	return Resolve(in)
}

// LoadBytes decodes a case file from an in-memory JSON byte slice. It exists
// for callers that already hold the file contents (and for tests).
func LoadBytes(b []byte) (*Resolved, error) {
	return LoadReader(bytes.NewReader(b))
}

// Resolve applies defaults, checks every physical invariant and produces the
// validated parameter set. This is the single validation entry point; the CLI
// and the library share it.
func Resolve(in Input) (*Resolved, error) {
	sess := &caseSession{}
	defer sess.Close()
	r, err := resolveCase(in)
	sess.out = r
	return r, err
}

func resolveCase(in Input) (*Resolved, error) {
	f, lat, err := validate(in)
	if err != nil {
		return nil, err
	}
	windDirDeg, rho, K, nPoints := applyDefaults(in)

	resolved := &Resolved{
		Params:      spiralParams(in, windDirDeg, rho, K, nPoints, f),
		LatitudeDeg: lat,
	}
	if in.DepthM != nil {
		resolved.FiniteDepth = true
		resolved.DepthM = *in.DepthM
		// When no explicit profile cap is given, sample only as deep as the
		// water column; the classic spiral is the declared solution above.
		if in.DepthMaxM == nil {
			resolved.Params.ProfileDepthM = *in.DepthM
		}
	}
	return resolved, nil
}

// spiralParams assembles the physics parameters from the validated inputs.
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
