package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSensitivityEndpoint(t *testing.T) {
	mux := New()
	payload := `{"case":{"tau":0.2,"latitude_deg":45,"K":0.05}}`
	req := httptest.NewRequest(http.MethodPost, "/api/sensitivity", strings.NewReader(payload))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestSweepEndpoint(t *testing.T) {
	mux := New()
	payload := `{"case":{"tau":0.2,"latitude_deg":45,"K":0.05},"param":"K","start":0.03,"stop":0.07,"step":0.02}`
	req := httptest.NewRequest(http.MethodPost, "/api/sweep", strings.NewReader(payload))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestHemisphereEndpoint(t *testing.T) {
	mux := New()
	payload := `{"tau":0.2,"wind_dir_deg":90,"latitude_deg":45,"K":0.05}`
	req := httptest.NewRequest(http.MethodPost, "/api/hemisphere", strings.NewReader(payload))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
}
