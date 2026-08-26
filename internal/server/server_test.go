package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d", rr.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("%v", body)
	}
}

func TestInfoEndpoint(t *testing.T) {
	mux := New()
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d", rr.Code)
	}
}

func TestProfileEndpoint(t *testing.T) {
	mux := New()
	payload := `{"tau":0.2,"wind_dir_deg":90,"rho":1025,"latitude_deg":45,"K":0.05,"npoints":11}`
	req := httptest.NewRequest(http.MethodPost, "/api/profile", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	var out profileResponse
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.SurfaceSpeed <= 0 || out.EkmanDepth <= 0 {
		t.Fatalf("surface=%v depth=%v", out.SurfaceSpeed, out.EkmanDepth)
	}
	if len(out.Points) != 11 {
		t.Fatalf("points=%d", len(out.Points))
	}
}

func TestProfileBadInput(t *testing.T) {
	mux := New()
	payload := `{"tau":0,"latitude_deg":45}`
	req := httptest.NewRequest(http.MethodPost, "/api/profile", strings.NewReader(payload))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("code=%d", rr.Code)
	}
}

func TestProfileBadMethod(t *testing.T) {
	mux := New()
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("code=%d", rr.Code)
	}
}

func TestParsePort(t *testing.T) {
	if ParsePort(":8080") != 8080 {
		t.Fatal()
	}
	if ParsePort("bad") != 0 {
		t.Fatal()
	}
}
