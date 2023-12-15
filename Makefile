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

ifeq (,$(shell go env GOARCH))
GOARCH=$(shell echo uname -p)
else
GOARCH=$(shell go env GOARCH)
endif

ifeq (darwin,$(GOOS))
GOTAGS = "-tags=dynamic"
else
GOTAGS =
endif

all: test

# Defines

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOIMPORTS ?= $(LOCALBIN)/goimports
GITCHGLOG ?= $(LOCALBIN)/git-chglog
export RPK       ?= $(LOCALBIN)/rpk

RPK_VERSION = v22.1.10
RPK_ZIP ?= https://github.com/redpanda-data/redpanda/releases/download/$(RPK_VERSION)/rpk-$(GOOS)-$(GOARCH).zip

.PHONY: rpk
rpk: $(RPK)
$(RPK): $(LOCALBIN)
	@test -s $(LOCALBIN)/rpk || $(call install,$(RPK),rpk,$(RPK_ZIP))

.PHONY: goimports
goimports: $(GOIMPORTS) ## Download goimports locally if necessary
$(GOIMPORTS): $(LOCALBIN)
	@test -s $(LOCALBIN)/goimports || GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@latest

.PHONY: chglog
chglog: $(GITCHGLOG) ## Download git-chglog locally if necessary
$(GITCHGLOG): $(LOCALBIN)
	@test -s $(LOCALBIN)/git-chglog || GOBIN=$(LOCALBIN) go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest


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
	go test $(GOTAGS) -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

.PHONY: integration-test
integration-test: rpk
	go test $(GOTAGS) -tags=integration -v -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out | grep total

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet $(GOTAGS) ./...

## scripts
define install
[ -f $(1) ] || { \
set -e;\
TMP_DIR=$$(mktemp -d);\
cd $$TMP_DIR ;\
wget -q $(3);\
unzip *.zip $(2);\
mv $(2) $(1);\
rm -rf $$TMP_DIR ;\
}
endef