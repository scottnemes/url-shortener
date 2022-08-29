package cache

import (
	"os"
	"testing"

	"example.com/url-shortener/internal/config"
	"example.com/url-shortener/internal/model"
)

// some global variables to avoid duplication
var configFileName string = "../../example/url_shortener.conf"
var verbose bool = true
var c config.Configuration = config.LoadConfig(configFileName, &verbose)

func TestGetCacheClient(t *testing.T) {
	_ = GetCacheClient(c.CacheHost, c.CachePort, c.CacheDB, c.CachePass)
	// if creating the connection fails, the function will fatal out prior to this
	t.Logf("PASSED creating cache connection. Expected: success, got: success")
}

func TestSetCachedUrl(t *testing.T) {
	testLog := "/tmp/TestSetCachedUrl.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	cacheClient := GetCacheClient(c.CacheHost, c.CachePort, c.CacheDB, c.CachePass)
	url := model.Url{Slug: "TEST1234", Target: "https://www.google.com"}
	SetCachedUrl(f, verbose, c.CacheExpireHours, cacheClient, url)
	if err != nil {
		t.Errorf("FAILED setting cached URL. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED setting cached URL. Expected: nil error, got: %v", err)
	}
}

func TestGetCachedUrl(t *testing.T) {
	testLog := "/tmp/TestGetCachedUrl.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	cacheClient := GetCacheClient(c.CacheHost, c.CachePort, c.CacheDB, c.CachePass)
	_, err = GetCachedUrl(f, verbose, cacheClient, "TEST1234")
	if err != nil {
		t.Errorf("FAILED getting cached URL. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED getting cached URL. Expected: nil error, got: %v", err)
	}
}

func TestDeleteCachedUrl(t *testing.T) {
	testLog := "/tmp/TestDeleteCachedUrl.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	cacheClient := GetCacheClient(c.CacheHost, c.CachePort, c.CacheDB, c.CachePass)
	err = DeleteCachedUrl(f, verbose, cacheClient, "TEST1234")
	if err != nil {
		t.Errorf("FAILED deleting cached URL. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED deleting cached URL. Expected: nil error, got: %v", err)
	}
}
