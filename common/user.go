package common

import (
	"github.com/mitchellh/mapstructure"
	"net/url"
)

type User struct {
	Email    string
	Password string
	ApiKey   string
}

// Returns index name to use in Elastic
func (this *User) GetIndexName() string {
	return this.ApiKey
}

func FindUserByEmail(email string) *User {
	return FindUserBy("email", email)
}

func FindUserByApiKey(apiKey string) *User {
	return FindUserBy("apiKey", apiKey)
}

func FindUserBy(key string, value string) *User {
	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"term": map[string]interface{}{
						key: map[string]interface{}{
							"value": value,
						},
					},
				},
			},
		},
	}

	searchResults, err := GetConnection().Search(query, []string{"users"}, []string{"user"}, url.Values{})

	if err != nil || searchResults.Hits.Total == 0 {
		return nil
	}

	user := &User{}
	mapstructure.Decode(searchResults.Hits.Hits[0].Source, user)

	return user
}
