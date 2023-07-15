package main

import (
	"errors"
	"fmt"
	"github.com/dapetoo/snippetbox/pkg/models"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
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

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//Parsing Form Data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//Retrieve the relevant data fields from the r.Postform map
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	//Map to hold any validation error
	formErrors := make(map[string]string)
	//Check that the fields aren't empty
	if strings.TrimSpace(title) == "" {
		formErrors["title"] = "This field cannot be empty"
	} else if utf8.RuneCountInString(title) > 100 {
		formErrors["title"] = "This field is too long, maximum is 100 characters"
	}
	//Validate the content field
	if strings.TrimSpace(content) == "" {
		formErrors["content"] = "This field cannot be empty"
	}
	//Check that the expires field is not blank
	if strings.TrimSpace(expires) == "" {
		formErrors["content"] = "This field cannot be empty"
	} else if expires != "365" && expires != "7" && expires != "1" {
		formErrors["expires"] = "The field is invalid"
	}

	if len(formErrors) > 0 {
		fmt.Sprintf()
		return
	}

	//Pass the data to the SnippetModel.Insert() method
	id, err := app.snippets.Insert(title, content, expires)
	log.Println(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
