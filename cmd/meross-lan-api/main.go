package main

import (
	"log"
)

var (
	config    Config
	apiServer APIServer
)

func main() {
	log.Println("Starting...")

	config = config.Load()
	apiServer.Start()
}
