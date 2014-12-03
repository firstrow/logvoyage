package common

import (
	"errors"
	"net/url"
	"sort"

	"github.com/Unknwon/com"
	"github.com/belogik/goes"
	"github.com/mitchellh/mapstructure"
	"github.com/xlab/handysort"
)

type User struct {
	Id        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Password  string     `json:"password"`
	ApiKey    string     `json:"apiKey"`
	Projects  []*Project `json:"projects"`
}

// Returns index name to use in Elastic
func (u *User) GetIndexName() string {
	return u.ApiKey
}

// Returns elastic search types
func (u *User) GetLogTypes() []string {
	t, err := GetTypes(u.GetIndexName())
	if err != nil {
		sort.Sort(handysort.Strings(t))
	}
	return t
}

////////////////////////
// Source groups
////////////////////////

// Source group represent group of log types.
// Each log type can be in various groups at the same time.
type Project struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Types       []string `json:"types"`
}

func (u *User) AddProject(p *Project) *User {
	if p.Id == "" {
		key := com.RandomCreateBytes(5)
		p.Id = string(key)
		u.Projects = append(u.Projects, p)
	} else {
		u.UpdateProject(p)
	}
	return u
}

func (u *User) UpdateProject(p *Project) {
	for key, g := range u.Projects {
		if p.Id == g.Id {
			u.Projects[key] = p
		}
	}
}

func (u *User) DeleteProject(id string) {
	for i, val := range u.Projects {
		if val.Id == id {
			copy(u.Projects[i:], u.Projects[i+1:])
			u.Projects[len(u.Projects)-1] = nil // or the zero value of T
			u.Projects = u.Projects[:len(u.Projects)-1]
			return
		}
	}
}

func (u *User) GetProject(id string) (*Project, error) {
	for _, val := range u.Projects {
		if val.Id == id {
			return val, nil
		}
	}
	return nil, errors.New("Source group not found")
}

////////////////////////
// Finders
////////////////////////

func FindUserByEmail(email string) (*User, error) {
	return FindUserBy("email", email)
}

func FindUserByApiKey(apiKey string) (*User, error) {
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

// Find user by any param.
// Returns err if ES can't perform/accept query,
// and nil if user not found.
func FindUserBy(key string, value string) (*User, error) {
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

	if err != nil {
		return nil, ErrSendingElasticSearchRequest
	}
	if searchResults.Hits.Total == 0 {
		return nil, nil
	}

	user := &User{}
	mapstructure.Decode(searchResults.Hits.Hits[0].Source, user)
	user.Id = searchResults.Hits.Hits[0].Id
	return user, nil
}
