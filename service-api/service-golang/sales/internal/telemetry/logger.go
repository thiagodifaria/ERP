// Logger concentra o ponto minimo de observabilidade local.
// Traces e metricas entram depois sem poluir o bootstrap.
package telemetry

import (
	"log"
	"os"
)

type Logger = log.Logger

func New(serviceName string) *Logger {
	return log.New(os.Stdout, serviceName+" ", log.Ldate|log.Ltime|log.LUTC)
}
