package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllFoods(t *testing.T) {
	req, err := http.NewRequest("GET", "/produce", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllFoods)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Error("wrong status code returned: expected %v, got %v", status, http.StatusOK)
	}
}