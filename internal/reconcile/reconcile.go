package reconcile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Invoice struct {
	ID          string  `json:"id"`
	Status      string  `json:"status"`
	SubmittedAt string  `json:"submitted_at"`
	KSeFNumber  *string `json:"ksef_number"`
	UPORef      *string `json:"upo_ref"`
}

type StateFile struct {
	Version  int       `json:"version"`
	Invoices []Invoice `json:"invoices"`
}

type Finding struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	InvoiceID string `json:"invoice_id"`
	Message   string `json:"message"`
}

type Report struct {
	SchemaVersion  string    `json:"schema_version"`
	GeneratedAt    string    `json:"generated_at"`
	ReadinessScore int       `json:"readiness_score"`
	Findings       []Finding `json:"findings"`
}

func LoadState(path string) (*StateFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var st StateFile
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	if st.Version == 0 {
		st.Version = 1
	}
	return &st, nil
}

func Analyze(st *StateFile) Report {
	findings := []Finding{}
	for _, inv := range st.Invoices {
		switch inv.Status {
		case "submitted", "pending":
			if inv.UPORef == nil || *inv.UPORef == "" {
				findings = append(findings, Finding{
					Code:      "MISSING_UPO",
					Severity:  "critical",
					InvoiceID: inv.ID,
					Message:   "invoice submitted without UPO reference",
				})
			}
		case "orphan":
			findings = append(findings, Finding{
				Code:      "ORPHAN_RECORD",
				Severity:  "critical",
				InvoiceID: inv.ID,
				Message:   "state marks invoice as orphan",
			})
		case "accepted":
			if inv.KSeFNumber == nil || *inv.KSeFNumber == "" {
				findings = append(findings, Finding{
					Code:      "MISSING_KSEF_NUMBER",
					Severity:  "high",
					InvoiceID: inv.ID,
					Message:   "accepted without ksef_number",
				})
			}
		}
	}
	score := 100
	for _, f := range findings {
		switch f.Severity {
		case "critical":
			score -= 35
		case "high":
			score -= 15
		default:
			score -= 5
		}
	}
	if score < 0 {
		score = 0
	}
	return Report{
		SchemaVersion:  "1.0",
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		ReadinessScore: score,
		Findings:       findings,
	}
}

func WriteReport(path string, rep Report) error {
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
