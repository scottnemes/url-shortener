package config

/*
	TODO:
	1. Load configuration from file
*/

const (
	ConfigDir    = "/etc/url-shortener"
	CounterFile  = "/tmp/counter_range.dat"
	dbHost       = "127.0.0.1"
	dbPort       = "27017"
	dbUser       = "mutiny"
	dbPass       = "password123"
	DBConnString = "mongodb://" + dbHost + ":" + dbPort
	DBDatabase   = "short_urls"
	DBCollection = "urls"
)
