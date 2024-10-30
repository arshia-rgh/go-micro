package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(request *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func Test_Authenticate(t *testing.T) {
	postBody := map[string]interface{}{
		"email":    "me@here.com",
		"password": "verysecret",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected the http.StatusAccepted but got %d", rr.Code)
	}
}
