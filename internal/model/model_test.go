package model

import (
	"context"
	"os"
	"testing"

	"example.com/url-shortener/internal/config"
)

// some global variables to avoid duplication
var configFileName string = "../../example/url_shortener.conf"
var verbose bool = true
var c config.Configuration = config.LoadConfig(configFileName, &verbose)

func TestGetDBClient(t *testing.T) {
	dbClient := GetDBClient(c.DBConnString)
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// if creating the connection fails, the function will fatal out prior to this
	t.Logf("PASSED creating database connection. Expected: success, got: success")
}

func TestInsertUrl(t *testing.T) {
	testLog := "/tmp/TestInsertUrl.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	dbClient := GetDBClient(c.DBConnString)
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	url := Url{Slug: "TEST1234", Target: "https://www.google.com"}
	err = InsertUrl(f, verbose, c.DBDatabase, c.DBCollection, dbClient, url)
	if err != nil {
		t.Errorf("FAILED inserting URL. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED inserting URL. Expected: nil error, got: %v", err)
	}
}

func TestGetUrl(t *testing.T) {
	testLog := "/tmp/TestGetUrl.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	dbClient := GetDBClient(c.DBConnString)
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	_, err = GetUrl(f, verbose, c.DBDatabase, c.DBCollection, dbClient, "TEST1234")
	if err != nil {
		t.Errorf("FAILED getting URL. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED getting URL. Expected: nil error, got: %v", err)
	}
}

func TestGetUrls(t *testing.T) {
	testLog := "/tmp/TestGetUrls.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	dbClient := GetDBClient(c.DBConnString)
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	_, err = GetUrls(f, verbose, c.DBDatabase, c.DBCollection, dbClient)
	if err != nil {
		t.Errorf("FAILED getting all URLs. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED getting all URLs. Expected: nil error, got: %v", err)
	}
}

func TestUpdateUrl(t *testing.T) {
	testLog := "/tmp/TestUpdateUrl.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	dbClient := GetDBClient(c.DBConnString)
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	url := Url{Slug: "TEST1234", Target: "https://www.reddit.com"}
	err = InsertUrl(f, verbose, c.DBDatabase, c.DBCollection, dbClient, url)
	if err != nil {
		t.Errorf("FAILED updating URL. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED updating URL. Expected: nil error, got: %v", err)
	}
}

func TestUpdateUrlHits(t *testing.T) {
	testLog := "/tmp/TestUpdateUrlHits.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	dbClient := GetDBClient(c.DBConnString)
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	err = UpdateUrlHits(f, verbose, c.DBDatabase, c.DBCollection, dbClient, "TEST1234")
	if err != nil {
		t.Errorf("FAILED updating URL hits. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED updating URL hits. Expected: nil error, got: %v", err)
	}
}

func TestDeleteUrl(t *testing.T) {
	testLog := "/tmp/TestDeleteUrl.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	dbClient := GetDBClient(c.DBConnString)
	// close the database connection before exit
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	err = DeleteUrl(f, verbose, c.DBDatabase, c.DBCollection, dbClient, "TEST1234")
	if err != nil {
		t.Errorf("FAILED deleting URL. Expected: nil error, got: %v", err)
	} else {
		t.Logf("PASSED deleting URL. Expected: nil error, got: %v", err)
	}
}
