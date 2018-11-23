VERSION=v1.2.2

.PHONY: release
release: tag build_darwin build_linux

.PHONY: tag
tag:
	git tag ${VERSION}

.PHONY: build_darwin
build_darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags \
			"-X github.com/sethpollack/dockerbox/version.Version=${VERSION} -X github.com/sethpollack/dockerbox/version.Commit=`git log -n 1 --pretty=format:"%h"`" \
			-o dockerbox_${VERSION}_darwin_amd64

.PHONY: build_linux
build_linux:
	GOOS=linux GOARCH=amd64 go build -ldflags \
			"-X github.com/sethpollack/dockerbox/version.Version=${VERSION} -X github.com/sethpollack/dockerbox/version.Commit=`git log -n 1 --pretty=format:"%h"`" \
			-o dockerbox_${VERSION}_linux_amd64
