package config

import (
	"os"
)

// define global const values
const (
	// Log level limited
	LogLevel = 3
)

// define global variable values
var (
	GlobalValueMap = map[string]interface{}{}
	// Root work dir
	Root = os.Getenv("PWD")
	// LogPath global for runner
	LogPath string
	// TestCommand global for runner
	TestCommand string
)

// SetValue .
func SetValue(key string, value interface{}) {
	GlobalValueMap[key] = value
}

// GetValue .
func GetValue(key string) interface{} {
	return GlobalValueMap[key]
}
