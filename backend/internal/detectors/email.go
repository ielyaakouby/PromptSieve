package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var emailRe = regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]{1,64}@[A-Z0-9.\-]{1,255}\.[A-Z]{2,24}\b`)

type emailDetector struct{}

func NewEmailDetector() models.Detector { return emailDetector{} }

func (emailDetector) Name() string { return "email" }

func (emailDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range emailRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		value, start, end, ok := validUTF8Slice(input, loc[0], loc[1])
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			"email", "personal_information", "email", models.SeverityMedium,
			start, end, value, "<REDACTED_EMAIL>",
			"Email address detected.",
			"Remove or replace personal email addresses before sending this prompt.",
			0.95,
		))
	}
	return findings, nil
}
