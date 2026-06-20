package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var (
	frIntlPhoneRe = regexp.MustCompile(`(?:\+|00)33[\s.\-]*[1-9](?:[\s.\-]?\d{2}){4}`)
	frNatPhoneRe  = regexp.MustCompile(`(?:\b|[^0-9])(0[1-9](?:[\s.\-]?\d{2}){4})\b`)
	e164PhoneRe   = regexp.MustCompile(`\+\d{1,3}[\s.\-]?\d(?:[\s.\-]?\d){6,13}`)
)

type phoneDetector struct{}

func NewPhoneDetector() models.Detector { return phoneDetector{} }

func (phoneDetector) Name() string { return "phone" }

func (phoneDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	seen := map[[2]int]struct{}{}
	var findings []models.Finding
	add := func(start, end int) {
		value, start, end, ok := validUTF8Slice(input, start, end)
		if !ok {
			return
		}
		digits := digitsOnly(value)
		if validFrenchPhone(digits) || (len(digits) >= 10 && validE164(digits) && (value[0] == '+' || nearby(input, start, end, 24, "tel", "tél", "phone", "mobile", "portable"))) {
			key := [2]int{start, end}
			if _, ok := seen[key]; ok {
				return
			}
			seen[key] = struct{}{}
			findings = append(findings, newFinding(
				"phone", "personal_information", "phone", models.SeverityMedium,
				start, end, value, "<REDACTED_PHONE>",
				"Phone number detected.",
				"Remove personal phone numbers before sending this prompt.",
				0.9,
			))
		}
	}

	for _, loc := range frIntlPhoneRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		add(loc[0], loc[1])
	}
	for _, loc := range frNatPhoneRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) >= 4 {
			add(loc[2], loc[3])
		}
	}
	for _, loc := range e164PhoneRe.FindAllStringIndex(input, maxFindingsPerDetector) {
		add(loc[0], loc[1])
	}
	return findings, nil
}
