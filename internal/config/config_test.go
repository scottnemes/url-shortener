package config

import (
	"testing"
)

/*
	Tests the LoadConfig function
*/
func TestLoadConfig(t *testing.T) {
	// move flag parsing to API, maybe new init func, pass in values to loadconfig func
	configFileName := "../../example/url_shortener.conf"
	verbose := true
	config := LoadConfig(configFileName, &verbose)
	if config.ConfigDir != "/etc/url_shortener" {
		t.Errorf("FAILED to load and validate configuration file. Expected: /etc/url_shortener, got: %v", config.ConfigDir)
	} else {
		t.Logf("PASSED loading and validating configuration file. Expected: /etc/url_shortener, got: %v", config.ConfigDir)
	}
}
