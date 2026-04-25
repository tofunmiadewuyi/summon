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
		// Also register with right option (altgr) so both option keys work.
		if ropt := withRightOption(parsed); ropt != nil {
			hook.Register(hook.KeyDown, ropt, func(e hook.Event) {
				summonApp(hk.App)
			})
		}
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

// withRightOption returns a copy of keys with "alt" replaced by "altgr" (right option).
// Returns nil if the combo doesn't use alt.
func withRightOption(keys []string) []string {
	found := false
	out := make([]string, len(keys))
	for i, k := range keys {
		if k == "alt" {
			out[i] = "altgr"
			found = true
		} else {
			out[i] = k
		}
	}
	if !found {
		return nil
	}
	return out
}

func summonApp(app string) error {
	script := fmt.Sprintf(`
		tell application "%s"
			reopen
			activate
		end tell
		tell application "System Events"
			tell process "%s" to set frontmost to true
		end tell`, app, app)
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

// carbonKeycodes maps macOS Carbon virtual keycodes to key names accepted by parser().
var carbonKeycodes = map[uint16]string{
	0: "a", 1: "s", 2: "d", 3: "f", 4: "h", 5: "g", 6: "z", 7: "x",
	8: "c", 9: "v", 11: "b", 12: "q", 13: "w", 14: "e", 15: "r",
	16: "y", 17: "t", 32: "u", 34: "i", 31: "o", 35: "p", 37: "l",
	38: "j", 39: "'", 40: "k", 41: ";", 43: ",", 44: "/", 45: "n",
	46: "m", 47: ".", 18: "1", 19: "2", 20: "3", 21: "4", 23: "5",
	22: "6", 26: "7", 28: "8", 25: "9", 29: "0",
	49: "space", 36: "enter", 48: "tab", 51: "delete", 53: "escape", 50: "`",
	123: "left", 124: "right", 125: "down", 126: "up",
	96: "f5", 97: "f6", 98: "f7", 99: "f3", 100: "f8", 101: "f9",
	103: "f11", 109: "f10", 111: "f12", 118: "f4", 120: "f2", 122: "f1",
}

// CaptureHotkey blocks until the user presses a combo with at least one modifier
// plus one non-modifier key. Returns the combo in summon format (e.g., "option+f").
func CaptureHotkey() (string, error) {
	fmt.Print("Press your desired hotkey combo... ")

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
			keyName, ok := carbonKeycodes[e.Rawcode]
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
