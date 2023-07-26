// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package detector

// fakeResourceGetter mocks internal.ResourceAtttributesGetter interface to retrieve env vars and metadata.
type fakeResourceGetter struct {
	envVars  map[string]string
	metaVars map[string]string
	fsPaths  map[string]string
}

func (g *fakeResourceGetter) EnvVar(name string) string {
	if g.envVars != nil {
		if v, ok := g.envVars[name]; ok {
			return v
		}
	}
	return ""
}

func (g *fakeResourceGetter) Metadata(path string) string {
	if g.metaVars != nil {
		if v, ok := g.metaVars[path]; ok {
			return v
		}
	}
	return ""
}

func (g *fakeResourceGetter) ReadAll(path string) string {
	if g.fsPaths != nil {
		if v, ok := g.fsPaths[path]; ok {
			return v
		}
	}
	return ""
}
