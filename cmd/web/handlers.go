package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	_, err := w.Write([]byte("Hello from Snippet box"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	_, err = fmt.Fprintf(w, "Display a specifc snippet with ID %d...", id)
	if err != nil {
		log.Println(err)
	}
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(405)
		http.Error(w, "Method Not Allowed", 405)
		_, err := w.Write([]byte("Method Not Allowed"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("Create  a new snippet"))
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}
