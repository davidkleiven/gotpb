package gotpb

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Link  string `yaml:"link"`
	Users []User `yaml:"users"`
}

type User struct {
	Email string `yaml:"email"`
}

func GetConf(fname string) Config {
	yamlFile, err := ioutil.ReadFile(fname)
	panicOnErr(err)

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	panicOnErr(err)
	return config
}
