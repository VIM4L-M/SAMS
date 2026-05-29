package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vimal/sams/internal/models"
)

func Rules(rules []models.Rule) http.HandlerFunc {
	type ruleDTO struct {
		ID          string `json:"id"`
		Category    string `json:"category"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Suggestion  string `json:"suggestion"`
		BaseWeight  int    `json:"baseWeight"`
	}

	dtos := make([]ruleDTO, len(rules))
	for i, r := range rules {
		dtos[i] = ruleDTO{
			ID:          r.ID,
			Category:    r.Category,
			Title:       r.Title,
			Description: r.Description,
			Suggestion:  r.Suggestion,
			BaseWeight:  r.BaseWeight,
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rules": dtos,
			"total": len(dtos),
		})
	}
}
