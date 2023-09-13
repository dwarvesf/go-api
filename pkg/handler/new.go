package handler

import (
	"log"
	"net/http"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/logger/monitor"
	"github.com/gin-gonic/gin"
)

// Handler for app
type Handler struct {
	log     *log.Logger
	cfg     config.Config
	monitor monitor.Tracer
}

// New will return an instance of Auth struct
func New(cfg config.Config, monitor monitor.Tracer) *Handler {

	return &Handler{
		log:     log.Default(),
		cfg:     cfg,
		monitor: monitor,
	}
}

// Healthz handler
// Return "OK"
func (h *Handler) Healthz(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("OK"))

}
