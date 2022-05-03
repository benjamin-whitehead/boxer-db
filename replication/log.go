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

const (
	COMMAND_TYPE_WRITE  = "WRITE"
	COMMAND_TYPE_DELETE = "DELETE"
	COMMAND_TYPE_READ   = "READ"
)

type Log struct {
	Entries []Entry
	File    *os.File
	HadData bool
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

	// Write log to file
	l.WriteLogToFile()

	return true, nil
}

var globalLog *Log = nil

func InitializeLog() {
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

func GetLog() *Log {
	if globalLog == nil {
		InitializeLog()
	}
	return globalLog
}

func (l *Log) WriteLogToFile() {
	log.Println("saving log to file...")
	// Write log to file
	jsonData, _ := json.Marshal(l.Entries)

	ioutil.WriteFile(config.GetConfig().LogFileLocation, jsonData, 0644)
	log.Println("log saved to file")
}
