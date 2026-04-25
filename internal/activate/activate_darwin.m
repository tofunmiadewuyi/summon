#import <AppKit/AppKit.h>
#include <ApplicationServices/ApplicationServices.h>

void activateAppNative(const char* cName) {
	NSString* appName = [NSString stringWithUTF8String:cName];
	dispatch_async(dispatch_get_main_queue(), ^{
		NSArray<NSRunningApplication*>* running = [[NSWorkspace sharedWorkspace] runningApplications];
		for (NSRunningApplication* app in running) {
			if ([app.localizedName caseInsensitiveCompare:appName] == NSOrderedSame) {
				[app unhide];
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
				[app activateWithOptions:NSApplicationActivateIgnoringOtherApps];
#pragma clang diagnostic pop
				return;
			}
		}
	});
}

void activateAppScript(const char* cName) {
	NSString* appName = [NSString stringWithUTF8String:cName];
	NSString* src = [NSString stringWithFormat:
		@"tell application \"%@\"\nreopen\nactivate\nend tell", appName];
	NSAppleScript* script = [[NSAppleScript alloc] initWithSource:src];
	[script executeAndReturnError:nil];
}
