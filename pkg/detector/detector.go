// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package detector

// Platform represents a GCP service platforms.
type Platform uint8

const (
	// UnknownPlatform is teh unknown platform.
	UnknownPlatform Platform = iota

	// GKE is the Kubernetes Engine platform.
	// TODO(zchee): not implemented yet.
	GKE

	// GCE is the Google Compute Engine platform.
	// TODO(zchee): not implemented yet.
	GCE

	// CloudRun is the Cloud Run platform.
	CloudRun

	// CloudRunJobs is the Cloud Run jobs platform.
	CloudRunJobs

	// CloudFunctions is the Cloud Functions platform.
	CloudFunctions

	// AppEngineStandard is the App Engine Standard 2nd platform.
	AppEngineStandard

	// AppEngineFlex is the App Engine Flex platform.
	AppEngineFlex
)

// Detector collects resource information for all GCP platforms
type Detector struct {
	attrs ResourceAttributesFetcher
}

func NewDetector(attrs ResourceAttributesFetcher) *Detector {
	return &Detector{
		attrs: attrs,
	}
}

// CloudPlatform returns the platform on which this program is running
func (d *Detector) CloudPlatform() Platform {
	d.attrs = fetcher

	switch {
	case d.isGCE(): // TODO(zchee): not implemented yet.
		return UnknownPlatform // GCE

	case d.isGKE(): // TODO(zchee): not implemented yet.
		return UnknownPlatform // GKE

	case d.isCloudRun():
		return CloudRun

	case d.isCloudRunJobs():
		return CloudRunJobs

	case d.isCloudFunctions():
		return CloudFunctions

	case d.isAppEngineStandard():
		return AppEngineStandard

	case d.isAppEngineFlex():
		return AppEngineFlex
	}

	return UnknownPlatform
}

// MetadataProvider contains the subset of the metadata.Client functions used
// by this resource Detector to allow testing with a fake implementation.
type MetadataProvider interface {
	ProjectID() (string, error)
	InstanceID() (string, error)
	Get(string) (string, error)
	InstanceName() (string, error)
	Zone() (string, error)
	InstanceAttributeValue(string) (string, error)
}

// OSProvider contains the subset of the os package functions used by.
type OSProvider interface {
	LookupEnv(string) (string, bool)
}
