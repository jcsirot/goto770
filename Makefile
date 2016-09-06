# Makefile for goto770
# -----------------------------------------------------------------
#
#        ENV VARIABLE
#
# -----------------------------------------------------------------

# go env vars
GO=$(firstword $(subst :, ,$(GOPATH)))
# list of pkgs for the project without vendor
PKGS=$(shell go list ./... | grep -v /vendor/)


# -----------------------------------------------------------------
#        Version
# -----------------------------------------------------------------

# version
VERSION=0.1.0
BUILDDATE=$(shell date -u '+%s')
BUILDHASH=$(shell git rev-parse --short HEAD)
VERSION_FLAG=-ldflags "-X main.Version=$(VERSION) -X main.GitHash=$(BUILDHASH) -X main.BuildStmp=$(BUILDDATE)"

# -----------------------------------------------------------------
#        Main targets
# -----------------------------------------------------------------

help:
	@echo
	@echo "----- BUILD ------------------------------------------------------------------------------"
	@echo "all                  clean and build the project"
	@echo "clean                clean the project"
	@echo "dependencies         download the dependencies"
	@echo "build                build all libraries and binaries"
	@echo "----- TESTS && LINT ----------------------------------------------------------------------"
	@echo "test                 test all packages"
	@echo "format               format all packages"
	@echo "lint                 lint all packages"
	@echo "----- OTHERS -----------------------------------------------------------------------------"
	@echo "help                 print this message"

all: clean build

clean:
	@go clean
	@rm -Rf .tmp
	@rm -Rf .DS_Store
	@rm -Rf *.log
	@rm -Rf *.out
	@rm -Rf *.lock
	@rm -Rf *.mem
	@rm -Rf *.test
	@rm -Rf build

dependencies:
	@echo
	@echo "----- DOWNLOADING -------------------------------------------------------------------------"
	@go get -u github.com/gorilla/mux
	@go get -u github.com/gorilla/context
	@go get -u github.com/urfave/negroni
	@go get -u github.com/urfave/cli
	@go get -u github.com/Sirupsen/logrus
	@go get -u gopkg.in/mgo.v2
	@go get -u github.com/tools/godep
	@go get -u github.com/golang/lint/golint
	@echo "----- DONE --------------------------------------------------------------------------------"

build: format
	@go build -v $(VERSION_FLAG) -o $(GO)/bin/goto770 goto770.go

format:
	@go fmt $(PKGS)

teardownTest:
	@$(shell docker kill handsongo-mongo-test 2&>/dev/null 1&>/dev/null)
	@$(shell docker rm handsongo-mongo-test 2&>/dev/null 1&>/dev/null)

setupTest: teardownTest
	@docker run -d --name handsongo-mongo-test -p "27017:27017" mongo:3.2

test: setupTest
	@export MONGODB_SRV=mongodb://$(DOCKER_IP)/spirits; go test -v $(PKGS); make teardownTest

bench:
	@go test -v -run TestSpiritHandlerGet -bench=. -memprofile=prof.mem github.com/sebastienfr/handsongo/web

benchTool: bench
	@echo "### TIP : type 'top 5' and 'list the first item'"
	@go tool pprof --alloc_space web.test prof.mem

lint:
	@golint dao/...
	@golint model/...
	@golint web/...
	@golint utils/...
	@golint ./.
	@go vet $(PKGS)

start:
	@handsongo -port 8020 -logl debug -logf text -statd 15s -db mongodb://$(DOCKER_IP)/spirits

stop:
	@killall handsongo

.PHONY: all test clean teardownTest setupTest
