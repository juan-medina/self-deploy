package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestRandom(t *testing.T) {
	server := NewServer()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/random", nil)
	server.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusOK
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
	body := response.Body.String()
	_, err := strconv.ParseInt(body, 10, 64)
	if err != nil {
		t.Errorf("we didn't get a number, got '%s', error was %s", body, err)
	}
}

func TestNotFound(t *testing.T) {
	server := NewServer()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/not/found", nil)
	server.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusNotFound
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
