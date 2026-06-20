package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var bearerRe = regexp.MustCompile(`(?i)\b(?:authorization\s*:\s*)?bearer\s+([A-Za-z0-9._\-+/=]{20,512})`)

type bearerDetector struct{}

func NewBearerDetector() models.Detector { return bearerDetector{} }

func (bearerDetector) Name() string { return "bearer-token" }

func (bearerDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range bearerRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) < 4 {
			continue
		}
		value, start, end, ok := validUTF8Slice(input, loc[2], loc[3])
		if !ok || isSafeAssignmentValue(value) {
			continue
		}
		token := "<REDACTED_SECRET>"
		fType := "bearer_token"
		if validJWT(value) {
			token = "<REDACTED_JWT>"
			fType = "jwt"
		}
		findings = append(findings, newFinding(
			"bearer-token", "auth_token", fType, models.SeverityHigh,
			start, end, strings.TrimSpace(value), token,
			"Bearer token detected.",
			"Remove authorization tokens before sending this prompt.",
			0.88,
		))
	}
	return findings, nil
}
