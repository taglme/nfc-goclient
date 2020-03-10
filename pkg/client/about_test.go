package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-goclient/pkg/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAboutGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/about")
		resp, err := json.Marshal(models.AppInfo{
			Name:    "Name",
			Version: "Version",
			Commit:  "Commit",
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

	// Use Client & URL from our local test server
	api := newAboutService(server.Client(), server.URL)
	body, err := api.Get()

	if err != nil {
		log.Fatal("Can't get about info", err)
	}

	assert.Equal(t, body.Name, "Name")
	assert.Equal(t, body.Version, "Version")
	assert.Equal(t, body.Commit, "Commit")
}

func TestAboutGetError(t *testing.T) {
	resp, err := json.Marshal(models.AppInfo{})
	if err != nil {
		log.Fatal("Can't marshall test model", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/about")
		rw.WriteHeader(500)
		_, err := rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return err", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := newAboutService(server.Client(), server.URL)
	body, err := api.Get()
	assert.NotNil(t, err)
	assert.Equal(t, body.Name, "")
}
