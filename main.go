package main

import (
	"log"

	"github.com/benjamin-whitehead/boxer-db/m/v2/api"
	"github.com/benjamin-whitehead/boxer-db/m/v2/config"
	"github.com/benjamin-whitehead/boxer-db/m/v2/db"
	"github.com/benjamin-whitehead/boxer-db/m/v2/replication"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitializeStore()
	configuration := config.GetConfig()

	router := gin.Default()

	api.InitializeAPIRoutes(router)

	// Initialize the log

	if configuration.Role == config.ROLE_LEADER {
		// TODO: Not sure if this should be for follower as well as leader, we will see
		log.Println("Initializing log...")
		replication.InitializeLog()
	}

	// TODO: Refactor
	if configuration.Role == config.ROLE_LEADER {
		log.Println("Starting as leader")
		router.Run("0.0.0.0:8080")
	} else {
		log.Println("Starting as follower")
		router.Run("0.0.0.0:8080")
	}

}
