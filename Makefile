# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

.PHONY: restart
restart: down up

.PHONY: up
up: build compose-up

.PHONY: down
down: compose-down

.PHONY: rebuild
rebuild: build compose-up

.PHONY: clean
clean:
	#################################
	######      Go clean       ######
	#################################

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

.PHONY: build
build:
	#################################
	######    Build Binary     ######
	#################################

	GOOS=linux CGO_ENABLED=0 go build -o release/vela-server github.com/go-vela/server/cmd/vela-server

.PHONY: compose-up
compose-up:
	#################################
	###### Docker Build/Start  ######
	#################################

	@docker-compose -f docker-compose.yml up -d --build # build and start app

.PHONY: compose-down
compose-down:
	#################################
	###### Docker Tear Down    ######
	#################################

	@docker-compose -f docker-compose.yml down
