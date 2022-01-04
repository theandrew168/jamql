package web

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/theandrew168/jamql/internal/core"
)

var (
	filterCount = 3
	stateCookieName = "state"
	tokenCookieName = "token"
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

// redirect user to spotify authorize w/ ID, scope, etc
func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	// simulate login when cfg.ClientID is unset
	if app.cfg.ClientID == "" {
		token, err := generateRandomString(32)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		app.session.Put(r, tokenCookieName, token)
		http.Redirect(w, r, "/jamql", 302)
		return
	}

	// generate state first since it can (rarely) cause errors
	state, err := generateRandomString(16)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// redirect user to Spotify's login service
	values := url.Values{}
	values.Set("response_type", "token")
	values.Set("client_id", app.cfg.ClientID)
	values.Set("scope", "playlist-modify-public")
	values.Set("redirect_uri", app.cfg.RedirectURI)
	values.Set("state", state)

	authURL, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	authURL.RawQuery = values.Encode()
	app.logger.Println(authURL.String())
}

// stores access_token in a cookie (URL param)
func (app *Application) handleCallback(w http.ResponseWriter, r *http.Request) {

}

// uses regular error responses (user isn't at the main app yet)
func (app *Application) handleJamQL(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"jamql.page.tmpl",
		"base.layout.tmpl",
	}

	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		FilterCount int
	}{
		filterCount,
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

	// handle no matching tracks
	if len(tracks) == 0 {
		app.clientErrorFlash(w, r, "No tracks match these filters!")
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

	// handle no matching tracks
	if len(tracks) == 0 {
		app.clientErrorFlash(w, r, "No tracks match these filters!")
		return
	}

	// join up all the filter values...
	values := []string{}
	for _, filter := range filters {
		values = append(values, strings.Title(filter.Value))
	}
	// to create a rough summary for the playlist name
	summary := strings.Join(values, " + ")

	// save tracks to a new playlist
	name := "JamQL Mix: " + summary
	desc := "A fresh mix generated by JamQL"
	err = app.storage.SaveTracks(r, tracks, name, desc)
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
			Op:	op,
			Value: value,
		}
		filters = append(filters, filter)
	}

	return filters
}

func generateRandomString(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
