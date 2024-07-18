package probe

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type Probe struct {
	Build string
}

func NewProbe(build string) *Probe {
	return &Probe{Build: build}
}

type LivenessResponse struct {
}

func (p *Probe) LivenessHandler(c *gin.Context) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}
	c.JSON(200, gin.H{
		"hostname":   host,
		"build":      p.Build,
		"status":     "up",
		"name":       utils.GetEnv("KUBERNETES_NAME", "dev"),
		"pod_ip":     utils.GetEnv("KUBERNETES_POD_IP", "localhost"),
		"node":       utils.GetEnv("KUBERNETES_NODE_NAME", "dev"),
		"namespace":  utils.GetEnv("KUBERNETES_NAMESPACE", "dev"),
		"GOMAXPROCS": utils.GetEnv("GOMAXPROCS", "dev"),
	})

}
