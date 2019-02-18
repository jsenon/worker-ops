#-----------------------------------------------------------------------------
# Global Variables
#-----------------------------------------------------------------------------

DOCKER_USER ?= $(DOCKER_USER)
DOCKER_PASS ?= 

DOCKER_BUILD_ARGS := --build-arg HTTP_PROXY=$(http_proxy) --build-arg HTTPS_PROXY=$(https_proxy)

APP_VERSION := latest
PACKAGE ?= $(shell go list ./... | grep config)
VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')


override LDFLAGS += \
  -X ${PACKAGE}.Version=${VERSION} \
  -X ${PACKAGE}.BuildDate=${BUILD_DATE} \
  -X ${PACKAGE}.GitCommit=${GIT_COMMIT} \


#-----------------------------------------------------------------------------
# BUILD
#-----------------------------------------------------------------------------

.PHONY: default build test publish build_local lint
default:  test lint build swagger

test:
	go test -v ./...
build_local:
	go mod tidy
	go build -ldflags '${LDFLAGS}' -o worker-ops  
build:
	go mod tidy
	docker build $(DOCKER_BUILD_ARGS) -t $(DOCKER_USER)/worker-ops:$(APP_VERSION)  .
lint:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	gometalinter ./... --exclude=vendor --deadline=60s

swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go generate


#-----------------------------------------------------------------------------
# PUBLISH
#-----------------------------------------------------------------------------

.PHONY: publish 

publish: 
	docker push $(DOCKER_USER)/worker-ops:$(APP_VERSION)

#-----------------------------------------------------------------------------
# CLEAN
#-----------------------------------------------------------------------------

.PHONY: clean 

clean:
	rm -rf worker-ops