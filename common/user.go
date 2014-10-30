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

func FindUserByEmail(email string) *User {
	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"term": map[string]interface{}{
						"email": map[string]interface{}{
							"value": email,
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
