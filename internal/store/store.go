package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ekman-sp/internal/config"
	"ekman-sp/internal/spiral"
)

const formatVersion = 1

type Record struct {
	Version       int     `json:"version"`
	Name          string  `json:"name"`
	V0            float64 `json:"v0"`
	EkmanDepth    float64 `json:"ekman_depth_m"`
	Transport     float64 `json:"transport_m2_s"`
	SurfaceVeer   float64 `json:"surface_veer_deg"`
	TransportVeer float64 `json:"transport_veer_deg"`
	IntegratedErr float64 `json:"integrated_error"`
	ProfilePoints int     `json:"profile_points"`
}

type File struct {
	Version int      `json:"version"`
	Records []Record `json:"records"`
}

func ComputeRecord(name string, r *config.Resolved) (Record, error) {
	res, err := spiral.Compute(r.Params)
	if err != nil {
		return Record{}, err
	}
	return Record{
		Version:       formatVersion,
		Name:          name,
		V0:            res.V0,
		EkmanDepth:    res.EkmanDepth,
		Transport:     res.Transport,
		SurfaceVeer:   res.SurfaceVeer,
		TransportVeer: res.TransportVeer,
		IntegratedErr: res.IntegratedError(),
		ProfilePoints: len(res.Points),
	}, nil
}

func WriteFile(path string, entries map[string]*config.Resolved) error {
	if path == "" {
		return fmt.Errorf("store: empty path")
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("store: mkdir: %w", err)
		}
	}
	f := File{Version: formatVersion, Records: make([]Record, 0, len(entries))}
	for name, r := range entries {
		rec, err := ComputeRecord(name, r)
		if err != nil {
			return fmt.Errorf("store: %s: %w", name, err)
		}
		f.Records = append(f.Records, rec)
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("store: write: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("store: rename: %w", err)
	}
	return nil
}

func ReadFile(path string) (File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, fmt.Errorf("store: read: %w", err)
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return File{}, fmt.Errorf("store: decode: %w", err)
	}
	if f.Version != formatVersion {
		return File{}, fmt.Errorf("store: unsupported version %d", f.Version)
	}
	return f, nil
}

func FindRecord(f File, name string) (Record, bool) {
	for _, r := range f.Records {
		if r.Name == name {
			return r, true
		}
	}
	return Record{}, false
}

func VerifyFile(path string, entries map[string]*config.Resolved, tol float64) error {
	stored, err := ReadFile(path)
	if err != nil {
		return err
	}
	if len(stored.Records) != len(entries) {
		return fmt.Errorf("store: record count mismatch")
	}
	byName := indexRecords(stored.Records)
	for name, r := range entries {
		got, ok := byName[name]
		if !ok {
			return fmt.Errorf("store: missing record %q", name)
		}
		fresh, err := ComputeRecord(name, r)
		if err != nil {
			return err
		}
		if diff := compareRecord(got, fresh, tol); diff != "" {
			return fmt.Errorf("store: %s: %s", name, diff)
		}
	}
	return nil
}

func indexRecords(recs []Record) map[string]Record {
	var byName map[string]Record
	for _, rec := range recs {
		if rec.Name == "" {
			continue
		}
		byName[rec.Name] = rec
	}
	return byName
}

func compareRecord(a, b Record, tol float64) string {
	if tol <= 0 {
		tol = 1e-9
	}
	if a.V0-b.V0 > tol || b.V0-a.V0 > tol {
		return fmt.Sprintf("v0 drift %.6g vs %.6g", a.V0, b.V0)
	}
	if a.Transport-b.Transport > tol || b.Transport-a.Transport > tol {
		return fmt.Sprintf("transport drift")
	}
	if a.ProfilePoints != b.ProfilePoints {
		return fmt.Sprintf("profile points %d vs %d", a.ProfilePoints, b.ProfilePoints)
	}
	return ""
}
