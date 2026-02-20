package models

import (
	"time"

	"github.com/pkg/errors"
)

type License struct {
	ID           string
	Owner        string
	Email        string
	Machine      string
	Type         string
	HostTier     string
	Start        time.Time
	End          time.Time
	Support      time.Time
	// Features is kept for backward compatibility with older API schema.
	// Prefer Plugins.
	Features []string
	Plugins  []string
	Applications []AppLicense
}

type LicenseResource struct {
	ID           string               `json:"id"`
	Owner        string               `json:"owner"`
	Email        string               `json:"email"`
	Machine      string               `json:"machine"`
	Type         string               `json:"type"`
	HostTier     string               `json:"host_tier"`
	Start        string               `json:"start"`
	End          string               `json:"end"`
	Support      string               `json:"support"`
	// Legacy field.
	Features []string `json:"features"`
	// Current field.
	Plugins []string `json:"plugins"`
	Applications []AppLicenseResource `json:"applications"`
}

type AppLicense struct {
	ID      string
	Name    string
	Type    string
	Start   time.Time
	End     time.Time
	Support time.Time
}
type AppLicenseResource struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Support string `json:"support"`
}
type LicenseMID struct {
	MID string `json:"mid"`
}

func (l AppLicense) IsActive() (ok bool) {
	if l.ID != "" {
		if l.End.IsZero() {
			ok = true
		} else {
			if time.Now().Before(l.End) {
				ok = true
			}
		}

	}
	return
}

func (l License) IsActive() (ok bool) {
	if l.ID != "" {
		if l.End.IsZero() {
			ok = true
		} else {
			if time.Now().Before(l.End) {
				ok = true
			}
		}

	}
	return
}

func (r *AppLicenseResource) ToAppLicense() (appLicense AppLicense, err error) {
	var licenseEnd, licenseStart, licenseSupport time.Time
	if r.End != "" {
		licenseEnd, err = time.Parse("2006-01-02", r.End)
		if err != nil {
			return appLicense, errors.Wrap(err, "Can't parse app license end time")
		}
	}
	if r.Start != "" {
		licenseStart, err = time.Parse("2006-01-02", r.Start)
		if err != nil {
			return appLicense, errors.Wrap(err, "Can't parse app license start time")
		}
	}
	if r.Support != "" {
		licenseSupport, err = time.Parse("2006-01-02", r.Support)
		if err != nil {
			return appLicense, errors.Wrap(err, "Can't parse app license support time")
		}
	}
	appLicense = AppLicense{
		ID:      r.ID,
		Name:    r.Name,
		Type:    r.Type,
		Start:   licenseStart,
		End:     licenseEnd,
		Support: licenseSupport,
	}
	return

}
func (r *LicenseResource) ToLicense() (license License, err error) {
	var licenseEnd, licenseStart, licenseSupport time.Time
	if r.End != "" {
		licenseEnd, err = time.Parse("2006-01-02", r.End)
		if err != nil {
			return license, errors.Wrap(err, "Can't parse license end time")
		}
	}
	if r.Start != "" {
		licenseStart, err = time.Parse("2006-01-02", r.Start)
		if err != nil {
			return license, errors.Wrap(err, "Can't parse license start time")
		}
	}
	if r.Support != "" {
		licenseSupport, err = time.Parse("2006-01-02", r.Support)
		if err != nil {
			return license, errors.Wrap(err, "Can't parse license support time")
		}
	}
	appLicenses := []AppLicense{}
	for _, appLicenseRes := range r.Applications {
		appLicense, err := appLicenseRes.ToAppLicense()
		if err != nil {
			return license, errors.Wrap(err, "Can't convert app license resource")
		}
		appLicenses = append(appLicenses, appLicense)
	}

	license = License{
		ID:           r.ID,
		Owner:        r.Owner,
		Email:        r.Email,
		Machine:      r.Machine,
		Type:         r.Type,
		HostTier:     r.HostTier,
		Start:        licenseStart,
		End:          licenseEnd,
		Support:      licenseSupport,
		Features:     r.Features,
		Plugins:      r.Plugins,
		Applications: appLicenses,
	}
	// Backward compatibility: if server still returns legacy features, map them to plugins.
	if len(license.Plugins) == 0 && len(license.Features) > 0 {
		license.Plugins = append([]string{}, license.Features...)
	}
	return

}
