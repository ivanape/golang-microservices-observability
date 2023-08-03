package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"strconv"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	// The most interesting part

	// initialize tracer
	tracer := otel.Tracer("broker-service")

	// get ctx, and span from tracer by starting it
	ctx, span := tracer.Start(r.Context(), "BrokerHandler")
	defer span.End()

	// creating request to send ctx to auth service, for auth service to catch context
	ctxReq, _ := http.NewRequest("GET", "http://localhost:8080/auth", nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(ctxReq.Header))

	// Send ctx to auth service
	client := &http.Client{}
	client.Do(ctxReq)

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))

	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	tracer := otel.Tracer("broker-tracer")
	ctx, span := tracer.Start(context.Background(), "operation-a")
	defer span.End()
	header := make(http.Header)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(header))

	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://auth:80/authenticate", bytes.NewBuffer(jsonData))

	request.Header = header

	ctx, requestSpan := tracer.Start(request.Context(), "request handling")
	defer requestSpan.End()

	// Now start a span for the DB operation
	_, dbSpan := tracer.Start(ctx, "DB operation")
	defer dbSpan.End()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {

		app.errorJSON(w, err)
		return

	}
	defer response.Body.Close()

	// make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New(strconv.Itoa(response.StatusCode)))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New(strconv.Itoa(response.StatusCode)))
		return
	}

	// Create a variable we'll read response.Body into

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}
