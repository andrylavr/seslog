GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-s -w -X main.GitBranch=${GIT_BRANCH} -X main.GitCommit=${GIT_COMMIT} -X main.BuildDate=`date -u +%Y-%m-%d.%H:%M:%S`"
CGO_ENABLED=0

build:
	@[ -d build ] || mkdir -p build
	go build ${LDFLAGS} -o build/seslog-server cmd/seslog-server/main.go
	@file  build/seslog-server
	@du -h build/seslog-server

build-json2ch:
	@[ -d build ] || mkdir -p build
	go build ${LDFLAGS} -o build/seslog-json2ch cmd/seslog-json2ch/main.go
	@file  build/seslog-json2ch
	@du -h build/seslog-json2ch

br:
	go build --race -o build/seslog-server -v -ldflags "-s" cmd/seslog-server/main.go

d:
	docker-compose -f dockerfiles/docker-compose.yml rm --force
	docker-compose -f dockerfiles/docker-compose.yml up --build

f:
	gofmt -l -s -w `find . -type f -name '*.go' -not -path "./*/vendor/*"`
	goimports -l -w `find . -type f -name '*.go' -not -path "./*/vendor/*"`

zip: build
	@[ -d build/seslog/resources ] || mkdir -p build/seslog/resources
	@cp config.json build/seslog/config.json
	@cp resources/regexes.yaml build/seslog/resources/regexes.yaml
	@cp build/seslog-server build/seslog/seslog-server
	@cp -r package/systemd build/seslog
	cd build && rm -f seslog.zip
	cd build && zip -r seslog.zip seslog

install:
	@mkdir -p /opt/seslog/resources
	cp build/seslog-server /opt/seslog/seslog-server
	cp resources/regexes.yaml /opt/seslog/resources/regexes.yaml
	cp package/systemd/seslog-server.service /etc/systemd/system/seslog-server.service
	/bin/systemctl daemon-reload
	/bin/systemctl enable seslog-server
	service seslog-server status

.PHONY: build