package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	health := r.Group("health")
	health.GET("/", healthHdl.Health)

	r.Use(authorizationMW.Authorize)

	volumes := r.Group("volumes")
	volumes.POST("/", volumeHdl.Create)
	volumes.GET("/", volumeHdl.GetAll)
	volumes.PUT("/:id", volumeHdl.Update)
	volumes.DELETE("/:id", volumeHdl.Delete)
	volumes.GET("/:id", volumeHdl.GetOne)

	entries := r.Group("entries")
	entries.POST("/", entryHdl.Create)
	entries.PUT("/:volumeName/*key", entryHdl.Update)
	entries.DELETE("/:volumeName/*key", entryHdl.Delete)
	entries.HEAD("/:volumeName/*key", entryHdl.Head)
	entries.GET("/:volumeName/*key", entryHdl.Get)
}
