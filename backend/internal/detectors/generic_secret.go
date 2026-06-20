package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var (
	assignmentRe = regexp.MustCompile(`(?i)\b([A-Za-z_][A-Za-z0-9_.-]{0,80})\s*[=:]\s*["']?([^\s"'#]{4,256})`)
	sensitiveKey = regexp.MustCompile(`(?i)(secret|token|password|passwd|pwd|api[_-]?key|private[_-]?key|access[_-]?token|auth[_-]?token|client[_-]?secret|database[_-]?password|db[_-]?password|database[_-]?url|connection[_-]?string)`)
)

var safeKeys = map[string]struct{}{
	"port": {}, "host": {}, "hostname": {}, "debug": {}, "env": {}, "node_env": {},
	"replicas": {}, "version": {}, "timeout": {}, "count": {}, "enabled": {},
	"path": {}, "home": {}, "lang": {}, "term": {}, "shell": {}, "tz": {},
	"ci": {}, "editor": {}, "log_level": {}, "loglevel": {}, "namespace": {},
}

type genericSecretDetector struct{}

func NewGenericSecretDetector() models.Detector { return genericSecretDetector{} }

func (genericSecretDetector) Name() string { return "generic-secret" }

func (genericSecretDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	for _, loc := range assignmentRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) < 6 {
			continue
		}
		key := input[loc[2]:loc[3]]
		value, start, end, ok := validUTF8Slice(input, loc[4], loc[5])
		if !ok {
			continue
		}
		keyNorm := strings.ToLower(strings.ReplaceAll(key, "-", "_"))
		if _, safe := safeKeys[keyNorm]; safe {
			continue
		}
		if keyNorm == "pwd" && isSafeAssignmentValue(value) {
			continue
		}
		if !sensitiveKey.MatchString(key) {
			continue
		}
		if isSafeAssignmentValue(value) {
			continue
		}
		fType, token := classifySecretKey(keyNorm)
		findings = append(findings, newFinding(
			"generic-secret", "credential", fType, models.SeverityHigh,
			start, end, strings.Trim(value, `"'`), token,
			"Sensitive assignment detected from variable name.",
			"Remove secrets from configuration snippets before sending this prompt.",
			0.86,
		))
	}
	return findings, nil
}

func classifySecretKey(key string) (string, string) {
	switch {
	case strings.Contains(key, "database_password") || strings.Contains(key, "db_password"):
		return "database_password", "<REDACTED_PASSWORD>"
	case strings.Contains(key, "password") || strings.Contains(key, "passwd") || strings.HasSuffix(key, "pwd"):
		return "password", "<REDACTED_PASSWORD>"
	case strings.Contains(key, "openai"):
		return "openai_api_key", "<REDACTED_OPENAI_API_KEY>"
	case strings.Contains(key, "github") || key == "gh_token":
		return "github_token", "<REDACTED_GITHUB_TOKEN>"
	default:
		return "generic_secret", "<REDACTED_SECRET>"
	}
}
