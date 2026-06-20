package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var dobRe = regexp.MustCompile(`(?i)(?:n[ée]e?\s+le|date\s+de\s+naissance|naissance|date\s+of\s+birth|\bdob\b|birth\s*date|born(?:\s+on)?)\s*[:\s]*(\d{1,2}[/\-.]\d{1,2}[/\-.]\d{2,4})`)

type dobDetector struct{}

func NewDOBDetector() models.Detector { return dobDetector{} }

func (dobDetector) Name() string { return "date-of-birth" }

func (dobDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range dobRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) < 4 {
			continue
		}
		value, start, end, ok := validUTF8Slice(input, loc[2], loc[3])
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			"date-of-birth", "personal_information", "date_of_birth", models.SeverityMedium,
			start, end, value, "<REDACTED_DATE_OF_BIRTH>",
			"Date of birth detected from contextual wording.",
			"Remove dates of birth before sending this prompt.",
			0.92,
		))
	}
	return findings, nil
}
