package models

type Edge struct {
	ID             string `json:"id"`
	Source         string `json:"source"`
	Target         string `json:"target"`
	ConnectionType string `json:"connectionType"`
	Protocol       string `json:"protocol"`
}
