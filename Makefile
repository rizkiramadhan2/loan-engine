#!/bin/bash

export NOW=$(shell date +"%Y-%m-%d")

PACKAGE = simple-app
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date +%FT%T%z)
ldflags = -X $(PACKAGE)/cmd.CommitHash=$(COMMIT_HASH) -X $(PACKAGE)/cmd.BuildDate=$(BUILD_DATE) -s -w

build:
	@echo "${NOW} == BUILDING HTTP SERVER API"
	@echo "${ldflags}"
	@CGO_ENABLED=0 go build -ldflags '$(ldflags)' -o http_api cmd/http-api/main.go

run: build
	@echo "${NOW} == RUNNING HTTP SERVER API"
	@./http_api

configure:
	@[[ -f .env ]] || cp .dev/.env.example .dev/.env

.PHONY: dev
dev:
	make configure
	make down
	@echo "${NOW} == RUNNING DOCKER"
	@docker-compose -f .dev/docker-compose.dev.yaml up -d

down:
	@docker-compose -f .dev/docker-compose.dev.yaml down

log-go:
	@docker logs --follow simple-app --tail 100

log-db:
	@docker logs --follow simple_app_db --tail 100

log-redis:
	@docker logs --follow simple_app_redis --tail 100

migrate:
	@.dev/migrator exec
	
generate-migration:
	@.dev/migrator create