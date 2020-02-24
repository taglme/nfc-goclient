package client

import (
	"github.com/taglme/nfc-client/pkg/models"
)

type SnippetFilter struct {
	Category models.SnippetCategory
	UsageID  string
}
type SnippetService interface {
	GetAll() ([]models.Snippet, error)
	GetFiltered(filter SnippetFilter) ([]models.Snippet, error)
}
