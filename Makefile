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
	$(eval WFVER := $(shell git for-each-ref --sort=creatordate --format '%(refname)' refs/tags | tail -1 | cut -d/ -f3))
	grep -q $(WFVER) info.plist || (echo "Update version in info.plist to $(WFVER)"; exit 1)
	zip -r package/Freddiebear.alfredworkflow info.plist icon.png download.sh

package: build workflow
	gzip -c target/freddiebear.amd64 > package/freddiebear.amd64.gz
	gzip -c target/freddiebear.arm64 > package/freddiebear.arm64.gz
