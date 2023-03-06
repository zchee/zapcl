// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package monitoredresource

import (
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"

	"github.com/zchee/zapcl/pkg/detector"
)

const (
	there               = "anyvalue"
	projectID           = "test-project"
	zoneID              = "test-region-zone"
	regionID            = "test-region"
	serviceName         = "test-service"
	version             = "1.0"
	instanceName        = "test-12345"
	qualifiedZoneName   = "projects/" + projectID + "/zones/" + zoneID
	qualifiedRegionName = "projects/" + projectID + "/regions/" + regionID
	funcSignature       = "test-cf-signature"
	funcTarget          = "test-cf-target"
	crConfig            = "test-cr-config"
	clusterName         = "test-k8s-cluster"
	podName             = "test-k8s-pod-name"
	containerName       = "test-k8s-container-name"
	namespaceName       = "test-k8s-namespace-name"
	instanceID          = "test-instance-12345"
)

// fakeResourceGetter mocks internal.ResourceAtttributesGetter interface to retrieve env vars and metadata
type fakeResourceGetter struct {
	envVars  map[string]string
	metaVars map[string]string
	fsPaths  map[string]string
}

// func (g *fakeResourceGetter) ProjectID() (string, error)    { return projectID, nil }
// func (g *fakeResourceGetter) InstanceID() (string, error)   { return instanceID, nil }
// func (g *fakeResourceGetter) InstanceName() (string, error) { return instanceName, nil }
// func (g *fakeResourceGetter) Zone() (string, error)         { return zoneID, nil }

// func (g *fakeResourceGetter) InstanceAttributeValue(s string) (string, error) { return g.Get(s) }

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

// setupDetectResource resets sync.Once on detectResource and enforces mocked resource attribute getter
func setupDetectedResource(envVars, metaVars, fsPaths map[string]string) {
	ResourceDetector.once = new(sync.Once)
	fake := &fakeResourceGetter{
		envVars:  envVars,
		metaVars: metaVars,
		fsPaths:  fsPaths,
	}
	ResourceDetector.attrs = fake
	ResourceDetector.pb = nil
}

func TestResourceDetection(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		metaVars map[string]string
		fsPaths  map[string]string
		want     *MonitoredResource
	}{
		{
			name: "CloudFunction",
			envVars: map[string]string{
				detector.EnvCloudFunctionsTarget:        funcTarget,
				detector.EnvCloudFunctionsSignatureType: funcSignature,
				detector.EnvCloudFunctionsKService:      serviceName,
				detector.EnvCloudRunRevision:            regionID,
			},
			metaVars: map[string]string{
				"":                   there,
				"project/project-id": projectID,
				"instance/region":    qualifiedRegionName,
			},
			want: &MonitoredResource{
				LogID: "cloudfunctions.googleapis.com%2Fcloud-functions",
				MonitoredResource: &mrpb.MonitoredResource{
					Type: "cloud_function",
					Labels: map[string]string{
						"project_id":    projectID,
						"region":        regionID,
						"function_name": serviceName,
					},
				},
			},
		},
		{
			name: "CloudRun",
			envVars: map[string]string{
				detector.EnvCloudRunConfig:   crConfig,
				detector.EnvCloudRunService:  serviceName,
				detector.EnvCloudRunRevision: version,
			},
			metaVars: map[string]string{
				"":                   there,
				"project/project-id": projectID,
				"instance/region":    qualifiedRegionName,
			},
			want: &MonitoredResource{
				LogID: "run.googleapis.com%2Fstdout",
				MonitoredResource: &mrpb.MonitoredResource{
					Type: "cloud_run_revision",
					Labels: map[string]string{
						"project_id":         projectID,
						"location":           regionID,
						"service_name":       serviceName,
						"revision_name":      version,
						"configuration_name": crConfig,
					},
				},
			},
		},
		{
			name: "CloudRunJobs",
			envVars: map[string]string{
				detector.EnvCloudRunJobsService:  serviceName,
				detector.EnvCloudRunJobsRevision: version,
			},
			metaVars: map[string]string{
				"":                   there,
				"project/project-id": projectID,
				"instance/region":    qualifiedRegionName,
			},
			want: &MonitoredResource{
				LogID: "run.googleapis.com%2Fstdout",
				MonitoredResource: &mrpb.MonitoredResource{
					Type: "cloud_run_job",
					Labels: map[string]string{
						"project_id": projectID,
						"location":   regionID,
						"job_name":   serviceName,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupDetectedResource(tt.envVars, tt.metaVars, tt.fsPaths)
			got := Detect()
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(mrpb.MonitoredResource{})); diff != "" {
				t.Errorf("got(-),want(+):\n%s", diff)
			}
		})
	}
}

// var benchmarkResultHolder *mrpb.MonitoredResource
//
// func BenchmarkDetectResource(b *testing.B) {
// 	var result *mrpb.MonitoredResource
//
// 	for n := 0; n < b.N; n++ {
// 		result = detectResource()
// 	}
//
// 	benchmarkResultHolder = result
// }
