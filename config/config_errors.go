package config

import (
	"fmt"
)

const (
	MISSING_ROLE_ERROR_MESSAGE                     = "node is missing a role"
	LEADER_MISSING_REPLICATION_NODES_ERROR_MESSAGE = "leader is missing replication nodes"
	LARGE_QUORUM_SIZE_ERROR_MESSAGE                = "quorum size is larger than number of total nodes"
	MISSING_QUORUM_SIZE_ERROR_MESSAGE              = "missing quorum size"
	INVALID_QUORUM_SIZE_ERROR_MESSAGE              = "invalid quorum size, did you specify a number?"
)

func ErrorMissingRole() error {
	return fmt.Errorf(MISSING_ROLE_ERROR_MESSAGE)
}

func ErrorMissingReplicationNodes() error {
	return fmt.Errorf(LEADER_MISSING_REPLICATION_NODES_ERROR_MESSAGE)
}

func ErrorLargeQuorumSize() error {
	return fmt.Errorf(LARGE_QUORUM_SIZE_ERROR_MESSAGE)
}

func ErrorMissingQuorumSize() error {
	return fmt.Errorf(MISSING_QUORUM_SIZE_ERROR_MESSAGE)
}

func ErrorInvalidQuorumSize() error {
	return fmt.Errorf(INVALID_QUORUM_SIZE_ERROR_MESSAGE)
}
