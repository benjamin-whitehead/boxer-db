package db

import "sync"

type BoxerKey struct {
	Key string
}

type BoxerValue struct {
	Value interface{}        `json:"value"`
	Meta  BoxerValueMetadata `json:"meta"`
}

type BoxerValueMetadata struct {
	Timestamp int64 `json:"timestamp"`
}

type BoxerStore struct {
	Container map[BoxerKey]BoxerValue
	Mutex     sync.Mutex
}
