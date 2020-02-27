package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/taglme/nfc-client/pkg/models"
)

type Client struct {
	Adapters AdapterService
	About    AboutService
	Events   EventService
	Snippets SnippetService
	Tags     TagService
	Runs     RunService
	Jobs     JobService
	Ws       WsService
	locale   string
}

func New(url string, locale string) *Client {
	httpClient := &http.Client{}

	//Set default locale in case of unknown locale
	_, ok := models.StringToLocale(locale)
	if !ok {
		locale = models.LocaleEn.String()
	}

	rt := withHeader(httpClient.Transport)
	rt.Set("Accept-Language", locale)
	httpClient.Transport = rt

	return &Client{
		Adapters: newAdapterService(httpClient, url),
		About:    newAboutService(httpClient, url),
		Events:   newEventService(httpClient, url),
		Snippets: newSnippetService(httpClient, url),
		Tags:     newTagService(httpClient, url),
		Runs:     newRunService(httpClient, url),
		Jobs:     newJobService(httpClient, url),
		Ws:       newWsService(url),
		locale:   locale,
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

type withHeaderStruct struct {
	http.Header
	rt http.RoundTripper
}

func withHeader(rt http.RoundTripper) withHeaderStruct {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeaderStruct{Header: make(http.Header), rt: rt}
}

func (h withHeaderStruct) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}
