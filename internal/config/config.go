// Package config provides a struct to store the applications config
package config

import (
	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/otelx"
)

// AppConfig stores all the config values for our application
var AppConfig struct {
	Logging  loggingx.Config
	Tracing  otelx.Config
	Generate GenerateConfig
}
