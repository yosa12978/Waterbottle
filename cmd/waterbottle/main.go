package main

import (
	"log"
	"os"

	"github.com/yosa12978/waterbottle/pkg/helpers"
)

func main() {

	if rootdir := os.Getenv("WATERBOTTLE_PATH"); rootdir == "" {
		log.Fatal("Please create WATERBOTTLE_PATH environment variable")
		os.Exit(1)
	}

	if err := helpers.InitCLI(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
