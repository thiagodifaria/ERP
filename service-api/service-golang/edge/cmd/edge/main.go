// Este arquivo inicia o servico edge e delega o bootstrap para a camada correta.
// Regra de negocio nao deve nascer aqui.
package main

import (
  "log"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/bootstrap"
)

func main() {
  app, err := bootstrap.NewApp()
  if err != nil {
    log.Fatalf("bootstrap error: %v", err)
  }

  if err := app.Run(); err != nil {
    app.Logger.Printf("server stopped with error: %v", err)
  }
}
