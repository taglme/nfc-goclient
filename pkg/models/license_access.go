package models

type JobCapabilities struct {
	AllowRepeat    bool `json:"allow_repeat"`
	AllowBatch     bool `json:"allow_batch"`
	AllowMultistep bool `json:"allow_multistep"`
}

type CreateJobConstraints struct {
	MaxRepeat               *int     `json:"max_repeat"`
	MaxSteps                *int     `json:"max_steps"`
	MaxQueuedJobsPerAdapter *int     `json:"max_queued_jobs_per_adapter"`
	AllowedScopes           []string `json:"allowed_scopes"`
	AllowedCommandScopes    []string `json:"allowed_command_scopes"`
}

type CreateJobRateLimit struct {
	MinIntervalMs *int `json:"min_interval_ms"`
	WindowMs      *int `json:"window_ms"`
	MaxInWindow   *int `json:"max_in_window"`
}

// LicenseAccess describes current host access level and allowed scopes as derived from license policies.
type LicenseAccess struct {
	MachineID            string   `json:"machine_id"`
	HostTier             string   `json:"host_tier"`
	AllowedScopes        []string `json:"allowed_scopes"`
	AllowedCommandScopes []string `json:"allowed_command_scopes"`

	CreateJobRateLimit   CreateJobRateLimit   `json:"create_job_rate_limit"`
	JobCapabilities      JobCapabilities      `json:"job_capabilities"`
	CreateJobConstraints CreateJobConstraints `json:"create_job_constraints"`
}
