package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var (
	githubPrefixRe = regexp.MustCompile(`\b(ghp_[A-Za-z0-9]{20,255}|github_pat_[A-Za-z0-9_]{20,255}|gho_[A-Za-z0-9]{20,255}|ghu_[A-Za-z0-9]{20,255}|ghs_[A-Za-z0-9]{20,255}|ghr_[A-Za-z0-9]{20,255})\b`)
	githubAssignRe = regexp.MustCompile(`(?i)\b(?:GITHUB_TOKEN|GH_TOKEN|GITHUB_PAT|GITHUB_OAUTH_TOKEN)\s*[=:]\s*["']?([^\s"']{8,256})`)
)

type githubDetector struct{}

func NewGitHubDetector() models.Detector { return githubDetector{} }

func (githubDetector) Name() string { return "github-token" }

func (githubDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
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
			"github-token", "source_control_token", "github_token", models.SeverityCritical,
			start, end, strings.Trim(value, `"'`), "<REDACTED_GITHUB_TOKEN>",
			"GitHub token detected.",
			"Revoke this token if it may have been exposed.",
			conf,
		))
	}
	for _, loc := range githubPrefixRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) >= 4 {
			add(loc[2], loc[3], 0.98)
		}
	}
	for _, loc := range githubAssignRe.FindAllStringSubmatchIndex(input, maxFindingsPerDetector) {
		if len(loc) >= 4 {
			add(loc[2], loc[3], 0.94)
		}
	}
	return findings, nil
}
