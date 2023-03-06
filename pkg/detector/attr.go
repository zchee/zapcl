// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package detector

import (
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"

	"cloud.google.com/go/compute/metadata"
)

// ResourceAttributesFetcher abstracts environment lookup methods to query for environment variables, metadata attributes and file content.
type ResourceAttributesFetcher interface {
	EnvVar(name string) string
	Metadata(path string) string
	ReadAll(path string) string
}

type resourceFetcher struct {
	mdClient *metadata.Client
}

var _ ResourceAttributesFetcher = (*resourceFetcher)(nil)

// EnvVar uses os.Getenv() to gets for environment variable by name.
func (g *resourceFetcher) EnvVar(name string) string {
	return os.Getenv(name)
}

// Metadata uses metadata package Client.Get() to lookup for metadata attributes by path.
func (g *resourceFetcher) Metadata(path string) string {
	val, err := g.mdClient.Get(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(val)
}

// ReadAll reads all content of the file as a string.
func (g *resourceFetcher) ReadAll(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return *(*string)(unsafe.Pointer(&data))
}

var fetcher = &resourceFetcher{
	mdClient: metadata.NewClient(&http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   1 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
		},
	}),
}

// ResourceAttributes provides read-only access to the ResourceAttributesFetcher interface implementation.
func ResourceAttributes() ResourceAttributesFetcher {
	return fetcher
}
