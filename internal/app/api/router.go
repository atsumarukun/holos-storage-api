package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	health := r.Group("health")
	health.GET("", healthHdl.Health)

	r.Use(authorizationMW.Authorize)

	volumes := r.Group("volumes")
	volumes.POST("", volumeHdl.Create)
	volumes.GET("", volumeHdl.GetAll)
	volumes.PUT("/:name", volumeHdl.Update)
	volumes.DELETE("/:name", volumeHdl.Delete)
	volumes.GET("/:name", volumeHdl.GetOne)

	entries := r.Group("entries")
	entries.POST("", entryHdl.Create)
	entries.GET("/:volumeName", entryHdl.Search)
	entries.POST("/:volumeName/*key", entryHdl.Copy)
	entries.PUT("/:volumeName/*key", entryHdl.Update)
	entries.DELETE("/:volumeName/*key", entryHdl.Delete)
	entries.HEAD("/:volumeName/*key", entryHdl.GetMeta)
	entries.GET("/:volumeName/*key", entryHdl.GetOne)
}
