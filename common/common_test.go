package common

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestItShouldReturnErrorIfApiKeyNotFound(t *testing.T) {
	logMessage := "0b1305-31-5f5b-5832-6a This is test logmessage"

	_, _, err := ExtractApiKey(logMessage)

	if err == nil {
		t.Fatal("It should return error")
	}
}

func TestExtractUserApiKeyAndTypeId(t *testing.T) {
	expectedKey := "0b137205-3291-5f5b-5832-ab2458b9936a"
	expectedType := "123"
	logMessage := "0b137205-3291-5f5b-5832-ab2458b9936a@123 This is test logmessage"

	key, logType, _ := ExtractApiKey(logMessage)

	if expectedKey != key {
		t.Fatal("Error extracting key")
	}
	if expectedType != logType {
		t.Fatal("Error extracting type")
	}

	// Test extraxt logType as string
	expectedType = "nginx_1"
	logMessage = "0b137205-3291-5f5b-5832-ab2458b9936a@nginx_1 This is test logmessage"

	key, logType, _ = ExtractApiKey(logMessage)

	if expectedKey != key {
		t.Fatal("Error extracting key 2")
	}
	if expectedType != logType {
		t.Fatal("Error extracting type 2")
	}

}

func TestRemoveApiKey(t *testing.T) {
	logMessage := "0b137205-3291-5f5b-5832-ab2458b9936a@2111 This is test logmessage"
	m := RemoveApiKey(logMessage)

	Convey("It should populate user from goes search response", t, func() {
		So(m, ShouldEqual, " This is test logmessage")
	})
}
