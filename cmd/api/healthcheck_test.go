package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	app := &application{}

	req, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.healthcheckHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, response["status"].(bool))

	respTimeStr := response["time"].(string)
	actualTime, err := time.Parse(time.RFC3339, respTimeStr)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	expectedTime := now.Round(time.Second)
	actualTime = actualTime.Round(time.Second)

	assert.WithinDuration(t, expectedTime, actualTime, time.Second)
}
