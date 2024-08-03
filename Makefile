init:
	mkdir -p target
	mkdir -p package

build: init
	GOOS=darwin GOARCH=amd64 go build -o target/freddiebear.amd64
	GOOS=darwin GOARCH=arm64 go build -o target/freddiebear.arm64

package: build
	zip -r package/freddiebear.amd64.zip target/freddiebear.amd64
	zip -r package/freddiebear.arm64.zip target/freddiebear.arm64

workflow: init
	zip -r package/Freddiebear.alfredworkflow info.plist icon.png

.PHONY: workflow init
