package accessibility

/*
#cgo LDFLAGS: -framework ApplicationServices
#include "accessibility_darwin.h"
*/
import "C"
import (
	"log"
	"time"
)

func Confirm() {
	if C.isAccessibilityEnabled(0) == 0 {
		C.isAccessibilityEnabled(1)
		log.Fatal("accessibility permission required. grant access in System Settings then restart summon.")
	}
}

func Wait() {
	if C.isAccessibilityEnabled(1) != 0 {
		return
	}
	log.Println("waiting for accessibility permission in System Settings...")
	for C.isAccessibilityEnabled(0) == 0 {
		time.Sleep(2 * time.Second)
	}
	log.Println("accessibility granted")
}
