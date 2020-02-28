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

func TestBuildQueryParams(t *testing.T) {
	id := "123"
	q := buildEventsQueryParams(&id, EventFilter{
		Name:    nil,
		SortBy:  nil,
		SortDir: nil,
		Offset:  nil,
		Limit:   nil,
	})

	assert.Equal(t, "?adapter_id=123", q)

	dir := "asc"
	by := "n"
	offset := 321
	l := 123

	q1 := buildEventsQueryParams(nil, EventFilter{
		Name:    nil,
		SortBy:  &by,
		SortDir: &dir,
		Offset:  &offset,
		Limit:   &l,
	})

	assert.Equal(t, "?sortby=n&sortdir=asc&offset=321&limit=123", q1)
}

func TestEventsGetAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/events")
		resp, err := json.Marshal(models.EventListResource{
			Total:  0,
			Length: 0,
			Limit:  0,
			Offset: 0,
			Items: []models.EventResource{{
				EventID:     "id",
				Name:        models.EventNameAdapterDiscovery.String(),
				AdapterID:   "adid",
				AdapterName: "adname",
				Data:        nil,
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

	api := newEventService(server.Client(), server.URL)
	body, pagInfo, err := api.GetAll()

	if err != nil {
		log.Fatal("Can't get events\n", err)
	}

	assert.Equal(t, "id", body[0].EventID)
	assert.Equal(t, "adid", body[0].AdapterID)
	assert.Equal(t, "adname", body[0].AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", body[0].CreatedAt.Format(time.RFC3339))
	assert.Equal(t, 0, pagInfo.Total)
	assert.Equal(t, 0, pagInfo.Offset)
	assert.Equal(t, 0, pagInfo.Limit)
	assert.Equal(t, 0, pagInfo.Length)
}

func TestEventsGetAll2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/events")
		resp, err := json.Marshal(models.ErrorResponse{
			Message: "err msg",
			Info:    "err info",
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(500)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newEventService(server.Client(), server.URL)
	_, _, err := api.GetAll()

	assert.NotNil(t, err)
}

func TestEventsGetFiltered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/events?adapter_id=1&limit=1", req.URL.String())
		resp, err := json.Marshal(models.EventListResource{
			Total:  0,
			Length: 0,
			Limit:  0,
			Offset: 0,
			Items: []models.EventResource{{
				EventID:     "id",
				Name:        models.EventNameAdapterDiscovery.String(),
				AdapterID:   "adid",
				AdapterName: "adname",
				Data:        nil,
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

	api := newEventService(server.Client(), server.URL)
	id := "1"
	l := 1
	body, pagInfo, err := api.GetFiltered(&id, EventFilter{
		Limit: &l,
	})

	if err != nil {
		log.Fatal("Can't get events\n", err)
	}

	assert.Equal(t, "id", body[0].EventID)
	assert.Equal(t, "adid", body[0].AdapterID)
	assert.Equal(t, "adname", body[0].AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", body[0].CreatedAt.Format(time.RFC3339))
	assert.Equal(t, 0, pagInfo.Total)
	assert.Equal(t, 0, pagInfo.Offset)
	assert.Equal(t, 0, pagInfo.Limit)
	assert.Equal(t, 0, pagInfo.Length)
}

func TestEventsAdd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/events")
		resp, err := json.Marshal(models.EventResource{
			EventID:     "eid",
			Name:        models.EventNameAdapterDiscovery.String(),
			AdapterID:   "id",
			AdapterName: "aid",
			Data:        nil,
			CreatedAt:   "2006-01-02T15:04:05Z",
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

	api := newEventService(server.Client(), server.URL)
	body, err := api.Add(models.NewEvent{
		Name:      models.EventNameAdapterDiscovery.String(),
		AdapterID: "id",
		Data:      nil,
	})

	if err != nil {
		log.Fatal("Can't post event\n", err)
	}

	assert.Equal(t, models.EventNameAdapterDiscovery, body.Name)
	assert.Equal(t, "id", body.AdapterID)
	assert.Equal(t, "eid", body.EventID)
	assert.Equal(t, "2006-01-02T15:04:05Z", body.CreatedAt.Format(time.RFC3339))
}

func TestEventsGet2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/events")
		resp, err := json.Marshal(models.ErrorResponse{
			Message: "err msg",
			Info:    "err info",
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(500)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	api := newEventService(server.Client(), server.URL)
	_, err := api.Add(models.NewEvent{
		Name:      "name",
		AdapterID: "id",
		Data:      nil,
	})

	assert.NotNil(t, err)
}
