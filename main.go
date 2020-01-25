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
		log.Fatal(err)
	}
}
