package gotpb

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Groups       map[string]string `yaml:"groups"`
	Users        []User            `yaml:"users"`
	MemoryMonths time.Duration     `yaml:"memoryMonts"`
	Db           string            `yaml:"db"`
}

func (c Config) UsersInGroup(group string) []User {
	users := []User{}
	for _, user := range c.Users {
		if user.Group == group {
			users = append(users, user)
		}
	}
	return users
}

type User struct {
	Email string `yaml:"email"`
	Group string `yaml:"group"`
}

func GetConf(fname string) Config {
	yamlFile, err := ioutil.ReadFile(fname)
	panicOnErr(err)

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	panicOnErr(err)
	return config
}

func ValidateConf(conf Config) bool {
	for _, user := range conf.Users {
		if _, ok := conf.Groups[user.Group]; !ok {
			return false
		}
	}
	return true
}
