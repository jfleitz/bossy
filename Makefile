MAKEFLAGS = -j1
SHELL := /bin/bash

# default target

all: build



###################
## BUILD SECTION ##
###################
export GO111MODULE := on
export CGO_ENABLED := 0
export GOPROXY     := direct
export GOSUMDB     := off

build:
	@echo "Building for local consumption"
	@go build

build_rpi:
	@echo "Building for Raspberry PI"
	@env GOOS=linux GOARCH=arm GOARM=7 go build


godeps:
	@echo "Getting go modules"
	@go mod download

clean: ## clean build output
	@rm -rf bin/