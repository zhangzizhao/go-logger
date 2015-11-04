package main

import (
	log "github.com/ianwoolf/go-logger/new"
)

func main() {
	logger := log.LogDir{
		Dir: "./log",
	}
	logger.Init()
	logger.Backend.Log(log.ERROR, []byte("test"))
	select {}
}
