package web

import (
	"bytes"
	"errors"
//	"fmt"
	"html/template"
	"net/http"

	"github.com/theandrew168/jamql/internal/core"
)

var (
	filterLimit = 3
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

	data := []bool{
		false,
		true,
		true,
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// uses flash messages for reporting errors
func (app *Application) handleSearch(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		app.serverErrorFlash(w, r, err)
		return
	}

	keys := r.PostForm["filter-key"]
	ops := r.PostForm["filter-op"]
	values1 := r.PostForm["filter-value-1"]
	values2 := r.PostForm["filter-value-2"]

	app.logger.Println(len(keys))
	app.logger.Println(keys)
	app.logger.Println(len(ops))
	app.logger.Println(ops)
	app.logger.Println(len(values1))
	app.logger.Println(values1)
	app.logger.Println(len(values2))
	app.logger.Println(values2)

	// TODO: rebuild filters from form data, add a helper func

	// convert form data into filters
	var filters []core.Filter
	for i := 0; i < filterLimit; i++ {
//		key := fmt.Sprintf("filter-key-%d", i)
//		op := fmt.Sprintf("filter-op-%d", i)
//		value := fmt.Sprintf("filter-value-%d", i)
//
//		// skip filters with empty values
//		if r.PostFormValue(value) == "" {
//			continue
//		}
//
//		filter := core.Filter{
//			Key:   r.PostFormValue(key),
//			Op:    r.PostFormValue(op),
//			Value: r.PostFormValue(value),
//		}
//		filters = append(filters, filter)
	}

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

	// convert form data into filters
	var filters []core.Filter
//	for i := 0; i < filterLimit; i++ {
//		key := fmt.Sprintf("filter-key-%d", i)
//		op := fmt.Sprintf("filter-op-%d", i)
//		value := fmt.Sprintf("filter-value-%d", i)
//
//		filter := core.Filter{
//			Key:   r.PostFormValue(key),
//			Op:    r.PostFormValue(op),
//			Value: r.PostFormValue(value),
//		}
//		filters = append(filters, filter)
//	}

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
