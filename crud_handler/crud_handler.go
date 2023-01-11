package crud_handler

import (
	"encoding/json"
	"fmt"
	"log"
	database "main/DataBase"
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
	router.Post("/records", h.HandlePost)
	router.Get("/records", h.HandleGetAll)
	router.Get("/records/{uuid}", h.HandleGet)
	router.Delete("/records/{uuid}", h.HandleDelete)
	router.Put("/records/{uuid}", h.HandlePut)
	log.Fatal(http.ListenAndServe(":8080", router))
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

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprint(w, "Created")

	result := new(repo.Record)
	result.Id = uuid.NewString()
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
	result.Result, err = tr.Transform(strings.NewReader(request.Input))
	if err != nil {
		http.Error(w, "Server Transformer error", http.StatusInternalServerError)
		return
	}
	result.Created_at = time.Now().Unix()

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

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "uuid")
	err := h.db.DeleteRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "Deleted")
}

func (h *Handler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	values, err := h.db.GetRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].Created_at > values[j].Created_at
	})
	enc := json.NewEncoder(w)
	err = enc.Encode(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "uuid")
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

func (h *Handler) HandlePut(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "uuid")
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
	transform_result, err := tr.Transform(strings.NewReader(request.Input))
	if err != nil {
		http.Error(w, "Server Transformer error", http.StatusInternalServerError)
		return
	}

	result := &repo.Record{
		Id:          id,
		Type:        request.Type,
		CaesarShift: request.CaesarShift,
		Result:      transform_result,
		Updated_at:  time.Now().Unix(),
	}
	err = h.db.UpdateRecord(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		result.Created_at = time.Now().Unix()
		err = h.db.NewRecord(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
