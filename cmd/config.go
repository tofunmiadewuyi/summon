package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tofunmiadewuyi/summon/internal/config"
)

func printConfig() {
	path, err := config.ConfigPath()
	if err != nil {
		log.Println("could not load config path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Could not find config file")
		} else {
			fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("------------------------------------------------------------------")
	fmt.Println(path)
	fmt.Println("------------------------------------------------------------------")
	fmt.Println(string(data))
}
