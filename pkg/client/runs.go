package client

import (
	"github.com/taglme/nfc-client/pkg/models"
)

type RunFilter struct {
	JobID   string
	Status  models.JobRunStatus
	SortBy  string
	SortDir string
}

type RunService interface {
	GetAll(adapterID string) ([]models.JobRun, error)
	GetFiltered(adapterID string, filter RunFilter) ([]models.JobRun, error)
	Get(adapterID string, runID string) (models.JobRun, error)
}
