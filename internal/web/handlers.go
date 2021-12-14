package web

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/theandrew168/jamql/internal/test"
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

func (app *Application) handleSearch(w http.ResponseWriter, r *http.Request) {
	// TODO: read form data
	// TODO: search for matching tracks

	ts, err := template.ParseFS(app.templates, "track.partial.tmpl")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// render tracks to a temp buffer
	var buf bytes.Buffer
	for _, track := range test.SampleTracks {
		err = ts.Execute(&buf, track)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// write all tracks at once
	w.Write(buf.Bytes())
}
