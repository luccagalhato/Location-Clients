package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

//Config ...
var Config struct {
	API struct {
		Port string `yaml:"port"`
	} `yaml:"api"`
	SQL struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Interval int    `yaml:"interval"`
	} `yaml:"sql"`
}

//LoadConfig ...
func LoadConfig() error {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, &Config)
}

//CreateConfigFile ...
func CreateConfigFile() {
	if _, err := os.Stat("config.yaml"); err == nil {
		fmt.Println("the 'config.yaml' already exists, do you really want to overwrite? (y/N)")
		var rsp string
		fmt.Scan(&rsp)
		if strings.ToLower(rsp) == "y" {
			writeFile()
		}
		return
	}
	writeFile()
}

func writeFile() {
	b, _ := yaml.Marshal(Config)
	ioutil.WriteFile("config.yaml", b, 0766)
}
