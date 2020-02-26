package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-client/pkg/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuildJobsQueryParams(t *testing.T) {
	l := 25
	q := buildRunsQueryParams(RunFilter{
		Limit: &l,
	})

	assert.Equal(t, "?limit=25", q)

	dir := "asc"
	by := "n"
	offset := 321
	status := models.JobStatusActive

	q1 := buildJobsQueryParams(JobFilter{
		SortBy:  &by,
		SortDir: &dir,
		Offset:  &offset,
		Status:  &status,
	})

	assert.Equal(t, "?status=active&sortby=n&sortdir=asc&offset=321", q1)
}

func TestJobsGetAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/adapters/id/jobs", req.URL.String())
		resp, err := json.Marshal(models.JobListResource{
			Total:  0,
			Length: 0,
			Limit:  0,
			Offset: 0,
			Items: []models.JobResource{{
				JobID:       "id",
				JobName:     "name",
				AdapterID:   "adid",
				AdapterName: "adname",
				CreatedAt:   "2006-01-02T15:04:05Z",
				Status:      models.JobStatusActive.String(),
				Steps: []models.JobStepResource{{
					Command: models.CommandRemovePassword.String(),
					Params:  nil,
				}},
			}},
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := newJobService(server.Client(), server.URL)
	body, err := api.GetAll("id")

	if err != nil {
		log.Fatal("Can't get runs\n", err)
	}

	assert.Equal(t, "id", body[0].JobID)
	assert.Equal(t, "name", body[0].JobName)
	assert.Equal(t, "adname", body[0].AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", body[0].CreatedAt.Format(time.RFC3339))
	assert.Equal(t, models.JobStatusActive, body[0].Status)
	assert.Equal(t, models.CommandRemovePassword, body[0].Steps[0].Command)
}

func TestJobsGetFiltered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/adapters/id/jobs?status=active&sortdir=asc", req.URL.String())
		resp, err := json.Marshal(models.JobListResource{
			Total:  0,
			Length: 0,
			Limit:  0,
			Offset: 0,
			Items: []models.JobResource{{
				JobID:       "id",
				JobName:     "name",
				AdapterID:   "adid",
				AdapterName: "adname",
				CreatedAt:   "2006-01-02T15:04:05Z",
				Status:      models.JobStatusActive.String(),
				Steps: []models.JobStepResource{{
					Command: models.CommandRemovePassword.String(),
					Params:  nil,
				}},
			}},
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := newJobService(server.Client(), server.URL)
	s := models.JobStatusActive
	sortDir := "asc"
	body, err := api.GetFiltered("id", JobFilter{
		Status: &s,
		SortDir: &sortDir,
	})

	if err != nil {
		log.Fatal("Can't get runs\n", err)
	}

	assert.Equal(t, "id", body[0].JobID)
	assert.Equal(t, "name", body[0].JobName)
	assert.Equal(t, "adname", body[0].AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", body[0].CreatedAt.Format(time.RFC3339))
	assert.Equal(t, models.JobStatusActive, body[0].Status)
	assert.Equal(t, models.CommandRemovePassword, body[0].Steps[0].Command)
}