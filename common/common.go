package common

import (
	"errors"
	"regexp"
)

const (
	format = "^([a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12})"
)

// Extracts api key from string
func ExtractApiKey(message string) (string, error) {
	re := regexp.MustCompile(format)
	result := re.Find([]byte(message))

	if result == nil {
		return "", errors.New("Error extracting apiKey")
	}

	return string(result), nil
}
