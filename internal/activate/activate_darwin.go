// Package activate brings a named macOS application to the foreground.
package activate

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit -framework ApplicationServices
#include "activate_darwin.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func Focus(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.activateAppNative(cs)
}
