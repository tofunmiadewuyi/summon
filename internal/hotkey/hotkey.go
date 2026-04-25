// Package hotkey defines the key registration behvaiour and app retrieval
package hotkey

import (
	"fmt"
	"os/exec"
	"strings"

	hook "github.com/robotn/gohook"
	"github.com/tofunmiadewuyi/summon/internal/config"
	serror "github.com/tofunmiadewuyi/summon/internal/serrors"
)

var hookRunning bool

func Register(cfg *config.Config) {
	// de-register all hooks
	if hookRunning {
		hook.End()
	}

	// confirm permissions
	confirmAccesibility()

	// register hooks
	hookRunning = true
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
		exists := appExists(hk.App)
		if !exists {
			fmt.Printf("app : %s is not installed on this computer", hk.App)
			continue
		}

		hook.Register(hook.KeyDown, parsed, func(e hook.Event) {
			summonApp(hk.App)
		})
	}

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()
}

func parser(keys string) ([]string, error) {
	modifiers := map[string]string{
		"option": "alt",
		"cmd":    "cmd",
		"ctrl":   "ctrl",
		"shift":  "shift",
	}

	// must not be empty
	if keys == "" {
		return nil, serror.ErrInvalidKeyCombo
	}

	parts := strings.Split(keys, "+")
	// must be a combo
	if len(parts) <= 1 {
		return nil, serror.ErrInvalidKeyCombo
	}

	var parsed []string

	// must start with a modifier
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
		_, ok = hook.Keycode[key]
		if ok {
			parsed = append(parsed, key)
		} else {
			return nil, serror.ErrInvalidKeyCombo
		}
	}

	return parsed, nil
}

func summonApp(app string) error {
	script := fmt.Sprintf(`tell application "%s" to activate`, app)
	return exec.Command("osascript", "-e", script).Run()
}

func appExists(app string) bool {
	script := fmt.Sprintf(`id of app "%s"`, app)
	err := exec.Command("osascript", "-e", script).Run()
	if err != nil {
		// does not exist
		return false
	}
	return true
}

// modifierRawcodes maps macOS Carbon virtual keycodes to summon modifier names.
var modifierRawcodes = map[uint16]string{
	58: "option", 61: "option", // Left/Right Option
	55: "cmd", 54: "cmd",       // Left/Right Cmd
	59: "ctrl", 62: "ctrl",     // Left/Right Ctrl
	56: "shift", 60: "shift",   // Left/Right Shift
}

// modifierOrder defines canonical ordering for building combo strings.
var modifierOrder = []string{"ctrl", "shift", "option", "cmd"}

// CaptureHotkey blocks until the user presses a combo with at least one modifier
// plus one non-modifier key. Returns the combo in summon format (e.g., "option+f").
func CaptureHotkey() (string, error) {
	fmt.Print("Press your desired hotkey combo... ")

	// Build reverse lookup: rawcode → key name (same names parser() accepts).
	reverseKeycode := make(map[uint16]string, len(hook.Keycode))
	for name, code := range hook.Keycode {
		reverseKeycode[code] = name
	}

	s := hook.Start()
	defer hook.End()

	activeModifiers := map[string]struct{}{}

	for e := range s {
		switch e.Kind {
		case hook.KeyDown, hook.KeyHold:
			if mod, ok := modifierRawcodes[e.Rawcode]; ok {
				activeModifiers[mod] = struct{}{}
				continue
			}
			if len(activeModifiers) == 0 {
				continue
			}
			// Resolve main key name.
			keyName, ok := reverseKeycode[e.Rawcode]
			if !ok {
				continue
			}
			// Build combo in canonical modifier order.
			var parts []string
			for _, mod := range modifierOrder {
				if _, held := activeModifiers[mod]; held {
					parts = append(parts, mod)
				}
			}
			parts = append(parts, keyName)
			fmt.Println()
			return strings.Join(parts, "+"), nil

		case hook.KeyUp:
			if mod, ok := modifierRawcodes[e.Rawcode]; ok {
				delete(activeModifiers, mod)
			}
		}
	}

	return "", serror.ErrCaptureFailed
}
