package main

import (
	"github.com/firstrow/logvoyage/common"
	"testing"
)

func TestItShouldReturnErrorIfApiKeyNotFound(t *testing.T) {
	logMessage := "0b1305-31-5f5b-5832-6a This is test logmessage"

	_, err := common.ExtractApiKey(logMessage)

	if err == nil {
		t.Fatal("It should return error")
	}
}

func TestExtractUserApiKeyFromString(t *testing.T) {
	apiKey := "0b137205-3291-5f5b-5832-ab2458b9936a"
	logMessage := "0b137205-3291-5f5b-5832-ab2458b9936a This is test logmessage"

	key, _ := common.ExtractApiKey(logMessage)

	if key != apiKey {
		t.Fatal("Error extracting key")
	}
}
