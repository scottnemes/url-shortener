package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Configuration struct {
	// General
	ConfigDir   string
	LogDir      string
	CounterFile string
	DebugMode   bool
	LogFile     string
	// Limits
	MaxSlugLen int
	// Gin
	GinPort string
	// TLS
	TlsCrt string
	TlsKey string
	// Database
	DBUser       string
	DBPass       string
	DBConnString string
	DBDatabase   string
	DBCollection string
	// Cache
	CacheEnabled     bool
	CacheHost        string
	CachePort        string
	CachePass        string
	CacheDB          int
	CacheExpireHours time.Duration
}

/*
	Processes initial command line flags. Loads configuration from disk, with the option to override the default config file with the -config command.
	Sets initial values for multiple variables that require a path in order to use the passed in config file if applicable.
*/
func LoadConfig(configFileName string, verbose *bool) Configuration {
	// if no config dir is provided, set a default to load
	if configFileName == "" {
		configFileName = "/etc/url_shortener/url_shortener.conf"
	}

	// open and decode configuration file
	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatalf("Cannot load configuration file (%v)", err)
	}
	defer configFile.Close()

	config := Configuration{}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	// set debug based on flag from above
	config.DebugMode = *verbose
	// build file paths based on relevent directories
	config.LogFile = fmt.Sprintf("%v/%v", config.LogDir, config.LogFile)
	config.CounterFile = fmt.Sprintf("%v/%v", config.ConfigDir, config.CounterFile)
	config.TlsCrt = fmt.Sprintf("%v/%v", config.ConfigDir, config.TlsCrt)
	config.TlsKey = fmt.Sprintf("%v/%v", config.ConfigDir, config.TlsKey)

	return config
}
