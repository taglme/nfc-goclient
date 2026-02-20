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

func TestLicenseGetAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/licenses/access", req.URL.String())
		resp, err := json.Marshal(models.LicenseResource{
			ID:       "LIC-1",
			Owner:    "Owner",
			Email:    "o@example.com",
			Machine:  "MID123",
			Type:     "standard",
			HostTier: "pro",
			Start:    "2024-01-01",
			End:      "2025-01-01",
			Support:  "2025-01-01",
			Plugins:  []string{"Nfcd"},
		})
		if err != nil {
			log.Fatal("Can't marshall test model", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return response", err)
		}
	}))
	defer server.Close()

	api := newLicenseService(server.Client(), server.URL)
	lic, err := api.GetAccess()
	assert.NoError(t, err)
	assert.Equal(t, "MID123", lic.Machine)
	assert.Equal(t, "pro", lic.HostTier)
	assert.Equal(t, "LIC-1", lic.ID)
}

func TestLicenseGetAccessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/licenses/access", req.URL.String())
		resp, err := json.Marshal(models.ErrorResponse{Message: "err msg", Info: "err info"})
		if err != nil {
			log.Fatal("Can't marshall test model", err)
		}
		rw.WriteHeader(500)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return response", err)
		}
	}))
	defer server.Close()

	api := newLicenseService(server.Client(), server.URL)
	_, err := api.GetAccess()
	assert.NotNil(t, err)
}
