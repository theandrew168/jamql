package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/go-chi/chi/v5"
	"github.com/golangcollege/sessions"
	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/jamql/internal/config"
	"github.com/theandrew168/jamql/internal/core"
	"github.com/theandrew168/jamql/internal/spotify"
	"github.com/theandrew168/jamql/internal/test"
	"github.com/theandrew168/jamql/internal/web"
)

//go:embed static
var staticFS embed.FS

//go:embed static/img/logo.webp
var logo []byte

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)

	conf := flag.String("conf", "jamql.conf", "app config file")
	flag.Parse()

	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Fatalln(err)
	}

	// create and configure session storage
	session := sessions.New([]byte(cfg.SecretKey))
	session.Lifetime = 1 * time.Hour
	session.HttpOnly = true
	session.Secure = true

	// use test storage when cfg.ClientID is unset
	var storage core.Storage
	if cfg.ClientID == "" {
		storage = test.NewStorage()
	} else {
		storage = spotify.NewStorage(session)
	}

	app := web.NewApplication(cfg, storage, session, logger)

	// setup http.Handler for static files
	static, _ := fs.Sub(staticFS, "static")
	staticServer := http.FileServer(http.FS(static))
	gzipStaticServer := gzhttp.GzipHandler(staticServer)

	// construct the top-level router
	r := chi.NewRouter()
	r.Mount("/", app.Router())
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/static/*", http.StripPrefix("/static", gzipStaticServer))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write(logo)
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	addr := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalln(err)
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	logger.Printf("started server on %s\n", addr)

	// kick off a goroutine to listen for SIGINT and SIGTERM
	shutdownError := make(chan error)
	go func() {
		// idle until a signal is caught
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// give the web server 5 seconds to shutdown gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown the web server and track any errors
		logger.Println("stopping server")
		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	err = srv.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalln(err)
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("stopped server")
}
