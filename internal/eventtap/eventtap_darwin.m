#include "eventtap_darwin.h"
#import <AppKit/AppKit.h>
#include <ApplicationServices/ApplicationServices.h>
#include <string.h>

static void activateAppScript(const char *cName) {
  NSString *appName = [NSString stringWithUTF8String:cName];
  NSString *src = [NSString
      stringWithFormat:@"tell application \"%@\"\nreopen\nactivate\nend tell",
                       appName];
  NSAppleScript *script = [[NSAppleScript alloc] initWithSource:src];
  [script executeAndReturnError:nil];
}

static void activateAppNative(const char *cName) {
  NSString *appName = [NSString stringWithUTF8String:cName];
  dispatch_async(dispatch_get_main_queue(), ^{
    NSArray<NSRunningApplication *> *running =
        [[NSWorkspace sharedWorkspace] runningApplications];
    for (NSRunningApplication *app in running) {
      if ([app.localizedName caseInsensitiveCompare:appName] == NSOrderedSame) {
        [app unhide];
        if (@available(macOS 14.0, *)) {
          [app activateFromApplication:[[NSWorkspace sharedWorkspace]
                                           frontmostApplication]
                               options:0];
        } else {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
          [app activateWithOptions:NSApplicationActivateIgnoringOtherApps];
#pragma clang diagnostic pop
        }
        return;
      }
    }
  });
}

#define MAX_BINDINGS 64

typedef struct {
  uint16_t keycode;
  CGEventFlags modifiers;
  char appName[256];
  void *script; // bridged NSAppleScript*, pre-compiled at registration
} TapBinding;

static TapBinding gBindings[MAX_BINDINGS];
static int gBindingCount = 0;
static CFMachPortRef gTap = NULL;
static CGKeyCode gLastSwallowed = 0xFFFF;

static CGEventFlags const kModMask =
    kCGEventFlagMaskAlternate | kCGEventFlagMaskCommand |
    kCGEventFlagMaskControl | kCGEventFlagMaskShift;

void clearTapBindings(void) {
  for (int i = 0; i < gBindingCount; i++) {
    if (gBindings[i].script) {
      CFBridgingRelease(gBindings[i].script);
      gBindings[i].script = NULL;
    }
  }
  gBindingCount = 0;
}

void addTapBinding(uint16_t keycode, uint32_t modifiers, const char *appName) {
  if (gBindingCount >= MAX_BINDINGS)
    return;
  gBindings[gBindingCount].keycode = keycode;
  gBindings[gBindingCount].modifiers = (CGEventFlags)modifiers;
  strlcpy(gBindings[gBindingCount].appName, appName, 256);
  // Pre-compile the AppleScript once; executeAndReturnError: at activation
  // time is pure IPC with no parsing overhead.
  NSString *name = [NSString stringWithUTF8String:appName];
  NSString *src = [NSString stringWithFormat:
      @"tell application \"%@\"\nreopen\nactivate\nend tell", name];
  NSAppleScript *script = [[NSAppleScript alloc] initWithSource:src];
  gBindings[gBindingCount].script = (void *)CFBridgingRetain(script);
  gBindingCount++;
}

static CGEventRef tapCallback(CGEventTapProxy proxy, CGEventType type,
                              CGEventRef event, void *refcon) {
  if (type == kCGEventTapDisabledByTimeout ||
      type == kCGEventTapDisabledByUserInput) {
    CGEventTapEnable(gTap, true);
    return event;
  }

  CGKeyCode keycode =
      (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);

  if (type == kCGEventKeyUp) {
    if (keycode == gLastSwallowed) {
      gLastSwallowed = 0xFFFF;
      return NULL;
    }
    return event;
  }

  CGEventFlags mods = CGEventGetFlags(event) & kModMask;

  for (int i = 0; i < gBindingCount; i++) {
    if (gBindings[i].keycode == keycode && gBindings[i].modifiers == mods) {
      gLastSwallowed = keycode;
      NSAppleScript *script = (__bridge NSAppleScript *)gBindings[i].script;
      dispatch_async(dispatch_get_global_queue(QOS_CLASS_USER_INTERACTIVE, 0), ^{
        [script executeAndReturnError:nil];
      });
      return NULL;
    }
  }

  gLastSwallowed = 0xFFFF;
  return event;
}

// — Capture ----------------------------------------------------------------

static volatile CGKeyCode gCaptureKeycode = 0;
static volatile uint64_t gCaptureModifiers = 0;
static volatile int gCaptureDone = 0;
static CFRunLoopRef gCaptureLoop = NULL;

static CGEventRef captureCallback(CGEventTapProxy proxy, CGEventType type,
                                  CGEventRef event, void *refcon) {
  if (type != kCGEventKeyDown)
    return event;

  CGEventFlags mods = CGEventGetFlags(event) & kModMask;
  if (mods == 0)
    return event;

  gCaptureKeycode =
      (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
  gCaptureModifiers = (uint64_t)mods;
  gCaptureDone = 1;
  CFRunLoopStop(gCaptureLoop);
  return event;
}

int captureNextCombo(uint16_t *keycode, uint64_t *modifiers) {
  gCaptureDone = 0;
  gCaptureLoop = CFRunLoopGetCurrent();

  CFMachPortRef tap = CGEventTapCreate(
      kCGHIDEventTap, kCGHeadInsertEventTap, kCGEventTapOptionDefault,
      CGEventMaskBit(kCGEventKeyDown), captureCallback, NULL);
  if (!tap)
    return 0;

  CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(NULL, tap, 0);
  CFRunLoopAddSource(gCaptureLoop, source, kCFRunLoopDefaultMode);
  CFRunLoopRun();
  CFRunLoopRemoveSource(gCaptureLoop, source, kCFRunLoopDefaultMode);
  CFRelease(source);
  CGEventTapEnable(tap, false);
  CFRelease(tap);

  if (!gCaptureDone)
    return 0;
  *keycode = gCaptureKeycode;
  *modifiers = gCaptureModifiers;
  return 1;
}

// — Persistent listener ----------------------------------------------------

void startEventTap(void) {
  if (gTap)
    return;

  gTap = CGEventTapCreate(
      kCGHIDEventTap, kCGHeadInsertEventTap, kCGEventTapOptionDefault,
      CGEventMaskBit(kCGEventKeyDown) | CGEventMaskBit(kCGEventKeyUp),
      tapCallback, NULL);
  if (!gTap)
    return;

  // CFMachPortCreateRunLoopSource wraps the raw Mach port into a run-loop
  // source. CFRunLoopAddSource attaches it to the main thread's run loop. Once
  // CFRunLoopRun() is spinning (in runloop_darwin.go), keydown events wake the
  // loop and invoke tapCallback.
  CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(NULL, gTap, 0);
  CFRunLoopAddSource(CFRunLoopGetMain(), source, kCFRunLoopCommonModes);
  CFRelease(source);
}
