package client

import (
	"encoding/json"
	"fmt"
	"github.com/taglme/nfc-client/pkg/models"
	"net/http"
)

type Client struct {
	Adapters AdapterService
	About    AboutService
	Events   EventService
	Snippets SnippetService
	Tags     TagService
	Runs     RunService
	Jobs     JobService
	/*
		TODO
		Ws WsService
	*/
}

func New(url string) *Client {
	httpClient := &http.Client{}

	return &Client{
		Adapters: newAdapterService(httpClient, url),
		About:    newAboutService(httpClient, url),
		Events:   newEventService(httpClient, url),
		Snippets: newSnippetService(httpClient, url),
		Tags:     newTagService(httpClient, url),
		Runs:     newRunService(httpClient, url),
		Jobs:     newJobService(httpClient, url),
	}
}

func handleHttpResponseCode(statusCode int, body []byte) (err error) {
	if statusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return
		}
		err = fmt.Errorf("server responded with an error. Error message: %s. Error info: %s", errorResponse.Message, errorResponse.Info)
		return err
	}

	return err
}
