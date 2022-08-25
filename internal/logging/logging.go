package logging

import (
	"log"
	"os"
)

/*
	Takes in the full path to the log file. Returns the file handler for use in other packages.
*/
func StartLogging(logFile string) *os.File {
	var err error
	var f *os.File
	f, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file (path: %v) (%v)", logFile, err)
	}
	return f
}
