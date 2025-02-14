package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	logger.Println("Hit the broker")

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		logger.Error("Error reading JSON: ", err)
		app.errorJSON(w, err)
		return
	}

	logger.WithContext(r.Context()).Printf("Received request with action: %s\n", requestPayload.Action)

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, r, requestPayload.Auth)
	//case "log":
	//	app.logItem(w, requestPayload.Log)
	default:
		app.errorJSON(w, errors.New("unknown action"))

	}
}

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request, a AuthPayload) {

	span, _ := opentracing.StartSpanFromContext(context.Background(), "authenticate")

	span.SetTag("service", "broker-service").
		SetTag("app", "example").
		SetTag("environment", "development")

	defer span.Finish()

	jsonData, _ := json.MarshalIndent(a, "", "\t")
	request, err := http.NewRequest("POST", "http://auth:80/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Inject the span's context into the HTTP headers
	err = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header),
	)
	if err != nil {
		logger.WithContext(r.Context()).Error("Error injecting span context: ", err)
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.WithContext(r.Context()).Error("Error calling auth service: ", err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		logger.WithContext(r.Context()).Error("Unauthorized")
		app.errorJSON(w, errors.New(strconv.Itoa(response.StatusCode)))
		return
	} else if response.StatusCode != http.StatusAccepted {
		logger.WithContext(r.Context()).Error("Error calling auth service: ", response.StatusCode)
		app.errorJSON(w, errors.New(strconv.Itoa(response.StatusCode)))
		return
	}

	// Create a variable we'll read response.Body into

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		logger.WithContext(r.Context()).Error("Error decoding response from auth service: ", err)
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		logger.WithContext(r.Context()).Error("Error from auth service: ", jsonFromService.Message)
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	logger.WithContext(r.Context()).Info("User authenticated: ", jsonFromService.Data)
	app.writeJSON(w, http.StatusAccepted, payload)

}
