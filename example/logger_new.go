package main

import (
	"runtime"
	"time"

	log "github.com/ianwoolf/go-logger/new"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger := log.LogDir{
		Dir:           "./log",
		FlushInterval: 2,   //s
		BufferSize:    256, // k
	}
	log.SetFall(true) // print all level to info
	log.SetConsole(true)
	logger.Init()

	for i := 10000; i > 0; i-- {
		// go logger.Backend.Log(log.ERROR, []byte("test"))
		logger.Info("test")
		logger.Debug("test")
		logger.Error("test")
		logger.Warn("test")
		logger.Fatal("test")
		time.Sleep(1000 * time.Millisecond)
	}

	select {}
}
