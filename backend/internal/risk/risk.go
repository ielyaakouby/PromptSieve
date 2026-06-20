package risk

import (
	"github.com/promptsieve/promptsieve/pkg/models"
)

type Score struct {
	Score  int
	Status string
}

func CalculateRisk(findings []models.Finding) Score {
	if len(findings) == 0 {
		return Score{Score: 0, Status: "SAFE"}
	}

	score := 0
	hasCritical := false
	hasHigh := false
	hasMedium := false

	for _, f := range findings {
		switch f.Severity {
		case models.SeverityCritical:
			score += 60
			hasCritical = true
		case models.SeverityHigh:
			score += 40
			hasHigh = true
		case models.SeverityMedium:
			score += 25
			hasMedium = true
		case models.SeverityLow:
			score += 10
		}
	}

	if score > 100 {
		score = 100
	}

	status := "REVIEW"
	switch {
	case hasCritical:
		if score < 70 {
			score = 70
		}
		status = "CRITICAL"
	case hasHigh:
		if score < 40 {
			score = 40
		}
		status = "SENSITIVE"
	case hasMedium:
		if score < 20 {
			score = 20
		}
		status = "REVIEW"
	default:
		status = "REVIEW"
	}

	return Score{
		Score:  score,
		Status: status,
	}
}
