package models

import "context"

type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type Finding struct {
	ID         string   `json:"id"`
	Detector   string   `json:"detector"`
	Category   string   `json:"category"`
	Type       string   `json:"type"`
	Severity   Severity `json:"severity"`
	Confidence float64  `json:"confidence"`

	Start int `json:"start"`
	End   int `json:"end"`

	Match            string `json:"-"`
	MaskedValue      string `json:"maskedValue"`
	ReplacementToken string `json:"replacementToken,omitempty"`
	Line             int    `json:"line,omitempty"`

	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
}

type Detector interface {
	Name() string
	Detect(ctx context.Context, input string) ([]Finding, error)
}
