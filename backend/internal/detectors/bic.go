package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var bicRe = regexp.MustCompile(`\b[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}(?:[A-Z0-9]{3})?\b`)

type bicDetector struct{}

func NewBICDetector() models.Detector { return bicDetector{} }

func (bicDetector) Name() string { return "bic" }

func (bicDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range bicRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		value, start, end, ok := validUTF8Slice(input, loc[0], loc[1])
		if !ok || !validBIC(value) {
			continue
		}
		if !nearby(input, start, end, 40, "bic", "swift", "bank identifier") {
			continue
		}
		findings = append(findings, newFinding(
			"bic", "financial_information", "bic", models.SeverityHigh,
			start, end, strings.ToUpper(value), "<REDACTED_BIC>",
			"BIC/SWIFT code detected.",
			"Remove bank identifier codes before sending this prompt.",
			0.9,
		))
	}
	return findings, nil
}
