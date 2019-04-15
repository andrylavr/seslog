GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
export GOPATH=${CURDIR}
LDFLAGS=-ldflags "-s -w -X main.GitBranch=${GIT_BRANCH} -X main.GitCommit=${GIT_COMMIT} -X main.BuildDate=`date -u +%Y-%m-%d.%H:%M:%S`"
CGO_ENABLED=0

build:
	@[ -d build ] || mkdir -p build
	go build ${LDFLAGS} -o build/seslog-server src/cmd/seslog-server/main.go
	@file  build/seslog-server
	@du -h build/seslog-server

d:
	docker-compose -f dockerfiles/docker-compose.yml rm --force
	docker-compose -f dockerfiles/docker-compose.yml up --build

f:
	gofmt -l -s -w `find . -type f -name '*.go' -not -path "./*/vendor/*"`
	goimports -l -w `find . -type f -name '*.go' -not -path "./*/vendor/*"`

.PHONY: build