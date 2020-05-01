package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Signer interface {
	Sign(req *http.Request) error
	SetMID(mid string)
	MID() string
}

type signer struct {
	appID  string
	secret *rsa.PrivateKey
	cert   string
	mid    string
}

func NewRSASigner(appID string, secret *rsa.PrivateKey, cert string) Signer {
	return &signer{
		appID:  appID,
		secret: secret,
		cert:   cert,
	}
}

type DesktopClaims struct {
	Certs []string `json:"x5c"`
	jwt.StandardClaims
	refAud string
}

func (s *signer) SetMID(mid string) {
	s.mid = mid
}

func (s *signer) MID() string {
	return s.mid
}

func (s *signer) Sign(req *http.Request) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	claims := DesktopClaims{
		[]string{s.cert},
		jwt.StandardClaims{
			Issuer:    s.appID,
			Audience:  "desktop:" + s.mid,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		"",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(s.secret)
	if err != nil {
		return fmt.Errorf("Request sign error: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+ss)
	return nil
}

func PrivateRSAKeyFromB64String(str string) (*rsa.PrivateKey, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("Key base64 decode error:", err)
	}
	pKey, err := x509.ParsePKCS1PrivateKey(decodedKey)
	if err != nil {
		return nil, fmt.Errorf("Key parse error:", err)
	}
	return pKey, nil
}
