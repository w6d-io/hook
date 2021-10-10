export GO111MODULE  := on
export PATH         := $(shell pwd)/bin:${PATH}
export NEXT_TAG     ?=

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

ifeq (,$(shell go env GOOS))
GOOS=$(shell echo $OS)
else
GOOS=$(shell go env GOOS)
endif

all: test

# Defines

## go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

# Rules

GOIMPORTS = $(shell pwd)/bin/goimports
bin/goimports: ## Download goimports locally if necessary
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)

GITCHGLOG = $(shell pwd)/bin/git-chglog
bin/git-chglog: ## Download git-chglog locally if necessary
	$(call go-get-toll,$(GITCHGLOG),github.com/git-chglog/git-chglog/cmd/git-chglog@latest)

# Formats the code
.PHONY: format
format: bin/goimports
	$(GOIMPORTS) -w -local github.com/w6d-io,gitlab.w6d.io/w6d http kafka *.go

# Changelog
.PHONY: changelog
changelog: bin/git-chglog
	$(GITCHGLOG) -o docs/CHANGELOG.md --next-tag $(NEXT_TAG)

.PHONY: test
test: fmt vet
	go test -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

