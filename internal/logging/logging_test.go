package logging

import (
	"os"
	"testing"
)

/*
	Tests the StartLogging function
*/
func TestStartLogging(t *testing.T) {
	testLog := "/tmp/TestStartLogging.log"
	f := StartLogging(testLog)
	defer f.Close()
	_, err := f.WriteString("Testing")
	if err != nil {
		t.Errorf("FAILED writing to %v. Expected: nil error, got: %v", testLog, err)
	} else {
		t.Logf("PASSED writing to %v.", testLog)
	}

	os.Remove(testLog)
}
