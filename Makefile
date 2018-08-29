# Copyright 2018 VMware, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Constants
SHELL := /bin/bash
GO := $(shell command -v go 2> /dev/null)
DOCKER := $(shell command -v docker 2> /dev/null)
GLIDE := $(GOBIN)/glide
GOIMPORTS := $(GOBIN)/goimports
PATH := $(GOBIN):$(PATH)
FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Prints the help message
# Usage: make help
.PHONY: help
help:
	@echo "Usage: make <TARGET> [<ARGUMENTS>]"
	@echo ""
	@echo "Available targets are:"
	@echo ""
	@echo "    build              Fix imports, formats files and builds the project locally."
	@echo "    docker-build       Build the docker image for the project."
	@echo "    docker-run         Run the docker container."
	@echo "    format             Fix imports and format files."
	@echo "    help               Show this help message."
	@echo "    setup              Set up the development environment."
	@echo "    test               Run all unit tests."
	@echo "    vendor-install     Install all dependencies to vendor directory."
	@echo "    vendor-get <pkg>   Add a new dependency to the vendor directory."
	@echo ""

# Checks for necessary variables
.PHONY: env
env:
ifndef GO
	$(error go is not installed)
endif
ifeq ($(GOPATH),)
	$(error GOPATH is not set)
endif
ifeq ($(GOBIN),)
	$(error GOBIN is not set)
endif
ifndef DOCKER
	$(error docker is not installed)
endif

# Configures the development environment
# Usage: make setup
.PHONY: setup
setup: env $(GLIDE) $(GOIMPORTS)
	@echo > /dev/null

# Installs Glide
$(GLIDE):
	curl https://glide.sh/get | sh
	@echo "Glide installed!"

# Installs goimports
$(GOIMPORTS):
	go get golang.org/x/tools/cmd/goimports
	@echo "goimports installed!"

# Fixes imports, formats files and builds the project
# # Usage: make build
.PHONY: build
build: format
	go build -v ./...
	cp wavefront/config/wavefront.yaml wavefront/operatorconfig/
	@echo "Build was successful!"

# Builds the docker image for the project
# # Usage: make docker-build
.PHONY: docker-build
docker-build: setup
	docker build -t vmware/wavefront-istio-mixer-adapter:latest .
	@echo "Docker image was built successfully!"

# Runs the docker container
# # Usage: make docker-run
.PHONY: docker-run
docker-run: setup
	docker run -it -p 8080:8080 vmware/wavefront-istio-mixer-adapter:latest

# Fixes imports and formats files
# Usage: make format
.PHONY: format
format: setup
	@$(GOIMPORTS) -w -l $(FILES)

# Installs dependencies from glide.lock to the vendor directory
# Usage: make vendor-install
.PHONY: vendor-install
vendor-install: setup
	@$(GLIDE) install --strip-vendor

# Adds a new dependency to glide.yaml, glide.lock and to the vendor directory
# Usage: make vendor-get <pkg>
# Example: make vendor-get pkg=github.com/foo/bar
.PHONY: vendor-get
vendor-get: setup
	@$(GLIDE) get $(pkg) --strip-vendor

# Runs unit tests
# Usage: make test
.PHONY: test
test:
	go test -v ./...
