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

func TestBuildRunQueryParams(t *testing.T) {
	l := 25
	q := buildRunsQueryParams(RunFilter{
		Limit: &l,
	})

	assert.Equal(t, "?limit=25", q)

	dir := "asc"
	by := "n"
	offset := 321
	jId := "123"
	status := models.JobRunStatusStarted

	q1 := buildRunsQueryParams(RunFilter{
		SortBy:  &by,
		SortDir: &dir,
		Offset:  &offset,
		Status:  &status,
		JobID:   &jId,
	})

	assert.Equal(t, "?status=started&job_id=123&sortby=n&sortdir=asc&offset=321", q1)
}

func TestRunsGetAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/adapters/id/runs", req.URL.String())
		resp, err := json.Marshal(models.JobRunListResource{
			Total:  0,
			Length: 0,
			Limit:  0,
			Offset: 0,
			Items: []models.JobRunResource{{
				JobID:       "id",
				JobName:     "name",
				AdapterID:   "adid",
				AdapterName: "adname",
				Results:     []models.StepResultResource{{Command: models.CommandGetDump.String()}},
				CreatedAt:   "2006-01-02T15:04:05Z",
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

	api := newRunService(server.Client(), server.URL)
	body, pagInfo, err := api.GetAll("id")

	if err != nil {
		log.Fatal("Can't get runs\n", err)
	}

	assert.Equal(t, "id", body[0].JobID)
	assert.Equal(t, "name", body[0].JobName)
	assert.Equal(t, "adname", body[0].AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", body[0].CreatedAt.Format(time.RFC3339))
	assert.Equal(t, models.CommandGetDump, body[0].Results[0].Command)
	assert.Equal(t, 0, pagInfo.Total)
	assert.Equal(t, 0, pagInfo.Offset)
	assert.Equal(t, 0, pagInfo.Limit)
	assert.Equal(t, 0, pagInfo.Length)
}

func TestRunsGetFiltered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/adapters/id/runs?job_id=id", req.URL.String())
		resp, err := json.Marshal(models.JobRunListResource{
			Total:  0,
			Length: 0,
			Limit:  0,
			Offset: 0,
			Items: []models.JobRunResource{{
				JobID:       "id",
				JobName:     "name",
				AdapterID:   "adid",
				AdapterName: "adname",
				Results:     []models.StepResultResource{{Command: models.CommandGetDump.String()}},
				CreatedAt:   "2006-01-02T15:04:05Z",
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

	api := newRunService(server.Client(), server.URL)
	jId := "id"
	body, pagInfo, err := api.GetFiltered("id", RunFilter{
		JobID: &jId,
	})

	if err != nil {
		log.Fatal("Can't get runs\n", err)
	}

	assert.Equal(t, "id", body[0].JobID)
	assert.Equal(t, "name", body[0].JobName)
	assert.Equal(t, "adname", body[0].AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", body[0].CreatedAt.Format(time.RFC3339))
	assert.Equal(t, models.CommandGetDump, body[0].Results[0].Command)
	assert.Equal(t, 0, pagInfo.Total)
	assert.Equal(t, 0, pagInfo.Offset)
	assert.Equal(t, 0, pagInfo.Limit)
	assert.Equal(t, 0, pagInfo.Length)
}
