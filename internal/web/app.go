package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golangcollege/sessions"

	"github.com/theandrew168/jamql/internal/config"
	"github.com/theandrew168/jamql/internal/core"
)

//go:embed templates
var templatesFS embed.FS

type Application struct {
	cfg config.Config

	templates fs.FS
	storage   core.Storage
	session   *sessions.Session
	logger    *log.Logger
}

func NewApplication(cfg config.Config, storage core.Storage, session *sessions.Session, logger *log.Logger) *Application {
	var templates fs.FS
	if strings.HasPrefix(os.Getenv("ENV"), "dev") {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		templates = os.DirFS("./internal/web/templates/")
	} else {
		// else use the embedded templates dir
		templates, _ = fs.Sub(templatesFS, "templates")
	}

	app := Application{
		cfg: cfg,

		templates: templates,
		storage:   storage,
		session:   session,
		logger:    logger,
	}

	// use the app's error handler for session errors
	session.ErrorHandler = app.serverErrorResponse

	return &app
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(app.session.Enable)
	r.Use(middleware.Recoverer)

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Get("/", app.handleIndex)
	r.Get("/login", app.handleLogin)
	r.Get("/callback", app.handleCallback)

	// TODO: require tok cookie, else redir to /login
	r.Get("/jamql", app.handleJamQL)
	r.Post("/search", app.handleSearch)
	r.Post("/save", app.handleSave)

	return r
}
