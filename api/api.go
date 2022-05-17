package api

import "github.com/gin-gonic/gin"

func InitializeAPIRoutes(router *gin.Engine) {

	// TODO: eventually try and get this to work with route groups
	router.GET("/api/v1/:key", GetKey)
	router.PUT("/api/v1/:key", checkRoleMiddleware(), PutKey)
	router.POST("/api/v1/:key", checkRoleMiddleware(), PutKey)
	router.DELETE("api/v1/:key", checkRoleMiddleware(), DeleteKey)

	router.POST("/api/v1/replication", checkRoleForReplicationMiddleware(), ReplicateLog)
	router.PUT("/api/v1/replication", checkRoleForReplicationMiddleware(), ReplicateLog)

	router.GET("/api/v1/log", GetLog)
	router.GET("/api/v1/role", GetRole)

}
