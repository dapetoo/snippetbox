package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	//FileServer to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//Register the File server with mux.Handle
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
