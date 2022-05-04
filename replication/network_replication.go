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
	if config.GetConfig().Role == config.ROLE_LEADER {
		addresses := config.GetConfig().ReplicationNodes

		// Get the latest entry in the log
		latestEntry := globalLog.Entries

		for _, address := range addresses {

			postBody, _ := json.Marshal(map[string]interface{}{
				"log": latestEntry,
			})

			responseBody := bytes.NewBuffer(postBody)

			requestUrl := fmt.Sprintf("%s/api/v1/replication", address)
			log.Println(requestUrl)

			response, err := http.Post(requestUrl, "application/json", responseBody)
			if err != nil {
				log.Fatal(err.Error())
			}

			body, _ := ioutil.ReadAll(response.Body)

			log.Println(string(body))

		}

	}

}
