package detectors

import (
	"context"
	"regexp"

	"github.com/promptsieve/promptsieve/pkg/models"
)

type RegexDetector struct {
	name        string
	regex       *regexp.Regexp
	category    string
	findingType string
	severity    models.Severity
	description string
	recomm      string
	token       string
	valueGroup  int
	confidence  float64
}

func NewRegexDetector(name, pattern, category, fType string, severity models.Severity, desc, rec, token string) *RegexDetector {
	return &RegexDetector{
		name:        name,
		regex:       regexp.MustCompile(pattern),
		category:    category,
		findingType: fType,
		severity:    severity,
		description: desc,
		recomm:      rec,
		token:       token,
		confidence:  0.99,
	}
}

func (d *RegexDetector) WithValueGroup(n int) *RegexDetector {
	d.valueGroup = n
	return d
}

func (d *RegexDetector) Name() string {
	return d.name
}

func (d *RegexDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	matches := d.regex.FindAllStringSubmatchIndex(input, maxFindingsPerDetector)
	for _, match := range matches {
		start, end := match[0], match[1]
		if d.valueGroup > 0 {
			g := d.valueGroup * 2
			if g+1 < len(match) && match[g] >= 0 {
				start, end = match[g], match[g+1]
			}
		}
		value, start, end, ok := validUTF8Slice(input, start, end)
		if !ok {
			continue
		}
		findings = append(findings, newFinding(
			d.name, d.category, d.findingType, d.severity,
			start, end, value, d.token, d.description, d.recomm, d.confidence,
		))
	}
	return findings, nil
}
