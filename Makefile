ZIP_NAME ?= "server-pdf.zip"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build ## build all

google-csv: ## download l10n csv from google
	make -C easydb-pdf-creator-plugin google_csv

build: ## build files in build folder
	mkdir -p build/server-pdf/webfrontend
	mkdir -p build/server-pdf/l10n
	cd html2pdf; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../build/server-pdf/html2pdf-linux-amd64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../build/server-pdf/html2pdf-linux-arm64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../build/server-pdf/html2pdf-windows-amd64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ../build/server-pdf/html2pdf-windows-arm64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../build/server-pdf/html2pdf-darwin-amd64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ../build/server-pdf/html2pdf-darwin-arm64.exe
	cp -r manifest.master.yml build/server-pdf/manifest.yml
	make -C easydb-pdf-creator-plugin build
	cp easydb-pdf-creator-plugin/build-info.json build/server-pdf
	cp -r easydb-pdf-creator-plugin/build/webfrontend/pdf-creator* build/server-pdf/webfrontend
	cp easydb-pdf-creator-plugin/l10n/pdf-creator.csv build/server-pdf/l10n

clean: ## clean
	rm -rf build || true
	rm -rf easydb-pdf-creator-plugin/build || true

test:
	cat apitest/payload.json | SERVER_PDF_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" build/server-pdf/html2pdf-darwin-arm64.exe > test.pdf

zip: build ## build zip file
	cd build && zip ${ZIP_NAME} -r server-pdf/

.PHONY: build test
