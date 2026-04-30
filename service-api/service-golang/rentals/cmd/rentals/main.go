package main

import (
	"log"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("rentals bootstrap error: %v", err)
	}
	defer func() {
		if app.Database != nil {
			_ = app.Database.Close()
		}
	}()

	if err := app.Run(); err != nil && err != http.ErrServerClosed {
		app.Logger.Printf("rentals runtime stopped with error: %v", err)
	}
}
