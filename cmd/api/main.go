// Command api is the service entrypoint.
package main

import (
	"os"

	"github.com/fermin/gophercraft/internal/application/command"
	"github.com/fermin/gophercraft/internal/application/query"
	"github.com/fermin/gophercraft/internal/infrastructure/clock"
	infraevent "github.com/fermin/gophercraft/internal/infrastructure/event"
	infrahandler "github.com/fermin/gophercraft/internal/infrastructure/handler"
	infralogger "github.com/fermin/gophercraft/internal/infrastructure/logger"
	inframetrics "github.com/fermin/gophercraft/internal/infrastructure/metrics"
	"github.com/fermin/gophercraft/internal/infrastructure/repository"
	"github.com/fermin/gophercraft/internal/infrastructure/uuid"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logPretty := os.Getenv("LOG_PRETTY") == "true" || os.Getenv("LOG_PRETTY") == "1"

	appLogger := infralogger.NewZerologLogger(logLevel, logPretty, infralogger.GlobalFieldsFromEnv())

	promMetrics := inframetrics.NewPrometheusRecorder()

	repo := repository.NewMemoryDummyRepository()
	publisher := infraevent.NewNoopPublisher(appLogger, promMetrics)
	createDummyHandler := command.NewCreateDummyHandler(repo, uuid.GoogleUUIDGenerator{}, clock.SystemClock{}, publisher)
	getDummyHandler := query.NewGetDummyHandler(repo)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	s, err := infrahandler.NewServer(appLogger, promMetrics, promMetrics)
	if err != nil {
		appLogger.Fatal("http server init", "error", err)
	}
	s.RegisterRoutes(
		infrahandler.NewCreateDummyGinHandler(createDummyHandler),
		infrahandler.NewGetDummyGinHandler(getDummyHandler),
	)
	if err = s.Run(addr); err != nil {
		appLogger.Fatal("http server", "error", err)
	}
}
