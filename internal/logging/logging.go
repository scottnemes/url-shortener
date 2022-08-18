package logging

import (
	"log"
	"os"

	"example.com/url-shortener/internal/config"
)

var F *os.File

func StartLogging() {
	var err error
	F, err = os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file (path: %v) (%v)", config.LogFile, err)
	}
}
