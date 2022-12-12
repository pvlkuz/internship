package crud_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Record struct {
	Id          string `json:"uuid"`
	Type        string `json:"type"`
	CaesarShift int    `json:"shift,omitempty"`
	Input       string `json:"input,omitempty"`
}

// Fake database just for example (fake id key especially for map_key to be real uuid)
var records = map[string]*Record{
	uuid.NewString(): &Record{Id: "1", Type: "reverse", Input: "123456789"},
	uuid.NewString(): &Record{Id: "2", Type: "base64", Input: "Man"},
	uuid.NewString(): &Record{Id: "3", Type: "caesar", CaesarShift: 1, Input: "zab"},
}

func HandleRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch r.Method {
	case "GET":
		record := records[vars["uuid"]]
		js, err := json.Marshal(record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(js))
		return

	case "PUT":
		record := records[vars["uuid"]]
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		record.Id = vars["uuid"]
		result, err := json.Marshal(record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(result))
	case "DELETE":
		delete(records, vars["uuid"])
		fmt.Fprint(w, "Status ", http.StatusNoContent, ": ", http.StatusText(http.StatusNoContent))
	}
}

func HandleRecords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		js, err := json.Marshal(records)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(js)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		record := new(Record)
		record.Id = uuid.NewString()
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if record.Type != "reverse" && record.Type != "caesar" && record.Type != "base64" {
			http.Error(w, "expected tranformation type field: reverse/caesar/base64", http.StatusInternalServerError)
			return
		}
		if record.Type == "caesar" && record.CaesarShift == 0 {
			http.Error(w, "expected shift field (not 0)", http.StatusInternalServerError)
			return
		}
		if record.Input == "" {
			http.Error(w, "expected input field", http.StatusInternalServerError)
			return
		}
		records[record.Id] = record

		result, err := json.Marshal(record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Status ", http.StatusCreated, ": ", http.StatusText(http.StatusCreated))
		log.SetOutput(w)
		log.Print(string(result))
	}
}
