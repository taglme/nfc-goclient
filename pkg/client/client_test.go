package client

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-goclient/pkg/models"
)

func TestHandleHttpResponseCode(t *testing.T) {
	err := handleHttpResponseCode(http.StatusOK, []byte("message"))

	assert.Nil(t, err)

	resp, err := json.Marshal(models.ErrorResponse{
		Message: "msg",
		Info:    "err",
	})
	if err != nil {
		log.Fatal("Can't marshal test model", err)
	}

	err = handleHttpResponseCode(http.StatusBadRequest, resp)
	assert.EqualError(t, err, "Server responded with an error: msg (err)")
}

func TestNew(t *testing.T) {
	c := New("http://my.url")
	assert.NotNil(t, c.Jobs)
	assert.NotNil(t, c.Runs)
	assert.NotNil(t, c.Snippets)
	assert.NotNil(t, c.Tags)
	assert.NotNil(t, c.About)
	assert.NotNil(t, c.Events)
	assert.NotNil(t, c.Adapters)
	assert.NotNil(t, c.Ws)

}
