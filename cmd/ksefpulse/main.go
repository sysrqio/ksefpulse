package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sysrqio/ksefpulse/internal/fa3"
	"github.com/sysrqio/ksefpulse/internal/reconcile"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "version":
		fmt.Println("ksefpulse", version)
	case "validate":
		runValidate(os.Args[2:])
	case "audit":
		runAudit(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "ksefpulse — KSeF reconciliation & FA(3) pre-audit\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n  ksefpulse audit --state <file> [--fa3 <file>] [--output report.json] [--min-score 70]\n  ksefpulse validate --fa3 <file>\n  ksefpulse version\n")
}

func runValidate(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	fa3Path := fs.String("fa3", "", "FA(3) XML path")
	_ = fs.Parse(args)
	if *fa3Path == "" {
		fmt.Fprintln(os.Stderr, "--fa3 required")
		os.Exit(1)
	}
	findings, err := fa3.ValidateFile(*fa3Path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, f := range findings {
		fmt.Printf("[%s] %s: %s\n", f.Severity, f.Code, f.Message)
	}
	if len(findings) > 0 {
		os.Exit(2)
	}
}

func runAudit(args []string) {
	fs := flag.NewFlagSet("audit", flag.ExitOnError)
	statePath := fs.String("state", "", "submission-state.json")
	fa3Path := fs.String("fa3", "", "optional FA(3) XML")
	output := fs.String("output", "reconciliation-report.json", "output report path")
	format := fs.String("format", "json", "json or table")
	minScore := fs.Int("min-score", 70, "minimum readiness score")
	_ = fs.Parse(args)
	if *statePath == "" {
		fmt.Fprintln(os.Stderr, "--state required")
		os.Exit(1)
	}
	st, err := reconcile.LoadState(*statePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rep := reconcile.Analyze(st)
	if *fa3Path != "" {
		ff, err := fa3.ValidateFile(*fa3Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, f := range ff {
			rep.Findings = append(rep.Findings, reconcile.Finding{
				Code:      f.Code,
				Severity:  f.Severity,
				InvoiceID: "",
				Message:   f.Message,
			})
			if f.Severity == "critical" {
				rep.ReadinessScore -= 25
			}
		}
		if rep.ReadinessScore < 0 {
			rep.ReadinessScore = 0
		}
	}
	if err := reconcile.WriteReport(*output, rep); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *format == "table" {
		fmt.Printf("KsefPulse readiness: %d/100 (%d findings)\n", rep.ReadinessScore, len(rep.Findings))
		for _, f := range rep.Findings {
			fmt.Printf("  - [%s] %s %s: %s\n", f.Severity, f.Code, f.InvoiceID, f.Message)
		}
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rep)
	}
	if rep.ReadinessScore < *minScore {
		os.Exit(2)
	}
}
