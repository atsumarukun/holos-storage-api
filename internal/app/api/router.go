package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	health := r.Group("health")
	health.GET("/", healthHandler.Health)
}
