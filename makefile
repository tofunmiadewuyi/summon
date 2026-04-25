build:
	go build -ldflags="-s -w -X main.version=dev" -o summon ./cmd &&  ./summon start

build_help:
	go build -ldflags="-s -w -X main.version=dev" -o summon ./cmd &&  ./summon help

build_config:
	go build -ldflags="-s -w -X main.version=dev" -o summon ./cmd &&  ./summon config

build_size:
	go build -ldflags="-s -w -X main.version=dev" -o summon ./cmd && ls -lh summon

release:
	@latest=$$(git tag --sort=-version:refname | head -1); \
	if [ -z "$$latest" ]; then next="v0.1.0"; \
	else \
		patch=$$(echo $$latest | cut -d. -f3); \
		prefix=$$(echo $$latest | cut -d. -f1-2); \
		next="$$prefix.$$((patch + 1))"; \
	fi; \
	echo "Tagging $$next"; \
	git tag $$next && git push origin $$next

