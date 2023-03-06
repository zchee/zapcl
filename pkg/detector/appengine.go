// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package detector

// List of App Engine Standard env vars:
//
// https://cloud.google.com/appengine/docs/standard/go/runtime#environment_variables
const (
	// EnvAppEngineEnv is the App Engine environment. Set to "standard".
	EnvAppEngineEnv = "GAE_ENV"
)

func (d *Detector) isAppEngineStandard() bool {
	env := d.attrs.EnvVar(EnvAppEngineEnv)

	return env == "standard"
}

// List of App Engine Flex env vars:
//
// https://cloud.google.com/appengine/docs/flexible/python/runtime#environment_variables
const (
	// EnvAppEngineFlexService is the service name specified in your application's app.yaml file, or if no service name is specified, it is set to default.
	EnvAppEngineFlexService = "GAE_SERVICE"

	// EnvAppEngineFlexVersion is the version label of the current application.
	EnvAppEngineFlexVersion = "GAE_VERSION"

	// EnvAppEngineFlexInstance is the name of the current instance.
	EnvAppEngineFlexInstance = "GAE_INSTANCE"
)

func (d *Detector) isAppEngineFlex() bool {
	service := d.attrs.EnvVar(EnvAppEngineFlexService)
	version := d.attrs.EnvVar(EnvAppEngineFlexVersion)
	instance := d.attrs.EnvVar(EnvAppEngineFlexInstance)

	return instance != "" && service != "" && version != ""
}
