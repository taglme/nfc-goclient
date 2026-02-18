package client

import (
	"net/http"
	"strings"

	"github.com/f2prateek/train"
	"github.com/pkg/errors"
)

const (
	HeaderXAppKey    = "X-App-Key"
	HeaderXUserToken = "X-User-Token"
)

type AuthHeaders struct {
	XAppKey       string
	XUserToken    string
	RequireAppKey bool
}

func NewAuthHeaders(xAppKey string, xUserToken string) train.Interceptor {
	return &AuthHeaders{
		XAppKey:       strings.TrimSpace(xAppKey),
		XUserToken:    strings.TrimSpace(xUserToken),
		RequireAppKey: true,
	}
}

func (i *AuthHeaders) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	if req == nil {
		return nil, errors.New("request is nil")
	}

	appKey := strings.TrimSpace(i.XAppKey)
	if i.RequireAppKey && appKey == "" {
		return nil, errors.New("X-App-Key is required")
	}
	if appKey != "" {
		req.Header.Set(HeaderXAppKey, appKey)
	}

	userToken := strings.TrimSpace(i.XUserToken)
	if userToken != "" {
		req.Header.Set(HeaderXUserToken, userToken)
	}

	return chain.Proceed(req)
}
