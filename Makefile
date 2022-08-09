.PHONY: all dep test lint build

WORKSPACE ?= $$(pwd)

GO_PKG_LIST := $(shell go list ./... | grep -v /vendor/)
SDK_VERSION := $(shell go list -m github.com/Axway/agent-sdk | cut -d ' ' -f 2 | cut -c 2-)

lint:
	@golint -set_exit_status ${GO_PKG_LIST}

dep:
	@echo "Resolving go package dependencies"
	@go mod tidy
	@echo "Package dependencies completed"

update-sdk:
	@echo "Updating SDK dependencies"
	@export GOFLAGS="-mod=mod" && go get "github.com/Axway/agent-sdk@main"


build: dep
	@echo "building discovery agent with sdk version $(SDK_VERSION)"
	export CGO_ENABLED=0
	export TIME=`date +%Y%m%d%H%M%S`
	@go build \
		-ldflags="-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildTime=${TIME}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildVersion=${VERSION}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildCommitSha=${COMMIT_ID}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.SDKBuildVersion=${SDK_VERSION}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentDescription=Layer7 Discovery Agent' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentName=Layer7DiscoveryAgent'" \
		-o bin/discovery main.go
	@echo "discovery agent binary placed at bin/discovery"

