package models

type Rule struct {
	ID          string     `yaml:"id"`
	Category    string     `yaml:"category"`
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	Suggestion  string     `yaml:"suggestion"`
	BaseWeight  int        `yaml:"baseWeight"`
	Conditions  Conditions `yaml:"conditions"`
}

type Conditions struct {
	SourceType     string            `yaml:"sourceType"`
	TargetType     string            `yaml:"targetType"`
	NodeType       string            `yaml:"nodeType"`
	RequiredNodes  []string          `yaml:"requiredNodes"`
	ForbiddenNodes []string          `yaml:"forbiddenNodes"`
	Properties     map[string]bool   `yaml:"properties"`
	ContextChecks  map[string]string `yaml:"contextChecks"`
	Direct         bool              `yaml:"direct"`
}
