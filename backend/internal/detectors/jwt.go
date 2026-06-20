package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var jwtCandidateRe = regexp.MustCompile(`\b[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{10,}\b`)

type jwtDetector struct{}

func NewJWTDetector() models.Detector { return jwtDetector{} }

func (jwtDetector) Name() string { return "jwt-token" }

func (jwtDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range jwtCandidateRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		value, start, end, ok := validUTF8Slice(input, loc[0], loc[1])
		if !ok || !validJWT(value) {
			continue
		}
		findings = append(findings, newFinding(
			"jwt-token", "auth_token", "jwt", models.SeverityHigh,
			start, end, value, "<REDACTED_JWT>",
			"JSON Web Token (JWT) detected.",
			"JWTs can contain sensitive auth information.",
			0.96,
		))
	}
	return findings, nil
}
