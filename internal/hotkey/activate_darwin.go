package hotkey

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit -framework ApplicationServices
#include "activate_darwin.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func deleteLastChar() {
	C.deleteLastChar()
}

func nativeActivate(app string) {
	cs := C.CString(app)
	defer C.free(unsafe.Pointer(cs))
	C.activateAppNative(cs)
}

func scriptActivate(app string) {
	cs := C.CString(app)
	defer C.free(unsafe.Pointer(cs))
	C.activateAppScript(cs)
}
