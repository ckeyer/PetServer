package main

type Config struct {
	App *AppConfig
}
type AppConfig struct {
	Port  int
	Debug bool
}

var config *Config
