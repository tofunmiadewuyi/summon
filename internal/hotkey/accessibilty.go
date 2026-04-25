package hotkey

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>

int isAccessibilityEnabled(int prompt) {
    NSDictionary *options = @{(id)kAXTrustedCheckOptionPrompt: @(prompt)};
    return AXIsProcessTrustedWithOptions((CFDictionaryRef)options);
}
*/
import "C"
import (
	"log"
	"time"
)

func confirmAccesibility() bool {
    if C.isAccessibilityEnabled(0) == 0 {
        C.isAccessibilityEnabled(1) // trigger the prompt
        log.Fatal("accessibility permission required. grant access in System Settings then restart summon.")
    }
    return true
}

func waitForAccessibility() {
    // prompt once
    if C.isAccessibilityEnabled(1) != 0 {
        return // already granted
    }
    log.Println("waiting for accessibility permission in System Settings...")
    // poll silently
    for C.isAccessibilityEnabled(0) == 0 {
        time.Sleep(2 * time.Second)
    }
    log.Println("accessibility granted")
}
