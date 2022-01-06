package web

import (
	"net/http"
)

// Based on:
// Let's Go - Chapter 11.6
func (app *Application) requireToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.session.Exists(r, "token") {
			http.Redirect(w, r, "/login", 302)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next(w, r)
	})
}
