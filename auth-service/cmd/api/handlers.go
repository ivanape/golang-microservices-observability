package main

import (
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"log"
	"net/http"
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
	)
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
		log.Println("Error happend: ", err)
		app.errorJSON(w, errors.New("invalid credentials"), 401)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusInternalServerError)
		return
	}

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
