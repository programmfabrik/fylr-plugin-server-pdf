ZIP_NAME ?= "server-pdf.zip"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build ## build all

build: ## build files in build folder
	cd html2pdf; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../build/server-pdf/html2pdf-linux-amd64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../build/server-pdf/html2pdf-linux-arm64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../build/server-pdf/html2pdf-windows-amd64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ../build/server-pdf/html2pdf-windows-arm64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../build/server-pdf/html2pdf-darwin-amd64.exe
	cd html2pdf; CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ../build/server-pdf/html2pdf-darwin-arm64.exe
	cp -r manifest.master.yml build/server-pdf/manifest.yml

clean: ## clean
	rm -rf build

zip: build ## build zip file
	cd build && zip ${ZIP_NAME} -r server-pdf/

.PHONY: build
