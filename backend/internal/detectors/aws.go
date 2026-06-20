package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var (
	awsAccessKeyRe = regexp.MustCompile(`\b(?:AKIA|ASIA)[A-Z0-9]{16}\b`)
	awsSecretRe    = regexp.MustCompile(`(?i)\bAWS[_-]?SECRET[_-]?ACCESS[_-]?KEY\s*[=:]\s*["']?([A-Za-z0-9/+=]{40})["']?`)
)

type awsDetector struct{}

func NewAWSDetector() models.Detector { return awsDetector{} }

func (awsDetector) Name() string { return "aws-credentials" }

func (awsDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range awsAccessKeyRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		value, start, end, ok := validUTF8Slice(input, loc[0], loc[1])
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			"aws-access-key", "cloud_credential", "aws_access_key", models.SeverityCritical,
			start, end, value, "<REDACTED_AWS_ACCESS_KEY>",
			"AWS access key ID detected.",
			"Remove or redact this credential before sharing.",
			0.99,
		))
	}
	for _, loc := range awsSecretRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) < 4 {
			continue
		}
		value, start, end, ok := validUTF8Slice(input, loc[2], loc[3])
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			"aws-secret-key", "cloud_credential", "aws_secret_access_key", models.SeverityCritical,
			start, end, strings.Trim(value, `"'`), "<REDACTED_AWS_SECRET_KEY>",
			"AWS secret access key detected from contextual assignment.",
			"Rotate this secret if it may have been exposed.",
			0.97,
		))
	}
	return findings, nil
}
