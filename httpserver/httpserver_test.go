package httpserver

import (
	"bytes"
	"encoding/json"

	"main/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

var TestRequest = models.TransformRequest{
	Type:        "caesar",
	CaesarShift: -3,
	Input:       "abc",
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

	data, err := json.Marshal(TestRequest)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/records", bytes.NewReader(data))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateRecord)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	res := models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	assert.ErrorIs(t, err, nil)

	assert.Equal(t, TestRecord, res)
}

func Test_GetAllRecords(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	req, err := http.NewRequest("GET", "/records", nil)
	require.ErrorIs(t, err, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAllRecords)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	res := []models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	assert.ErrorIs(t, err, nil)

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

	req, err := http.NewRequest("GET", "/records/1111", nil)
	require.ErrorIs(t, err, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetRecord)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	res := models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	assert.ErrorIs(t, err, nil)

	assert.Equal(t, TestRecord, res)
}

func Test_UpdateRecord(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	data, err := json.Marshal(TestRequest)
	require.NoError(t, err)
	req, err := http.NewRequest("PUT", "/records/1111", bytes.NewReader(data))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateRecord)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	res := models.Record{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(&res)
	assert.ErrorIs(t, err, nil)

	assert.Equal(t, TestRecord, res)
}

func Test_DeleteRecord(t *testing.T) {
	service := new(MockService)
	h := NewHandler(service)

	req, err := http.NewRequest("DELETE", "/records/1111", nil)
	require.ErrorIs(t, err, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteRecord)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)

}
