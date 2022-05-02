package replication

import (
	"time"

	"github.com/benjamin-whitehead/boxer-db/m/v2/db"
)

// For replication, I have decided to go with Write-Ahead Logging (WAL)
// More information can be found at: https://martinfowler.com/articles/patterns-of-distributed-systems/wal.html

const (
	COMMAND_TYPE_WRITE  = "WRITE"
	COMMAND_TYPE_DELETE = "DELETE"
	COMMAND_TYPE_READ   = "READ"
)

type Log struct {
	Entries []Entry
	Index   int64 // Possibly delete?
}

type Entry struct {
	EntryKey    db.BoxerKey
	EntryValue  db.BoxerValue
	CommandType string
	Timestamp   int64
}

func (l *Log) AppendLog(key db.BoxerKey, value db.BoxerValue, commandType string) (bool, error) {
	if commandType == COMMAND_TYPE_WRITE {
		l.Entries = append(l.Entries, Entry{key, value, commandType, time.Now().UnixNano()})
	}
	return true, nil
}

var GlobalLog *Log = nil

func InitializeLog() {
	if GlobalLog == nil {
		GlobalLog = &Log{
			Entries: make([]Entry, 0),
			Index:   0,
		}
	}
}
