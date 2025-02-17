package main

import (
	"authentication/obs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
	otelchimetric "github.com/riandyrn/otelchi/metric"
)

func (app *Config) routes(metricConfig otelchimetric.BaseConfig) http.Handler {
	r := chi.NewRouter()

	// specify who is allowed to connect
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}),
		otelchi.Middleware(obs.DefaultServiceTags["service"], otelchi.WithChiRoutes(r)),
		otelchimetric.NewRequestDurationMillis(metricConfig),
		otelchimetric.NewRequestInFlight(metricConfig),
		otelchimetric.NewResponseSizeBytes(metricConfig),
	)

	r.Use(middleware.Heartbeat("/ping"))
	r.Post("/authenticate", app.Authenticate)
	return r
}
