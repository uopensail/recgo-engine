# ================================
# Project variables
# ================================
CURDIR       := $(shell pwd)
OS           := $(shell go env GOOS)
ARCH         := $(shell go env GOARCH)
PROJECT_NAME := recgo-engine
PUBLISHDIR   := $(CURDIR)/dist

# ================================
# Git build metadata
# ================================
GITCOMMIT    := $(shell git rev-parse --short=7 HEAD)
GITBRANCH    := $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
BUILD_TIME   := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# ================================
# Go build flags
# Inject git info & build time into main.go variables
# ================================
GOLDFLAGS    := -X main.gitCommitInfo=$(GITCOMMIT).$(GITBRANCH) -X main.buildTime=$(BUILD_TIME)
GOFLAGS      := -ldflags "$(GOLDFLAGS)"

# ================================
# Docker registry (optional)
# ================================
DOCKER_REGISTRY ?= $(ACR_CONTAINER_REGISTRY_SERVER)/uopensail

.PHONY: all clean build run docker-build docker-push \
        pre prod dev test local

# ================================
# Default target
# ================================
all: dev

# ================================
# Clean build artifacts
# ================================
clean:
	rm -rf $(PUBLISHDIR)

# ================================
# Build binary
# ================================
build: clean
	@mkdir -p $(PUBLISHDIR)
	go build -o $(PUBLISHDIR)/$(PROJECT_NAME) $(GOFLAGS)

# ================================
# Run locally (build & execute)
# ================================
run: build
	cd $(PUBLISHDIR) && ./$(PROJECT_NAME)

# ================================
# Environment builds (pre/prod/dev/test/local)
# Will copy config files and run.sh script to dist/
# ================================
pre: ENV=pre
prod: ENV=prod
dev: ENV=dev
test: ENV=test
local: ENV=local

pre prod dev test local: build
	@mkdir -p $(PUBLISHDIR)/conf
	@if [ ! -d "conf/$(ENV)" ]; then \
		echo "ERROR: Config directory conf/$(ENV) does not exist."; \
		exit 1; \
	fi
	cp -aRf conf/$(ENV)/* $(PUBLISHDIR)/conf
	cp run.sh $(PUBLISHDIR)/

# ================================
# Build Docker image
# ================================
docker-build: build
	@mv $(PUBLISHDIR)/$(PROJECT_NAME) $(PUBLISHDIR)/main
	docker build \
		--build-arg GIT_HASH=$(GITCOMMIT) \
		--build-arg GIT_TAG=$(GITBRANCH) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(DOCKER_REGISTRY)/$(PROJECT_NAME):$(GITBRANCH)-$(GITCOMMIT) \
		-t $(DOCKER_REGISTRY)/$(PROJECT_NAME):latest \
		$(PUBLISHDIR)

# ================================
# Push Docker image
# ================================
docker-push: docker-build
	docker push $(DOCKER_REGISTRY)/$(PROJECT_NAME):$(GITBRANCH)-$(GITCOMMIT)
	docker push $(DOCKER_REGISTRY)/$(PROJECT_NAME):latest
