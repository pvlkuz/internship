package httpserver

import (
	"encoding/json"
	"fmt"
	"strings"

	"main/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (mock *MockService) CreateRecord(request models.TransformRequest) (*models.Record, error) {
	return &TestRecord, nil
}
func (mock *MockService) GetRecord(id string) (*models.Record, error) {
	return &TestRecord, nil
}
func (mock *MockService) GetAllRecords() ([]models.Record, error) {
	res := []models.Record{
		{
			ID:          "11",
			Type:        "reverse",
			CaesarShift: 0,
			Result:      "321",
		},
		{
			ID:          "12",
			Type:        "caesar",
			CaesarShift: -3,
			Result:      "xyz",
		},
	}

	return res, nil
}
func (mock *MockService) UpdateRecord(id string, request models.TransformRequest) *models.Record {
	return &TestRecord
}
func (mock *MockService) DeleteRecord(id string) error {
	return nil
}

var TestRequest = []models.TransformRequest{
	{Type: "caesar", CaesarShift: -3, Input: "abc"},
	{Type: "reverse", CaesarShift: 0, Input: "54321"},
	{Type: "base64", CaesarShift: 0, Input: "Man"},
}

var TestRecord = models.Record{
	ID:          "1111",
	Type:        "reverse",
	CaesarShift: 0,
	Result:      "12345",
}

func Test_CreateRecord(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	req, err := http.NewRequest("POST", "/records", strings.NewReader(fmt.Sprintf(`{"type":"%s", "input":"%s", "shift":%d}`, TestRequest[0].Type, TestRequest[0].Input, TestRequest[0].CaesarShift)))
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.NewRecord)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusCreated)
	}
	res := models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	if err != nil {
		t.Errorf("decoding error")
	}
	assert.Equal(t, TestRecord, res)
}

func Test_GetAllRecords(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	req, err := http.NewRequest("GET", "/records", strings.NewReader(""))
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAllRecords)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
	res := []models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	if err != nil {
		t.Errorf("decoding error")
	}

	assert.Equal(t, []models.Record{
		{
			ID:          "11",
			Type:        "reverse",
			CaesarShift: 0,
			Result:      "321",
		},
		{
			ID:          "12",
			Type:        "caesar",
			CaesarShift: -3,
			Result:      "xyz",
		},
	}, res)
}

func Test_GetRecord(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	req, err := http.NewRequest("GET", "/records/1111", strings.NewReader(""))
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetRecord)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
	res := models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	if err != nil {
		t.Errorf("decoding error")
	}

	assert.Equal(t, TestRecord, res)
}

func Test_UpdateRecord(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	req, err := http.NewRequest("PUT", "/records/1111", strings.NewReader(fmt.Sprintf(`{"type":"%s", "input":"%s", "shift":%d}`, TestRequest[0].Type, TestRequest[0].Input, TestRequest[0].CaesarShift)))
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateRecord)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
	res := models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	if err != nil {
		t.Errorf("decoding error")
	}

	assert.Equal(t, TestRecord, res)
}

func Test_DeleteRecord(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	MyURL := fmt.Sprintf("http://localhost/records/%s", "1111")
	req, err := http.NewRequest("DELETE", MyURL, nil)
	if err != nil {
		t.Fatalf("failer to create request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteRecord)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusNoContent)
	}

}
