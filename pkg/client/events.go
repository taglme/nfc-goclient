package client

import (
	"github.com/taglme/nfc-client/pkg/models"
)

type EventFilter struct {
	Name    models.EventName
	SortBy  string
	SortDir string
}

type EventService interface {
	GetAll() ([]models.Event, error)
	GetFiltered(adapterID string, filter EventFilter) ([]models.Event, error)
	Add(event models.NewEvent) (models.Event, error)
}
