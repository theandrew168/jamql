package web

import (
	"html/template"
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"index.page.tmpl",
		"base.layout.tmpl",
	}

	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleJamQL(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"jamql.page.tmpl",
		"base.layout.tmpl",
		"filter.partial.tmpl",
	}

	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := []struct {
		ID     int
		Hidden bool
	}{
		{1, false},
		{2, true},
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
