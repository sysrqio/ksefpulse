package fa3

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Finding struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func ValidateFile(path string) ([]Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ValidateXML(data)
}

func ValidateXML(data []byte) ([]Finding, error) {
	findings := []Finding{}
	trim := strings.TrimSpace(string(data))
	if trim == "" {
		return nil, fmt.Errorf("empty fa3 document")
	}
	if !strings.HasPrefix(trim, "<") {
		return nil, fmt.Errorf("not xml")
	}
	var doc any
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("xml parse: %w", err)
	}
	lower := strings.ToLower(trim)
	if !strings.Contains(lower, "faktura") && !strings.Contains(lower, "invoice") {
		findings = append(findings, Finding{
			Code:     "FA3_ROOT_HINT_MISSING",
			Severity: "high",
			Message:  "expected FA(3) invoice root element hint (Faktura/Invoice)",
		})
	}
	return findings, nil
}
