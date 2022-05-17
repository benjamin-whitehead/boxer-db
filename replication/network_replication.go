package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/benjamin-whitehead/boxer-db/m/v2/config"
)

// Replicate sends the HTTP request to the replica nodes with the updated log
func Replicate(address string) {
	// Create a new body to send to the replica node
	postBody, _ := json.Marshal(map[string]interface{}{
		"log": globalLog.Entries,
	})
	responseBody := bytes.NewBuffer(postBody)
	// Create the endpoint that the request will be sent to
	requestUrl := fmt.Sprintf("%s/api/v1/replication", address)
	// Make the post request to the replica node, log any errors
	_, err := http.Post(requestUrl, "application/json", responseBody)
	if err != nil {
		log.Println(err.Error())
	}
}

// ReplicateLog replicates the entire log to the replica nodes
func ReplicateLog() {
	// If the node is not a leader, we don't want to
	if config.GetConfig().Role != config.ROLE_LEADER {
		log.Println("error: not leader")
		return
	}
	// create a new wait group
	var wg sync.WaitGroup
	// Get the addressess of the replica nodes
	addresses := config.GetConfig().ReplicationNodes
	// Get the size of the quorum, minus one since the leader counts as a replica
	quorum := config.GetConfig().QuorumSize - 1
	// Since we are waiting for a minimum of quorum responses, add the number of replica nodes to the wait group
	wg.Add(quorum)
	// Iterate over the replicae nodes, and send a request to each one asynchronously
	var index int
	for index = 0; index < quorum; index++ {
		go func(address string) {
			defer wg.Done()
			Replicate(address)
		}(addresses[index])
	}
	// Iterate over the remaining replica nodes, and send a request to each one asynchronously
	for ; index < len(addresses); index++ {
		go func(address string) {
			Replicate(address)
		}(addresses[index])
	}
	// Wait for all the requests to finish
	wg.Wait()
}
