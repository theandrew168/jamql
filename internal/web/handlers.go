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
	filterCount = 3
)

// uses regular error responses (user isn't at the main app yet)
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

// uses regular error responses (user isn't at the main app yet)
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

	// TODO: loop range filterCount
	data := []struct {
		ID     int
		Hidden bool
	}{
		{0, false},
		{1, true},
		{2, true},
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// be sure to call r.ParseForm() before using this helper
func parseFiltersForm(r *http.Request) []core.Filter {
	var filters []core.Filter
	for i := 0; i < filterCount; i++ {
		keyName := fmt.Sprintf("filter-key-%d", i)
		opName := fmt.Sprintf("filter-op-%d", i)
		value1Name := fmt.Sprintf("filter-value1-%d", i)
		value2Name := fmt.Sprintf("filter-value2-%d", i)

		key := r.PostFormValue(keyName)
		op := r.PostFormValue(opName)
		value1 := r.PostFormValue(value1Name)
		value2 := r.PostFormValue(value2Name)

		// flip "year contains" to "year between"
		if key == "year" && op == "contains" {
			op = "between"
		}

		// ignore filters with missing fields
		if value1 == "" {
			continue
		}
		if key == "year" && op == "between" && value2 == "" {
			continue
		}

		// rebuild "year between" value if necessary
		var value string
		if key == "year" && op == "between" {
			value = value1 + "-" + value2
		} else {
			value = value1
		}

		filter := core.Filter{
			Key:   key,
			Op:    op,
			Value: value,
		}
		filters = append(filters, filter)
	}

	return filters
}

// uses flash messages for reporting errors
func (app *Application) handleSearch(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		app.serverErrorFlash(w, r, err)
		return
	}

	filters := parseFiltersForm(r)

	// handle empty filters
	if len(filters) == 0 {
		app.clientErrorFlash(w, r, "Try filling in some filters first!")
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
		app.serverErrorFlash(w, r, err)
		return
	}

	ts, err := template.ParseFS(app.templates, "track.partial.tmpl")
	if err != nil {
		app.serverErrorFlash(w, r, err)
		return
	}

	// render tracks to a temp buffer
	var buf bytes.Buffer
	for _, track := range tracks {
		err = ts.Execute(&buf, track)
		if err != nil {
			app.serverErrorFlash(w, r, err)
			return
		}
	}

	// write all tracks at once
	w.Write(buf.Bytes())
}

// uses flash messages for reporting errors
func (app *Application) handleSave(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		app.serverErrorFlash(w, r, err)
		return
	}

	filters := parseFiltersForm(r)

	// search for matching tracks
	tracks, err := app.storage.SearchTracks(r, filters)
	if err != nil {
		// redirect to /login if unauthorized (token expired)
		if errors.Is(err, core.ErrUnauthorized) {
			http.Redirect(w, r, "/login", 303)
			return
		}
		app.serverErrorFlash(w, r, err)
		return
	}

	// save tracks to a new playlist
	err = app.storage.SaveTracks(r, tracks, "TODO Title", "TODO Description")
	if err != nil {
		app.serverErrorFlash(w, r, err)
		return
	}

	ts, err := template.ParseFS(app.templates, "flash-success.partial.tmpl")
	if err != nil {
		app.serverErrorFlash(w, r, err)
		return
	}

	message := "Playlist created!"
	err = ts.Execute(w, message)
	if err != nil {
		app.serverErrorFlash(w, r, err)
		return
	}
}
