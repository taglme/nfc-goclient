package client

import (
	"github.com/taglme/nfc-client/pkg/models"
)

type TagService interface {
	GetAll(adapterID string) ([]models.Tag, error)
	Get(adapterID string, tagID string) (models.Tag, error)
}
