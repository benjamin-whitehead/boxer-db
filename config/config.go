package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	ROLE_FOLLOWER = "FOLLOWER"
	ROLE_LEADER   = "LEADER"
)

var globalConfig *Configuration = nil

// Configuration struct will contain necessary values needed for the program to run
// This will contain the CLI flags and environment variables
type Configuration struct {
	Role             string
	ReplicationNodes []string
	QuorumSize       int
}

func handleEnvironmentVariables() error {
	roleEnv := os.Getenv("NODE_ROLE")

	if roleEnv == "" {
		return ErrorMissingRole()
	}

	if roleEnv == ROLE_LEADER {
		globalConfig.Role = ROLE_LEADER
	} else {
		globalConfig.Role = ROLE_FOLLOWER
	}

	replicationNodes := os.Getenv("REPLICATION_NODES")
	if roleEnv == ROLE_LEADER && replicationNodes == "" {
		return ErrorMissingReplicationNodes()
	}

	followerNodes := strings.Split(replicationNodes, ",")
	globalConfig.ReplicationNodes = followerNodes

	quorumSize := os.Getenv("QUORUM_SIZE")
	if roleEnv == ROLE_LEADER && quorumSize == "" {
		return ErrorMissingQuorumSize()
	}

	size, err := strconv.Atoi(quorumSize)
	if roleEnv == ROLE_LEADER && err != nil {
		return ErrorInvalidQuorumSize()
	}

	globalConfig.QuorumSize = size

	return nil
}

func GetConfig() *Configuration {
	if globalConfig == nil {
		globalConfig = &Configuration{}

		// Handle any additional setup here
		err := handleEnvironmentVariables()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return globalConfig
}
