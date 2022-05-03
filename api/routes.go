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
	value, err := db.GlobalStore.Get(db.BoxerKey{Key: key})
	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
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

	replication.GetLog().AppendLog(boxerKey, boxerValue, replication.COMMAND_TYPE_WRITE)

	db.GlobalStore.Put(boxerKey, boxerValue)
	c.Status(http.StatusOK)
}

// DeleteKey deletes a value associated with a key
// Example: DELETE https://localhost:8080/api/v1/hello
func DeleteKey(c *gin.Context) {
	key := c.Param("key")
	err := db.GlobalStore.Delete(db.BoxerKey{Key: key})
	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		c.Status(http.StatusOK)
	}
}

// GetRole return's the role of the node
func GetRole(c *gin.Context) {
	role := config.GetConfig().Role
	c.JSON(http.StatusOK, RoleResponse{Role: role})
}
