#!/usr/bin/env make

.PHONY: start dev build clean docker-build docker-console build-console

export NODE_ENV=development

# ---------------------------------------------------------------------------------------------------------------------
# CONFIG
# ---------------------------------------------------------------------------------------------------------------------

DOCKER_IMAGE_VERSION=0.1.0
DOCKER_IMAGE_TAG=orange-judge/$(DOCKER_IMAGE_VERSION)

# ---------------------------------------------------------------------------------------------------------------------
# UTILS
# ---------------------------------------------------------------------------------------------------------------------

clean:
	rm server

# ---------------------------------------------------------------------------------------------------------------------
# DEVELOPMENT
# ---------------------------------------------------------------------------------------------------------------------

dev:
	yarn dev

# ---------------------------------------------------------------------------------------------------------------------
# PRODUCTION
# ---------------------------------------------------------------------------------------------------------------------

build:
	go build -o server.exe -v main.go

start:
	./server -d

# ---------------------------------------------------------------------------------------------------------------------
# DOCKER
# ---------------------------------------------------------------------------------------------------------------------

docker-build:
	@docker build -t $(DOCKER_IMAGE_TAG) .

docker-console:
	docker-compose run orange-judge /bin/bash

build-console: docker-build docker-console

docker-up:
	docker-compose up