package sanitizer

import (
	"sort"

	"github.com/promptsieve/promptsieve/pkg/models"
)

// Sanitize replaces each finding range with its replacement token.
// It uses the provided findings as the only source of truth and does not rescan.
func Sanitize(input string, findings []models.Finding) string {
	if len(findings) == 0 || input == "" {
		return input
	}

	usable := make([]models.Finding, 0, len(findings))
	for _, finding := range findings {
		if finding.Start < 0 || finding.End > len(input) || finding.Start >= finding.End {
			continue
		}
		usable = append(usable, finding)
	}
	if len(usable) == 0 {
		return input
	}

	sort.SliceStable(usable, func(i, j int) bool {
		if usable[i].Start == usable[j].Start {
			if usable[i].End != usable[j].End {
				return usable[i].End > usable[j].End
			}
			return replacementOf(usable[i]) < replacementOf(usable[j])
		}
		return usable[i].Start < usable[j].Start
	})

	resolved := make([]models.Finding, 0, len(usable))
	for _, finding := range usable {
		if len(resolved) == 0 || finding.Start >= resolved[len(resolved)-1].End {
			resolved = append(resolved, finding)
			continue
		}
		prev := resolved[len(resolved)-1]
		if prefer(finding, prev) {
			resolved[len(resolved)-1] = finding
		}
	}

	out := make([]byte, 0, len(input))
	cursor := 0
	for _, finding := range resolved {
		out = append(out, input[cursor:finding.Start]...)
		out = append(out, replacementOf(finding)...)
		cursor = finding.End
	}
	out = append(out, input[cursor:]...)
	return string(out)
}

func replacementOf(finding models.Finding) string {
	if finding.ReplacementToken != "" {
		return finding.ReplacementToken
	}
	if finding.MaskedValue != "" {
		return finding.MaskedValue
	}
	return "<REDACTED>"
}

func prefer(candidate, current models.Finding) bool {
	candLen := candidate.End - candidate.Start
	currLen := current.End - current.Start
	if candLen != currLen {
		return candLen > currLen
	}
	if candidate.Confidence != current.Confidence {
		return candidate.Confidence > current.Confidence
	}
	return severityWeight(candidate.Severity) > severityWeight(current.Severity)
}

func severityWeight(s models.Severity) int {
	switch s {
	case models.SeverityCritical:
		return 4
	case models.SeverityHigh:
		return 3
	case models.SeverityMedium:
		return 2
	case models.SeverityLow:
		return 1
	default:
		return 0
	}
}
