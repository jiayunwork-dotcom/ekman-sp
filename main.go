package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"ekman-sp/internal/config"
	"ekman-sp/internal/report"
	"ekman-sp/internal/server"
	"ekman-sp/internal/spiral"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

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
	case "serve":
		return runServe(args[1:], stderr)
	case "help", "-h", "--help":
		writeUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "ekman-sp: unknown command %q\n", args[0])
		writeUsage(stderr)
		return 2
	}
}

func runServe(args []string, stderr io.Writer) int {
	addr := ":8080"
	if len(args) > 0 {
		addr = args[0]
	}
	log.Printf("ekman-sp listening on %s", server.FormatAddr(addr))
	if err := server.ListenAndServe(server.Config{Addr: addr}); err != nil {
		fmt.Fprintf(stderr, "ekman-sp: %v\n", err)
		return 1
	}
	return 0
}

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

func writeUsage(w io.Writer) {
	fmt.Fprintln(w, "ekman-sp — steady-state ocean surface Ekman spiral calculator")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  ekman-sp serve [addr]                 HTTP API (default :8080)")
	fmt.Fprintln(w, "  ekman-sp profile <case-file.json>   load a case and print the spiral report")
	fmt.Fprintln(w, "  ekman-sp help                       show this help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "example:")
	fmt.Fprintln(w, "  ekman-sp serve :8080")
	fmt.Fprintln(w, "  ekman-sp profile example/midlat-wind.json")
}
