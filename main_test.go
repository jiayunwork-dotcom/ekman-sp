package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCase(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "case.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write case: %v", err)
	}
	return path
}

func runCLI(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errBuf bytes.Buffer
	code = run(args, &out, &errBuf)
	return code, out.String(), errBuf.String()
}

func TestRunProfileOK(t *testing.T) {
	path := writeCase(t, `{"tau":0.2,"wind_dir_deg":90,"rho":1025,"latitude_deg":45,"K":0.05,"npoints":21}`)
	code, out, errOut := runCLI(t, "profile", path)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 (stderr: %s)", code, errOut)
	}
	for _, want := range []string{"Ekman depth De", "surface speed |V0|", "transport |Me|", "depth profile"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing %q", want)
		}
	}
	if errOut != "" {
		t.Errorf("stderr not empty: %s", errOut)
	}
}

func TestRunProfileMissingFile(t *testing.T) {
	code, _, errOut := runCLI(t, "profile", "/nonexistent/case.json")
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(errOut, "cannot read case file") {
		t.Errorf("stderr = %q, want read error", errOut)
	}
}

func TestRunProfileMalformedJSON(t *testing.T) {
	path := writeCase(t, `{"tau":`)
	code, _, errOut := runCLI(t, "profile", path)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(errOut, "cannot parse") {
		t.Errorf("stderr = %q, want parse error", errOut)
	}
}

func TestRunProfileZeroDensity(t *testing.T) {
	path := writeCase(t, `{"tau":0.2,"rho":0,"latitude_deg":45}`)
	code, _, errOut := runCLI(t, "profile", path)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(errOut, "rho") {
		t.Errorf("stderr = %q, want density error", errOut)
	}
}

func TestRunProfileEquator(t *testing.T) {
	path := writeCase(t, `{"tau":0.2,"latitude_deg":0}`)
	code, _, errOut := runCLI(t, "profile", path)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(errOut, "equator") {
		t.Errorf("stderr = %q, want equator error", errOut)
	}
}

func TestRunNoArgs(t *testing.T) {
	code, out, errOut := runCLI(t)
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if errOut == "" {
		t.Error("usage must go to stderr")
	}
	if out != "" {
		t.Errorf("stdout should be empty, got %q", out)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	code, _, errOut := runCLI(t, "bogus")
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if !strings.Contains(errOut, "unknown command") {
		t.Errorf("stderr = %q, want unknown command", errOut)
	}
}

func TestRunProfileTooFewArgs(t *testing.T) {
	code, _, errOut := runCLI(t, "profile")
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if !strings.Contains(errOut, "usage: ekman-sp profile") {
		t.Errorf("stderr = %q, want usage hint", errOut)
	}
}

func TestRunHelp(t *testing.T) {
	for _, arg := range []string{"help", "-h", "--help"} {
		code, out, errOut := runCLI(t, arg)
		if code != 0 {
			t.Errorf("%s: exit code = %d, want 0", arg, code)
		}
		if !strings.Contains(out, "usage:") {
			t.Errorf("%s: stdout = %q, want usage", arg, out)
		}
		if errOut != "" {
			t.Errorf("%s: stderr = %q, want empty", arg, errOut)
		}
	}
}

func TestRunProfileSouthHemisphere(t *testing.T) {
	path := writeCase(t, `{"tau":0.2,"wind_dir_deg":90,"latitude_deg":-45}`)
	code, out, errOut := runCLI(t, "profile", path)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 (stderr: %s)", code, errOut)
	}
	if !strings.Contains(out, "45.0 deg (NE)") {
		t.Errorf("southern surface heading missing from report:\n%s", out)
	}
}
