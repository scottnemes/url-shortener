package util

import (
	"os"
	"testing"
)

/*
	Tests the GetNewRange and GetAndIncrease functions
*/
func TestGetAndIncrease(t *testing.T) {
	testLog := "/tmp/TestGetNewRange.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	c := Counter{}

	// tests getting a new range
	c.GetAndIncrease(f, true)
	if c.Counter != 1000000 || c.CounterEnd != 2000000 {
		t.Errorf("FAILED getting new counter range. Expected: 1000000/2000000, got: %v/%v", c.Counter, c.CounterEnd)
	} else {
		t.Logf("PASSED getting new counter range. Expected: 1000000/2000000, got: %v/%v", c.Counter, c.CounterEnd)
	}

	// tests increasing an existing range
	c.GetAndIncrease(f, true)
	if c.Counter != 1000001 || c.CounterEnd != 2000000 {
		t.Errorf("FAILED incrementing existing counter range. Expected: 1000001/2000000, got: %v/%v", c.Counter, c.CounterEnd)
	} else {
		t.Logf("PASSED incrementing existing counter range. Expected: 1000001/2000000, got: %v/%v", c.Counter, c.CounterEnd)
	}

	os.Remove(testLog)
}

/*
	Tests the SaveCounterRange and LoadCounterRange functions
*/
func TestLoadCounterRange(t *testing.T) {
	testLog := "/tmp/TestGetNewRange.log"
	testFile := "/tmp/TestLoadCounterRange.dat"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()

	c := Counter{}
	c.Counter = 100
	c.CounterEnd = 200

	c.SaveCounterRange(f, true, testFile)
	c.LoadCounterRange(f, true, testFile)

	if c.Counter != 100 || c.CounterEnd != 200 {
		t.Errorf("FAILED saving and load counter range from file. Expected: 100/200, got: %v/%v", c.Counter, c.CounterEnd)
	} else {
		t.Logf("PASSED saving and loading counter range from file. Expected: 100/200, got: %v/%v", c.Counter, c.CounterEnd)
	}

	os.Remove(testFile)
	os.Remove(testLog)
}

func TestFileExists(t *testing.T) {
	testFile := "/tmp/TestFileExists.log"
	f, err := os.OpenFile(testFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()

	fileExists := FileExists(testFile)

	if fileExists != true {
		t.Errorf("FAILED testing if a file exists. Expected: true, got: %v", fileExists)
	} else {
		t.Logf("PASSED testing if a file exists. Expected: true, got: %v", fileExists)
	}

	os.Remove(testFile)
}

func TestGenerateUrlSlug(t *testing.T) {
	testLog := "/tmp/TestGetNewRange.log"
	f, err := os.OpenFile(testLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("FAILED opening log file. Expected: nil error, got: %v", err)
	}
	defer f.Close()
	c := Counter{}
	c.Counter = 1000000
	slug := GenerateUrlSlug(f, true, &c)
	if slug != "4C92" {
		t.Errorf("FAILED generating slug. Expected: 4C92, got: %v", slug)
	} else {
		t.Logf("PASSED generating slug. Expected 4C92, got: %v", slug)
	}
}

func TestIsValidSlug(t *testing.T) {
	// test a valid slug
	slug := "4C92"
	isValid := IsValidSlug(7, slug)
	if isValid != true {
		t.Errorf("FAILED validating valid slug. Expected: true, got: %v", isValid)
	} else {
		t.Logf("PASSED validating valid slug. Expected true, got: %v", isValid)
	}

	// test an invalid slug (non-base62 char)
	slug = "4C92$"
	isValid = IsValidSlug(7, slug)
	if isValid != false {
		t.Errorf("FAILED validating invalid slug (non-base62 char). Expected: false, got: %v", isValid)
	} else {
		t.Logf("PASSED validating invalid slug (non-base62 char). Expected false, got: %v", isValid)
	}

	// test an invalid slug (too long)
	slug = "4C92"
	isValid = IsValidSlug(3, slug)
	if isValid != false {
		t.Errorf("FAILED validating invalid slug (too long). Expected: false, got: %v", isValid)
	} else {
		t.Logf("PASSED validating invalid slug (too long). Expected false, got: %v", isValid)
	}

	// test an invalid slug (empty slug)
	slug = ""
	isValid = IsValidSlug(7, slug)
	if isValid != false {
		t.Errorf("FAILED validating invalid slug (empty slug). Expected: false, got: %v", isValid)
	} else {
		t.Logf("PASSED validating invalid slug (empty slug). Expected false, got: %v", isValid)
	}
}

func TestIsValidUrl(t *testing.T) {
	// test a valid URL
	url := "https://www.google.com"
	isValid := IsValidUrl(url)
	if isValid != true {
		t.Errorf("FAILED validating valid url. Expected: true, got: %v", isValid)
	} else {
		t.Logf("PASSED validating valid url. Expected true, got: %v", isValid)
	}

	// test an invalid URL (no protocol)
	url = "www.google.com"
	isValid = IsValidUrl(url)
	if isValid != false {
		t.Errorf("FAILED validating invalid url (no protocol). Expected: false, got: %v", isValid)
	} else {
		t.Logf("PASSED validating invalid url (no protocol). Expected false, got: %v", isValid)
	}
}
