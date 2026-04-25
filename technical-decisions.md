# Technical Decisions

## Hotkey listening: gohook → native CGEventTap

Initially used [gohook](https://github.com/robotn/gohook) to listen for global hotkey combos. gohook wraps CGEventTap internally but exposes a Go-level API.

The problem: gohook intercepts events in listen-only mode — it can observe them but not consume them. When a hotkey fired, the key event still reached the focused app. Pressing `option+s` in a terminal would summon Slack AND type `ß` into the terminal. We worked around this by synthesising a backspace after activation (`deleteLastChar`), but that required checking whether a text field was focused first (via the Accessibility API) to avoid sending spurious backspaces elsewhere. It was a hack on top of a hack.

Replacing gohook with a direct `CGEventTapCreate` call gave us an active tap — one that sits at the front of the HID event pipeline and can return `NULL` to swallow events entirely. The hotkey combo never reaches any app. No stray character, no backspace, no Accessibility check needed.

The secondary benefit: one less dependency. gohook has had maintenance gaps and introduced CGo build friction on its own. Everything it did we now do directly.

## App activation: osascript subprocess → NSAppleScript → native Cocoa

**First version: `exec.Command("osascript", ...)`**

The initial implementation shelled out to the `osascript` binary to run an AppleScript that called `activate` on the target app. Simple to write, but each invocation spawned a new process, waited for it to start, compiled the AppleScript, executed it, and exited. The round trip was ~200ms — perceptible, especially compared to tools like Aerospace which switch instantly.

**Second version: in-process NSAppleScript**

Replaced the subprocess with `NSAppleScript` called directly via CGo. Same AppleScript source string, but compiled and executed in-process. Eliminated the process spawn overhead. Still compiled a script on every activation.

**Third version: native Cocoa (`NSRunningApplication`)**

Dropped AppleScript entirely. `[[NSWorkspace sharedWorkspace] runningApplications]` returns every running process directly. We walk the list, match by name, call `unhide` + `activateWithOptions:`. No script, no compilation, no interpreter — just Cocoa API calls.

This required two things to work correctly from a background daemon:

1. **`dispatch_async(dispatch_get_main_queue(), ...)`** — AppKit calls must run on the main thread. Dispatching to the main queue ensures this.
2. **`CFRunLoopRun()` on the main OS thread** — `dispatch_get_main_queue()` delivers blocks to the main thread's run loop. A CLI tool has no run loop by default, so we added one explicitly. Go's `runtime.LockOSThread()` pins the main goroutine to the main OS thread before `CFRunLoopRun()` is called, ensuring dispatch_async has somewhere to deliver work.
