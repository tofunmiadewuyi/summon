// Package hotkey wires together accessibility, event tap, and activation.
package hotkey

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/tofunmiadewuyi/summon/internal/activate"
	"github.com/tofunmiadewuyi/summon/internal/config"
	"github.com/tofunmiadewuyi/summon/internal/eventtap"
	serror "github.com/tofunmiadewuyi/summon/internal/serrors"
)

func Register(cfg *config.Config) {
	eventtap.SetCallback(activate.Focus)
	eventtap.Clear()
	seen := map[string]bool{}
	for _, hk := range cfg.Hotkeys {
		if seen[hk.Keys] {
			fmt.Printf("warning: duplicate hotkey %s — only the last binding will work", hk.Keys)
		}
		seen[hk.Keys] = true
		fmt.Printf("registering: %s → %s\n", hk.Keys, hk.App)
		parsed, err := parser(hk.Keys)
		if err != nil {
			fmt.Printf("invalid hotkey combo: %s for summoning %s\n", hk.Keys, hk.App)
			continue
		}
		if !appExists(hk.App) {
			fmt.Printf("app: %s is not installed on this computer\n", hk.App)
			continue
		}
		eventtap.Register(parsed, hk.App)
	}
	eventtap.Start()
}

func CaptureHotkey() (string, error) {
	fmt.Print("Press your desired hotkey combo... ")
	combo, err := eventtap.CaptureCombo()
	if err != nil {
		return "", err
	}
	fmt.Println()
	return combo, nil
}

func parser(keys string) ([]string, error) {
	modifiers := map[string]string{
		"option": "alt",
		"cmd":    "cmd",
		"ctrl":   "ctrl",
		"shift":  "shift",
	}

	if keys == "" {
		return nil, serror.ErrInvalidKeyCombo
	}

	parts := strings.Split(keys, "+")
	if len(parts) <= 1 {
		return nil, serror.ErrInvalidKeyCombo
	}

	var parsed []string

	val, ok := modifiers[parts[0]]
	if ok {
		parsed = append(parsed, val)
	} else {
		return nil, serror.ErrInvalidKeyCombo
	}

	for _, key := range parts[1:] {
		val, ok := modifiers[key]
		if ok {
			parsed = append(parsed, val)
			continue
		}
		if eventtap.ValidKeyName(key) {
			parsed = append(parsed, key)
		} else {
			return nil, serror.ErrInvalidKeyCombo
		}
	}

	return parsed, nil
}

func appExists(app string) bool {
	script := fmt.Sprintf(`id of app "%s"`, app)
	err := exec.Command("osascript", "-e", script).Run()
	return err == nil
}
