// Package service manages the running of summons in the background
package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/tofunmiadewuyi/summon/internal/config"
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
	<key>StandardOutPath</key>
	<string>{{.LogPath}}</string>
	<key>StandardErrorPath</key>
	<string>{{.LogPath}}</string>
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

func LaunchdStart() {
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine binary path: %v\n", err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine home directory: %v\n", err)
		os.Exit(1)
	}
	logPath := filepath.Join(home, "Library", "Logs", "summon.log")

	path, err := plistPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine plist path: %v\n", err)
		os.Exit(1)
	}

	// Unload any existing instance before writing the new plist so a fresh
	// daemon always starts with the current binary.
	exec.Command("launchctl", "unload", path).Run()

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
		LogPath    string
	}{plistLabel, binaryPath, logPath}); err != nil {
		fmt.Fprintf(os.Stderr, "could not render plist: %v\n", err)
		os.Exit(1)
	}

	if out, err := exec.Command("launchctl", "load", path).CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "launchctl load failed: %s\n", out)
		os.Exit(1)
	}

	fmt.Printf("summon started and will run on login (logs: %s)\n", logPath)
}

func LaunchdStop() {
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

func LaunchdStatus() {
	// running = launchctl list returns exit 0
	running := exec.Command("launchctl", "list", plistLabel).Run() == nil

	if running {
		fmt.Println("status:  running")
	} else {
		fmt.Println("status:  stopped")
	}

	cfgPath, err := config.ConfigPath()
	if err == nil {
		fmt.Println("config: ", cfgPath)
	}

	cfg, err := config.Load()
	if err == nil {
		fmt.Printf("bindings: %d\n", len(cfg.Hotkeys))
	}
}
