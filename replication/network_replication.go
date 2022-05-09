package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/benjamin-whitehead/boxer-db/m/v2/config"
)

func Replicate(address string) {
	postBody, _ := json.Marshal(map[string]interface{}{
		"log": globalLog.Entries,
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

func ReplicateLog() {
	if config.GetConfig().Role != config.ROLE_LEADER {
		log.Println("error: not leader")
		return
	}

	var wg sync.WaitGroup

	addresses := config.GetConfig().ReplicationNodes

	quorum := config.GetConfig().QuorumSize - 1

	var index int

	wg.Add(quorum)
	for index = 0; index < quorum; index++ {

		go func(address string) {
			defer wg.Done()
			Replicate(address)
		}(addresses[index])

	}

	for ; index < len(addresses); index++ {
		go func(address string) {
			Replicate(address)
		}(addresses[index])
	}
	wg.Wait()

}
