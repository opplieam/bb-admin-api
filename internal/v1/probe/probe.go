package probe

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type Store interface {
	HealthCheck() (bool, error)
}

type Handler struct {
	Build string
	Store Store
}

func NewHandler(build string, store Store) *Handler {
	return &Handler{
		Build: build,
		Store: store,
	}
}

func (h *Handler) LivenessHandler(c *gin.Context) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}
	c.JSON(200, gin.H{
		"hostname":   host,
		"build":      h.Build,
		"status":     "up",
		"name":       utils.GetEnv("KUBERNETES_NAME", "dev"),
		"pod_ip":     utils.GetEnv("KUBERNETES_POD_IP", "localhost"),
		"node":       utils.GetEnv("KUBERNETES_NODE_NAME", "dev"),
		"namespace":  utils.GetEnv("KUBERNETES_NAMESPACE", "dev"),
		"GOMAXPROCS": utils.GetEnv("GOMAXPROCS", "dev"),
	})

}

func (h *Handler) ReadinessHandler(c *gin.Context) {
	healthy, err := h.Store.HealthCheck()
	if err != nil || !healthy {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{
		"status": "up",
	})
}
