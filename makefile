LD := -ldflags="-s -w -X main.version=dev"
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

replace_build:
	go build $(LD) -o $(shell which summon) ./cmd && summon start

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

