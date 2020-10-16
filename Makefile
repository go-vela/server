# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

# capture the current date we build the application from
BUILD_DATE = $(shell date +%Y-%m-%dT%H:%M:%SZ)

# check if a git commit sha is already set
ifndef GITHUB_SHA
	# capture the current git commit sha we build the application from
	GITHUB_SHA = $(shell git rev-parse HEAD)
endif

# check if a git tag is already set
ifndef GITHUB_TAG
	# capture the current git tag we build the application from
	GITHUB_TAG = $(shell git describe --tag --abbrev=0)
endif

# check if a go version is already set
ifndef GOLANG_VERSION
	# capture the current go version we build the application from
	GOLANG_VERSION = $(shell go version | awk '{ print $$3 }')
endif

# create a list of linker flags for building the golang application
LD_FLAGS = -X github.com/go-vela/server/version.Commit=${GITHUB_SHA} -X github.com/go-vela/server/version.Date=${BUILD_DATE} -X github.com/go-vela/server/version.Go=${GOLANG_VERSION} -X github.com/go-vela/server/version.Tag=${GITHUB_TAG}

# The `clean` target is intended to clean the workspace
# and prepare the local changes for submission.
#
# Usage: `make clean`
.PHONY: clean
clean: tidy vet fmt fix

# The `restart` target is intended to destroy and
# create the local Docker compose stack.
#
# Usage: `make restart`
.PHONY: restart
restart: down up

# The `up` target is intended to create
# the local Docker compose stack.
#
# Usage: `make up`
.PHONY: up
up: build pull compose-up

# The `down` target is intended to destroy
# the local Docker compose stack.
#
# Usage: `make down`
.PHONY: down
down: compose-down

# The `tidy` target is intended to clean up
# the Go module files (go.mod & go.sum).
#
# Usage: `make tidy`
.PHONY: tidy
tidy:
	@echo
	@echo "### Tidying Go module"
	@go mod tidy

# The `vet` target is intended to inspect the
# Go source code for potential issues.
#
# Usage: `make vet`
.PHONY: vet
vet:
	@echo
	@echo "### Vetting Go code"
	@go vet ./...

# The `fmt` target is intended to format the
# Go source code to meet the language standards.
#
# Usage: `make fmt`
.PHONY: fmt
fmt:
	@echo
	@echo "### Formatting Go Code"
	@go fmt ./...

# The `fix` target is intended to rewrite the
# Go source code using old APIs.
#
# Usage: `make fix`
.PHONY: fix
fix:
	@echo
	@echo "### Fixing Go Code"
	@go fix ./...

# The `test` target is intended to run
# the tests for the Go source code.
#
# Usage: `make test`
.PHONY: test
test:
	@echo
	@echo "### Testing Go Code"
	@go test -race ./...

# The `test-cover` target is intended to run
# the tests for the Go source code and then
# open the test coverage report.
#
# Usage: `make test-cover`
.PHONY: test-cover
test-cover:
	@echo
	@echo "### Creating test coverage report"
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo
	@echo "### Opening test coverage report"
	@go tool cover -html=coverage.out

# The `build` target is intended to compile
# the Go source code into a binary.
#
# Usage: `make build`
.PHONY: build
build:
	@echo
	@echo "### Building release/vela-server binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-ldflags '${LD_FLAGS}' \
		-o release/vela-server \
		github.com/go-vela/server/cmd/vela-server

# The `build-static` target is intended to compile
# the Go source code into a statically linked binary.
#
# Usage: `make build-static`
.PHONY: build-static
build-static:
	@echo
	@echo "### Building static release/vela-server binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-ldflags '-s -w -extldflags "-static" ${LD_FLAGS}' \
		-o release/vela-server \
		github.com/go-vela/server/cmd/vela-server

# The `build-static-ci` target is intended to compile
# the Go source code into a statically linked binary
# when used within a CI environment.
#
# Usage: `make build-static-ci`
.PHONY: build-static-ci
build-static-ci:
	@echo
	@echo "### Building CI static release/vela-server binary"
	@go build -a \
		-ldflags '-s -w -extldflags "-static" ${LD_FLAGS}' \
		-o release/vela-server \
		github.com/go-vela/server/cmd/vela-server

# The `check` target is intended to output all
# dependencies from the Go module that need updates.
#
# Usage: `make check`
.PHONY: check
check: check-install
	@echo
	@echo "### Checking dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -update

# The `check-direct` target is intended to output direct
# dependencies from the Go module that need updates.
#
# Usage: `make check-direct`
.PHONY: check-direct
check-direct: check-install
	@echo
	@echo "### Checking direct dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -direct

# The `check-full` target is intended to output
# all dependencies from the Go module.
#
# Usage: `make check-full`
.PHONY: check-full
check-full: check-install
	@echo
	@echo "### Checking all dependencies for updates"
	@go list -u -m -json all | go-mod-outdated

# The `check-install` target is intended to download
# the tool used to check dependencies from the Go module.
#
# Usage: `make check-install`
.PHONY: check-install
check-install:
	@echo
	@echo "### Installing psampaz/go-mod-outdated"
	@go get -u github.com/psampaz/go-mod-outdated

# The `bump-deps` target is intended to upgrade
# non-test dependencies for the Go module.
#
# Usage: `make bump-deps`
.PHONY: bump-deps
bump-deps: check
	@echo
	@echo "### Upgrading dependencies"
	@go get -u ./...

# The `bump-deps-full` target is intended to upgrade
# all dependencies for the Go module.
#
# Usage: `make bump-deps-full`
.PHONY: bump-deps-full
bump-deps-full: check
	@echo
	@echo "### Upgrading all dependencies"
	@go get -t -u ./...

# The `pull` target is intended to pull all
# images for the local Docker compose stack.
#
# Usage: `make pull`
.PHONY: pull
pull:
	@echo
	@echo "### Pulling images for docker-compose stack"
	@docker-compose pull

# The `compose-up` target is intended to build and create
# all containers for the local Docker compose stack.
#
# Usage: `make compose-up`
.PHONY: compose-up
compose-up:
	@echo
	@echo "### Creating containers for docker-compose stack"
	@docker-compose -f docker-compose.yml up -d --build

# The `compose-down` target is intended to destroy
# all containers for the local Docker compose stack.
#
# Usage: `make compose-down`
.PHONY: compose-down
compose-down:
	@echo
	@echo "### Destroying containers for docker-compose stack"
	@docker-compose -f docker-compose.yml down
