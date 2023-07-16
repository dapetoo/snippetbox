package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dapetoo/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s, - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Defer func which will always run in the event of a panic
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// SessionLoad loads and saves session data for current request
func (app *application) SessionLoad(next http.Handler) http.Handler {
	return app.session.LoadAndSave(next)
}

// Ensure unauthorized user does not have access to create snippet
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//If the user is not authenticated, redirect to login page
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		//Otherwise set the "Cache-Control: no-store header so that pages require authentication arent s
		//stored in the user's browser caches
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
	//If user s

}

// NoSurf Middleware to prevent a CSRF attack
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check if a authenticatedUserID value exists in the session
		exists := app.session.Exists(r.Context(), "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		//Fetch the details of the current user from the database. If no matching record is found, or the current user account is
		//deactivated, remove the (invalid) authenticatedUserID value from the session and call the next handler
		// in the chain as normal
		user, err := app.users.Get(app.session.GetInt(r.Context(), "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r.Context(), "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
