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

	// TODO: Refactor this setup to functions

	configuration := config.GetConfig()
	log.Println("Initializing log...")
	replication.GetLog()

	// InitializeLog needs to be called before InitializeStore
	db.InitializeStore()

	if replication.GetLog().HadData {
		log.Println("Restoring log...")
		// Had data, load the data into the store
		for _, entry := range replication.GetLog().Entries {
			if entry.CommandType == replication.COMMAND_TYPE_WRITE {
				db.GlobalStore.Put(entry.EntryKey, entry.EntryValue)
			}
			if entry.CommandType == replication.COMMAND_TYPE_DELETE {
				db.GlobalStore.Delete(entry.EntryKey)
			}
		}
	}

	router := gin.Default()

	api.InitializeAPIRoutes(router)

	// Initialize the log

	// TODO: Refactor
	if configuration.Role == config.ROLE_LEADER {
		log.Println("Starting as leader")
		router.Run("0.0.0.0:8080")
	} else {
		log.Println("Starting as follower")
		router.Run("0.0.0.0:8080")
	}

}
