// Command ekman-sp is a command-line calculator for the steady-state ocean
// surface Ekman spiral. It reads a wind-stress case file (JSON) and prints
// the surface current magnitude and direction, the Ekman depth, the
// depth-integrated Ekman transport and the turning velocity profile.
//
// Usage:
//
//	ekman-sp profile example/midlat-wind.json
//	ekman-sp help
//
// The profile subcommand loads the JSON case, validates every physical input
// (density, eddy viscosity, Coriolis parameter, optional finite water depth)
// and renders the report. Any invalid input is printed to standard error and
// exits non-zero; nothing is silently ignored.
package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"ekman-sp/internal/config"
	"ekman-sp/internal/report"
	"ekman-sp/internal/spiral"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the CLI against the given arguments and output sinks. It
// returns the process exit code: 0 on success, 1 on any input or computation
// error, 2 on usage errors. Keeping this separable from main makes the CLI
// behaviour directly testable.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writeUsage(stderr)
		return 2
	}
	switch args[0] {
	case "profile":
		if len(args) != 2 {
			fmt.Fprintln(stderr, "usage: ekman-sp profile <case-file.json>")
			return 2
		}
		return runProfile(args[1], stdout, stderr)
	case "help", "-h", "--help":
		writeUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "ekman-sp: unknown command %q\n", args[0])
		writeUsage(stderr)
		return 2
	}
}

// runProfile loads, validates and renders a single case file.
func runProfile(path string, stdout, stderr io.Writer) int {
	resolved, err := config.Load(path)
	if err != nil {
		fmt.Fprintf(stderr, "ekman-sp: %v\n", err)
		return 1
	}

	result, err := spiral.Compute(resolved.Params)
	if err != nil {
		fmt.Fprintf(stderr, "ekman-sp: %v\n", err)
		return 1
	}

	meta := report.Meta{
		Tau:         resolved.Params.Tau,
		WindHeading: resolved.Params.WindHeading,
		Rho:         resolved.Params.Rho,
		F:           resolved.Params.F,
		K:           resolved.Params.K,
		Latitude:    resolved.LatitudeDeg,
		HasLatitude: !math.IsNaN(resolved.LatitudeDeg),
		FiniteDepth: resolved.FiniteDepth,
		DepthM:      resolved.DepthM,
		NPoints:     resolved.Params.NPoints,
	}
	if err := report.Render(stdout, &result, meta); err != nil {
		fmt.Fprintf(stderr, "ekman-sp: cannot write report: %v\n", err)
		return 1
	}
	return 0
}

// writeUsage prints the command synopsis to w.
func writeUsage(w io.Writer) {
	fmt.Fprintln(w, "ekman-sp — steady-state ocean surface Ekman spiral calculator")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  ekman-sp profile <case-file.json>   load a case and print the spiral report")
	fmt.Fprintln(w, "  ekman-sp help                       show this help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "example:")
	fmt.Fprintln(w, "  ekman-sp profile example/midlat-wind.json")
}
