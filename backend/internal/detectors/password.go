package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var passwordRe = regexp.MustCompile(`(?i)(?:mot\s+de\s+passe|passphrase|password|passwd|\bpwd\b)\s*(?:is|:|=|:=)\s*["']?([^\s"']{4,128})`)

type passwordDetector struct{}

func NewPasswordDetector() models.Detector { return passwordDetector{} }

func (passwordDetector) Name() string { return "generic-password" }

func (passwordDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range passwordRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) < 4 {
			continue
		}
		value, start, end, ok := validUTF8Slice(input, loc[2], loc[3])
		if !ok || isSafeAssignmentValue(value) {
			continue
		}
		findings = append(findings, newFinding(
			"generic-password", "credential", "password", models.SeverityHigh,
			start, end, strings.Trim(value, `"'`), "<REDACTED_PASSWORD>",
			"Potential cleartext password detected.",
			"Remove or redact passwords before sharing.",
			0.93,
		))
	}
	return findings, nil
}
