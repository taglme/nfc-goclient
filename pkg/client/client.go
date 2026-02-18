package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/f2prateek/train"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//Client represents main client structure. It is used to communicate with server API
type Client struct {
	Adapters AdapterService
	About    AboutService
	Events   EventService
	Snippets SnippetService
	Tags     TagService
	Runs     RunService
	Jobs     JobService
	Licenses LicenseService
	Ws       WsService
}

type Options struct {
	// BaseURL is the HTTP base URL, e.g. "http://127.0.0.1:3011" or "https://api.example.com".
	// If scheme is omitted, http is assumed.
	BaseURL string

	// WSURL is the websocket base URL, e.g. "ws://127.0.0.1:3011" or "wss://api.example.com".
	// If empty, derived from BaseURL (http->ws, https->wss).
	WSURL string

	// XAppKey is required for authenticated API access.
	XAppKey string
	// XUserToken is optional.
	XUserToken string

	// Locale sets Accept-Language header if non-empty.
	Locale string

	// Interceptors are applied in the given order.
	Interceptors []train.Interceptor
}

//New create new client to communicate with server API
func New(hostOrBaseURL string, interceptors ...train.Interceptor) *Client {
	opts := Options{BaseURL: hostOrBaseURL, Interceptors: interceptors}
	return NewWithOptions(opts)
}

func NewWithOptions(opts Options) *Client {
	baseHTTP := normalizeHTTPBaseURL(opts.BaseURL)
	baseWS := strings.TrimSpace(opts.WSURL)
	if baseWS == "" {
		baseWS = deriveWSBaseURL(baseHTTP)
	} else {
		baseWS = normalizeWSBaseURL(baseWS)
	}

	interceptors := make([]train.Interceptor, 0, len(opts.Interceptors)+2)
	if strings.TrimSpace(opts.Locale) != "" {
		interceptors = append(interceptors, NewLocalizer(strings.TrimSpace(opts.Locale)))
	}
	if strings.TrimSpace(opts.XAppKey) != "" || strings.TrimSpace(opts.XUserToken) != "" {
		interceptors = append(interceptors, NewAuthHeaders(opts.XAppKey, opts.XUserToken))
	}
	interceptors = append(interceptors, opts.Interceptors...)

	transport := train.Transport(interceptors...)
	httpClient := &http.Client{Transport: transport}

	return &Client{
		Adapters: newAdapterService(httpClient, baseHTTP),
		About:    newAboutService(httpClient, baseHTTP),
		Events:   newEventService(httpClient, baseHTTP),
		Snippets: newSnippetService(httpClient, baseHTTP),
		Tags:     newTagService(httpClient, baseHTTP),
		Runs:     newRunService(httpClient, baseHTTP),
		Jobs:     newJobService(httpClient, baseHTTP),
		Licenses: newLicenseService(httpClient, baseHTTP),
		Ws:       newWsService(baseWS),
	}
}

func normalizeHTTPBaseURL(input string) string {
	s := strings.TrimSpace(input)
	if s == "" {
		return "http://127.0.0.1:3011"
	}
	if strings.Contains(s, "://") {
		return strings.TrimRight(s, "/")
	}
	return strings.TrimRight("http://"+s, "/")
}

func normalizeWSBaseURL(input string) string {
	s := strings.TrimSpace(input)
	if s == "" {
		return "ws://127.0.0.1:3011"
	}
	if strings.Contains(s, "://") {
		return strings.TrimRight(s, "/")
	}
	return strings.TrimRight("ws://"+s, "/")
}

func deriveWSBaseURL(httpBase string) string {
	u, err := url.Parse(strings.TrimSpace(httpBase))
	if err != nil || u == nil {
		return normalizeWSBaseURL(httpBase)
	}
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	switch scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	default:
		// Best-effort fallback.
		u.Scheme = "ws"
	}
	return strings.TrimRight(u.String(), "/")
}

func handleHttpResponseCode(statusCode int, body []byte) (err error) {
	if statusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return
		}
		err = fmt.Errorf("Server responded with an error: %s (%s)", errorResponse.Message, errorResponse.Info)
		return err
	}

	return err
}
