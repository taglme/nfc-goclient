package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/taglme/nfc-goclient/pkg/models"
)

var licensePath = "/licenses/apps/"

type LicenseChecker interface {
	Start()
	Stop()
	License() (models.ApplicationLicense, bool)
	IsActive() bool
}

type licenseChecker struct {
	host     string
	appID    string
	interval time.Duration
	callback func(bool)
	license  models.ApplicationLicense
	started  bool
	stopSig  chan bool
}

func NewLicenseChecker(host string, interval time.Duration, callback func(bool), appID string) LicenseChecker {
	return &licenseChecker{
		host:     host,
		appID:    appID,
		interval: interval,
		callback: callback,
	}
}

func (l *licenseChecker) License() (al models.ApplicationLicense, ok bool) {
	al = l.license
	if l.license.ID != "" {
		ok = true
	}
	return
}

func (l *licenseChecker) IsActive() (active bool) {
	if l.license.ID != "" {
		if l.license.End == "" {
			active = true
		} else {
			licenseEnd, err := time.Parse("2006-01-02", l.license.End)
			if err == nil {
				if time.Now().Before(licenseEnd) {
					active = true
				}
			}

		}

	}
	return
}

func (l *licenseChecker) Start() {
	if !l.started {
		fmt.Println("Start check license")
		stopSig := make(chan bool)
		l.stopSig = stopSig
		go func() {
			for {
				select {
				case <-stopSig:
					return
				default:
					l.license = models.ApplicationLicense{}
					resp, err := http.Get(l.host + licensePath + l.appID)
					if err != nil {
						break
					}
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						break
					}
					resp.Body.Close()
					var appLicense models.ApplicationLicense
					err = json.Unmarshal(body, &appLicense)
					if err != nil {
						break
					}
					if appLicense.ID == "" {
						break
					}
					l.license = appLicense

				}
				if l.callback != nil {
					go l.callback(l.IsActive())
				}
				time.Sleep(l.interval)

			}
		}()
		l.started = true
	}

}
func (l *licenseChecker) Stop() {
	if l.started {
		fmt.Println("Stop check license")
		l.stopSig <- true
		l.started = false
	}

}
