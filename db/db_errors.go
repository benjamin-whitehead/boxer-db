package db

import (
	"fmt"
)

const (
	KEY_NOT_FOUND_ERROR_MESSAGE = "key %s not found"
)

func ErrorKeyNotFound(key BoxerKey) error {
	return fmt.Errorf(fmt.Sprintf(KEY_NOT_FOUND_ERROR_MESSAGE, key.Key))
}
