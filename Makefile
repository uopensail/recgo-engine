CURDIR:=$(shell pwd)
OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)
.PHONY: build clean run

GITCOMMITHASH := $(shell git rev-parse --short=7 HEAD)
GITBRANCHNAME := $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

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

build-docker-image: build
	mv dist/recgo-engine dist/main
	docker build \
		--build-arg GIT_HASH=$(GITCOMMITHASH) \
		--build-arg GIT_TAG=$(GITBRANCHNAME) \
		--build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		-t $(ACR_CONTAINER_REGISTRY_SERVER)/uopensail/${PROJECT_NAME}:${GITBRANCHNAME}-${GITCOMMITHASH} \
		-t $(ACR_CONTAINER_REGISTRY_SERVER)/uopensail/${PROJECT_NAME}:latest \
		.

push-docker-image: build-docker-image
	docker push $(ACR_CONTAINER_REGISTRY_SERVER)/uopensail/${PROJECT_NAME}:${GITBRANCHNAME}-${GITCOMMITHASH}
	docker push $(ACR_CONTAINER_REGISTRY_SERVER)/uopensail/${PROJECT_NAME}:latest

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