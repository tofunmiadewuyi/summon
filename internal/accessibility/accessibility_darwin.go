package accessibility

/*
#cgo LDFLAGS: -framework ApplicationServices
#include "accessibility_darwin.h"
*/
import "C"
import "log"

// Confirm silently checks first; only shows the system prompt if not yet granted,
// then exits so the user can restart after granting.
func Confirm() {
	if C.isAccessibilityEnabled(0) != 0 {
		return
	}
	C.isAccessibilityEnabled(1)
	log.Fatal("accessibility permission required. grant access in System Settings then restart summon.")
}
