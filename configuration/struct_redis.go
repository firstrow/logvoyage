package configuration

import "strconv"

// RedisStruct define de struct of indexes names to ElasticSearch
type RedisStruct struct {
    Host	string 	`yaml:"host"`
    Port	int		`yaml:"port"`
}

// SetDefaults set the defaults values of all keys has not configured
func (s *RedisStruct) SetDefaults() {
    if len(s.Host) == 0  { s.Host = "localhost" }
    if s.Port <= 0       { s.Port = 6379 }
}

// GetURI returns the uri to create conecction
func (s *RedisStruct) GetURI() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}