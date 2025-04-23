package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	health := r.Group("health")
	health.GET("/", healthHdl.Health)

	r.Use(authorizationMW.Authorize)

	volumes := r.Group("volumes")
	volumes.POST("/", volumeHdl.Create)
}
