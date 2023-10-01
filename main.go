package main

import (
	"log"

	"github.com/abrucker235/go-chi-starter/cmd"
)

func main() {
	if err := cmd.RootCMD.Execute(); err != nil {
		log.Fatal(err)
	}
}
