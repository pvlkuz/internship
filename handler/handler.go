package handler

import (
	"encoding/json"
	"log"
	"main/repo"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service interface {
	CreateRecord(request repo.TransformRequest) (*repo.Record, error)
	GetRecord(id string) (*repo.Record, error)
	GetAllRecords() ([]repo.Record, error)
	UpdateRecord(id string, request repo.TransformRequest) *repo.Record
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

func (h *Handler) RunServer() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/records", h.NewRecord)
	router.Get("/records", h.GetAllRecords)
	router.Get("/records/{id}", h.GetRecord)
	router.Delete("/records/{id}", h.DeleteRecord)
	router.Put("/records/{id}", h.UpdateRecord)

	//nolint:exhaustivestruct, exhaustruct
	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("server listenig error")
	}
}

func ResponseWithJSON(w http.ResponseWriter, statusCode int, record interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(record)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewRecord(w http.ResponseWriter, r *http.Request) {
	var request *repo.TransformRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = request.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := h.service.CreateRecord(*request)
	if err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, nil)
	}

	ResponseWithJSON(w, http.StatusCreated, result)
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllRecords(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.GetAllRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if values == nil {
		ResponseWithJSON(w, http.StatusNoContent, nil)
		return
	}

	ResponseWithJSON(w, http.StatusOK, values)
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.service.GetRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJSON(w, http.StatusOK, result)
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	request := new(repo.TransformRequest)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = request.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.service.UpdateRecord(id, *request)

	ResponseWithJSON(w, http.StatusOK, result)
}
