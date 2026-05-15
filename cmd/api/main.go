// Command api is the service entrypoint.
package main

import (
	"log"
	"os"

	"github.com/fermin/gophercraft/internal/application/command"
	"github.com/fermin/gophercraft/internal/application/query"
	"github.com/fermin/gophercraft/internal/infrastructure/clock"
	infraevent "github.com/fermin/gophercraft/internal/infrastructure/event"
	infrahandler "github.com/fermin/gophercraft/internal/infrastructure/handler"
	"github.com/fermin/gophercraft/internal/infrastructure/repository"
	"github.com/fermin/gophercraft/internal/infrastructure/uuid"
)

func main() {
	repo := repository.NewMemoryDummyRepository()
	_ = command.NewCreateDummyHandler(repo, uuid.GoogleUUIDGenerator{}, clock.SystemClock{}, infraevent.NoopPublisher{})
	_ = query.NewGetDummyHandler(repo)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	s, err := infrahandler.NewServer()
	if err != nil {
		log.Fatalf("http server init: %v", err)
	}
	s.RegisterRoutes()
	if err = s.Run(addr); err != nil {
		log.Fatalf("http server: %v", err)
	}
}
