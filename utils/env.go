package utils

import "os"

// List all key to fetch environment variables
const (
	KEY_ENVIRONMENT = "DIONYSOS_ENVIRONMENT"
)

// defaults represents all default key-value for environment variables.
var defaults = map[string]string{
	// Describe the execution environment of the server.
	// Can have values PROD or DEV.
	// By default, it is on DEV.
	KEY_ENVIRONMENT: "DEV",
}

// LoadEnvironment set defaults values on environment variables if they are unset.
func LoadEnvironment() {
	for key, defaultValue := range defaults {
		if _, isset := os.LookupEnv(key); !isset {
			os.Setenv(key, defaultValue)
		}
	}
}
