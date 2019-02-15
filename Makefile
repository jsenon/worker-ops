#-----------------------------------------------------------------------------
# Global Variables
#-----------------------------------------------------------------------------

DOCKER_USER ?= $(DOCKER_USER)
DOCKER_PASS ?= 

DOCKER_BUILD_ARGS := --build-arg HTTP_PROXY=$(http_proxy) --build-arg HTTPS_PROXY=$(https_proxy)

APP_VERSION := latest

#-----------------------------------------------------------------------------
# BUILD
#-----------------------------------------------------------------------------

.PHONY: default build test publish build_local lint
default: depend test lint build swagger

depend:
	go get -u -v golang.org/x/vgo
	vgo generate ./...
test:
	go test -v ./...
build_local:
	vgo build ./...
build:
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