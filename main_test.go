package main

import (
	"bytes"
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
		t.Errorf("wrong status code returned: expected %v, got %v", status, http.StatusOK)
	}
}

func TestGetFoodByCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/produce/ABC-DEF", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodByCode)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("wrong status code returned: expected %v, got %v", status, http.StatusNotFound)
	}
}

func TestAddFood(t *testing.T) {
	var payload = []byte(`[{"code":"1234-5678-AAAA-BBBB", "name":"Celery", "price":"1.99"}, {"code":"9101-1121-CCCC-DDDD", "name": "Garlic", "price": "1.65"}]`)

	req, err := http.NewRequest("POST", "/groceries", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(AddFood)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("wrong status code returned: expected %v, got %v", status, http.StatusCreated)
	}
}

func TestDeleteFood(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/groceries/ABC-DEF", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteFood)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("wrong status code returned: expected %v, got %v", status, http.StatusNotFound)
	}
}