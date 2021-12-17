package web

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/theandrew168/jamql/internal/core"
)

var (
	filterLimit = 1
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
		{0, false},
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleSearch(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// convert form data into filters
	var filters []core.Filter
	for i := 0; i < filterLimit; i++ {
		key := fmt.Sprintf("filter-key-%d", i)
		op := fmt.Sprintf("filter-op-%d", i)
		value := fmt.Sprintf("filter-value-%d", i)

		filter := core.Filter{
			Key:   r.PostFormValue(key),
			Op:    r.PostFormValue(op),
			Value: r.PostFormValue(value),
		}
		filters = append(filters, filter)
	}

	// handle empty filters
	empty := true
	for _, filter := range filters {
		if filter.Value != "" {
			empty = false
		}
	}

	if empty {
		err = app.flashError(w, "Try filling in some filters first!")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	// search for matching tracks
	tracks, err := app.storage.SearchTracks(r, filters)
	if err != nil {
		// redirect to /login if unauthorized (token expired)
		if errors.Is(err, core.ErrUnauthorized) {
			http.Redirect(w, r, "/login", 303)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	ts, err := template.ParseFS(app.templates, "track.partial.tmpl")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// render tracks to a temp buffer
	var buf bytes.Buffer
	for _, track := range tracks {
		err = ts.Execute(&buf, track)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// write all tracks at once
	w.Write(buf.Bytes())
}

func (app *Application) handleSave(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// convert form data into filters
	var filters []core.Filter
	for i := 0; i < filterLimit; i++ {
		key := fmt.Sprintf("filter-key-%d", i)
		op := fmt.Sprintf("filter-op-%d", i)
		value := fmt.Sprintf("filter-value-%d", i)

		filter := core.Filter{
			Key:   r.PostFormValue(key),
			Op:    r.PostFormValue(op),
			Value: r.PostFormValue(value),
		}
		filters = append(filters, filter)
	}

	// search for matching tracks
	tracks, err := app.storage.SearchTracks(r, filters)
	if err != nil {
		// redirect to /login if unauthorized (token expired)
		if errors.Is(err, core.ErrUnauthorized) {
			http.Redirect(w, r, "/login", 303)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// save tracks to a new playlist
	err = app.storage.SaveTracks(r, tracks, "TODO Title", "TODO Description")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.flashSuccess(w, "Playlist created!")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) flashSuccess(w http.ResponseWriter, message string) error {
	ts, err := template.ParseFS(app.templates, "flash-success.partial.tmpl")
	if err != nil {
		return err
	}

	err = ts.Execute(w, message)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) flashError(w http.ResponseWriter, message string) error {
	ts, err := template.ParseFS(app.templates, "flash-error.partial.tmpl")
	if err != nil {
		return err
	}

	err = ts.Execute(w, message)
	if err != nil {
		return err
	}

	return nil
}
