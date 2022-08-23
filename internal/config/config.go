package config

import (
	"encoding/json"
	"flag"
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

func LoadConfig() Configuration {
	var configDir string
	// process flags
	verbose := flag.Bool("v", false, "Enable debug output")
	flag.StringVar(&configDir, "config", "", "Path to configuration directory")
	flag.Parse()

	// if no config dir is provided, set a default to load
	if configDir == "" {
		configDir = "/etc/url_shortener"
	}

	configFileName := fmt.Sprintf("%v/url_shortener.conf", configDir)

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

	// persist config dir in case one was passed in
	config.ConfigDir = configDir
	// set debug based on flag from above
	config.DebugMode = *verbose
	// build file paths based on relevent directories
	config.LogFile = fmt.Sprintf("%v/%v", config.LogDir, config.LogFile)
	config.CounterFile = fmt.Sprintf("%v/%v", config.ConfigDir, config.CounterFile)
	config.TlsCrt = fmt.Sprintf("%v/%v", config.ConfigDir, config.TlsCrt)
	config.TlsKey = fmt.Sprintf("%v/%v", config.ConfigDir, config.TlsKey)

	return config
}
