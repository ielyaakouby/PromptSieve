package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var cardRe = regexp.MustCompile(`\b(?:\d[ -]?){13,19}\b`)

type cardDetector struct{}

func NewCardDetector() models.Detector { return cardDetector{} }

func (cardDetector) Name() string { return "payment-card" }

func (cardDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range cardRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		value, start, end, ok := validUTF8Slice(input, loc[0], loc[1])
		if !ok {
			continue
		}
		digits := digitsOnly(value)
		if !validLuhn(digits) {
			continue
		}
		findings = append(findings, newFinding(
			"payment-card", "financial_information", "payment_card", models.SeverityCritical,
			start, end, value, "<REDACTED_PAYMENT_CARD>",
			"Payment card number detected and Luhn-validated.",
			"Never include payment card numbers in prompts.",
			0.98,
		))
	}
	return findings, nil
}
