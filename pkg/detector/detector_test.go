// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package detector

import (
	"testing"
)

func TestCloudPlatformAppEngineStandard(t *testing.T) {
	d := NewDetector(&fakeResourceGetter{
		envVars: map[string]string{
			EnvAppEngineFlexService: "foo",
			EnvAppEngineEnv:         "standard",
		},
	})
	platform := d.CloudPlatform()

	if platform != AppEngineStandard {
		t.Fatalf("got %d but want %d", platform, AppEngineStandard)
	}
}

func TestCloudPlatformAppEngineFlex(t *testing.T) {
	d := NewDetector(&fakeResourceGetter{
		envVars: map[string]string{
			EnvAppEngineFlexService:  "foo",
			EnvAppEngineFlexVersion:  "001",
			EnvAppEngineFlexInstance: "foo",
		},
	})
	platform := d.CloudPlatform()

	if platform != AppEngineFlex {
		t.Fatalf("got %d but want %d", platform, AppEngineFlex)
	}
}

func TestCloudPlatformCloudRun(t *testing.T) {
	d := NewDetector(&fakeResourceGetter{
		envVars: map[string]string{
			EnvCloudRunService:  "foo",
			EnvCloudRunRevision: "foo-001",
			EnvCloudRunConfig:   "foo",
		},
	})
	platform := d.CloudPlatform()

	if platform != CloudRun {
		t.Fatalf("got %d but want %d", platform, CloudRun)
	}
}

func TestCloudPlatformCloudRunJobs(t *testing.T) {
	d := NewDetector(&fakeResourceGetter{
		envVars: map[string]string{
			EnvCloudRunJobsService:  "foo",
			EnvCloudRunJobsRevision: "foo-001",
		},
	})
	platform := d.CloudPlatform()

	if platform != CloudRunJobs {
		t.Fatalf("got %d but want %d", platform, CloudRunJobs)
	}
}

func TestCloudPlatformCloudFunctions(t *testing.T) {
	d := NewDetector(&fakeResourceGetter{
		envVars: map[string]string{
			EnvCloudFunctionsTarget:        "foo",
			EnvCloudFunctionsSignatureType: "foo",
			EnvCloudFunctionsKService:      "foo",
			EnvCloudFunctionsKRevision:     "foo",
		},
	})
	platform := d.CloudPlatform()

	if platform != CloudFunctions {
		t.Fatalf("got %d but want %d", platform, CloudFunctions)
	}
}
