package client

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-goclient/pkg/models"
)

func TestTagsGetAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/adapters/id/tags?type=bluetooth", req.URL.String())
		resp, err := json.Marshal([]models.TagShortResource{{
			TagID:       "id",
			Kind:        "kind",
			Href:        "href",
			Type:        models.TagTypeBluetooth.String(),
			AdapterID:   "id",
			AdapterName: "name",
			Uid:         "za7VbYcSQU2zRgGQXQAm/g==",
		}})
		if err != nil {
			log.Fatal("Can't marshall test model", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return err", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := newTagService(server.Client(), server.URL)

	tagType := models.TagTypeBluetooth
	body, err := api.GetAll("id", &tagType)

	if err != nil {
		log.Fatal("Can't get tags", err)
	}

	assert.Equal(t, "name", body[0].AdapterName)
	assert.Equal(t, "id", body[0].AdapterID)
	assert.Equal(t, models.TagTypeBluetooth, body[0].Type)
	assert.Equal(t, "id", body[0].TagID)
}

func TestTagsGetOne(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/adapters/id/tags/bluetooth", req.URL.String())
		resp, err := json.Marshal(models.TagResource{
			TagID:       "id",
			Type:        models.TagTypeBluetooth.String(),
			AdapterID:   "id",
			AdapterName: "name",
			Vendor:      "vendor",
			Product:     "product",
		})
		if err != nil {
			log.Fatal("Can't marshall test model", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return err", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := newTagService(server.Client(), server.URL)
	body, err := api.Get("id", models.TagTypeBluetooth.String())

	if err != nil {
		log.Fatal("Can't get tags", err)
	}

	assert.Equal(t, "name", body.AdapterName)
	assert.Equal(t, "id", body.AdapterID)
	assert.Equal(t, models.TagTypeBluetooth, body.Type)
	assert.Equal(t, "id", body.TagID)
	assert.Equal(t, "product", body.Product)
	assert.Equal(t, "vendor", body.Vendor)
}
