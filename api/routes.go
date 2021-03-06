package api

import (
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

// ReplicateLog replicates the log to followers
func ReplicateLog(c *gin.Context) {
	// Bind the POST body to the request struct
	var request ReplicationRequest
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Iterate over the log and add each entry to the key value store
	var err error
	for _, logEntry := range request.Log {
		// If the entry is a WRITE command, write to the store
		if logEntry.CommandType == replication.COMMAND_TYPE_WRITE {
			db.GlobalStore.Put(logEntry.EntryKey, logEntry.EntryValue)
			replication.GetLog().AppendLog(logEntry.EntryKey, logEntry.EntryValue, replication.COMMAND_TYPE_WRITE)
		}
		// If the entry is a DELETE command, delete from the store
		if logEntry.CommandType == replication.COMMAND_TYPE_DELETE {
			err = db.GlobalStore.Delete(logEntry.EntryKey)
			replication.GetLog().AppendLog(logEntry.EntryKey, db.BoxerValue{}, replication.COMMAND_TYPE_DELETE)
		}
		// If the entry is a READ command, READ from the store
		if logEntry.CommandType == replication.COMMAND_TYPE_READ {
			_, err = db.GlobalStore.Get(logEntry.EntryKey)
			replication.GetLog().AppendLog(logEntry.EntryKey, db.BoxerValue{}, replication.COMMAND_TYPE_READ)
		}

	}
	// If an error occurred, then an error occurred with the replication, return not found
	if err != nil {
		c.Status(http.StatusNotFound)
	}
	// Otherwise, return ok
	c.Status(http.StatusOK)
}

// GetLog returns the log for this node
func GetLog(c *gin.Context) {
	log := replication.GetLog().Entries
	c.JSON(http.StatusOK, log)
}
