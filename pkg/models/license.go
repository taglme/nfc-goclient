package models

type ApplicationLicense struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Support string `json:"support"`
}
