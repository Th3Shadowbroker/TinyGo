package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	DefaultConfiguration = Configuration{
		Service: ServiceConfiguration{
			Port: 8080,
		},

		Database: DatabaseConfiguration{
			Host:     "localhost",
			Port:     27017,
			User:     "",
			Password: "",
			Database: "tiny",
		},
	}
)

type Configuration struct {
	Service ServiceConfiguration `json:"service"`

	Database DatabaseConfiguration `json:"database"`
}

type ServiceConfiguration struct {
	Port int `json:"port"`
}

type DatabaseConfiguration struct {
	Host string `json:"host"`

	Port int `json:"port"`

	User string `json:"user"`

	Password string `json:"password"`

	Database string `json:"database"`
}

func LoadConfiguration(filename string) Configuration {
	var config Configuration

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config = DefaultConfiguration
		bytes, _ := json.MarshalIndent(DefaultConfiguration, "", "  ")
		ioutil.WriteFile(filename, bytes, 0644)
	} else {
		bytes, _ := ioutil.ReadFile(filename)
		json.Unmarshal(bytes, &config)
	}

	return config
}

func (c *DatabaseConfiguration) ConnectionString() string {
	if len(c.User) == 0 && len(c.Password) == 0 {
		return fmt.Sprintf("mongodb://%s:%d/%s", c.Host, c.Port, c.Database)
	} else {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	}
}
