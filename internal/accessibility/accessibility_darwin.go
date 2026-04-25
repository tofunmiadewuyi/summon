package accessibility

/*
#cgo LDFLAGS: -framework ApplicationServices
#include "accessibility_darwin.h"
*/
import "C"
import (
	"log"
	"os"
	"time"
)

// Confirm silently checks first. If not granted, shows the system prompt,
// then polls until granted and restarts via os.Exit(0) so launchd re-runs cleanly.
func Confirm() {
	if C.isAccessibilityEnabled(0) != 0 {
		return
	}
	C.isAccessibilityEnabled(1)
	log.Println("accessibility permission required — grant access in System Settings")
	for C.isAccessibilityEnabled(0) == 0 {
		time.Sleep(2 * time.Second)
	}
	os.Exit(0) // launchd will restart the process now that permission is granted
}
