package main

import (
	"broker/obs"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

const webPort = "80"

type Config struct{}

var logger *logrus.Logger

func main() {
	app := Config{}
	obs.InitTracer(obs.DefaultServiceTags["service"])

	logger = logrus.New()
	// Configure the Loki hook
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithFormatter(&logrus.JSONFormatter{}).
		WithStaticLabels(obs.DefaultServiceTags)

	lokiWebHookUrl := os.Getenv("LOKI_WEBHOOK_URL")

	hook := lokirus.NewLokiHookWithOpts(
		lokiWebHookUrl,
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel)

	logger.Hooks.Add(hook)

	logger.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
