package main

import (
	"log"
	"os"

	"github.com/fermin/gophercraft/internal/application/command"
	"github.com/fermin/gophercraft/internal/application/query"
	infrahandler "github.com/fermin/gophercraft/internal/infrastructure/handler"
	"github.com/fermin/gophercraft/internal/infrastructure/repository"
)

func main() {
	repo := repository.NewMemoryDummyRepository()
	_ = command.NewCreateDummyHandler(repo)
	_ = query.NewGetDummyHandler(repo)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	s := infrahandler.NewServer()
	s.RegisterRoutes()
	if err := s.Run(addr); err != nil {
		log.Fatalf("http server: %v", err)
	}
}
