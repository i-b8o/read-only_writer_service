package main

import (
	"context"
	"log"
	"regulations_writable_service/internal/app"
	"regulations_writable_service/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetConfig()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
