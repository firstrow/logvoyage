package configuration

// IndexesStruct define de struct of indexes names to ElasticSearch
type IndexesStruct struct {
    User    string  `yaml:"user"`
}

// SetDefaults set the defaults values of all keys has not configured
func (s *IndexesStruct) SetDefaults() {
    if len(s.User) == 0  { s.User = "users" }
}