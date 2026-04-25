package eventtap

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit -framework ApplicationServices -framework CoreFoundation
#include "eventtap_darwin.h"
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"strings"
	"unsafe"

	serror "github.com/tofunmiadewuyi/summon/internal/serrors"
)

// carbonKeycodes maps macOS Carbon virtual keycodes to summon key names.
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

var reverseCarbon map[string]uint16

func init() {
	reverseCarbon = make(map[string]uint16, len(carbonKeycodes))
	for code, name := range carbonKeycodes {
		reverseCarbon[name] = code
	}
}

// ValidKeyName reports whether name is a recognised non-modifier key.
func ValidKeyName(name string) bool {
	_, ok := reverseCarbon[name]
	return ok
}

// CGEventFlags values for each modifier (stable Apple constants).
var modifierFlags = map[string]uint32{
	"shift":  0x00020000,
	"ctrl":   0x00040000,
	"alt":    0x00080000,
	"cmd":    0x00100000,
}

func Register(parsed []string, app string) {
	var mods uint32
	var keyName string
	for _, part := range parsed {
		if f, ok := modifierFlags[part]; ok {
			mods |= f
		} else {
			keyName = part
		}
	}
	keycode, ok := reverseCarbon[keyName]
	if !ok {
		return
	}
	cs := C.CString(app)
	defer C.free(unsafe.Pointer(cs))
	C.addTapBinding(C.uint16_t(keycode), C.uint32_t(mods), cs)
}

func Clear() {
	C.clearTapBindings()
}

func Start() {
	C.startEventTap()
}

var captureModifierNames = []struct {
	flag uint64
	name string
}{
	{0x00040000, "ctrl"},
	{0x00020000, "shift"},
	{0x00080000, "option"},
	{0x00100000, "cmd"},
}

func CaptureCombo() (string, error) {
	var keycode C.uint16_t
	var modifiers C.uint64_t
	if C.captureNextCombo(&keycode, &modifiers) == 0 {
		return "", serror.ErrCaptureFailed
	}
	keyName, ok := carbonKeycodes[uint16(keycode)]
	if !ok {
		return "", serror.ErrCaptureFailed
	}
	var parts []string
	for _, m := range captureModifierNames {
		if uint64(modifiers)&m.flag != 0 {
			parts = append(parts, m.name)
		}
	}
	parts = append(parts, keyName)
	return strings.Join(parts, "+"), nil
}
