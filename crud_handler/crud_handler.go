package crud_handler

import (
	"encoding/json"
	"fmt"
	"main/transformer"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TransformRequest struct {
	Type        string `json:"type"`
	CaesarShift int    `json:"shift,omitempty"`
	Input       string `json:"input,omitempty"`
}

type TransformResult struct {
	Id          string `json:"uuid"`
	Type        string `json:"type"`
	CaesarShift int    `json:"shift,omitempty"`
	Result      string `json:"result,omitempty"`
	Created_at  int64  `json:"created"`
	Updated_at  int64  `json:"updated,omitempty"`
}

// Fake database just for example
var temp1, temp2, temp3 = uuid.NewString(), uuid.NewString(), uuid.NewString()
var TransformRequests = map[string]*TransformRequest{
	temp1: &TransformRequest{Type: "reverse", Input: "123456789"},
	temp2: &TransformRequest{Type: "base64", Input: "Man"},
	temp3: &TransformRequest{Type: "caesar", CaesarShift: 1, Input: "zab"},
}
var TransformResults = map[string]*TransformResult{
	temp1: &TransformResult{Id: temp1, Type: "reverse", Result: "987654321", Created_at: time.Now().Unix()},
	temp2: &TransformResult{Id: temp2, Type: "base64", Result: "TWFu", Created_at: time.Now().Unix() + 1},
	temp3: &TransformResult{Id: temp3, Type: "caesar", CaesarShift: 1, Result: "abc", Created_at: time.Now().Unix() + 2},
}

func SetJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	request := new(TransformRequest)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if request.Type != "reverse" && request.Type != "caesar" && request.Type != "base64" {
		http.Error(w, "expected tranformation type field: reverse/caesar/base64", http.StatusInternalServerError)
		return
	}
	if request.Type == "caesar" && request.CaesarShift == 0 {
		http.Error(w, "expected shift field (not 0)", http.StatusInternalServerError)
		return
	}
	if request.Input == "" {
		http.Error(w, "expected input field", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Status ", http.StatusCreated, ": ", http.StatusText(http.StatusCreated))

	result := new(TransformResult)
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
	TransformResults[result.Id] = result

	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "uuid")
	delete(TransformResults, id)
	fmt.Fprint(w, "Status ", http.StatusNoContent, ": ", http.StatusText(http.StatusNoContent))
}

func HandleGetAll(w http.ResponseWriter, r *http.Request) {
	values := make([]*TransformResult, 0, len(TransformResults))
	for _, value := range TransformResults {
		values = append(values, value)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].Created_at > values[j].Created_at
	})
	enc := json.NewEncoder(w)
	err := enc.Encode(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	result := TransformResults[chi.URLParam(r, "uuid")]
	enc := json.NewEncoder(w)
	err := enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandlePut(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "uuid")
	request := new(TransformRequest)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if request.Type != "reverse" && request.Type != "caesar" && request.Type != "base64" {
		http.Error(w, "expected tranformation type field: reverse/caesar/base64", http.StatusInternalServerError)
		return
	}
	if request.Type == "caesar" && request.CaesarShift == 0 {
		http.Error(w, "expected shift field (not 0)", http.StatusInternalServerError)
		return
	}
	if request.Input == "" {
		http.Error(w, "expected input field", http.StatusInternalServerError)
		return
	}

	result, ok := TransformResults[id]
	if !ok {
		result = new(TransformResult)
		result.Id = id
		TransformResults[result.Id] = result
	}
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

	if result.Created_at != 0 {
		result.Updated_at = time.Now().Unix()
	} else {
		result.Created_at = time.Now().Unix()
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
