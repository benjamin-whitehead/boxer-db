package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/benjamin-whitehead/boxer-db/m/v2/config"
)

func ReplicateLog() {
	addresses := config.GetConfig().ReplicationNodes

	// Get the latest entry in the log
	latestEntry := globalLog.Entries[len(globalLog.Entries)-1]

	postBody, _ := json.Marshal(map[string]interface{}{
		"Log": latestEntry,
	})

	responseBody := bytes.NewBuffer(postBody)

	requestUrl := fmt.Sprintf("%s/api/v1/replication/", addresses[0])

	response, err := http.Post(requestUrl, "application/json", responseBody)
	if err != nil {
		log.Fatal(err.Error())
	}

	body, _ := ioutil.ReadAll(response.Body)

	log.Println(string(body))

}
