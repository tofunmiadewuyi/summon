package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tofunmiadewuyi/summon/internal/config"
	"github.com/tofunmiadewuyi/summon/internal/hotkey"
)

func add() {
	reader := bufio.NewReader(os.Stdin)
	var combo string

	for {
		var err error
		combo, err = hotkey.CaptureHotkey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to capture hotkey: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Detected: %s\n", combo)
		fmt.Printf("Add binding for %s? [y/N] ", combo)

		answer, _ := reader.ReadString('\n')
		runes := []rune(strings.TrimSpace(strings.ToLower(answer)))
		if len(runes) > 0 && runes[len(runes)-1] == 'y' {
			break
		}
		fmt.Println("Try again.")
	}

	fmt.Print("Enter app name: ")
	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)
	if appName == "" {
		fmt.Fprintln(os.Stderr, "app name cannot be empty")
		os.Exit(1)
	}

	if err := config.AppendBinding(combo, appName); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save binding: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added: %s → %s\n", combo, appName)
}
