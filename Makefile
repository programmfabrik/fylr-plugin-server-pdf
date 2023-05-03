ZIP_NAME ?= "server-pdf.zip"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build ## build all

build: ## build files in build folder
	cd html2pdf; GOOS=linux go build -o ../build/server-pdf/html2pdf-linux.exe
	cd html2pdf; GOOS=windows go build -o ../build/server-pdf/html2pdf-windows.exe
	cd html2pdf; GOOS=darwin GOARCH=amd64 go build -o ../build/server-pdf/html2pdf-darwin.exe
	cp -r manifest.master.yml build/server-pdf/manifest.yml

clean: ## clean
	rm -rf build

zip: build ## build zip file
	cd build && zip ${ZIP_NAME} -r server-pdf/

.PHONY: build
