package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-client/pkg/models"
	"log"
	"net/http"
	"testing"
)

func TestHandleHttpResponseCode(t *testing.T) {
	err := handleHttpResponseCode(http.StatusOK, []byte("message"))

	assert.Nil(t, err)

	resp, err := json.Marshal(models.ErrorResponse{
		Message: "msg",
		Info:    "err",
	})
	if err != nil {
		log.Fatal("Can't marshall test model\n", err)
	}

	err = handleHttpResponseCode(http.StatusBadRequest, resp)
	assert.EqualError(t, err, "server responded with an error. Error message: msg. Error info: err")
}

func TestNew(t *testing.T) {
	c := New("http://my.url", "en")

	assert.Equal(t, "en", c.locale)
	assert.NotNil(t, c.Jobs)
	assert.NotNil(t, c.Runs)
	assert.NotNil(t, c.Snippets)
	assert.NotNil(t, c.Tags)
	assert.NotNil(t, c.About)
	assert.NotNil(t, c.Events)
	assert.NotNil(t, c.Adapters)
	assert.NotNil(t, c.Ws)

	c = New("http://my.url", "")
	assert.Equal(t, "en", c.locale)
}
