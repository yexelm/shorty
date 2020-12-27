package config

import "os"

// GetEnv fallbacks env variable to default value if it was not set
func GetEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}
