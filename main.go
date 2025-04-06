package main

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/app"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"log"
)

func main() {
	repository.FillSlices()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newApp, err := app.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = newApp.Start()
	if err != nil {
		log.Fatal(err)
	}
}
