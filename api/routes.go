package api

import (
	"log"
	"net/http"
	"time"

	"github.com/benjamin-whitehead/boxer-db/m/v2/config"
	"github.com/benjamin-whitehead/boxer-db/m/v2/db"
	"github.com/benjamin-whitehead/boxer-db/m/v2/replication"
	"github.com/gin-gonic/gin"
)

// GetKey returns a value associated with a key
// Example: GET https://localhost:8080/api/v1/hello
func GetKey(c *gin.Context) {
	key := c.Param("key")
	boxerKey := db.BoxerKey{Key: key}
	value, err := db.GlobalStore.Get(db.BoxerKey{Key: key})
	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		replication.GetLog().AppendLog(boxerKey, db.BoxerValue{}, replication.COMMAND_TYPE_READ)
		c.JSON(http.StatusOK, value)
	}
}

// PutKey sets a value associated with a key
// Example: PUT https://localhost:8080/api/v1/hello
// {
// 	"value": "world"
// }
func PutKey(c *gin.Context) {
	var request Request
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	key := c.Param("key")
	value := request.Value

	boxerKey := db.BoxerKey{Key: key}
	boxerValue := db.BoxerValue{Value: value, Meta: db.BoxerValueMetadata{Timestamp: time.Now().UnixNano()}}

	db.GlobalStore.Put(boxerKey, boxerValue)
	replication.GetLog().AppendLog(boxerKey, boxerValue, replication.COMMAND_TYPE_WRITE)

	c.Status(http.StatusOK)
}

// DeleteKey deletes a value associated with a key
// Example: DELETE https://localhost:8080/api/v1/hello
func DeleteKey(c *gin.Context) {
	key := c.Param("key")
	boxerKey := db.BoxerKey{Key: key}
	err := db.GlobalStore.Delete(db.BoxerKey{Key: key})

	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		replication.GetLog().AppendLog(boxerKey, db.BoxerValue{}, replication.COMMAND_TYPE_DELETE)
		c.Status(http.StatusOK)
	}
}

// GetRole return's the role of the node
func GetRole(c *gin.Context) {
	role := config.GetConfig().Role
	c.JSON(http.StatusOK, RoleResponse{Role: role})
}

func ReplicateLog(c *gin.Context) {

	log.Println("HERE!")

	var request ReplicationRequest
	if err := c.BindJSON(&request); err != nil {
		log.Println("ERROR: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	var err error

	log.Println(len(request.Log))
	for _, logEntry := range request.Log {

		log.Println(logEntry.CommandType)

		if logEntry.CommandType == replication.COMMAND_TYPE_WRITE {
			db.GlobalStore.Put(logEntry.EntryKey, logEntry.EntryValue)
			replication.GetLog().AppendLog(logEntry.EntryKey, logEntry.EntryValue, replication.COMMAND_TYPE_WRITE)
		}
		if logEntry.CommandType == replication.COMMAND_TYPE_DELETE {
			err = db.GlobalStore.Delete(logEntry.EntryKey)
			replication.GetLog().AppendLog(logEntry.EntryKey, db.BoxerValue{}, replication.COMMAND_TYPE_DELETE)
		}
		if logEntry.CommandType == replication.COMMAND_TYPE_READ {
			_, err = db.GlobalStore.Get(logEntry.EntryKey)
			replication.GetLog().AppendLog(logEntry.EntryKey, db.BoxerValue{}, replication.COMMAND_TYPE_READ)
		}

	}
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	c.Status(http.StatusOK)

}
