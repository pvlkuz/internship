package handler

import (
	"main/transformer"
	"net/http"
	"strconv"
)

func ReverseHandler(w http.ResponseWriter, r *http.Request) {
	result, err := transformer.NewReverseTransformer().Transform(r.Body)
	if err != nil {
		http.Error(w, "Server ReverseHandler Transformer error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, "server write output error", http.StatusInternalServerError)
		return
	}
}

func CaesarHandler(w http.ResponseWriter, r *http.Request) {
	shiftStr := r.URL.Query().Get("shift")
	shift, err := strconv.Atoi(shiftStr)
	if err != nil {
		http.Error(w, "No integer shift given", http.StatusInternalServerError)
		return
	}
	result, err := transformer.NewCaesarTransformer(shift).Transform(r.Body)
	if err != nil {
		http.Error(w, "Server CaesarHandler Transformer error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, "server write output error", http.StatusInternalServerError)
		return
	}
}

func Base64Handler(w http.ResponseWriter, r *http.Request) {
	result, err := transformer.NewBase64Transformer().Transform(r.Body)
	if err != nil {
		http.Error(w, "Server Bae64Handler Transformer error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, "server write output error", http.StatusInternalServerError)
		return
	}
}
