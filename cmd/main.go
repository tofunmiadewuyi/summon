package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: summon <cmd>")
		return
	}

	switch os.Args[1] {
	case "run":
		start()

	case "start":
		launchdStart()

	case "stop":
		launchdStop()

	case "add":
		add()

	case "config":
		printConfig()

	case "upgrade":
		upgrade()

	case "version", "-v":
		fmt.Println(version)

	case "help", "-h":
		fmt.Printf("Usage: %s <command>\n", program)
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  start            Install summon as a login item and start it")
		fmt.Println("  stop             Stop summon and remove it from login items")
		fmt.Println("  add              Add a new hotkey binding interactively")
		fmt.Println("  config           Print the config file path and contents")
		fmt.Printf("  upgrade          Upgrade %s to the latest release\n", program)
		fmt.Println("  version          Print the current version")
		fmt.Println("  help             Show this help message")

	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Run 'dbq help' for usage.")
	}
}
