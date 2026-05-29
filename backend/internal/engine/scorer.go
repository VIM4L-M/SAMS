package engine

import (
	"github.com/vimal/sams/internal/models"
)

func contextMultiplier(ctx models.Context) float64 {
	switch ctx.TrafficLevel {
	case "low":
		return 0.5
	case "medium":
		return 1.0
	case "high":
		return 1.5
	case "massive":
		return 2.0
	default:
		return 1.0
	}
}

func calculateScore(
	errors []models.Issue,
	warnings []models.Issue,
	rules []models.Rule,
	ctx models.Context,
) int {
	ruleWeights := make(map[string]int)
	for _, r := range rules {
		ruleWeights[r.ID] = r.BaseWeight
	}

	multiplier := contextMultiplier(ctx)
	penalty := 0.0

	for _, issue := range errors {
		w := ruleWeights[issue.RuleID]
		if w == 0 {
			w = 10
		}
		penalty += float64(w) * multiplier
	}
	for _, issue := range warnings {
		w := ruleWeights[issue.RuleID]
		if w == 0 {
			w = 5
		}
		// Warnings carry half the penalty of errors
		penalty += float64(w) * multiplier * 0.5
	}

	score := int(100.0 - penalty)
	if score < 0 {
		score = 0
	}
	return score
}
