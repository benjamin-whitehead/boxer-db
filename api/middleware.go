package api

import (
	"net/http"

	"github.com/benjamin-whitehead/boxer-db/m/v2/config"
	"github.com/gin-gonic/gin"
)

func checkRoleMiddleware() gin.HandlerFunc {
	role := config.GetConfig().Role

	return func(c *gin.Context) {
		if role != config.ROLE_LEADER {
			c.AbortWithStatusJSON(http.StatusForbidden, InvalidRequestToFollowerResponse{
				Message: "This action is only available to the leader",
			})
		}
		c.Next()
	}
}
