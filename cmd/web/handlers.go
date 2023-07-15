package main

import (
	"errors"
	"fmt"
	"github.com/dapetoo/snippetbox/pkg/forms"
	"github.com/dapetoo/snippetbox/pkg/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	//Use the Snippet object's Get method to retrieve the data for a specific record based on it's ID
	s, err := app.snippets.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) showSnippet2(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	snippetID, err := strconv.Atoi(id)
	if err != nil || snippetID < 1 {
		app.notFound(w)
		return
	}

	// Use the Snippet object's Get method to retrieve the data for a specific record based on its ID
	s, err := app.snippets.GetByID(snippetID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//Retrieve the value of the "flash" key
	flash := app.session.GetString(r.Context(), "flash")

	app.render(w, r, "show.page.tmpl", &templateData{
		Flash:   flash,
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//Parsing Form Data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Create a new forms.Form struct
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	//If the form isnt valid redisplay the template passing in the form.Form object as the data
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	//Pass the data to the SnippetModel.Insert() method
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Put() to add a string value(Flash message)
	app.session.Put(r.Context(), "flash", "Snippet successfully created")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
