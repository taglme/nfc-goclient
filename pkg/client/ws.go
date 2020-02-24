package client

import (
	"github.com/taglme/nfc-client/pkg/models"
)

type WsService interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	OnTagDiscovery(func(adapterID string, tag models.Tag))
	OnTagRelease(func(adapterID string, tag models.Tag))
	OnAdapterDiscovery(func(adapterID string, adapter models.Adapter))
	OnAdapterRelease(func(adapterID string, adapter models.Adapter))
	OnRunStarted(func(adapterID string, jobRun models.JobRun))
	OnRunSuccess(func(adapterID string, jobRun models.JobRun))
	OnRunError(func(adapterID string, jobRun models.JobRun))
	OnJobSubmitted(func(adapterID string, job models.Job))
	OnJobRemoved(func(adapterID string, job models.Job))
	OnJobActivated(func(adapterID string, job models.Job))
	OnJobPending(func(adapterID string, job models.Job))
	OnJobFinished(func(adapterID string, job models.Job))
	SetLocale(locale string) error
	ConnString() string
}
