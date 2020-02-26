package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/taglme/nfc-client/pkg/models"
)

type AdapterFilter struct {
	AdapterType models.AdapterType
}
type AdapterService interface {
	GetAll() ([]models.Adapter, error)
	GetFiltered(adapterType models.AdapterType) ([]models.Adapter, error)
	Get(adapterID string) (models.Adapter, error)
}

type adapterService struct {
	url    string
	client *http.Client
}

func newAdapterService(client *http.Client, url string) AdapterService {
	return &adapterService{
		url:    url,
		client: client,
	}
}

func (s *adapterService) GetAll() ([]models.Adapter, error) {
	return s.GetFiltered(0)
}

func (s *adapterService) GetFiltered(adapterType models.AdapterType) (adapters []models.Adapter, err error) {
	var adaptersListResource models.AdapterListResource
	targetURL := s.url + "/adapters"
	adapterTypeStr := adapterType.String()
	if adapterTypeStr != "" {
		targetURL = targetURL + "?type=" + adapterTypeStr
	}
	resp, err := s.client.Get(targetURL)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return
		}
		err = fmt.Errorf("Server responded with error. Error message: %s. Error info: %s", errorResponse.Message, errorResponse.Info)
		return
	}

	err = json.Unmarshal(body, &adaptersListResource)
	if err != nil {
		return
	}
	for _, adapterShortResource := range adaptersListResource {
		adapters = append(adapters, adapterShortResource.ToAdapter())
	}
	return
}
func (s *adapterService) Get(adapterID string) (adapter models.Adapter, err error) {
	var adapterResource models.AdapterResource
	if adapterID == "" {
		err = errors.New("Adapter ID not set")
		return
	}

	targetURL := s.url + "/adapters/" + adapterID
	resp, err := s.client.Get(targetURL)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return
		}
		err = fmt.Errorf("Server responded with error. Error message: %s. Error info: %s", errorResponse.Message, errorResponse.Info)
		return
	}

	err = json.Unmarshal(body, &adapterResource)
	if err != nil {
		return
	}
	adapter = adapterResource.ToAdapter()
	return
}
