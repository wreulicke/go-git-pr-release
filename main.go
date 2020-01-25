package main

import (
	"log"
	"os"
)

func mainInternal() error {
	return Run(os.Args)
}

func main() {
	if err := mainInternal(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
