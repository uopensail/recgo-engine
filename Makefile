CURDIR:=$(shell pwd)
OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)
.PHONY: build clean run

GITCOMMITHASH := $(shell git rev-parse --short=7 HEAD)
GITBRANCHNAME := $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
GIT_TAG := $(shell  git describe --tags --abbrev=0 2>/dev/null || echo "untagged")

GOLDFLAGS += -X handler.__GITCOMMITINFO__=$(GITCOMMITHASH).${GITBRANCHNAME}
GOFLAGS = -ldflags "$(GOLDFLAGS)"

PUBLISHDIR=${CURDIR}/dist
PROJECT_NAME=recgo-engine

all: dev

clean:
	rm -rf ${PUBLISHDIR}/
	mkdir -pv ${PUBLISHDIR}/conf

build: clean
	go build -o ${PUBLISHDIR}/${PROJECT_NAME} $(GOFLAGS)

docker-image: build
	docker build -t ${PROJECT_NAME}:$@ .
	docker build \
		--build-arg GIT_HASH=$GIT_HASH \
		--build-arg GIT_TAG=$GIT_TAG \
		--build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		-t registry.uopensail.com/${PROJECT_NAME}:${GIT_TAG}-${GIT_HASH} \
		-t registry.uopensail.com/${PROJECT_NAME}:latest \
		.

docker-image-push: docker-image
	echo "$(DOCKER_REGISTRY_PASSWORD)" | docker login registry.uopensail.com -u $(DOCKER_REGISTRY_USER) --password-stdin
	docker push registry.uopensail.com/${PROJECT_NAME}:${GIT_TAG}-${GIT_HASH}
	docker push registry.uopensail.com/${PROJECT_NAME}:latest

pre: build
	cp -aRf conf/$@/* ${PUBLISHDIR}/conf
	cp run.sh ${PUBLISHDIR}/

prod: build
	cp -aRf conf/$@/* ${PUBLISHDIR}/conf
	cp run.sh ${PUBLISHDIR}/

dev: build
	cp -aRf conf/$@/* ${PUBLISHDIR}/conf
	cp run.sh ${PUBLISHDIR}/
test: build
	cp -aRf conf/$@/* ${PUBLISHDIR}/conf
	cp run.sh ${PUBLISHDIR}/
local: build
	cp -aRf conf/$@/* ${PUBLISHDIR}/conf
	cd dist && ./${PROJECT_NAME}