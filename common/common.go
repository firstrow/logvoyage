package common

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

const (
	ApiKeyFormat = `^([a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12})@([a-z0-9\_]{1,20})`
)

var (
	ErrExtractingKey = errors.New("Error extracting key and type")
)

// Extracts api key and log type from string
func ExtractApiKey(message string) (string, string, error) {
	re := regexp.MustCompile(ApiKeyFormat)
	result := re.FindAllStringSubmatch(message, -1)

	if result == nil {
		return "", "", ErrExtractingKey
	}

	return result[0][1], result[0][2], nil
}

// Removes api key and log type from string
func RemoveApiKey(message string) string {
	re := regexp.MustCompile(ApiKeyFormat)
	result := re.ReplaceAll([]byte(message), []byte(""))

	return string(result)
}

// Builds full path to application based on $GOPATH
func AppPath(elem ...string) string {
	app_path := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "firstrow", "logvoyage")
	dir_path := append([]string{app_path}, elem...)
	return filepath.Join(dir_path...)
}
