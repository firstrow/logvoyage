package configuration

import (
    "io/ioutil"
    "fmt"

    "gopkg.in/yaml.v2"
    "github.com/cosiner/gohper/errors"
)

var (
    // Cfg is a configuration parsed from yaml file
    Cfg Config

    // AlternativeConfPath used from testing to read a diferent configuration file
    AlternativeConfPath string

    initialized bool
)

// Config define a struct used to readded the configuration settings
type Config struct {
    Debug       bool            `yaml:"debug"`
    Indexes     IndexesStruct   `yaml:"indexes"`
    Redis       RedisStruct     `yaml:"redis"`
}

// SetDefaults set the defaults values of all keys has not configured
func (cfg *Config) SetDefaults() *Config {
    cfg.Indexes.SetDefaults()
    cfg.Redis.SetDefaults()
    return cfg
}

// CreateConfFile print to terminal the yaml struct to generate a new default configuration file
func CreateConfFile() {
    c := Config{
        Debug: false,
        Indexes: IndexesStruct{
            User: "users",
        },
        Redis: RedisStruct{
            Host: "localhost",
            Port: 6379,
        },
    }
    b, err := yaml.Marshal(c)
    errors.Fatalln(err)
    fmt.Println(string(b))
}

// ReadConf open the configuration file e parse the content to a ney `Config` variable
func ReadConf(str ...string) *Config {
    if ! initialized {
        var filePath string

        if len(str) == 0 && AlternativeConfPath != "" {
            filePath = AlternativeConfPath
        }else if len(str) == 0 {
            filePath = "/etc/logvoyage.yml"
        }else{
            filePath = str[0]
        }

        yamlFile, err := ioutil.ReadFile(filePath)
        errors.Fatalln(err)

        err = yaml.Unmarshal(yamlFile, &Cfg)
        errors.Fatalln(err)

        Cfg.SetDefaults()
        initialized = true
    }
    return &Cfg
}
