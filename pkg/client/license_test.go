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
		resp, err := json.Marshal(models.LicenseAccess{
			MachineID:            "MID123",
			HostTier:             "pro",
			AllowedScopes:        []string{"nfc.read", "nfc.write"},
			AllowedCommandScopes: []string{"run.create"},
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
	access, err := api.GetAccess()
	assert.NoError(t, err)
	assert.Equal(t, "MID123", access.MachineID)
	assert.Equal(t, "pro", access.HostTier)
	assert.Equal(t, []string{"nfc.read", "nfc.write"}, access.AllowedScopes)
	assert.Equal(t, []string{"run.create"}, access.AllowedCommandScopes)
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
