GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
BINARY_FOLDER=./build

Version := $(shell git describe --tags)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/updiver/cli/cmd/cli/cmd.Version=$(Version) -X github.com/updiver/cli/cmd/cli/cmd.GitCommit=$(GitCommit)"

.PHONY: all clean-all clean-binaries build-all build-% run-% run-tests run-tests-cover run-tests-cover-profile build-all-platforms

default:
	${MAKE} run-cli

all:
		$(MAKE) run

clean-all:
		$(MAKE) clean-binaries
clean-binaries:
		-rm -f $(BINARY_FOLDER)/*

build-all:
		$(MAKE) build-cli
build-%:
		$(GOBUILD) -ldflags=${LDFLAGS} -o $(BINARY_FOLDER)/$* -v ./cmd/$*
		chmod +x $(BINARY_FOLDER)/$*

build-all-platforms:
	mkdir -p build/
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags $(LDFLAGS) -o build/dumper-cli-darwin-m1-${Version} ./cmd/cli
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -o build/dumper-cli-amd64-${Version} ./cmd/cli
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags $(LDFLAGS) -o build/dumper-cli-darwin-${Version} ./cmd/cli
	GOARM=7 GOARCH=arm CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -o build/dumper-cli-arm-${Version} ./cmd/cli
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -o build/dumper-cli-arm64-${Version} ./cmd/cli
	GOOS=windows CGO_ENABLED=0 go build -a -ldflags $(LDFLAGS) -o build/dumper-cli-${Version}.exe ./cmd/cli

run-%:
		$(MAKE) build-$*
		$(BINARY_FOLDER)/$* --help

run-tests:
	$(GOCMD) test -count=1 ./... -v
run-tests-cover:
	$(GOCMD) test -count=1 ./... -v -cover
run-tests-cover-profile:
	$(GOCMD) test -count=1 ./... -v -coverprofile cover.out
	$(GOCMD) tool cover -html=cover.out