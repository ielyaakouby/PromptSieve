package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var ibanRe = regexp.MustCompile(`\b[A-Z]{2}\d{2}(?:[ ]?[A-Z0-9]){11,30}\b`)

type ibanDetector struct{}

func NewIBANDetector() models.Detector { return ibanDetector{} }

func (ibanDetector) Name() string { return "iban" }

func (ibanDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range ibanRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		value, start, end, ok := validUTF8Slice(input, loc[0], loc[1])
		if !ok || !validIBAN(value) {
			continue
		}
		findings = append(findings, newFinding(
			"iban", "financial_information", "iban", models.SeverityHigh,
			start, end, value, "<REDACTED_IBAN>",
			"IBAN detected and checksum-validated.",
			"Remove account identifiers before sending this prompt.",
			0.97,
		))
	}
	return findings, nil
}
