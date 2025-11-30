package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	healthCheckPath = "/health"
)

type HealthCheck struct{}

func (h *HealthCheck) AddRoutes(r *gin.Engine) {
	r.GET(healthCheckPath, h.handler)
}

func (h *HealthCheck) handler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
