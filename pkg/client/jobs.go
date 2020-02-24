package client

import (
	"github.com/taglme/nfc-client/pkg/models"
)

type JobFilter struct {
	Status  models.JobStatus
	SortBy  string
	SortDir string
}
type JobService interface {
	GetAll(adapterID string) ([]models.Job, error)
	GetFiltered(adapterID string, filter JobFilter) ([]models.Job, error)
	Delete(adapterID string, jobID string) error
	DeleteAll(adapterID string) error
	Get(adapterID string, jobID string) (models.Job, error)
	Add(adapterID string, job models.NewJob) (models.Job, error)
	UpdateStatus(adapterID string, jobID string, status models.JobStatus) (models.Job, error)
}
