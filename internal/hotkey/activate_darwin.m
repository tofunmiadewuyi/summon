#import <AppKit/AppKit.h>
#include <ApplicationServices/ApplicationServices.h>

void deleteLastChar(void) {
  CGEventRef down = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)51, true);
  CGEventSetFlags(down, 0);
  CGEventPost(kCGHIDEventTap, down);
  CFRelease(down);
  CGEventRef up = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)51, false);
  CGEventSetFlags(up, 0);
  CGEventPost(kCGHIDEventTap, up);
  CFRelease(up);
}

static BOOL hasFocusedTextField(void) {
  NSRunningApplication *front =
      [[NSWorkspace sharedWorkspace] frontmostApplication];
  if (!front)
    return NO;

  AXUIElementRef appEl = AXUIElementCreateApplication(front.processIdentifier);
  AXUIElementRef focused = NULL;
  AXError err = AXUIElementCopyAttributeValue(
      appEl, kAXFocusedUIElementAttribute, (CFTypeRef *)&focused);
  CFRelease(appEl);
  if (err != kAXErrorSuccess || !focused)
    return NO;

  CFStringRef role = NULL;
  err = AXUIElementCopyAttributeValue(focused, kAXRoleAttribute,
                                      (CFTypeRef *)&role);
  CFRelease(focused);
  if (err != kAXErrorSuccess || !role)
    return NO;

  BOOL isText = CFEqual(role, kAXTextFieldRole) ||
                CFEqual(role, kAXTextAreaRole) ||
                CFEqual(role, CFSTR("AXComboBox")) ||
                CFEqual(role, CFSTR("AXSearchField"));
  CFRelease(role);
  return isText;
}

void activateAppNative(const char *cName) {
  NSString *appName = [NSString stringWithUTF8String:cName];
  dispatch_async(dispatch_get_main_queue(), ^{
    if (hasFocusedTextField()) {
      deleteLastChar();
    }
    NSArray<NSRunningApplication *> *running =
        [[NSWorkspace sharedWorkspace] runningApplications];
    for (NSRunningApplication *app in running) {
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

void activateAppScript(const char *cName) {
  NSString *appName = [NSString stringWithUTF8String:cName];
  NSString *src = [NSString
      stringWithFormat:@"tell application \"%@\"\nreopen\nactivate\nend tell",
                       appName];
  NSAppleScript *script = [[NSAppleScript alloc] initWithSource:src];
  [script executeAndReturnError:nil];
}
