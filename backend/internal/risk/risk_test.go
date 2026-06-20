package risk

import (
	"testing"

	"github.com/promptsieve/promptsieve/pkg/models"
)

func TestEmptyIsSafe(t *testing.T) {
	got := CalculateRisk(nil)
	if got.Status != "SAFE" || got.Score != 0 {
		t.Fatalf("%+v", got)
	}
}

func TestEmailOnlyIsReview(t *testing.T) {
	got := CalculateRisk([]models.Finding{{Severity: models.SeverityMedium}})
	if got.Status != "REVIEW" {
		t.Fatalf("got %s", got.Status)
	}
}

func TestPasswordIsSensitive(t *testing.T) {
	got := CalculateRisk([]models.Finding{{Severity: models.SeverityHigh}})
	if got.Status != "SENSITIVE" {
		t.Fatalf("got %s", got.Status)
	}
}

func TestCriticalIsHighRisk(t *testing.T) {
	got := CalculateRisk([]models.Finding{{Severity: models.SeverityCritical}})
	if got.Status != "CRITICAL" || got.Score < 70 {
		t.Fatalf("%+v", got)
	}
}

func TestAnyFindingIsNeverSafe(t *testing.T) {
	got := CalculateRisk([]models.Finding{{Severity: models.SeverityLow}})
	if got.Status == "SAFE" {
		t.Fatal("low finding must not be SAFE")
	}
}

func TestScoreClamped(t *testing.T) {
	findings := make([]models.Finding, 8)
	for i := range findings {
		findings[i].Severity = models.SeverityCritical
	}
	got := CalculateRisk(findings)
	if got.Score != 100 {
		t.Fatalf("score=%d", got.Score)
	}
}
