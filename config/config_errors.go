package config

import (
	"fmt"
)

const (
	MISSING_ROLE_ERROR_MESSAGE                     = "node is missing a role"
	LEADER_MISSING_REPLICATION_NODES_ERROR_MESSAGE = "leader is missing replication nodes"
)

func ErrorMissingRole() error {
	return fmt.Errorf(MISSING_ROLE_ERROR_MESSAGE)
}

func ErrorMissingReplicationNodes() error {
	return fmt.Errorf(LEADER_MISSING_REPLICATION_NODES_ERROR_MESSAGE)
}
