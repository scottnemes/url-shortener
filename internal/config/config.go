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
}

const (
	configFileName = "/etc/url_shortener/url_shortener.conf"
)

func LoadConfig() Configuration {
	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal(err)
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
