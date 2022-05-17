package replication

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/benjamin-whitehead/boxer-db/m/v2/config"
	"github.com/benjamin-whitehead/boxer-db/m/v2/db"
)

// For replication, I have decided to go with Write-Ahead Logging (WAL)
// More information can be found at: https://martinfowler.com/articles/patterns-of-distributed-systems/wal.html

// These are the types of commands that can be written to the log
const (
	COMMAND_TYPE_WRITE  = "WRITE"
	COMMAND_TYPE_DELETE = "DELETE"
	COMMAND_TYPE_READ   = "READ"
)

// Log represents the WAL log that is persisted to disk, as well as replicated to follower nodes
type Log struct {
	Entries []Entry  // The entries in the log
	File    *os.File // The file that the log is persisted to
	HadData bool     // Whether or not the log had data in it
}

// Entry represents a single entry in the WAL log
type Entry struct {
	EntryKey    db.BoxerKey   `json:"key"`          // The key of the entry
	EntryValue  db.BoxerValue `json:"value"`        // The value of the entry
	CommandType string        `json:"command_type"` // The type of command that was executed
	Timestamp   int64         `json:"timestamp"`    // The timestamp of the log write
}

var globalLog *Log = nil

// WriteEntry writes an entry to the log
func (l *Log) AppendLog(key db.BoxerKey, value db.BoxerValue, commandType string) (bool, error) {

	l.Entries = append(l.Entries, Entry{key, value, commandType, time.Now().UnixNano()})

	// Write log to file
	return l.writeLogToFile()
}

// InitializeLog initializes the log, and loads it from disk if it exists
func initializeLog() {
	if globalLog == nil {
		globalLog = &Log{
			Entries: make([]Entry, 0),
		}
		// Check if the file exists at the path
		file, err := os.Open(config.GetConfig().LogFileLocation)
		if os.IsNotExist(err) {
			log.Println("Log file does not exist, creating new one")
			// File does not exist, create it
			file, err = os.Create(config.GetConfig().LogFileLocation)
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			log.Println("Log file exists, reading from it")
			byteValue, _ := ioutil.ReadAll(file)
			var entries []Entry
			json.Unmarshal(byteValue, &entries)
			globalLog.Entries = entries
			globalLog.HadData = true
		}
		globalLog.File = file
	}
}

// GetLog returns the global log
func GetLog() *Log {
	if globalLog == nil {
		initializeLog()
	}
	return globalLog
}

// WriteLogToFile writes the log to disk
func (l *Log) writeLogToFile() (bool, error) {
	log.Println("saving log to file...")
	jsonData, err := json.Marshal(l.Entries)
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(config.GetConfig().LogFileLocation, jsonData, 0644)
	if err != nil {
		return false, err
	}
	log.Println("log saved to file")

	// Replicate the log to follower nodes
	ReplicateLog()

	return true, nil
}
