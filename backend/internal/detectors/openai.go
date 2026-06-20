package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var (
	openaiKeyRe    = regexp.MustCompile(`\b(sk-proj-[A-Za-z0-9_-]{20,}|sk-[A-Za-z0-9]{20,})\b`)
	openaiAssignRe = regexp.MustCompile(`(?i)\bOPENAI_(?:API_)?KEY\s*[=:]\s*["']?([^\s"']{8,256})`)
)

type openaiDetector struct{}

func NewOpenAIDetector() models.Detector { return openaiDetector{} }

func (openaiDetector) Name() string { return "openai-api-key" }

func (openaiDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	seen := map[[2]int]struct{}{}
	var findings []models.Finding
	add := func(start, end int, conf float64) {
		value, start, end, ok := validUTF8Slice(input, start, end)
		if !ok || isSafeAssignmentValue(value) {
			return
		}
		key := [2]int{start, end}
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		findings = append(findings, newFinding(
			"openai-api-key", "api_key", "openai_api_key", models.SeverityCritical,
			start, end, strings.Trim(value, `"'`), "<REDACTED_OPENAI_API_KEY>",
			"OpenAI-style API key detected.",
			"Revoke and replace this key if it may have been exposed.",
			conf,
		))
	}
	for _, loc := range openaiKeyRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) >= 4 {
			add(loc[2], loc[3], 0.98)
		}
	}
	for _, loc := range openaiAssignRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) >= 4 {
			add(loc[2], loc[3], 0.94)
		}
	}
	return findings, nil
}
