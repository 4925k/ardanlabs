package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "expvar"

	"github.com/4925k/practical_go_foundations/nlp"
	"github.com/4925k/practical_go_foundations/nlp/stemmer"
	"github.com/gorilla/mux"
)

func main() {
	// LOGGING
	logger := log.New(log.Writer(), "nlp: ", log.LstdFlags|log.Lshortfile)
	s := Server{
		logger: logger, // DEPENDENCY INJECTION
	}

	// ROUTING

	// /health is an exact mathc
	// /health/ is a prefix match
	// http.HandleFunc("/health", healthHandler)
	// http.HandleFunc("/tokenize", tokenizeHandler)

	// USING A 3rd PARTY ROUTER
	r := mux.NewRouter()
	r.HandleFunc("/health", s.healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/tokenize", s.tokenizeHandler).Methods(http.MethodPost)
	r.HandleFunc("/stem/{word}", s.stemHandler).Methods(http.MethodGet)
	http.Handle("/", r)

	// SERVER
	addr := ":8080"
	s.logger.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

type Server struct {
	logger *log.Logger
}

func (s *Server) stemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	stem := stemmer.Stem(word)

	fmt.Fprintln(w, stem)
}

// healthHandler returns "OK"
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// tokenizeHandler tokenizes text
func (s *Server) tokenizeHandler(w http.ResponseWriter, r *http.Request) {
	// VALIDATION
	// NOTE: handled by gorilla mux
	// if r.Method != http.MethodPost {
	// 	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	// 	return
	// }

	// REQUEST PARSING
	defer r.Body.Close()
	rdr := io.LimitReader(r.Body, 1_000_000)
	data, err := io.ReadAll(rdr)
	if err != nil {
		s.logger.Printf("error reading request body: %v", err)
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
