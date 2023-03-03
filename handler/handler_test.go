package handler

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"main/cache"
	"main/models"
	"main/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testTime = time.Now()

type MockDB struct {
	mock.Mock
}

func (mock *MockDB) CreateRecord(r *models.Record) error {
	return nil
}
func (mock *MockDB) GetRecord(id string) (models.Record, error) {
	result := models.Record{
		ID:          "1111",
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
		CreatedAt:   testTime,
	}
	return result, nil
}
func (mock *MockDB) GetAllRecords() ([]models.Record, error) {
	result := []models.Record{
		{
			ID:          uuid.NewString(),
			Type:        "reverse",
			CaesarShift: 0,
			Result:      "54321",
			CreatedAt:   testTime,
		},
		{
			ID:          uuid.NewString(),
			Type:        "caesar",
			CaesarShift: -3,
			Result:      "xyz",
			CreatedAt:   testTime,
		},
	}
	return result, nil
}
func (mock *MockDB) UpdateRecord(r *models.Record) error {
	return nil
}
func (mock *MockDB) DeleteRecord(id string) error {
	return nil
}

var NewRecordRequestTable = []models.TransformRequest{
	{Type: "caesar", CaesarShift: -3, Input: "abc"},
	{Type: "reverse", CaesarShift: 0, Input: "54321"},
	{Type: "base64", CaesarShift: 0, Input: "Man"},
}
var NewRecordResultTable = []string{
	"xyz", "12345", "TWFu",
}

func Test_NewRecordHandler(t *testing.T) {
	db := new(MockDB)
	service := service.NewService(db, nil)
	h := NewHandler(service)

	for i, test := range NewRecordRequestTable {
		req, err := http.NewRequest("POST", "/records", strings.NewReader(fmt.Sprintf(`{"type":"%s", "input":"%s", "shift":%d}`, test.Type, test.Input, test.CaesarShift)))
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

		res := new(models.Record)
		dec := json.NewDecoder(rr.Body)
		err = dec.Decode(&res)
		if err != nil {
			t.Errorf("decoding error")
		}
		if res.Type != test.Type || res.CaesarShift != test.CaesarShift || res.Result != NewRecordResultTable[i] {
			t.Errorf("Mismatch result=%s, %d, %s", res.Type, res.CaesarShift, res.Result)
		}
	}
}

var GetAllRecordsTestTable = []models.TransformRequest{
	{Type: "reverse", CaesarShift: 0, Input: "54321"},
	{Type: "caesar", CaesarShift: -3, Input: "xyz"},
}

func Test_GetAllRecordsHandler(t *testing.T) {
	db := new(MockDB)
	service := service.NewService(db, nil)
	h := NewHandler(service)

	req, err := http.NewRequest("GET", "/records", nil)
	if err != nil {
		t.Fatalf("failer to create request: %s", err)
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
		t.Errorf("decoding out array error")
	}
	for i, test := range res {
		if test.Type != GetAllRecordsTestTable[i].Type || test.CaesarShift != GetAllRecordsTestTable[i].CaesarShift || test.Result != GetAllRecordsTestTable[i].Input {
			t.Errorf("Mismatch result=%s, %d, %s", test.Type, test.CaesarShift, test.Result)
		}
	}
}

func Test_GetRecordHandler(t *testing.T) {
	db := new(MockDB)
	cache := cache.NewInMemoCache()
	service := service.NewService(db, cache)
	h := NewHandler(service)

	reverseTs := httptest.NewServer(http.HandlerFunc(h.GetRecord))
	defer reverseTs.Close()
	MyURL := fmt.Sprintf("%s/%s", reverseTs.URL, "1111")
	res, err := http.Get(MyURL)
	//http.Post(reverseTs.URL, "application/json", strings.NewReader(in))
	if err != nil {
		t.Errorf("error in GetAll request")
	}

	result := models.Record{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&result)
	if err != nil {
		t.Errorf("decoding out array error")
	}
	ExpectedResult := models.Record{
		ID:          "1111",
		Type:        "reverse",
		CaesarShift: 0,
		Result:      "54321",
		CreatedAt:   testTime,
	}

	assert.Equal(t, ExpectedResult.ID, result.ID)
	assert.Equal(t, ExpectedResult.Type, result.Type)
	assert.Equal(t, ExpectedResult.CaesarShift, result.CaesarShift)
	assert.Equal(t, ExpectedResult.Result, result.Result)
}

var UpdateRecordTestTable = []models.TransformRequest{
	{Type: "caesar", CaesarShift: -3, Input: "abc"},
	{Type: "reverse", CaesarShift: 0, Input: "54321"},
	{Type: "base64", CaesarShift: 0, Input: "Man"},
}
var UpdateRecordResultTable = []string{
	"xyz", "12345", "TWFu",
}

func Test_UpdateRecord(t *testing.T) {
	db := new(MockDB)
	cache := cache.NewInMemoCache()
	service := service.NewService(db, cache)
	h := NewHandler(service)

	MyURL := fmt.Sprintf("http://localhost/records/%s", "1111")
	for i, test := range UpdateRecordTestTable {
		req, err := http.NewRequest("PUT", MyURL, strings.NewReader(fmt.Sprintf(`{"type":"%s", "input":"%s", "shift":%d}`, test.Type, test.Input, test.CaesarShift)))
		if err != nil {
			t.Fatalf("failer to create request: %s", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.UpdateRecord)
		handler.ServeHTTP(rr, req)

		res := new(models.Record)
		dec := json.NewDecoder(rr.Body)
		err = dec.Decode(&res)
		if err != nil {
			t.Errorf("decoding error")
		}
		if res.Type != test.Type || res.CaesarShift != test.CaesarShift || res.Result != UpdateRecordResultTable[i] {
			t.Errorf("Mismatch result=%s, %d, %s", res.Type, res.CaesarShift, res.Result)
		}
	}
}

func Test_DeleteRecord(t *testing.T) {
	db := new(MockDB)
	cache := cache.NewInMemoCache()
	service := service.NewService(db, cache)
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
