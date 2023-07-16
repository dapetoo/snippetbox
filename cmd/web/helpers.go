package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Server Error
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	err = app.errorLog.Output(2, trace)
	if err != nil {
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Client Error
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Not found Error
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	//Flash Message to the templateData if one exists
	td.Flash = app.session.PopString(r.Context(), "flash")
	//Add the Authentication status to the template data
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	//Retrieve the appropriate template from the cache
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s doesn't exist", name))
		return
	}

	//Initialize a new buffer
	buf := new(bytes.Buffer)

	//Execute the template set passing in any dynamic data
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Write the content of the buffer to the http.ResponseWriter
	_, err = buf.WriteTo(w)
	if err != nil {
		return
	}
}

// Return true if the current user is from authenticated user, otherwise return false
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r.Context(), "authenticatedUserID")
}
