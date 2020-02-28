package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-client/pkg/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdapterService_GetAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/adapters")
		resp, err := json.Marshal([]models.AdapterResource{{
			AdapterID: "123",
			Kind:      "kind",
			Href:      "href",
			Name:      "name",
			Type:      models.AdapterTypeBarcode.String(),
			Driver:    "driver",
		}})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newAdapterService(server.Client(), server.URL)
	body, err := api.GetAll()
	if err != nil {
		log.Fatal("Can't get adapters\n", err)
	}

	assert.Equal(t, "123", body[0].AdapterID)
	assert.Equal(t, "name", body[0].Name)
	assert.Equal(t, models.AdapterTypeBarcode, body[0].Type)
}

func TestAdapterService_GetAll2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/adapters")
		resp, err := json.Marshal(models.ErrorResponse{
			Message: "err msg",
			Info:    "err info",
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(500)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newAdapterService(server.Client(), server.URL)
	_, err := api.GetAll()
	assert.NotNil(t, err)
}

func TestAdapterService_GetFiltered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/adapters?type=barcode")
		resp, err := json.Marshal([]models.AdapterResource{{
			AdapterID: "123",
			Kind:      "kind",
			Href:      "href",
			Name:      "name",
			Type:      models.AdapterTypeBarcode.String(),
			Driver:    "driver",
		}})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newAdapterService(server.Client(), server.URL)
	aType := models.AdapterTypeBarcode
	body, err := api.GetFiltered(&aType)
	if err != nil {
		log.Fatal("Can't get adapters\n", err)
	}

	assert.Equal(t, "123", body[0].AdapterID)
	assert.Equal(t, "name", body[0].Name)
	assert.Equal(t, models.AdapterTypeBarcode, body[0].Type)
}

func TestAdapterService_GetFiltered2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/adapters?type=barcode")
		resp, err := json.Marshal(models.ErrorResponse{
			Message: "err msg",
			Info:    "err info",
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(500)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newAdapterService(server.Client(), server.URL)
	aType := models.AdapterTypeBarcode
	_, err := api.GetFiltered(&aType)
	assert.NotNil(t, err)
}

func TestAdapterService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/adapters/123")
		resp, err := json.Marshal(models.AdapterResource{
			AdapterID: "123",
			Kind:      "kind",
			Href:      "href",
			Name:      "name",
			Type:      models.AdapterTypeBarcode.String(),
			Driver:    "driver",
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newAdapterService(server.Client(), server.URL)
	body, err := api.Get("123")
	if err != nil {
		log.Fatal("Can't get adapter\n", err)
	}

	assert.Equal(t, "123", body.AdapterID)
	assert.Equal(t, "name", body.Name)
	assert.Equal(t, models.AdapterTypeBarcode, body.Type)
	assert.Equal(t, "driver", body.Driver)
}

func TestAdapterService_Get2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/adapters/123")
		resp, err := json.Marshal(models.ErrorResponse{
			Message: "err msg",
			Info:    "err info",
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(500)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newAdapterService(server.Client(), server.URL)
	_, err := api.Get("123")
	assert.NotNil(t, err)
}
