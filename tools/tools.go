// Package tools manages tools using during development.

//go:build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/incu6us/goimports-reviser/v3"
	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
