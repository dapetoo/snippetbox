package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
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
	err := ts.Execute(buf, td)
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
