package main

import (
	"fmt"
	"log"

	"github.com/fermin/gophercraft/internal/application/command"
	"github.com/fermin/gophercraft/internal/application/query"
	infrahandler "github.com/fermin/gophercraft/internal/infrastructure/handler"
	"github.com/fermin/gophercraft/internal/infrastructure/repository"
)

func main() {
	repo := repository.NewMemoryExampleRepository()
	_ = command.NewCreateExampleHandler(repo)
	_ = query.NewGetExampleHandler(repo)
	_ = infrahandler.NewServer()

	fmt.Println("gophercraft api — startup")
	log.Println("wire HTTP listener and routes in a follow-up")
}
