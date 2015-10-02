package main

import (
	logpkg "github.com/ckeyer/go-log"
)

var log *logpkg.Logger

func init() {
	log = logpkg.GetDefaultLogger()
}
