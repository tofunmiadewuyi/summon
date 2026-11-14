VERSION ?= dev
LD := -ldflags="-s -w -X main.version=$(VERSION)"
APP := Summon.app
CGO_LDFLAGS := -Wl,-no_warn_duplicate_libraries
export CGO_LDFLAGS

build:
	go build $(LD) -o summon ./cmd

build_run:
	go build $(LD) -o summon ./cmd && ./summon run

build_help:
	go build $(LD) -o summon ./cmd && ./summon help

build_config:
	go build $(LD) -o summon ./cmd && ./summon config

build_size:
	go build $(LD) -o summon ./cmd && ls -lh summon

# Build Summon.app (binary + Info.plist + icon), ad-hoc signed
app:
	rm -rf $(APP)
	mkdir -p $(APP)/Contents/MacOS $(APP)/Contents/Resources
	cp assets/Info.plist $(APP)/Contents/
	cp assets/AppIcon.icns $(APP)/Contents/Resources/
	go build $(LD) -o $(APP)/Contents/MacOS/summon ./cmd
	codesign --force --sign - $(APP)

# Install the local build to ~/Applications, link it into ~/.local/bin, restart the agent
replace_build: app
	rm -rf ~/Applications/$(APP)
	mkdir -p ~/Applications ~/.local/bin
	cp -R $(APP) ~/Applications/
	ln -sf ~/Applications/$(APP)/Contents/MacOS/summon ~/.local/bin/summon
	~/.local/bin/summon start

release:
	@latest=$$(git tag --sort=-version:refname | head -1); \
	if [ -z "$$latest" ]; then next="v0.0.1"; \
	else \
		patch=$$(echo $$latest | cut -d. -f3); \
		prefix=$$(echo $$latest | cut -d. -f1-2); \
		next="$$prefix.$$((patch + 1))"; \
	fi; \
	echo "Tagging $$next"; \
	git tag $$next && git push origin $$next

