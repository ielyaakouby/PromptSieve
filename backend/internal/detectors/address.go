package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var addressRe = regexp.MustCompile(`(?i)\b\d{1,4}\s+(?:rue|avenue|ave|boulevard|bd|chemin|route|impasse|place|all[eé]e|cours|quai)\s+[A-Za-zÀ-ÿ0-9''\-][A-Za-zÀ-ÿ0-9''\- ]{1,50}(?:,\s*)?(?:\d{5}\s+[A-Za-zÀ-ÿ][A-Za-zÀ-ÿ''\- ]{1,40})?`)

type addressDetector struct{}

func NewAddressDetector() models.Detector { return addressDetector{} }

func (addressDetector) Name() string { return "address" }

func (addressDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range addressRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		value, start, end, ok := validUTF8Slice(input, loc[0], loc[1])
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			"address", "personal_information", "address", models.SeverityMedium,
			start, end, value, "<REDACTED_ADDRESS>",
			"Postal or street address detected.",
			"Remove physical addresses before sending this prompt.",
			0.88,
		))
	}
	return findings, nil
}
