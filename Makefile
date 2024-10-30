# set Go specific env variables
# export GOARCH = amd64
# export GOOS = linux
# export CGO_ENABLED = 0
# export GO111MODULE = on
# export GOPROXY = direct

# TAG=v1.0.0
# ARCH=amd64
# OS=linux
# GIT_COMMIT=git-$(shell git rev-parse --short HEAD)
BINARY_NAME=dbctl

# VERSION := $(GIT_TAG)
# LDFLAGS:=-X main.GitVersion=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)
# BUILD_FLAGS=-ldflags '$(LDFLAGS)'

# Build go app
build:
	echo "Building dbctl"
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
