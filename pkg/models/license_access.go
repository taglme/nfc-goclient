package models

// LicenseAccess describes current host access level and allowed scopes as derived from license policies.
type LicenseAccess struct {
	MachineID            string   `json:"machine_id"`
	HostTier             string   `json:"host_tier"`
	AllowedScopes        []string `json:"allowed_scopes"`
	AllowedCommandScopes []string `json:"allowed_command_scopes"`
}
