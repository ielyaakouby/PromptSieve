package scanner

import (
	"sort"

	"github.com/promptsieve/promptsieve/pkg/models"
)

var typeSpecificity = map[string]int{
	"private_key":           100,
	"openai_api_key":        90,
	"github_token":          90,
	"aws_secret_access_key": 90,
	"aws_access_key":        85,
	"payment_card":          85,
	"iban":                  85,
	"jwt":                   80,
	"cvv":                   75,
	"card_expiry":           70,
	"database_password":     65,
	"password":              60,
	"bic":                   55,
	"email":                 50,
	"phone":                 50,
	"date_of_birth":         50,
	"address":               45,
	"person_name":           40,
	"bearer_token":          30,
	"generic_secret":        20,
}

func Deduplicate(findings []models.Finding) []models.Finding {
	if len(findings) == 0 {
		return findings
	}

	valid := make([]models.Finding, 0, len(findings))
	for _, finding := range findings {
		if finding.Start < 0 || finding.End <= finding.Start {
			continue
		}
		valid = append(valid, finding)
	}
	if len(valid) < 2 {
		return valid
	}

	sort.SliceStable(valid, func(i, j int) bool {
		if valid[i].Start != valid[j].Start {
			return valid[i].Start < valid[j].Start
		}
		si, sj := specificity(valid[i]), specificity(valid[j])
		if si != sj {
			return si > sj
		}
		if valid[i].Confidence != valid[j].Confidence {
			return valid[i].Confidence > valid[j].Confidence
		}
		return (valid[i].End - valid[i].Start) > (valid[j].End - valid[j].Start)
	})

	kept := make([]models.Finding, 0, len(valid))
	for _, finding := range valid {
		replaced := false
		skip := false
		for i, existing := range kept {
			if !overlaps(existing, finding) {
				continue
			}
			if better(finding, existing) {
				kept[i] = finding
				replaced = true
			} else {
				skip = true
			}
			break
		}
		if !replaced && !skip {
			kept = append(kept, finding)
		}
	}
	return kept
}

func overlaps(a, b models.Finding) bool {
	return a.Start < b.End && b.Start < a.End
}

func specificity(f models.Finding) int {
	if v, ok := typeSpecificity[f.Type]; ok {
		return v
	}
	return 10
}

func better(candidate, current models.Finding) bool {
	cs, cu := specificity(candidate), specificity(current)
	if cs != cu {
		return cs > cu
	}
	if candidate.Confidence != current.Confidence {
		return candidate.Confidence > current.Confidence
	}
	return (candidate.End - candidate.Start) > (current.End - current.Start)
}
