# fylr plugins are built by fylr-build-plugin, the build driver that knows how
# a fylr plugin is put together (compile, assemble build/, zip, seal, loca).
# This Makefile is a thin shim for muscle memory — all logic lives in the
# tool. @latest always resolves the tool's newest release, so plugins pick up
# fixes without being touched; an incompatible tool change would come as a new
# major version (import path .../v2), which is the only event that changes
# this line.
#
# Tools needed:
#   go       runs fylr-build-plugin + compiles html2pdf — https://go.dev/dl/
FYLR_BUILD_PLUGIN ?= go run github.com/programmfabrik/fylr-build-plugin@latest

# The tool itself reads NO environment variables — everything is passed as
# flags. The release workflow's RELEASE_TAG env is translated into a flag
# right here. The zip's NAME is not passed: fylr-build-plugin always names
# it <plugin.name>.zip, the naming rule for every fylr plugin release.
RELEASE_FLAGS = $(if $(RELEASE_TAG),-release "$(RELEASE_TAG)")

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build ## build all

build: ## build the plugin into build/<name>/ — loadable by fylr via plugin.paths
	$(FYLR_BUILD_PLUGIN) build $(RELEASE_FLAGS)

zip: ## build the release zip
	$(FYLR_BUILD_PLUGIN) zip $(RELEASE_FLAGS)

check: ## validate the build tree against the manifest
	$(FYLR_BUILD_PLUGIN) check

clean: ## clean build files
	$(FYLR_BUILD_PLUGIN) clean

test: ## render apitest/payload.json to test.pdf with the local darwin binary
	cat apitest/payload.json | SERVER_PDF_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" build/server-pdf/html2pdf-darwin-arm64.exe > test.pdf

.PHONY: help all build zip check clean test
