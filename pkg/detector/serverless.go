// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package detector

// List of Cloud Run env vars:
//
//	https://cloud.google.com/run/docs/container-contract#env-vars
const (
	// EnvCloudRunService is the name of the Cloud Run service being run.
	EnvCloudRunService = "K_SERVICE"

	// EnvCloudRunRevision is the name of the Cloud Run revision being run.
	EnvCloudRunRevision = "K_REVISION"

	// EnvCloudRunConfig is	the name of the Cloud Run configuration that created the revision.
	EnvCloudRunConfig = "K_CONFIGURATION"
)

func (d *Detector) isCloudRun() bool {
	config := d.attrs.EnvVar(EnvCloudRunConfig)
	// note that this envvar is also present in Cloud Function environments
	service := d.attrs.EnvVar(EnvCloudRunService)
	revision := d.attrs.EnvVar(EnvCloudRunRevision)

	return config != "" && service != "" && revision != ""
}

// List of Cloud Run jobs env vars:
//
//	https://cloud.google.com/run/docs/container-contract#jobs-env-vars
const (
	// EnvCloudRunJobsService is the name of the Cloud Run job being run.
	EnvCloudRunJobsService = "CLOUD_RUN_JOB"

	// EnvCloudRunJobsRevision is the name of the Cloud Run execution being run.
	EnvCloudRunJobsRevision = "CLOUD_RUN_EXECUTION"

	// EnvCloudRunJobsTaskIndex for each task, this will be set to a unique value between 0 and the number of tasks minus 1.
	EnvCloudRunJobsTaskIndex = "CLOUD_RUN_TASK_INDEX"

	// cloudRunJobsTaskAttempt is the number of times this task has been retried.
	//
	// Starts at 0 for the first attempt; increments by 1 for every successive retry, up to the maximum retries value.
	EnvCloudRunJobsTaskAttempt = "CLOUD_RUN_TASK_ATTEMPT"

	// cloudRunJobsRevisionEnv is the number of tasks defined in the --tasks parameter.
	EnvCloudRunJobsTaskCount = "CLOUD_RUN_TASK_COUNT"
)

func (d *Detector) isCloudRunJobs() bool {
	service := d.attrs.EnvVar(EnvCloudRunJobsService)
	revision := d.attrs.EnvVar(EnvCloudRunJobsRevision)

	return service != "" && revision != ""
}

// List of Cloud Functions newer runtimes env vars:
//
//	https://cloud.google.com/functions/docs/configuring/env-var#newer_runtimes
const (
	// EnvCloudFunctionsTarget is the function to be executed.
	EnvCloudFunctionsTarget = "FUNCTION_TARGET"

	// EnvCloudFunctionsSignatureType is the type of the function: http for HTTP functions, and event for event-driven functions.
	EnvCloudFunctionsSignatureType = "FUNCTION_SIGNATURE_TYPE"

	// EnvCloudFunctionsKService is the name of the function resource.
	//
	// Note that this envvar is also present in Cloud Run environments.
	EnvCloudFunctionsKService = "K_SERVICE"

	// EnvCloudFunctionsKRevision is the version identifier of the function.
	//
	// Note that this envvar is also present in Cloud Run environments.
	EnvCloudFunctionsKRevision = "K_REVISION"
)

func (d *Detector) isCloudFunctions() bool {
	target := d.attrs.EnvVar(EnvCloudFunctionsTarget)
	signatureType := d.attrs.EnvVar(EnvCloudFunctionsSignatureType)
	service := d.attrs.EnvVar(EnvCloudFunctionsKService)
	revision := d.attrs.EnvVar(EnvCloudFunctionsKRevision)

	return target != "" && signatureType != "" && service != "" && revision != ""
}
