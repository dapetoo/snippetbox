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

	//Initialize Middleware
	mux.Use(app.SessionLoad)

	//Routes
	mux.Get("/", app.home)
	mux.Get("/snippet/create", app.createSnippetForm)
	mux.Get("/snippet/{id}", app.showSnippet2)
	mux.Post("/snippet/create", app.createSnippet)
	mux.Get("/user/signup", app.signupUserForm)
	mux.Post("/user/signup", app.signupUser)
	mux.Get("/user/login", app.loginUserForm)
	mux.Post("/user/login", app.loginUser)
	mux.Post("/user/logout", app.logoutUser)

	//FileServer to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//Register the File server with mux.Handle
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
