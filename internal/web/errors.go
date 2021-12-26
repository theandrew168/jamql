package web

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
)

var (
	ErrUnauthorized = errors.New("core: unauthorized")
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, code int, tmpl string) {
	files := []string{
		tmpl,
		"base.layout.tmpl",
	}

	// attempt to parse error template
	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	// render template to a temp buffer
	var buf bytes.Buffer
	err = ts.Execute(&buf, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	// write the status and error message
	w.WriteHeader(code)
	w.Write(buf.Bytes())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 404, "404.page.tmpl")
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 405, "405.page.tmpl")
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// skip 2 frames to identify original caller
	app.logger.Output(2, err.Error())
	app.errorResponse(w, r, 500, "500.page.tmpl")
}

func (app *Application) errorFlash(w http.ResponseWriter, r *http.Request, code int, message string) {
	ts, err := template.ParseFS(app.templates, "flash-error.partial.tmpl")
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	// render template to a temp buffer
	var buf bytes.Buffer
	err = ts.Execute(&buf, message)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	// write the status and error message
	w.WriteHeader(code)
	w.Write(buf.Bytes())
}

func (app *Application) clientErrorFlash(w http.ResponseWriter, r *http.Request, message string) {
	app.errorFlash(w, r, 400, message)
}

func (app *Application) serverErrorFlash(w http.ResponseWriter, r *http.Request, err error) {
	// skip 2 frames to identify original caller
	app.logger.Output(2, err.Error())
	app.errorFlash(w, r, 500, "Internal server error")
}
