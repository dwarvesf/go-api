package handler

import (
	"log"
	"net/http"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/gin-gonic/gin"
)

// Handler for app
type Handler struct {
	log *log.Logger
	cfg config.Config
}

// New will return an instance of Auth struct
func New(cfg config.Config) *Handler {

	return &Handler{
		log: log.Default(),
		cfg: cfg,
	}
}

// Healthz handler
// Return "OK"
func (h *Handler) Healthz(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("OK"))

}
