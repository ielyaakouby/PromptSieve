package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

// Person detection is contextual only. PromptSieve does not ship a NER model.
var (
	personLabelRe = regexp.MustCompile(`(?i)\b(?:client|customer|name|nom|pr[eé]nom|monsieur|mme|mr\.?|mrs\.?)\s*[:\s]+([A-ZÀ-Ÿ][a-zà-ÿ''\-]{1,30}(?:\s+[A-ZÀ-Ÿ][a-zà-ÿ''\-]{1,30}){1,2})`)
	personNearDOB = regexp.MustCompile(`([A-ZÀ-Ÿ][a-zà-ÿ''\-]{1,30}\s+[A-ZÀ-Ÿ][a-zà-ÿ''\-]{1,30})[\s\S]{0,80}(?:n[ée]e?\s+le|born\b|date\s+of\s+birth|\bdob\b)`)
)

type personDetector struct{}

func NewPersonDetector() models.Detector { return personDetector{} }

func (personDetector) Name() string { return "person-name" }

func (personDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	seen := map[[2]int]struct{}{}
	var findings []models.Finding
	add := func(start, end int, conf float64) {
		value, start, end, ok := validUTF8Slice(input, start, end)
		if !ok {
			return
		}
		key := [2]int{start, end}
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		findings = append(findings, newFinding(
			"person-name", "personal_information", "person_name", models.SeverityMedium,
			start, end, value, "<REDACTED_PERSON_NAME>",
			"Person name detected from contextual signals.",
			"Remove personal names before sending this prompt.",
			conf,
		))
	}

	for _, loc := range personLabelRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) >= 4 {
			add(loc[2], loc[3], 0.82)
		}
	}
	for _, loc := range personNearDOB.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) >= 4 {
			add(loc[2], loc[3], 0.8)
		}
	}
	return findings, nil
}
