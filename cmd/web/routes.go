package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//Using Chi as a Router
	mux := chi.NewRouter()
	mux.Get("/", app.home)
	mux.Get("/snippet/create", app.createSnippet)
	mux.Post("/snippet/create", app.createSnippetForm)
	mux.Get("/snippet/{id}", app.showSnippet2)

	//FileServer to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//Register the File server with mux.Handle
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
