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

func (c *Configuration) Init() {
	c.LogFile = fmt.Sprintf("%v/%v", c.ConfigDir, c.LogFile)
	c.CounterFile = fmt.Sprintf("%v/%v", c.ConfigDir, c.CounterFile)
	c.TlsCrt = fmt.Sprintf("%v/%v", c.ConfigDir, c.TlsCrt)
	c.TlsKey = fmt.Sprintf("%v/%v", c.ConfigDir, c.TlsKey)

	// process flags
	verbose := flag.Bool("v", false, "Enable debug output")
	flag.Parse()
	c.DebugMode = *verbose

	// process env variables
	if os.Getenv("USHORT_LOG_FILE") != "" {
		c.LogFile = os.Getenv("USHORT_LOG_FILE")
	}
}

func LoadConfig() Configuration {
	// check if a custom configuration file path is set
	var configFileName string
	if os.Getenv("USHORT_CONFIG_FILE") != "" {
		configFileName = os.Getenv("USHORT_CONFIG_FILE")
	} else {
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
	return config
}
