package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/4925k/practical_go_foundations/nlp"
)

func main() {
	// ROUTING

	// /health is an exact mathc
	// /health/ is a prefix match
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/tokenize", tokenizeHandler)

	// SERVER
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// healthHandler returns "OK"
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// tokenizeHandler tokenizes text
func tokenizeHandler(w http.ResponseWriter, r *http.Request) {
	// VALIDATION
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// REQUEST PARSING
	defer r.Body.Close()
	rdr := io.LimitReader(r.Body, 1_000_000)
	data, err := io.ReadAll(rdr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// VALIDATION
	if len(data) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	text := string(data)

	// tokenize
	tokens := nlp.Tokenize(text)

	// RESPONSE
	resp := map[string]any{
		"tokens": tokens,
	}

	ret, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}
