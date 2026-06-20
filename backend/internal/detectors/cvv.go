package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var cvvRe = regexp.MustCompile(`(?i)(?:cvv2?|cvc2?|cid|security\s*code|code\s+s[eé]curit[eé]|cryptogramme)\s*[:#=\s]+(\d{3,4})\b`)

type cvvDetector struct{}

func NewCVVDetector() models.Detector { return cvvDetector{} }

func (cvvDetector) Name() string { return "cvv" }

func (cvvDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range cvvRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) < 4 {
			continue
		}
		value, start, end, ok := validUTF8Slice(input, loc[2], loc[3])
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			"cvv", "financial_information", "cvv", models.SeverityCritical,
			start, end, value, "<REDACTED_CVV>",
			"Card security code detected from contextual label.",
			"Never include CVV/CVC values in prompts.",
			0.95,
		))
	}
	return findings, nil
}
