package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTransactions(t *testing.T) {
	router := setupRouter()
	req, err := http.NewRequest("GET", "/api/users/1/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, req)
	if status := recoder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	t.Log(recoder.Body.String())
}

func TestGetTransactionByAccountID(t *testing.T) {
	router := setupRouter()
	req, err := http.NewRequest("GET", "/api/users/1/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}
	query := req.URL.Query()
	query.Add("account_id", "1")
	req.URL.RawQuery = query.Encode()
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, req)
	if status := recoder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	t.Log(recoder.Body.String())
}

func TestCreateTransaction(t *testing.T) {
	var jsonStr = []byte(`{
		"account_id": 1,
		"amount": 20.5,
		"transaction_type": "deposit"
	}`)
	router := setupRouter()
	req, err := http.NewRequest("POST", "/api/users/1/transactions", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, req)
	if status := recoder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	t.Log(recoder.Body.String())
}

func TestUpdateTransaction(t *testing.T) {
	var jsonStr = []byte(`{
		"amount": 1.5,
		"transaction_type": "deposit"
	}`)
	router := setupRouter()
	req, err := http.NewRequest("PATCH", "/api/users/1/transactions/8", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, req)
	if status := recoder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	t.Log(recoder.Body.String())
}

func TestDeleteTransaction(t *testing.T) {
	router := setupRouter()
	req, err := http.NewRequest("DELETE", "/api/users/1/transactions/8", nil)
	if err != nil {
		t.Fatal(err)
	}
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, req)
	if status := recoder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	t.Log(recoder.Body.String())
}
