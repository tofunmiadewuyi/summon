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
	combo, err := hotkey.CaptureHotkey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to capture hotkey: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Detected: %s\n", combo)
	fmt.Printf("Add binding for %s? [y/N] ", combo)

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(answer)) != "y" {
		fmt.Println("Cancelled.")
		return
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
