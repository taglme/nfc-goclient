package models

import (
	"fmt"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type JobRun struct {
	RunID       string
	JobID       string
	JobName     string
	Status      JobRunStatus
	AdapterID   string
	AdapterName string
	Tag         Tag
	Results     []StepResult
	CreatedAt   time.Time
}

func (j Job) ToJobRun() JobRun {
	id, _ := uuid.NewV4()
	createdAt := time.Now().UTC()
	stepResults := make([]StepResult, len(j.Steps))
	for i, _ := range stepResults {
		stepResults[i] = j.Steps[i].ToStepResult()
	}

	jobRun := JobRun{
		RunID:       id.String(),
		JobID:       j.JobID,
		JobName:     j.JobName,
		Status:      JobRunStatusStarted,
		AdapterID:   j.AdapterID,
		AdapterName: j.AdapterName,
		Tag:         Tag{},
		Results:     stepResults,
		CreatedAt:   createdAt,
	}
	return jobRun
}

type JobRunResource struct {
	RunID       string               `json:"run_id"`
	Kind        string               `json:"kind"`
	Href        string               `json:"href"`
	JobID       string               `json:"job_id"`
	JobName     string               `json:"job_name"`
	Status      string               `json:"status"`
	AdapterID   string               `json:"adapter_id"`
	AdapterName string               `json:"adapter_name"`
	Tag         TagResource          `json:"tag"`
	Results     []StepResultResource `json:"results"`
	CreatedAt   string               `json:"created_at"`
}

func (jr JobRun) ToResource() JobRunResource {
	var stepResultResources []StepResultResource
	for _, stepResult := range jr.Results {
		stepResultResources = append(stepResultResources, stepResult.ToResource())
	}
	resource := JobRunResource{
		RunID:       jr.RunID,
		Kind:        "JobRun",
		Href:        fmt.Sprintf(`/adapters/%s/runs/%s`, jr.AdapterID, jr.RunID),
		JobID:       jr.JobID,
		JobName:     jr.JobName,
		Status:      jr.Status.String(),
		AdapterID:   jr.AdapterID,
		AdapterName: jr.AdapterName,
		Tag:         jr.Tag.ToResource(),
		Results:     stepResultResources,
		CreatedAt:   jr.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
	}
	return resource
}

type JobRunListResource struct {
	Total  int
	Length int
	Limit  int
	Offset int
	Items  []JobRunResource
}

type JobRunStatus int

const (
	JobRunStatusStarted JobRunStatus = iota + 1
	JobRunStatusSuccess
	JobRunStatusError
)

func StringToJobRunStatus(s string) (JobRunStatus, bool) {
	switch s {
	case JobRunStatusStarted.String():
		return JobRunStatusStarted, true
	case JobRunStatusSuccess.String():
		return JobRunStatusSuccess, true
	case JobRunStatusError.String():
		return JobRunStatusError, true
	}
	return 0, false
}

func (jobRunStatus JobRunStatus) String() string {
	names := [...]string{
		"unknown",
		"started",
		"success",
		"error",
	}
	if jobRunStatus < JobRunStatusStarted || jobRunStatus > JobRunStatusError {
		return names[0]
	}
	return names[jobRunStatus]
}

type CommandStatus int

const (
	CommandStatusSuccess CommandStatus = iota + 1
	CommandStatusError
)

func StringToCommandStatus(s string) (CommandStatus, bool) {
	switch s {
	case CommandStatusSuccess.String():
		return CommandStatusSuccess, true
	case CommandStatusError.String():
		return CommandStatusError, true
	}
	return 0, false
}

func (commandStatus CommandStatus) String() string {
	names := [...]string{
		"unknown",
		"success",
		"error",
	}
	if commandStatus < CommandStatusSuccess || commandStatus > CommandStatusError {
		return names[0]
	}
	return names[commandStatus]
}

type StepResult struct {
	Command Command
	Params  CommandParams
	Output  CommandOutput
	Status  CommandStatus
	Message string
}
type StepResultResource struct {
	Command string                `json:"command"`
	Params  CommandParamsResource `json:"params"`
	Output  CommandOutputResource `json:"output"`
	Status  string                `json:"status"`
	Message string                `json:"message"`
}

func (stepResult StepResult) ToResource() StepResultResource {
	var params CommandParamsResource
	var output CommandOutputResource

	if stepResult.Params != nil {
		params = stepResult.Params.ToResource()
	}
	if stepResult.Output != nil {
		output = stepResult.Output.ToResource()
	}

	resource := StepResultResource{
		Command: stepResult.Command.String(),
		Params:  params,
		Output:  output,
		Status:  stepResult.Status.String(),
		Message: stepResult.Message,
	}
	return resource
}
func (stepResultResource StepResultResource) ToStepResult() (StepResult, error) {
	var params CommandParams
	var output CommandOutput
	var err error

	if stepResultResource.Params != nil {
		params, err = stepResultResource.Params.ToParams()
		if err != nil {
			return StepResult{}, err
		}
	}
	if stepResultResource.Output != nil {
		output, err = stepResultResource.Output.ToOutput()
		if err != nil {
			return StepResult{}, err
		}
	}

	command, _ := StringToCommand(stepResultResource.Command)
	status, _ := StringToCommandStatus(stepResultResource.Status)

	resource := StepResult{
		Command: command,
		Params:  params,
		Output:  output,
		Status:  status,
	}
	return resource, nil
}
