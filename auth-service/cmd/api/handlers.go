package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	tracer := opentracing.GlobalTracer()

	// Extract the span's context from the HTTP headers.
	spanContext, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)

	// Start a new span with the extracted span context as its parent.
	span := tracer.StartSpan(
		"auth",
		opentracing.ChildOf(spanContext),
	).SetTag("service", "auth-service").
		SetTag("app", "example").
		SetTag("environment", "development")

	defer span.Finish()

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusConflict)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		logger.WithContext(r.Context()).Error("Error getting user by email: ", err)
		app.errorJSON(w, errors.New("invalid credentials"), 401)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil && !valid {
		logger.WithContext(r.Context()).
			WithFields(logrus.Fields{
				"spanID":  span.Context().(jaeger.SpanContext).SpanID().String(),
				"traceID": span.Context().(jaeger.SpanContext).TraceID().String(),
			}).
			Error("Error validating password: ", err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusInternalServerError)
		return
	}

	if err != nil {
		logger.WithContext(r.Context()).Error("Error validating password: ", err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	logger.WithContext(r.Context()).
		WithFields(logrus.Fields{
			"spanID":  span.Context().(jaeger.SpanContext).SpanID().String(),
			"traceID": span.Context().(jaeger.SpanContext).TraceID().String(),
		}).
		Info("User logged in: ", user.Email)

	app.writeJSON(w, http.StatusAccepted, payload)

}
