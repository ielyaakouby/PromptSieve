package detectors

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/promptsieve/promptsieve/pkg/models"
)

const maxFindingsPerDetector = 500

func newFinding(
	detector, category, fType string,
	severity models.Severity,
	start, end int,
	match, token, desc, rec string,
	confidence float64,
) models.Finding {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	return models.Finding{
		ID:               fmt.Sprintf("%s:%d:%d", detector, start, end),
		Detector:         detector,
		Category:         category,
		Type:             fType,
		Severity:         severity,
		Confidence:       confidence,
		Start:            start,
		End:              end,
		Match:            match,
		MaskedValue:      maskForType(fType, match),
		ReplacementToken: token,
		Description:      desc,
		Recommendation:   rec,
	}
}

func maskForType(fType, value string) string {
	switch fType {
	case "password", "database_password", "cvv":
		return "********"
	case "email":
		return maskEmail(value)
	case "phone":
		return maskKeepSuffix(value, 2)
	case "payment_card":
		return maskKeepSuffix(digitsOnly(value), 4)
	case "iban":
		return maskKeepPrefixSuffix(strings.ReplaceAll(value, " ", ""), 2, 4)
	case "aws_access_key":
		return maskKeepPrefixSuffix(value, 4, 4)
	case "openai_api_key", "github_token":
		return maskKeepPrefixSuffix(value, 4, 4)
	case "jwt":
		return "eyJ**************"
	default:
		if len(value) <= 4 {
			return "********"
		}
		return maskKeepPrefixSuffix(value, 2, 2)
	}
}

func maskEmail(value string) string {
	at := strings.LastIndex(value, "@")
	if at <= 0 {
		return "***@***"
	}
	return "***" + value[at:]
}

func maskKeepSuffix(value string, suffix int) string {
	if suffix <= 0 || len(value) <= suffix {
		return strings.Repeat("*", max(4, len(value)))
	}
	return strings.Repeat("*", len(value)-suffix) + value[len(value)-suffix:]
}

func maskKeepPrefixSuffix(value string, prefix, suffix int) string {
	if len(value) <= prefix+suffix {
		return strings.Repeat("*", max(4, len(value)))
	}
	return value[:prefix] + strings.Repeat("*", len(value)-prefix-suffix) + value[len(value)-suffix:]
}

func digitsOnly(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func nearby(input string, start, end, radius int, needles ...string) bool {
	from := start - radius
	if from < 0 {
		from = 0
	}
	to := end + radius
	if to > len(input) {
		to = len(input)
	}
	window := strings.ToLower(input[from:to])
	for _, needle := range needles {
		if strings.Contains(window, strings.ToLower(needle)) {
			return true
		}
	}
	return false
}

func validUTF8Slice(input string, start, end int) (string, int, int, bool) {
	if start < 0 || end > len(input) || start >= end {
		return "", 0, 0, false
	}
	if !utf8.ValidString(input[start:end]) {
		return "", 0, 0, false
	}
	return input[start:end], start, end, true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
