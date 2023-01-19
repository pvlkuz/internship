package crud_handler

import (
	"encoding/json"
	"log"
	database "main/data-base"
	"main/repo"
	"main/transformer"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Handler struct {
	db database.RecordDB
}

func NewHandler(db database.RecordDB) *Handler {
	return &Handler{
		db: db,
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

	//log.Fatal(http.ListenAndServe(":8080", router))
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

	result := new(repo.Record)
	result.ID = uuid.NewString()
	result.Type = request.Type
	result.CaesarShift = request.CaesarShift
	var tr transformer.Transformer
	switch {
	case request.Type == "reverse":
		tr = transformer.NewReverseTransformer()
	case request.Type == "caesar":
		tr = transformer.NewCaesarTransformer(request.CaesarShift)
	case request.Type == "base64":
		tr = transformer.NewBase64Transformer()
	}
	result.Result, err = tr.Transform(strings.NewReader(request.Input), false)
	if err != nil {
		http.Error(w, "Server Transformer error", http.StatusInternalServerError)
		return
	}
	result.CreatedAt = time.Now().Unix()

	err = h.db.NewRecord(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.db.DeleteRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllRecords(w http.ResponseWriter, r *http.Request) {
	values, err := h.db.GetRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].CreatedAt > values[j].CreatedAt
	})
	enc := json.NewEncoder(w)
	err = enc.Encode(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := h.db.GetRecord(id)
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

	var tr transformer.Transformer
	switch {
	case request.Type == "reverse":
		tr = transformer.NewReverseTransformer()
	case request.Type == "caesar":
		tr = transformer.NewCaesarTransformer(request.CaesarShift)
	case request.Type == "base64":
		tr = transformer.NewBase64Transformer()
	}
	transform_result, err := tr.Transform(strings.NewReader(request.Input), false)
	if err != nil {
		http.Error(w, "Server Transformer error", http.StatusInternalServerError)
		return
	}

	result, err := h.db.GetRecord(id)
	result.Type = request.Type
	result.CaesarShift = request.CaesarShift
	result.Result = transform_result
	result.UpdatedAt = time.Now().Unix()
	if err != nil {
		result.ID = id
		result.CreatedAt = time.Now().Unix()
		err = h.db.NewRecord(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = h.db.UpdateRecord(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
