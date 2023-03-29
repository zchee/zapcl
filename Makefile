# -----------------------------------------------------------------------------
# global

.DEFAULT_GOAL := help

# -----------------------------------------------------------------------------
# go

GO_PATH ?= $(shell go env GOPATH)
GO_OS ?= $(shell go env GOOS)
GO_ARCH ?= $(shell go env GOARCH)

PACKAGES := $(subst $(GO_PATH)/src/,,$(CURDIR))
CGO_ENABLED ?= 0
GO_BUILDTAGS=osusergo,netgo,static
GO_LDFLAGS=-s -w
ifeq (${GO_OS},linux)
GO_LDFLAGS+=-d
endif
GO_LDFLAGS+=-buildid= "-extldflags=-static"
GO_FLAGS ?= -tags='${GO_BUILDTAGS}' -ldflags='${GO_LDFLAGS}'

GO_PACKAGES = $(shell go list ./...)
GO_TEST_PACKAGES ?= $(shell go list -f='{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)

GO_TEST ?= ${TOOLS_BIN}/gotestsum -f standard-verbose --
GO_TEST_FLAGS ?= -race -count=1 -shuffle=on
GO_TEST_FUNC ?= .
GO_COVERAGE_OUT ?= coverage.out
GO_BENCH_FLAGS ?= -benchmem
GO_BENCH_FUNC ?= .
GO_LINT_PACKAGE ?= ./...
GO_LINT_FLAGS ?=

TOOLS_DIR := ${CURDIR}/tools
TOOLS_BIN := ${TOOLS_DIR}/bin
TOOLS = $(shell cd ${TOOLS_DIR}; go list -tags=tools -f='{{ join .Imports " " }}')

JOBS = $(shell getconf _NPROCESSORS_CONF)

# -----------------------------------------------------------------------------
# defines

define target
@printf "+ $(patsubst ,$@,$(1))\\n" >&2
endef

# -----------------------------------------------------------------------------
# target

##@ test, bench, coverage

.PHONY: test
test: CGO_ENABLED=1
test: GO_FLAGS=-tags='$(subst ${space},${comma},${GO_BUILDTAGS})'
test: tools/bin/gotestsum  ## Runs package test including race condition.
	$(call target)
	@CGO_ENABLED=${CGO_ENABLED} GOTESTSUM_FORMAT=standard-verbose ${GO_TEST} ${GO_TEST_FLAGS} -run=${GO_TEST_FUNC} $(strip ${GO_FLAGS}) ${GO_TEST_PACKAGES}

.PHONY: coverage
coverage: CGO_ENABLED=1
coverage: GO_FLAGS=-tags='$(subst ${space},${comma},${GO_BUILDTAGS})'
coverage: tools/bin/gotestsum  ## Takes packages test coverage.
	$(call target)
	@CGO_ENABLED=${CGO_ENABLED} ${GO_TEST} ${GO_TEST_FLAGS} -covermode=atomic -coverpkg=./... -coverprofile=${GO_COVERAGE_OUT} $(strip ${GO_FLAGS}) ${GO_PACKAGES}

##@ fmt, lint

.PHONY: fmt
fmt: tools/bin/goimports-reviser tools/bin/gofumpt  ## Run goimports-reviser and gofumpt.
	$(call target)
	@${TOOLS_BIN}/goimports-reviser -project-name ${PACKAGES} ./...
	@${TOOLS_BIN}/gofumpt -extra -w $(shell find $$PWD -iname "*.go" -not -iname "*pb.go" -not -iwholename "*vendor*")

.PHONY: lint
lint: lint/golangci-lint  ## Run all linters.

.PHONY: lint/golangci-lint
lint/golangci-lint: tools/bin/golangci-lint .golangci.yaml  ## Run golangci-lint.
	$(call target)
	@${TOOLS_BIN}/golangci-lint -j ${JOBS} run $(strip ${GO_LINT_FLAGS}) ${GO_LINT_PACKAGE}

##@ tools

define build_tool
@cd ${TOOLS_DIR}; \
for t in ${TOOLS}; do \
	if [ -z '$1' ] || [ $$(basename $${t%%/v*}) = '$1' ]; then \
		echo "Install $$t ..." >&2; \
		GOBIN=${TOOLS_BIN} CGO_ENABLED=0 go install -v -mod=readonly ${GO_FLAGS} "$${t}"; \
	fi \
done
endef

tools/bin/%: ${TOOLS_DIR}/go.mod ${TOOLS_DIR}/go.sum
tools/bin/%:  ## Install an individual dependency tool
	$(call build_tool,$(@F))

.PHONY: tools
tools:  ## Install tools
	$(call build_tool)

##@ clean

.PHONY: clean
clean:  ## Cleanups binaries and extra files in the package.
	$(call target)
	@rm -rf *.out *.test *.prof trace.txt ${TOOLS_BIN}

##@ help

.PHONY: help
help:  ## Show this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[33m<target>\033[0m\n"} /^[a-zA-Z_0-9\/%_-]+:.*?##/ { printf "  \033[1;32m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ misc

.PHONY: todo
todo:  ## Print the all of (TODO|BUG|XXX|FIXME|NOTE) in packages.
	@grep -E '(TODO|BUG|XXX|FIXME)(\(.+\):|:)' $(shell find . -type f -name '*.go' -and -not -iwholename '*vendor*')

.PHONY: env/%
env/%: ## Print the value of MAKEFILE_VARIABLE. Use `make env/GO_FLAGS` or etc.
	@echo $($*)

