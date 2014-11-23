package common

import (
	"github.com/belogik/goes"
	"github.com/mitchellh/mapstructure"
	"net/url"
)

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	ApiKey    string `json:"apiKey"`
}

// Returns index name to use in Elastic
func (this *User) GetIndexName() string {
	return this.ApiKey
}

func (this *User) GetLogTypes() []string {
	userLogTypes, _ := GetTypes(this.GetIndexName())
	return userLogTypes
}

func FindUserByEmail(email string) *User {
	return FindUserBy("email", email)
}

func FindUserByApiKey(apiKey string) *User {
	return FindUserBy("apiKey", apiKey)
}

func (this *User) Save() {
	doc := goes.Document{
		Index:  "users",
		Type:   "user",
		Id:     this.Id,
		Fields: this,
	}
	extraArgs := make(url.Values, 0)
	GetConnection().Index(doc, extraArgs)
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
	user.Id = searchResults.Hits.Hits[0].Id

	return user
}
