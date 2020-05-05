package client

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-goclient/pkg/models"
)

func TestBuildSnippetsQueryParams(t *testing.T) {
	id := "123"
	q := buildSnippetsQueryParams(SnippetFilter{
		UsageID: &id,
	})

	assert.Equal(t, "?usage_id=123", q)

	c := "tag"
	id1 := "123"
	cat, ok := models.StringToSnippetCategory(c)
	if !ok {
		log.Printf("Can't convert snippet resource category")
	}

	q1 := buildSnippetsQueryParams(SnippetFilter{
		Category: &cat,
		UsageID:  &id1,
	})

	assert.Equal(t, "?usage_id=123&category=tag", q1)
}

func TestSnippetsGetAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/snippets")
		resp, err := json.Marshal([]models.SnippetResource{{
			Name:        "name",
			Category:    "tag",
			UsageID:     "id",
			UsageName:   "uname",
			Description: "dec",
			Code:        "code",
		}})
		if err != nil {
			log.Fatal("Can't marshall test model", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return err", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := newSnippetService(server.Client(), server.URL)
	body, err := api.GetAll()

	if err != nil {
		log.Fatal("Can't get events", err)
	}

	cat, ok := models.StringToSnippetCategory("tag")
	if !ok {
		log.Printf("Can't convert snippet resource category")
	}

	assert.Equal(t, "id", body[0].UsageID)
	assert.Equal(t, cat, body[0].Category)
	assert.Equal(t, "name", body[0].Name)
	assert.Equal(t, "code", body[0].Code)
}

func TestSnippetsGetFiltered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/snippets?usage_id=123")
		resp, err := json.Marshal([]models.SnippetResource{{
			Name:        "name",
			Category:    "tag",
			UsageID:     "id",
			UsageName:   "uname",
			Description: "dec",
			Code:        "code",
		}})
		if err != nil {
			log.Fatal("Can't marshall test model", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return err", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := newSnippetService(server.Client(), server.URL)
	id := "123"
	body, err := api.GetFiltered(SnippetFilter{
		UsageID: &id,
	})

	if err != nil {
		log.Fatal("Can't get events", err)
	}

	cat, ok := models.StringToSnippetCategory("tag")
	if !ok {
		log.Printf("Can't convert snippter resource category")
	}

	assert.Equal(t, "id", body[0].UsageID)
	assert.Equal(t, cat, body[0].Category)
	assert.Equal(t, "name", body[0].Name)
	assert.Equal(t, "code", body[0].Code)
}
