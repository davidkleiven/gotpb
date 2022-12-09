package gotpb

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Groups map[string]string `yaml:"groups"`
	Users  []User            `yaml:"users"`
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
