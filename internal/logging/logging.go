package logging

import (
	"log"
	"os"
)

var F *os.File

func StartLogging(logFile string) {
	var err error
	F, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file (path: %v) (%v)", logFile, err)
	}
}
