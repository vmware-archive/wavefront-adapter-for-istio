# Copyright 2018 VMware, Inc.
# SPDX-License-Identifier: Apache-2.0
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
	@echo "    dep-add <pkg>      Add a package to go.mod."
	@echo "    dep-tidy           Clean up unused dependencies."
	@echo "    docker-build       Build the docker image for the project."
	@echo "    docker-run         Run the docker container."
	@echo "    format             Fix imports and format files."
	@echo "    helm-pack          Create a Helm configuration package."
	@echo "    helm-print         Dry run and print the Helm manifest."
	@echo "    helm-generate      Generate the manifest from Helm configuration."
	@echo "    help               Show this help message."
	@echo "    setup              Set up the development environment."
	@echo "    test               Run all unit tests."
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
setup: env $(GOIMPORTS)
	@echo > /dev/null

# Installs goimports
$(GOIMPORTS):
	go get golang.org/x/tools/cmd/goimports
	@echo "goimports installed!"

# Fixes imports, formats files and builds the project
# Usage: make build
.PHONY: build
build: format
	go build -v ./...
	cp wavefront/config/wavefront.yaml install/wavefront/templates/
	sed -i '' -E 's/namespace: istio-system/namespace: {{ .Values.namespaces.istio }}/' install/wavefront/templates/wavefront.yaml
	sed -i '' -E 's/- metric/- metric.{{ .Values.namespaces.adapter }}/' install/wavefront/templates/wavefront.yaml
	@echo "Build was successful!"

# Builds the docker image for the project
# Usage: make docker-build
.PHONY: docker-build
docker-build: build
	# Sometimes docker build may fail in ubuntu due to network driver issue, in that case pass --network=host to docker build command.
	docker build -t vmware/wavefront-adapter-for-istio:latest .
	@echo "Docker image was built successfully!"

# Runs the docker container
# Usage: make docker-run
.PHONY: docker-run
docker-run: setup
	docker run -it -p 8000:8000 vmware/wavefront-adapter-for-istio:latest

# Dry-runs and prints the Helm manifest
# Usage: make helm-print
.PHONY: helm-print
helm-print:
	helm install --dry-run --debug install/wavefront/ --generate-name

# Creates a Helm configuration package
# Usage: make helm-pack
.PHONY: helm-pack
helm-pack:
	helm package install/wavefront/

# Generates the manifest from Helm configuration
# Usage: make helm-generate
.PHONY: helm-generate
helm-generate:
	@rm -f install/config.yaml
	helm template install/wavefront > install/config.yaml

# Fixes imports and formats files
# Usage: make format
.PHONY: format
format: setup
	@$(GOIMPORTS) -w -l $(FILES)

# Runs unit tests
# Usage: make test
.PHONY: test
test: build
	go test -v ./...

# Adds a package to go.mod
# Usage: make dep-add pkg=istio.io/istio@1.0.4
.PHONY: dep-add
dep-add:
	go mod edit -require $(pkg)

# Cleans up unused dependencies
# Usage: make dep-tidy
.PHONY: dep-tidy
dep-tidy:
	go mod tidy
