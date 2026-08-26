package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ekman-sp/internal/config"
	"ekman-sp/internal/spiral"
)

type sensitivityRequest struct {
	Case config.Input `json:"case"`
}

type sweepRequest struct {
	Case  config.Input `json:"case"`
	Param string       `json:"param"`
	Start float64      `json:"start"`
	Stop  float64      `json:"stop"`
	Step  float64      `json:"step"`
}

type compareRequest struct {
	Left  config.Input `json:"left"`
	Right config.Input `json:"right"`
}

func handleSensitivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, fmt.Errorf("POST only"))
		return
	}
	var req sensitivityRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	resolved, err := config.Resolve(req.Case)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	sens, err := spiral.AnalyzeSensitivity(resolved.Params)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sens)
}

func handleSweep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, fmt.Errorf("POST only"))
		return
	}
	var req sweepRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	resolved, err := config.Resolve(req.Case)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	param := req.Param
	if param == "" {
		param = "K"
	}
	if req.Step <= 0 {
		req.Step = 0.01
	}
	if req.Stop <= req.Start {
		req.Stop = req.Start + req.Step*10
	}
	sweep, err := spiral.SweepParameter(resolved.Params, param, req.Start, req.Stop, req.Step)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sweep)
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, fmt.Errorf("POST only"))
		return
	}
	var req compareRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	left, err := config.Resolve(req.Left)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	right, err := config.Resolve(req.Right)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	cmp, err := spiral.Compare(left.Params, right.Params)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cmp)
}

func handleHemisphereCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, fmt.Errorf("POST only"))
		return
	}
	var in config.Input
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	resolved, err := config.Resolve(in)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	cmp, err := spiral.CompareHemispheres(resolved.Params)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cmp)
}
