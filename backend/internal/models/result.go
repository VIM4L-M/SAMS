package models

type ValidationResult struct {
	Errors   []Issue      `json:"errors"`
	Warnings []Issue      `json:"warnings"`
	Passed   []PassedRule `json:"passed"`
	Score    int          `json:"score"`
}

type Issue struct {
	RuleID        string   `json:"ruleId"`
	Category      string   `json:"category"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	AffectedNodes []string `json:"affectedNodes"`
	AffectedEdges []string `json:"affectedEdges"`
	Suggestion    string   `json:"suggestion"`
}

type PassedRule struct {
	RuleID   string `json:"ruleId"`
	Category string `json:"category"`
	Title    string `json:"title"`
}

type ValidationRequest struct {
	Nodes   []Node  `json:"nodes"`
	Edges   []Edge  `json:"edges"`
	Context Context `json:"context"`
}

type Context struct {
	TrafficLevel   string `json:"trafficLevel"`
	ReadWriteRatio string `json:"readWriteRatio"`
	UserBase       string `json:"userBase"`
	TeamSize       string `json:"teamSize"`
	Stage          string `json:"stage"`
}

type ValidationResponse struct {
	Results  ValidationResult `json:"results"`
	Metadata Metadata         `json:"metadata"`
}

type Metadata struct {
	RulesChecked int   `json:"rulesChecked"`
	TimeMs       int64 `json:"timeMs"`
}
