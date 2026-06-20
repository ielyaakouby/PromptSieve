package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var expiryRe = regexp.MustCompile(`(?i)(?:exp(?:ir(?:y|ation|e))?|valid\s*thru|date\s+d['’]expiration)\s*[:\s]+(\d{1,2}\s*[/\-]\s*\d{2,4})`)

type expiryDetector struct{}

func NewCardExpiryDetector() models.Detector { return expiryDetector{} }

func (expiryDetector) Name() string { return "card-expiry" }

func (expiryDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range expiryRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) < 4 {
			continue
		}
		value, start, end, ok := validUTF8Slice(input, loc[2], loc[3])
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			"card-expiry", "financial_information", "card_expiry", models.SeverityHigh,
			start, end, value, "<REDACTED_CARD_EXPIRY>",
			"Payment card expiration date detected from contextual wording.",
			"Remove card expiration dates before sending this prompt.",
			0.9,
		))
	}
	return findings, nil
}
