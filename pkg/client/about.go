package client

import (
	"github.com/taglme/nfc-client/pkg/models"
)

type AboutService interface {
	Get() (models.AppInfo, error)
}
