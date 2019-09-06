PACKAGE  = stash.kopano.io/kgol/kweb
PACKAGE_NAME = kopano-$(shell basename $(PACKAGE))

# Tools

GO      ?= go
GOFMT   ?= gofmt
DEP     ?= dep
GOLINT  ?= golint

GO2XUNIT ?= go2xunit

CHGLOG ?= git-chglog

# Cgo
CGO_ENABLED ?= 0

# Variables
ARGS    ?=
PWD     := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null | sed 's/^v//' || \
			cat $(CURDIR)/.version 2> /dev/null || echo 0.0.0-unreleased)
PKGS     = $(or $(PKG),$(shell $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS) 2>/dev/null)
CMDS     = $(or $(CMD),$(addprefix cmd/,$(notdir $(shell find "$(PWD)/cmd/" -type d))))
TIMEOUT  = 30

export CGO_ENABLED

# Build

.PHONY: all
all: fmt vendor | $(CMDS) $(PLUGINS)

plugins: fmt vendor | $(PLUGINS)

.PHONY: $(CMDS)
$(CMDS): vendor ; $(info building $@ ...) @
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-trimpath \
		-tags release \
		-ldflags '-s -w -X $(PACKAGE)/version.Version=$(VERSION) -X $(PACKAGE)/version.BuildDate=$(DATE) -extldflags -static' \
		-o bin/$(notdir $@) $(PACKAGE)/$@

# Helpers

.PHONY: lint
lint: vendor ; $(info running golint ...)	@
	@ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret

.PHONY: vet
vet: vendor ; $(info running go vet ...)	@
	@ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GO) vet $$pkg)" || ret=1 ; \
	done ; exit $$ret

.PHONY: fmt
fmt: ; $(info running gofmt ...)	@
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

.PHONY: check
check: ; $(info checking dependencies ...) @
	@$(DEP) check && echo OK

# Tests

TEST_TARGETS := test-default test-bench test-short test-race test-verbose
.PHONY: $(TEST_TARGETS)
test-bench:   ARGS=-run=_Bench* -test.benchmem -bench=.
test-short:   ARGS=-short
test-race:    ARGS=-race
test-race:    CGO_ENABLED=1
test-verbose: ARGS=-v
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test

.PHONY: test
test: vendor ; $(info running $(NAME:%=% )tests ...)	@
	@CGO_ENABLED=$(CGO_ENABLED) $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

TEST_XML_TARGETS := test-xml-default test-xml-short test-xml-race
.PHONY: $(TEST_XML_TARGETS)
test-xml-short: ARGS=-short
test-xml-race:  ARGS=-race
test-xml-race:  CGO_ENABLED=1
$(TEST_XML_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_XML_TARGETS): test-xml

.PHONY: test-xml
test-xml: vendor ; $(info running $(NAME:%=% )tests ...)	@
	@mkdir -p test
	2>&1 CGO_ENABLED=$(CGO_ENABLED) $(GO) test -timeout $(TIMEOUT)s $(ARGS) -v $(TESTPKGS) | tee test/tests.output
	$(shell test -s test/tests.output && $(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml)

# Dep

Gopkg.lock: Gopkg.toml ; $(info updating dependencies ...)
	@$(DEP) ensure -v -update
	@touch $@

vendor: Gopkg.lock ; $(info retrieving dependencies ...)
	@$(DEP) ensure -v -vendor-only
	@touch $@

# Dist

.PHONY: licenses
licenses: ; $(info building licenses files ...)
	$(CURDIR)/scripts/go-license-ranger.py > $(CURDIR)/3rdparty-LICENSES.md

3rdparty-LICENSES.md: licenses

.PHONY: dist
dist: 3rdparty-LICENSES.md ; $(info building dist tarball ...)
	@rm -rf "dist/${PACKAGE_NAME}-${VERSION}"
	@mkdir -p "dist/${PACKAGE_NAME}-${VERSION}"
	@mkdir -p "dist/${PACKAGE_NAME}-${VERSION}/scripts"
	@cd dist && \
	cp -avf ../LICENSE.txt "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../README.md "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../3rdparty-LICENSES.md "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../bin/* "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../scripts/kopano-kwebd.binscript "${PACKAGE_NAME}-${VERSION}/scripts" && \
	cp -avf ../scripts/kopano-kwebd.service "${PACKAGE_NAME}-${VERSION}/scripts" && \
	cp -avf ../scripts/kwebd.cfg "${PACKAGE_NAME}-${VERSION}/scripts" && \
	cp -avf ../scripts/robots.txt "${PACKAGE_NAME}-${VERSION}/scripts" && \
	cp -avf ../scripts/favicon.ico "${PACKAGE_NAME}-${VERSION}/scripts" && \
	tar --owner=0 --group=0 -czvf ${PACKAGE_NAME}-${VERSION}.tar.gz "${PACKAGE_NAME}-${VERSION}" && \
	cd ..

.PHONE: changelog
changelog: ; $(info updating changelog ...)
	$(CHGLOG) --output CHANGELOG.md $(ARGS)

# Rest

.PHONY: clean
clean: ; $(info cleaning ...)	@
	@rm -rf bin
	@rm -rf test/test.*

.PHONY: version
version:
	@echo $(VERSION)
