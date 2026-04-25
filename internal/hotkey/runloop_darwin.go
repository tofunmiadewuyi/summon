package hotkey

/*
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"

// RunMainLoop runs the CoreFoundation run loop on the calling thread.
// Must be called from the main OS thread so dispatch_get_main_queue() works.
func RunMainLoop() {
	C.CFRunLoopRun()
}
