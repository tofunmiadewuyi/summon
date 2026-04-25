package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const plistLabel = "com.summon"

var plistTemplate = template.Must(template.New("plist").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>{{.Label}}</string>
	<key>ProgramArguments</key>
	<array>
		<string>{{.BinaryPath}}</string>
		<string>run</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
</dict>
</plist>
`))

func plistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", plistLabel+".plist"), nil
}

func launchdStart() {
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine binary path: %v\n", err)
		os.Exit(1)
	}

	path, err := plistPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine plist path: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "could not create LaunchAgents directory: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not write plist: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := plistTemplate.Execute(f, struct {
		Label      string
		BinaryPath string
	}{plistLabel, binaryPath}); err != nil {
		fmt.Fprintf(os.Stderr, "could not render plist: %v\n", err)
		os.Exit(1)
	}

	if out, err := exec.Command("launchctl", "load", path).CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "launchctl load failed: %s\n", out)
		os.Exit(1)
	}

	fmt.Println("summon started and will run on login")
}

func launchdStop() {
	path, err := plistPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine plist path: %v\n", err)
		os.Exit(1)
	}

	if out, err := exec.Command("launchctl", "unload", path).CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "launchctl unload failed: %s\n", out)
		os.Exit(1)
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "could not remove plist: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("summon stopped and removed from login items")
}
