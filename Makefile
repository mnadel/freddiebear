clean:
	rm -rf target
	rm -rf package

init:
	mkdir -p target
	mkdir -p package

build: init
	GOOS=darwin GOARCH=amd64 go build -o target/freddiebear.amd64
	GOOS=darwin GOARCH=arm64 go build -o target/freddiebear.arm64

workflow:
	zip -r package/Freddiebear.alfredworkflow info.plist icon.png

package: build workflow
	gzip -c target/freddiebear.amd64 > package/freddiebear.amd64.gz
	gzip -c target/freddiebear.arm64 > package/freddiebear.arm64.gz
