package httpserver

import (
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service interface {
	CreateRecord(request models.TransformRequest) (*models.Record, error)
	GetRecord(id string) (*models.Record, error)
	GetAllRecords() ([]models.Record, error)
	UpdateRecord(id string, request models.TransformRequest) *models.Record
	DeleteRecord(id string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RunServer() error {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Post("/records", h.NewRecord)
	router.Get("/records", h.GetAllRecords)
	router.Get("/records/{id}", h.GetRecord)
	router.Delete("/records/{id}", h.DeleteRecord)
	router.Put("/records/{id}", h.UpdateRecord)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server listenig error: %w", err)
	}

	return nil
}

func ResponseWithJSON(w http.ResponseWriter, statusCode int, record interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(record)

	if err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) NewRecord(w http.ResponseWriter, r *http.Request) {
	var request *models.TransformRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseWithJSON(w, http.StatusBadRequest, err.Error())
	}

	err = request.Validate()
	if err != nil {
		ResponseWithJSON(w, http.StatusUnprocessableEntity, err.Error())
	}

	result, err := h.service.CreateRecord(*request)
	if err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
	}

	ResponseWithJSON(w, http.StatusCreated, result)
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteRecord(id)
	if err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllRecords(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.GetAllRecords()
	if err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
	}

	ResponseWithJSON(w, http.StatusOK, values)
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.service.GetRecord(id)
	if err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
	}

	ResponseWithJSON(w, http.StatusOK, result)
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	request := new(models.TransformRequest)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseWithJSON(w, http.StatusBadRequest, err.Error())
	}

	err = request.Validate()
	if err != nil {
		ResponseWithJSON(w, http.StatusUnprocessableEntity, err.Error())
	}

	result := h.service.UpdateRecord(id, *request)

	ResponseWithJSON(w, http.StatusOK, result)
}
