package models

type Node struct {
	ID   string   `json:"id"`
	Type string   `json:"type"`
	Data NodeData `json:"data"`
}

type NodeData struct {
	Label      string          `json:"label"`
	Subtype    string          `json:"subtype"`
	Properties map[string]bool `json:"properties"`
}
