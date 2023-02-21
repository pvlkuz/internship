package crud_handler

import (
	"encoding/json"
	"log"
	"main/repo"
	"main/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	service service.ServiceInterface
}

type DBLayer interface {
	NewRecord(r *repo.Record) error
	GetRecord(id string) (repo.Record, error)
	GetRecords() ([]repo.Record, error)
	UpdateRecord(r *repo.Record) error
	DeleteRecord(id string) error
}

func NewHandler(service service.ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

type TransformRequest struct {
	Type        string `json:"type"`
	CaesarShift int    `json:"shift,omitempty"`
	Input       string `json:"input,omitempty"`
}

func (h *Handler) RunServer() {
	router := chi.NewRouter()
	router.Use(SetJSONContentType)
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
		log.Fatal("server listenig error")
	}
}

func SetJSONContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func CheckValidRequest(request *TransformRequest) string {
	if request.Type != "reverse" && request.Type != "caesar" && request.Type != "base64" {
		return "expected tranformation type field: reverse/caesar/base64"
	}
	if request.Type == "caesar" && request.CaesarShift == 0 {
		return "expected shift field (not 0)"
	}
	if request.Input == "" {
		return "expected input field"
	}
	return ""
}

func (h *Handler) NewRecord(w http.ResponseWriter, r *http.Request) {
	request := new(TransformRequest)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	invalid := CheckValidRequest(request)
	if invalid != "" {
		http.Error(w, invalid, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	result := h.service.NewRecord(service.TransformRequest(*request))
	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.DeleteRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllRecords(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.GetRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := h.service.GetRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	request := new(TransformRequest)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	invalid := CheckValidRequest(request)
	if invalid != "" {
		http.Error(w, invalid, http.StatusBadRequest)
		return
	}
	result := h.service.UpdateRecord(id, service.TransformRequest(*request))

	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
