# summon

Summon any app instantly with a global hotkey — from anywhere, any desktop.

## Why

I tried AeroSpace for window management but it mangled apps like WhatsApp that have a minimum width — tiling just doesn't play nice with them. I noticed some performance issues too. Then I tried Raycast, but it only focuses apps on the current desktop, which defeats the whole point. I want to hit a key and have my app — wherever it is.

So I built summon. It uses osascript under the hood, keeping it as macOS-native as possible. No tiling, no launcher UI. Just a hotkey and your app.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/tofunmiadewuyi/summon/main/install.sh | bash
```

Or download manually from [releases](https://github.com/tofunmiadewuyi/summon/releases).

## Usage

```bash
summon add       # add a new hotkey binding interactively
summon start     # install summon as a login item and start it
summon stop      # stop summon and remove it from login items
summon status    # show whether summon is running and its config
summon config    # print the config file path and contents
summon upgrade   # upgrade to the latest release
```

## Adding a binding

```
$ summon add
Press your desired hotkey combo...
Detected: option+s
Add binding for option+s? [y/N] y
Enter app name: Slack
Added: option+s → Slack
```

If a binding isn't working, check that the app name is exactly right:

```bash
osascript -e 'name of app "Safari"'
```

## Config

Bindings are stored at `~/.config/summon/config.toml`:

```toml
[[binding]]
keys = "option+f"
app = "Finder"

[[binding]]
keys = "option+s"
app = "Slack"
```

Summon hot-reloads the config — edit the file and changes apply immediately without restarting.

## Supported modifiers

`option`, `cmd`, `ctrl`, `shift`

## Requirements

- macOS
- Accessibility permission (summon will prompt on first run)
