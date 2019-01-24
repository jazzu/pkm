package tools

import "os"

func GetEnvParam(key string, defaultValue string) string {
	result := os.Getenv(key)
	if len(result) == 0 {
		return defaultValue
	}
	return result
}
