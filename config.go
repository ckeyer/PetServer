package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	App *AppConfig
}
type AppConfig struct {
	Port   int
	Debug  bool
	ImgDir string
}

var config *Config

// LoadConf ...
func LoadConf(path string) (err error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Error(err)
		return
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(err)
		return
	}

	config = new(Config)
	err = json.Unmarshal(bs, config)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
