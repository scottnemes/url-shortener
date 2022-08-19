package config

/*
	TODO:
	1. Load configuration from file
*/

const (
	// General
	ConfigDir   = "/etc/url_shortener"
	CounterFile = "/tmp/counter_range.dat"
	DebugMode   = true
	LogFile     = "/tmp/url_shortener.log"
	// Limits
	MaxSlugLen = 7
	// Gin
	GinPort = "8443"
	// TLS
	TlsCrt = "/tmp/localhost.crt"
	TlsKey = "/tmp/localhost.key"
	// Database
	dbHost       = "127.0.0.1"
	dbPort       = "27017"
	dbUser       = "mutiny"
	dbPass       = "password123"
	DBConnString = "mongodb://" + dbHost + ":" + dbPort
	DBDatabase   = "short_urls"
	DBCollection = "urls"
	// Cache
	CacheEnabled     = true
	CacheHost        = "localhost"
	CachePort        = "6379"
	CachePass        = ""
	CacheDB          = 0
	CacheExpireHours = 1
)
