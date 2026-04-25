#include <ApplicationServices/ApplicationServices.h>
#include "accessibility_darwin.h"

int isAccessibilityEnabled(int prompt) {
	CFStringRef  keys[] = { kAXTrustedCheckOptionPrompt };
	CFBooleanRef vals[] = { prompt ? kCFBooleanTrue : kCFBooleanFalse };
	CFDictionaryRef opts = CFDictionaryCreate(NULL,
		(const void**)keys, (const void**)vals, 1,
		&kCFTypeDictionaryKeyCallBacks,
		&kCFTypeDictionaryValueCallBacks);
	int result = AXIsProcessTrustedWithOptions(opts);
	CFRelease(opts);
	return result;
}
