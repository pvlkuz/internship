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

type ServiceInterface interface {
	CreateRecord(request repo.TransformRequest) *repo.Record
	GetRecord(id string) (*repo.Record, error)
	GetRecords() ([]repo.Record, error)
	UpdateRecord(id string, request repo.TransformRequest) *repo.Record
	DeleteRecord(id string) error
}

type Handler struct {
	service ServiceInterface
}

type DBLayer interface {
	NewRecord(r *repo.Record) error
	GetRecord(id string) (repo.Record, error)
	GetRecords() ([]repo.Record, error)
	UpdateRecord(r *repo.Record) error
	DeleteRecord(id string) error
}

func NewHandler(service ServiceInterface) *Handler {
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

func ResponseWithJSON(w http.ResponseWriter, result *repo.Record, results *[]repo.Record, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var err error

	if result != nil {
		err = json.NewEncoder(w).Encode(result)
	} else {
		err = json.NewEncoder(w).Encode(results)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewRecord(w http.ResponseWriter, r *http.Request) {
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

	result := h.service.CreateRecord(*request)

	ResponseWithJSON(w, result, nil, http.StatusCreated)
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
	values, err := h.service.GetRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJSON(w, nil, &values, http.StatusOK)
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.service.GetRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJSON(w, result, nil, http.StatusOK)
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

	ResponseWithJSON(w, result, nil, http.StatusOK)
}
