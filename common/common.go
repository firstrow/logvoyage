package common

import (
	"errors"
	"regexp"
)

const (
	ApiKeyFormat = "^([a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12})"
)

// Extracts api key from string
func ExtractApiKey(message string) (string, error) {
	re := regexp.MustCompile(ApiKeyFormat)
	result := re.Find([]byte(message))

	if result == nil {
		return "", errors.New("Error extracting apiKey")
	}

	return string(result), nil
}

// Extracts api key from string
func RemoveApiKey(message string) string {
	re := regexp.MustCompile(ApiKeyFormat)
	result := re.ReplaceAll([]byte(message), []byte(""))

	return string(result)
}
