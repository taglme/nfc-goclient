package models

type PageInfo struct {
	Total  int `json:"total"`
	Length int `json:"length"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
