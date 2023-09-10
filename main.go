package main

import (
	"chia-goths/internal"
	"chia-goths/internal/apps"
	"chia-goths/internal/apps/about"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "embed"
)

func main() {
	internal.LoadEnv()

	configLogger()

	c := chi.NewRouter()
	c.Use(middleware.Logger)
	c.Use(middleware.Recoverer)
	c.Use(middleware.Compress(5))

	main := apps.NewAppConfig(c, "/apps/about")
	main.InitApp(&about.App{})

	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/apps/about", http.StatusFound)
	})

	log.Info().Str("listenAddr", internal.EnvVars.ListenAddr).Msg("starting server")
	if err := http.ListenAndServe(internal.EnvVars.ListenAddr, csrf.Protect(internal.EnvVars.CSRFKey)(c)); err != nil {
		panic(fmt.Errorf("failed to listen and serve: %w", err))
	}
}

func configLogger() {
	if internal.EnvVars.DevMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("dev mode enabled")
	}

	// set chi middleware logger to zerolog
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger:  &log.Logger,
			NoColor: !internal.EnvVars.DevMode,
		})
}
