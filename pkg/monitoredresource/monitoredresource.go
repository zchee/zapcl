// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package monitoredresource

import (
	"strings"
	"sync"

	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"

	"github.com/zchee/zap-cloudlogging/pkg/detector"
)

type Type string

type Label map[string]string

// List of Monitored resource types.
//
//	https://cloud.google.com/logging/docs/api/v2/resource-list
const (
	// aiplatform.googleapis.com/Endpoint
	// Vertex AI Endpoint
	// A Vertex AI API Endpoint where Models are deployed into it.
	//
	// resource_container: The identifier of the GCP Project owning the Endpoint.
	// location: The region in which the service is running.
	// endpoint_id: The ID of the Endpoint.

	// aiplatform.googleapis.com/Featurestore
	// Vertex AI Feature Store	A Vertex AI Feature Store.
	//
	// resource_container: The identifier of the GCP Project owning the Featurestore.
	// location: The region in which the service is running.
	// featurestore_id: The ID of the Featurestore.

	// aiplatform.googleapis.com/IndexEndpoint
	// Matching Engine Index Endpoint	An Endpoint to which Matching Engine Indexes are deployed.
	//
	// resource_container: The identifier of the GCP Project owning the Index.
	// location: The region in which the service is running.
	// index_endpoint_id: The ID of the index endpoint.

	// aiplatform.googleapis.com/PipelineJob
	// Vertex Pipelines Job	A Vertex Pipelines Job.
	//
	// resource_container: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The region in which the service is running.
	// pipeline_job_id: The ID of the PipelineJob.

	// alloydb.googleapis.com/Instance
	// AlloyDB instance	Monitored resource representing an AlloyDB instance.
	//
	// resource_container: The identifier of the GCP project associated with this resource.
	// location: The Google Cloud region in which the AlloyDB instance is running.
	// cluster_id: AlloyDB cluster identifier.
	// instance_id: AlloyDB instance identifier.

	// api
	// Produced API	An API provided by the producer.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// service: The API service name, such as "cloudsql.googleapis.com".
	// method: The API method, such as "disks.list".
	// version: The API version, such as "v1".
	// location: The service specific notion of location. This can be the name of a zone, region, or "global".

	// apigateway.googleapis.com/Gateway
	// API Gateway	Fully managed API Gateway.
	//
	// resource_container: The identifier of the GCP Project owning the Gateway.
	// location: The region in which the Gateway is running.
	// gateway_id: The ID of the Gateway.

	// apigee.googleapis.com/Environment
	// Apigee environment	Monitored resource for Apigee environment.
	//
	// resource_container: The GCP project ID that writes to this monitored resource.
	// org: An organization is a container for all the objects in an Apigee account.
	// env: An environment is a runtime execution context for the proxies in an organization.
	// location: Location where the Apigee infrastructure is provisioned.

	// app_script_function
	// Apps Script Function	An Apps Script function.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// invocation_type: The invocation type.
	// function_name: The function name.

	// assistant_action
	// Google Assistant Action	An Action in a Google Assistant App.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// version_id: Stringified version ID of the assistant agent.
	// action_id: Action ID of the assistant agent.

	// audited_resource
	// Audited Resource	A Google Cloud resource that produces an audit log.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// service: The name of the API service generating the audit log.
	// method: The name of the API method generating the audit log.

	// autoscaler
	// Autoscaler	An autoscaler for a single managed instance group.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The zone or region for the autoscaler.
	// autoscaler_id: The identifier for the autoscaler.
	// autoscaler_name: The name of the autoscaler.
	// instance_group_manager_id: The identifier for the managed instance group scaled by the given autoscaler.
	// instance_group_manager_name: The name of the managed instance group scaled by the givenautoscaler.

	// aws_alb_load_balancer
	// Amazon ALB Load Balancer	A load balancer in Amazon ALB.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// name: The name of the load balancer.
	// region: The AWS region for the load balancer. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the load balancer.

	// aws_cloudfront_distribution
	// Amazon CloudFront CDN	A CloudFront content distribution network.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// distribution_id: The CloudFront distribution identifier assigned by AWS.
	// region: The AWS region for the CloudFront distribution. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the CDN.

	// aws_dynamodb_table
	// Amazon DynamoDB Table	A table in Amazon DynamoDB.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// table: The table name.
	// region: The AWS region for the table. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the table.

	// aws_ebs_volume
	// Amazon EBS Volume	An Amazon EC2 Elastic Block Storage volume.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// volume_id: The EBS volume identifier assigned by AWS.
	// region: The AWS region for the volume. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the volume.

	// aws_ec2_instance
	// Amazon EC2 Instance	A VM instance in Amazon EC2.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// instance_id: The VM instance identifier assigned by AWS.
	// aws_account: The AWS account number under which the VM is running.
	// region: The AWS region in which the VM is running. Supported AWS region values are listed by service at http://docs.aws.amazon.com/general/latest/gr/rande.html. The value supplied for this label must be prefixed with 'aws:' (for example, 'aws:us-east-1' is a valid value while 'us-east-1' is not).

	// aws_elasticache_cluster
	// Amazon Elasticache Cluster	A cache cluster in Amazon Elasticache.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// cluster_id: The cluster identifier.
	// region: The AWS region for the cluster. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the cluster.

	// aws_elb_load_balancer
	// Amazon Elastic Load Balancer	A load balancer in Amazon Elastic Load Balancer.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// region: The AWS region for the load balancer. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// name: The name of the load balancer.
	// aws_account: The AWS account number for the load balancer.

	// aws_emr_cluster
	// Amazon EMR Cluster	A cluster in Amazon Elastic MapReduce.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// cluster_id: The cluster identifier.
	// region: The AWS region for the cluster. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the cluster.

	// aws_kinesis_stream
	// Amazon Kinesis Stream	A stream in Amazon Kinesis.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// stream_name: The stream name.
	// region: The AWS region for the stream. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the stream.

	// aws_lambda_function
	// Amazon Lambda Function	A function in Amazon Lambda.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// function_name: The function name.
	// region: The AWS region for the function. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the function.

	// aws_rds_database
	// Amazon RDS Database	A database in Amazon Relational Database Service.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// name: The database name.
	// region: The AWS region for the database. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the database.

	// aws_redshift_cluster
	// Amazon Redshift Cluster	A cluster in Amazon Redshift.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// cluster_identifier: The cluster name.
	// region: The AWS region for the cluster. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the cluster.

	// aws_s3_bucket
	// Amazon S3 Bucket	A bucket in Amazon S3.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// bucket_name: The bucket name.
	// region: The AWS region for the bucket. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the bucket.

	// aws_ses
	// Amazon SES Region	An Amazon region with Amazon Simple Email Service enabled.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// region: The AWS region. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the SES region.

	// aws_sns_topic
	// Amazon SNS Topic	A topic in Amazon SNS.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// topic: The topic name.
	// region: The AWS region for the topic. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the topic.

	// aws_sqs_queue
	// Amazon SQS Queue	A queue in Amazon Simple Queue Service.
	//
	// project_id: The identifier of the GCP project under which data is stored for the AWS account specified in the aws_account label, such as "my-project".
	// queue: The queue name.
	// region: The AWS region for the queue. The format of this field is "aws:{region}", where supported values for {region} are listed at http://docs.aws.amazon.com/general/latest/gr/rande.html.
	// aws_account: The AWS account number for the queue.

	// bigquery_biengine_model
	// BigQuery BI Engine Model	BigQuery BI Engine Model.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the resource.
	// model_id: The identifier of the BI model.

	// bigquery_dataset
	// BigQuery Dataset	A dataset in BigQuery.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// dataset_id: The name of the BigQuery dataset.

	// bigquery_dts_config
	// BigQuery DTS Config	A BigQuery Data Transfer Service configuration.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the resource
	// config_id: The id of the DTS configuration.

	// bigquery_dts_run
	// BigQuery DTS Run	A BigQuery Data Transfer Service Run.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the resource
	// config_id: The name of the DTS config that created the run.
	// run_id: The unique resource name of the BigQuery DTS run.

	// bigquery_project
	// BigQuery Project	BigQuery Project.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: Location of the resource.

	// bigquery_resource
	// BigQuery	A BigQuery resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".

	// bigquery_table
	// BigQuery Table	An individual BigQuery table.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// dataset_id: The name of the BigQuery dataset.
	// table_id: The name of the BigQuery table.

	// billing_account
	// Cloud Billing Account	A Cloud Billing Account.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// account_id: The unique id of the billing account.

	// build
	// Cloud Build	A build in Cloud Build.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// build_id: The unique id of the build.
	// build_trigger_id: The unique id of the build trigger.
	// certificatemanager.googleapis.com/Project
	// Certificate Manager project	Certificate Manager project.
	// resource_container: The GCP container associated with the resource.
	// location: GCP location.

	// client_auth_config_brand
	// OAuth2 Brand	Consent screen data shown to users during three-legged OAuth2 flows.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// brand_id: The unique id of the brand.
	// client_auth_config_client
	// OAuth2 Client	A client used in OAuth2 flows.
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// client_id: The unique id of the client.

	// cloud_composer_environment
	// Cloud Composer Environment	A Composer environment runs the managed Apache Airflow service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Cloud Composer location in which the environment is running.
	// environment_name: The user-specified environment name.

	// cloud_dataproc_batch
	// Cloud Dataproc Batch	A Dataproc batch execution.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Cloud Dataproc region to which the batch was submitted.
	// batch_id: The user-specified batch id.
	// cloud_dataproc_cluster

	// Cloud Dataproc Cluster	A Dataproc cluster with separate cluster name and id labels.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// cluster_name: The user-specified cluster name.
	// cluster_uuid: The generated cluster id.
	// region: The Cloud Dataproc region in which the cluster is running.
	// cloud_dataproc_job
	// Cloud Dataproc Job	A Dataproc job execution.
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// region: The Cloud Dataproc region to which the job was submitted.
	// job_id: The user-specified job id.
	// job_uuid: The generated job uuid.

	// cloud_debugger_resource
	// Cloud Debugger	A Google Cloud Debugger resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// app: The application to which the debugger is attached.

	// CloudFunction is a function in Google Cloud Functions.
	//
	//  project_id
	// The identifier of the GCP project associated with this resource, such as "my-project".
	//
	//  function_name
	// The short function name.
	//
	//  region
	// The region in which the function is running.
	CloudFunction Type = "cloud_function"

	// CloudRunJob is a job in Cloud Run.
	//
	//  project_id
	// The identifier of the GCP project associated with this resource, such as "my-project".
	//
	//  job_name
	// Name of the monitored job.
	//
	//  location
	// Region where the job exists.
	CloudRunJob Type = "cloud_run_job"

	// CloudRunRevision is a revision in Cloud Run.
	//
	//  project_id
	// The identifier of the GCP project associated with this resource, such as "my-project".
	//
	//  service_name
	// Name of the service.
	//
	//  revision_name
	// Name of the monitored revision.
	//
	//  location
	// Region where the service is running.
	//
	//  configuration_name
	// Name of the configuration which created the monitored revision.
	CloudRunRevision Type = "cloud_run_revision"

	// cloud_scheduler_job
	// Cloud Scheduler Job	A Cloud Scheduler Job.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The region of the job.
	// job_id: Identifier of the job.

	// cloud_tasks_queue
	// Cloud Tasks Queue	A queue in Cloud Tasks.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// queue_id: The name of the queue.
	// target_type: The target type the queue is dispatching to.
	// location: The zone or region where the application is running.
	// clouddeploy.googleapis.com/DeliveryPipeline
	// Cloud Deploy Delivery Pipeline	A Cloud Deploy Delivery Pipeline.
	// resource_container: The identifier of the Google Cloud project associated with this resource.
	// location: The Google Cloud location where the resource resides.
	// pipeline_id: ID of the delivery pipeline resource.

	// cloudiot_device
	// Cloud IoT Device	A Device in Google Cloud IoT.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// device_num_id: The unique numeric identifier of the device.
	// device_registry_id: The user-defined string identifier of the device registry.
	// location: The cloud region of the device registry.

	// cloudiot_device_registry
	// Cloud IoT Registry	A Device Registry in Google Cloud IoT.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// device_registry_id: The user-defined string identifier of the device registry.
	// location: The cloud region of the device registry.

	// cloudkms_cryptokey
	// Cloud KMS CryptoKey	Cryptographic key in the KMS.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The region Crypto Key belongs to.
	// key_ring_id: Key Ring the Crypto Key belongs to.
	// crypto_key_id: Crypto Key Identifier.

	// cloudkms_cryptokeyversion
	// Cloud KMS CryptoKeyVersion	Version of a cryptographic key.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The region Crypto Key Version belongs to.
	// key_ring_id: Key Ring the Crypto Key Version belongs to.
	// crypto_key_id: Crypto Key the Crypto Key Version belongs to.
	// crypto_key_version_id: Crypto Key Version Identifier.

	// cloudkms_keyring
	// Cloud KMS Key Ring	Collection of cryptographic keys.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The region Key Ring belongs to.
	// key_ring_id: Key Ring Identifier.

	// cloudml_model_version
	// Cloud ML Model Version	A Google Cloud ML model version.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// model_id: An immutable identifier for a model.
	// version_id: An immutable identifier for a version.
	// region: Cloud ML region.

	// cloudsql_database
	// Cloud SQL Database	A database hosted in Google Cloud SQL.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// database_id: The ID of the database.
	// region: The Google Cloud SQL region in which the database is running.
	// cloudvolumesgcp-api.netapp.com/CloudVolume
	// Monitored Resource for NetApp CVS	Monitored Resource for NetApp CVS.
	// resource_container: Project information.
	// location: Region/Zone information.
	// volume_id: ID of the volume.
	// service_type: Service type of the volume or replication relationship.
	// name: Name of the volume or replication relationship.

	// consumed_api
	// Consumed API	An API used by customers.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as 'my-project'.
	// service: The API service name, such as 'cloudsql.googleapis.com'.
	// method: The API method name, such as 'disks.list'.
	// version: The API version, such as 'v1'.
	// location: The service specific notion of location. This can be a name of a zone or region. If a service does not have any notion of zones then 'global' can be used.
	// credential_id: The client credential ID, such as an API key ID or the OAuth client ID.

	// container
	// GKE Container	A Google Kubernetes Engine (GKE) container instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// cluster_name: An immutable name for the cluster the container is running in.
	// namespace_id: Immutable ID of the cluster namespace the container is running in.
	// instance_id: Immutable ID of the GCE instance the container is running in.
	// pod_id: Immutable ID of the pod the container is running in.
	// container_name: Immutable name of the container.
	// zone: The GCE zone in which the instance is running.

	// csr_repository
	// Cloud Source Repository	A repository in Google Cloud Source Repositories.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// name: The name of the repository.

	// dataflow_step
	// Dataflow Step	A step in a Dataflow job.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// job_id: The ID of the job.
	// step_id: The ID of the step.
	// job_name: The name of the job.
	// region: The region in which the job is running.
	// datamigration.googleapis.com/MigrationJob
	// Database migration service migration job	Database migration service migration job.
	// resource_container: The resource container (project ID).
	// location: The location.
	// migration_job_id: The migration job ID.

	// dataplex.googleapis.com/Environment
	// Cloud Dataplex Environment	An Environment within a Cloud Dataplex Lake.
	//
	// resource_container: The identifier of GCP project associated with this resource.
	// location: The GCP region associated with this resource.
	// lake_id: The identifier of the Lake resource containing this resource.
	// environment_id: The identifier of this Environment resource.

	// dataplex.googleapis.com/Lake
	// Cloud Dataplex Lake	A Cloud Dataplex Lake.
	//
	// resource_container: The identifier of GCP project associated with this resource.
	// location: The GCP region associated with this resource.
	// lake_id: The identifier of this Lake resource.

	// dataplex.googleapis.com/Task
	// Cloud Dataplex Task	A Task within a Cloud Dataplex Lake.
	//
	// resource_container: The identifier of GCP project associated with this resource.
	// location: The GCP region associated with this resource.
	// lake_id: The identifier of the Lake resource containing this resource.
	// task_id: The identifier of this Task resource.

	// dataplex.googleapis.com/Zone
	// Cloud Dataplex Zone	A Zone within a Cloud Dataplex Lake.
	//
	// resource_container: The identifier of GCP project associated with this resource.
	// location: The GCP region associated with this resource.
	// lake_id: The identifier of the Lake resource containing this resource.
	// zone_id: The identifier of this Zone resource.

	// dataproc_cluster
	// Dataproc Cluster	A Dataproc cluster.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// cluster_id: The cluster ID, concatenated from the cluster name and uuid
	// zone: The GCE zone in which the instance is running.

	// datastore_database
	// Cloud Datastore Database	A Cloud Datastore database.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// database_id: The unique id of the database.

	// datastore_index
	// Cloud Datastore Index	A Cloud Datastore index.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// database_id: The database the index belongs to.
	// index_id: The unique id of the index.

	// datastream.googleapis.com/Stream
	// Datastream Stream	A Datastream stream.
	//
	// resource_container: The resource container (project ID).
	// location: The location.
	// stream_id: The stream ID.

	// deployment
	// Deployment	A Deployment Manager deployment.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// name: Name that uniquely identifies the deployment within a project.

	// deployment_manager_manifest
	// Deployment Manager Manifest	A Deployment Manager manifest which is used to specify the contents of a deployment.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// manifest_name: Name that uniquely identifies the manifest within a project.
	// deployment_name: Name of the deployment.

	// deployment_manager_operation
	// Deployment Manager Operation	A Deployment Manager operation.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// operation_name: Name that uniquely identifies the operation within a project.

	// deployment_manager_resource
	// Deployment Manager Resource	Deployment Manager's record of Google Cloud Platform resources in a Deployment, such as a VM or a bucket.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// resource_name: Name of the resource, unique within a deployment.
	// deployment_name: Name of the deployment.

	// deployment_manager_type
	// Deployment Manager Type	A Deployment Manager type.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// name: Name that uniquely identifies the type within a project.

	// dns_managed_zone
	// Managed DNS Zone	A ManagedZone in the Google Cloud DNS service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// zone_name: The name of the ManagedZone.
	// location: The location field is provided for compatibility with other GCP services. Its value is always set to 'global'

	// dns_policy
	// Cloud DNS Policy	A Policy in the Google Cloud DNS service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// policy_name: The name of the Policy.
	// location: The location field is provided for compatibility with other GCP services. Its value is always set to 'global'

	// dns_query
	// Cloud DNS Query	A DNS query to a private DNS handled by the Google Cloud DNS service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// target_name: The DNS name managed by Cloud DNS to be resolved (e.g. the zone name, policy name, internal domain name). External names will have the value "external"
	// location: The GCP zone where the DNS request was received (e.g. us-east1, us-west1).
	// target_type: The target of the resolution of the DNS query (e.g. public-zone, private-zone, external).
	// source_type: Source of the query (e.g. gce-vm, internet).

	// firebase_domain
	// Firebase Hosting Site Domain	A domain from which a Firebase Hosting site is serving traffic.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as 'my-project'.
	// site_name: The name of a Firebase Hosting site, that is the subdomain in .web.app.
	// domain_name: The default subdomain (on web.app or firebaseapp.com) or custom domain from which content was served.

	// firebase_namespace
	// Firebase Realtime Database	A Firebase Realtime Database.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// table_name: The name of the database.
	// location: The location of the database.
	// fleetengine.googleapis.com/Fleet
	// Fleet Engine On Demand Rides and Deliveries	A top-level resource for Fleet Engine On Demand Rides and Deliveries metrics and logs.
	// resource_container: The identifier of the GCP container associated with Fleet.
	// location: The region in which the Fleet Engine instance is running.

	// folder
	// Google Folder	A Google Cloud Platform folder.
	//
	// folder_id: Numeric id of the folder.

	// GAEApp is an application running in Google App Engine (GAE).
	//
	//  project_id
	// The identifier of the GCP project associated with this resource, such as "my-project".
	//
	//  module_id
	// The service/module name.
	//
	//  version_id
	// The version name.
	//
	//  zone
	// The GAE zone where the application is running.
	GAEApp Type = "gae_app"

	// gateway_scope
	// Gateway Scope	GatewayScope represents a set of Gateways with the same merged configs.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The location of the control plane
	// scope: The name of the gateway_scope

	// gce_autoscaler
	// GCE Autoscaler	A Google Compute Engine (GCE) autoscaler.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// autoscaler_id: Unique identifier of the autoscaler.
	// location: GCE zone or region where the autoscaler is running.

	// gce_backend_bucket
	// GCE Backend Bucket	A Google Compute Engine (GCE) backend bucket.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// backend_bucket_id: Unique identifier of the backend bucket.

	// gce_backend_service
	// Compute Engine Backend Service	A Compute Engine backend service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// backend_service_id: Unique identifier of the backend service.
	// location: Global or Compute Engine region containing the backend service

	// gce_client_ssl_policy
	// GCE Client SSL Policy	A Google Compute Engine (GCE) client SSL policy.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// client_ssl_policy_id: Unique identifier of the client SSL policy.

	// gce_commitment
	// GCE Committed Use Discount	A Google Compute Engine (GCE) committed use discount.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// commitment_id: Unique identifier of the committed use discount.
	// location: GCE region where the committed use discount is active.

	// gce_disk
	// Disk	A disk belonging to a Compute Engine instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// disk_id: Unique identifier of the disk.
	// zone: The Compute Engine zone where the disk resides.

	// gce_firewall_rule
	// GCE Firewall Rule	A Google Compute Engine (GCE) firewall rule.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// firewall_rule_id: Unique identifier of the firewall rule.

	// gce_forwarding_rule
	// GCE Forwarding Rule	A Google Compute Engine (GCE) Forwarding Rule.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// forwarding_rule_id: Unique identifier of the forewarding rule.
	// region: GCE region where the forwarding rule resides.

	// gce_health_check
	// GCE Health Check	A Google Compute Engine (GCE) health check.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// health_check_id: Unique identifier of the health check.

	// gce_image
	// GCE Image	A Google Compute Engine (GCE) image resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// image_id: Unique numerical identifier of the image.

	// gce_instance
	// VM Instance	A virtual machine instance hosted in Compute Engine.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// instance_id: The numeric VM instance identifier assigned by Compute Engine.
	// zone: The Compute Engine zone in which the VM is running.

	// gce_instance_group
	// GCE Instance Group	A Google Compute Engine (GCE) instance group resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// instance_group_id: The unique numerical identifier of the instance group.
	// instance_group_name: The unique user provided name of the instance group.
	// location: GCE zone containing the instance group.

	// gce_instance_group_manager
	// GCE Instance Group Manager	A Google Compute Engine (GCE) instance group manager resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// instance_group_manager_id: The unique numerical identifier of the instance group manager.
	// instance_group_manager_name: The unique user provided name of the instance group manager.
	// location: GCE zone or region where the instance group manager is located.

	// gce_instance_template
	// GCE Instance Template	A Google Compute Engine (GCE) instance template resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// instance_template_id: The unique numerical identifier of the instance template.
	// instance_template_name: The unique user provided name of the instance template.

	// gce_license
	// GCE License	A Google Compute Engine (GCE) license.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// license_id: Unique identifier of the license.

	// gce_network
	// GCE Network	A Google Compute Engine (GCE) network.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// network_id: Unique identifier of the network.

	// gce_network_endpoint_group
	// Network Endpoint Group	A Compute Engine network endpoint group resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// zone: The name of the zone where the network endpoint group is located.
	// network_endpoint_group_id: The ID of the network endpoint group.

	// gce_network_region
	// Network Region	A region of a Compute Engine network.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// network_id: The ID of the Compute Engine network.
	// region: The name of the network region.
	// gce_node_group
	// Node Group	A Compute Engine node group.
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// node_group_id: Unique identifier of the node group.
	// zone: Zone of the node group.

	// gce_node_template
	// Node Template	A Compute Engine node template.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// node_template_id: Unique identifier of the node template.
	// region: Region of the node template.

	// gce_operation
	// GCE Operation	A Google Compute Engine (GCE) operation resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// operation_name: The unique user provided name of the operation.
	// location: Location of the resource.

	// gce_packet_mirroring
	// GCE Packet Mirroring	A Google Compute Engine (GCE) packet mirroring.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// packet_mirroring_id: Unique identifier of the packet mirroring.
	// region: Region of the packet mirroring.
	// gce_project
	// GCE Project	A Google Compute Engine (GCE) project resource.
	// project_id: GCE specific numeric identifier of the GCE project resource.

	// gce_reserved_address
	// GCE Reserved Address	A Google Compute Engine reserved address.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// reserved_address_id: Unique identifier of the reserved address.
	// location: Global or GCE region containing the reserved address

	// gce_resource_policy
	// Resource Policy	A Compute Engine resource policy.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// resource_policy_id: Unique identifier of the resource policy.
	// region: Region of the resource policy.

	// gce_route
	// GCE Route	A Google Compute Engine (GCE) route.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// route_id: Unique identifier of the route.

	// gce_router
	// GCE Router	A Google Compute Engine (GCE) router.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// router_id: Unique identifier of the router.
	// region: Region of the router.

	// gce_snapshot
	// GCE Snapshot	A Google Compute Engine (GCE) snapshot.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// snapshot_id: Unique identifier of the snapshot.

	// gce_ssl_certificate
	// GCE SSL Certificate	A Google Compute Engine (GCE) SSL certificate.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// ssl_certificate_id: The unique numerical identifier of the SSL certificate.
	// ssl_certificate_name: The unique user provided name of the SSL Certificate.

	// gce_subnetwork
	// Subnetwork	A Compute Engine subnetwork.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// subnetwork_id: The unique numerical identifier of the subnetwork.
	// subnetwork_name: The unique user provided name of the subnetwork.
	// location: Location of the resource.

	// gce_target_http_instance
	// GCE Target HTTP Instance	A Google Compute Engine (GCE) target http instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// target_http_instance_id: Unique identifier of the target http instance.
	// zone: GCE zone where the target http instance resides.

	// gce_target_http_proxy
	// GCE Target HTTP Proxy	A Google Compute Engine (GCE) target http proxy.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// target_http_proxy_id: Unique identifier of the target http proxy.

	// gce_target_https_proxy
	// GCE Target HTTPS Proxy	A Google Compute Engine (GCE) target https proxy.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// target_https_proxy_id: Unique identifier of the target https proxy.

	// gce_target_pool
	// GCE Target Pool	A Google Compute Engine (GCE) target pool.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// target_pool_id: Unique identifier of the target pool.
	// zone: GCE zone where the pool resides.

	// gce_target_ssl_proxy
	// GCE Target SSL Proxy	A Google Compute Engine (GCE) target SSL proxy.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// target_ssl_proxy_id: Unique identifier of the target ssl proxy.

	// gce_url_map
	// GCE URL Map	A Google Compute Engine (GCE) URL map.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// url_map_id: Unique identifier of the url map.

	// gcs_bucket
	// GCS Bucket	A Google Cloud Storage (GCS) bucket.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// bucket_name: An immutable name of the bucket.
	// location: Location of the bucket.

	// generic_node
	// Generic Node	A generic node identifies a machine or other computational resource for which no more specific resource type is applicable. The label values must uniquely identify the node.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The GCP or AWS region in which data about the resource is stored. For example, "us-east1-a" (GCP) or "aws:us-east-1a" (AWS).
	// namespace: A namespace identifier, such as a cluster name.
	// node_id: A unique identifier for the node within the namespace, such as a hostname or IP address.

	// generic_task
	// Generic Task	A generic task identifies an application process for which no more specific resource is applicable, such as a process scheduled by a custom orchestration system. The label values must uniquely identify the task.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The GCP or AWS region in which data about the resource is stored. For example, "us-east1-a" (GCP) or "aws:us-east-1a" (AWS).
	// namespace: A namespace identifier, such as a cluster name.
	// job: An identifier for a grouping of related tasks, such as the name of a microservice or distributed batch job.
	// task_id: A unique identifier for the task within the namespace and job, such as a replica index identifying the task within the job.

	// genomics_dataset
	// Genomics Dataset	A dataset in the Google Genomics service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// dataset_id: Unique identifier of the dataset.

	// genomics_operation
	// Genomics Operation	A long running operation in the Google Genomics service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// operation_id: Unique identifier of the long running operation.

	// gke_cluster
	// GKE Cluster Operations	A Google Kubernetes Engine (GKE) Cluster. It contains events and audit logs about cluster operations.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// cluster_name: The name of the GKE Cluster.
	// location: The location in which the GKE Cluster is running.

	// gke_nodepool
	// GKE Node Pool Operations	A Google Kubernetes Engine (GKE) Node Pool. It contains audit logs about Node Pool operations.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// nodepool_name: The name of the GKE Node Pool.
	// location: The location in which the GKE Cluster is running.
	// cluster_name: The name of the GKE Cluster to which this Node Pool belongs.

	// gkebackup.googleapis.com/BackupPlan
	// GKE Backup Plan	A backup plan provides configuration, location, and management functions for a sequence of backups.
	//
	// resource_container: The identifier of the Google Cloud container associated with the resource.
	// location: The Google Cloud location where this backupPlan resides.
	// backup_plan_id: The name of the backupPlan.

	// gkebackup.googleapis.com/RestorePlan
	// GKE Restore Plan	A restore plan defines the configuration of a series of restore operations to be performed against backups which belong to the specified backup plan.
	//
	// resource_container: The identifier of the Google Cloud container associated with the resource.
	// location: The Google Cloud location where this restorePlan resides.
	// restore_plan_id: The name of the restorePlan.

	// global
	// Global	A resource type used to indicate that a log is not associated with any specific resource.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".

	// healthcare_annotation_store
	// Healthcare Annotation Store	A Cloud Healthcare Annotation store containing Annotation records.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the dataset that contains the Annotation store.

	// dataset_id: The ID of the dataset.
	//
	// annotation_store_id: The ID of the Annotation store.

	// healthcare_consent_store
	// Healthcare Consent Store	A Cloud Healthcare Consent store containing consent records.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the dataset that contains the Consent store.
	// dataset_id: The ID of the dataset.
	// consent_store_id: The ID of the Consent store.

	// healthcare_dataset
	// Healthcare Dataset	A Cloud Healthcare dataset.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the dataset.
	// dataset_id: The ID of the dataset.

	// healthcare_dicom_store
	// Healthcare DICOM Store	A Cloud Healthcare DICOM store containing DICOM instances.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the dataset that contains the DICOM store.
	// dataset_id: The ID of the dataset.
	// dicom_store_id: The ID of the DICOM store.

	// healthcare_fhir_store
	// Healthcare FHIR Store	A Cloud Healthcare FHIR store containing FHIR resources representing electronic medical information.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the dataset that contains the FHIR store.
	// dataset_id: The ID of the dataset.
	// fhir_store_id: The ID of the FHIR store.

	// healthcare_hl7v2_store
	// Healthcare HL7v2 Store	A Cloud Healthcare HL7v2 store containing clinical messages.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The Google Cloud location of the dataset that contains the HL7v2 store.
	// dataset_id: The ID of the dataset.
	// hl7v2_store_id: The ID of the HL7v2 store.

	// http_external_regional_lb_rule
	// HTTP/S External Regional Load Balancing Rule	A resource descriptor for HTTP/S External Regional load balancing behavior.
	//
	// project_id: The identifier of the Google Cloud project associated with this resource, such as 'my-project'.
	// network_name: The name of the customer network in which the Load Balancer resides.
	// region: The region under which the Load Balancer is defined.
	// url_map_name: The name of the urlmap.
	// forwarding_rule_name: The name of the forwarding rule.
	// target_proxy_name: The name of the target HTTP/S proxy.
	// matched_url_path_rule: The prefix of URL defined in urlmap tree. 'UNMATCHED' for the sink default rule.
	// backend_target_name: The name of the backend target or service.
	// backend_target_type: The type of the backend target. Can be 'BACKEND_SERVICE', or 'UNKNOWN' if the backend wasn't assigned.
	// backend_name: The name of the backend group. Can be '' if the backend wasn't assigned.
	// backend_type: The type of the backend group. Can be 'INSTANCE_GROUP', 'NETWORK_ENDPOINT_GROUP', or 'UNKNOWN' if the backend wasn't assigned.
	// backend_scope: The scope of the backend group. Can be 'UNKNOWN' if the backend wasn't assigned.
	// backend_scope_type: The type of the scope of the backend group. Can be 'ZONE', 'REGION', or 'UNKNOWN' in case the backend wasn't assigned.

	// http_load_balancer
	// Cloud HTTP Load Balancer	A Cloud HTTP Load Balancer Instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// forwarding_rule_name: The name of the forwarding rule.
	// url_map_name: The name of the urlmap.
	// target_proxy_name: The name of the target proxy.
	// backend_service_name: The name of the backend service.
	// zone: The zone in which the load balancer is running.

	// iam_role
	// IAM Role	An IAM role.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// role_name: The name of the IAM custom role; this labelappears only on custom roles.(e.g., roles/[CUSTOM_ROLE],organizations/123456/roles/[CUSTOM_ROLE],projects/myproject/roles/[CUSTOM_ROLE]).

	// identitytoolkit_project
	// Project	An Identity Toolkit project.
	//
	// project_id: The identifier of the GCP project associated with this resource.

	// identitytoolkit_tenant
	// Identity Toolkit Tenant	An Identity Toolkit tenant.
	//
	// project_id: The identifier of the GCP project associated with this resource.
	// tenant_name: The name of the tenant.

	// ids.googleapis.com/Endpoint
	// IDS Endpoint	A Cloud IDS Endpoint.
	//
	// resource_container: The identifier of the GCP project owning the Endpoint.
	// location: The zone of the IDS Endpoint.
	// id: The ID of the Endpoint.

	// istio_control_plane
	// Istio Control Plane	An Istio Control Plane is an instance of a service that provides xDS and related functionality to a set of managed Istio proxies.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// mesh_uid: Unique identifier for an Istio service mesh.
	// location: The physical location in which the workload for the Control Plane is located.
	// revision: Immutable revision of Istio managed by the Control Plane.
	// build_id: Immutable build tag for the instance of the Control Plane.
	// owner: Immutable name of the owner of the Control Plane.

	// k8s_cluster
	// Kubernetes Cluster	A Kubernetes cluster. It contains Kubernetes audit logs from the cluster.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The physical location of the cluster.
	// cluster_name: The name of the cluster.

	// k8s_container
	// Kubernetes Container	A Kubernetes container instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The physical location of the cluster that contains the container.
	// cluster_name: The name of the cluster that the container is running in.
	// namespace_name: The name of the namespace that the container is running in.
	// pod_name: The name of the pod that the container is running in.
	// container_name: The name of the container.

	// k8s_control_plane_component
	// Kubernetes Control Plane Component	A Kubernetes Control Plane component.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The physical location of the cluster that contains the control plane component.
	// cluster_name: The name of the cluster that the control plane component is running in.
	// component_name: The name of the control plane component.
	// component_location: The physical location where the control plane component is running.

	// k8s_node
	// Kubernetes Node	A Kubernetes node instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The physical location of the cluster that contains the node.
	// cluster_name: The name of the cluster that the node is a part of.
	// node_name: The name of the node.

	// k8s_pod
	// Kubernetes Pod	A Kubernetes pod instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The physical location of the cluster that contains the pod.
	// cluster_name: The name of the cluster that the pod is running in.
	// namespace_name: The name of the namespace that the pod is running in.
	// pod_name: The name of the pod.

	// k8s_service
	// Kubernetes Service	A Kubernetes Service instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The physical location of the cluster that contains the service.
	// cluster_name: The name of the cluster that the service is running in.
	// namespace_name: The name of the namespace that the service is running in.
	// service_name: The name of the service.

	// l4_proxy_rule
	// Layer 4 Proxying Rule for TCP/UDP/SSL Traffic	A resource descriptor for TCP/SSL/UDP Internal Regional load balancing behavior.
	//
	// project_id: The identifier of the Google Cloud project associated with this resource, such as 'my-project'.
	// network_name: The name of the customer network in which the Load Balancer resides.
	// region: The region under which the Load Balancer is defined.
	// load_balancing_scheme: The load balancing scheme associated with the forwarding rule, one of [INTERNAL_MANAGED, EXTERNAL_MANAGED].
	// protocol: The protocol associated with the traffic processed by the proxy, one of [TCP, UDP, SSL, UNKNOWN].
	// forwarding_rule_name: The name of the forwarding rule.
	// target_proxy_name: The name of the target proxy.
	// backend_target_name: The name of the backend target or service.
	// backend_target_type: The type of the backend target, one of ['BACKEND_SERVICE'; 'UNKNOWN' - if the backend wasn't assigned].
	// backend_name: The name of the backend group. Can be '' if the backend wasn't assigned.
	// backend_type: The type of the backend group, one of ['INSTANCE_GROUP'; 'NETWORK_ENDPOINT_GROUP'; 'UNKNOWN' - if the backend wasn't assigned].
	// backend_scope: The scope of the backend group. Can be 'UNKNOWN' if the backend wasn't assigned.
	// backend_scope_type: The type of the scope of the backend group, one of ['ZONE'; 'REGION'; 'UNKNOWN' - in case the backend wasn't assigned].

	// livestream.googleapis.com/Channel
	// Live Stream API Channel	A Live Stream API Channel.
	//
	// resource_container: The identifier of the GCP project associated with this channel resource.
	// location: The GCP location where the channel resource resides.
	// channel_id: ID of the channel resource.

	// loadbalancing.googleapis.com/ExternalNetworkLoadBalancerRule
	// Google Cloud External Network Load Balancer Rule	A set of definitions for multi protocol network load balancing behavior.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// region: The Google Cloud Platform region of the backend instance that connected to network load balancing forwarding rule.
	// backend_network_name: The network name of the NIC of the instance that received the Net LB flow.
	// backend_target_type: The type of the backend target that handled the connection.
	// backend_service_name: The name of the backend service that handled the connection.
	// primary_target_pool: The name of the primary target pool.
	// target_pool: The name of the target pool.
	// forwarding_rule_name: The name of the forwarding rule.
	// backend_group_name: The name of the backend group that handled the connection.
	// backend_group_type: The type of the backend group that handled the connection.
	// backend_group_scope: The scope (zone or region) of the backend group that handled the connection.
	// backend_subnetwork_name: The name of the subnetwork of the instance that handled the connection.
	// backend_zone: The zone of the endpoint (VM instance) that handled the connection.

	// loadbalancing.googleapis.com/InternalNetworkLoadBalancerRule
	// Google Cloud Internal Network Load Balancer Rule	A set of definitions for multi protocol internal load balancing behavior.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// region: The Google Cloud Platform region of the backend instance that connected to network load balancing forwarding rule.
	// backend_network_name: The network name of the NIC of the instance that received the Net LB flow.
	// backend_service_name: The name of the backend service that handled the connection.
	// forwarding_rule_name: The name of the forwarding rule.
	// backend_group_name: The name of the backend group that handled the connection.
	// backend_group_type: The type of the backend group that handled the connection.
	// backend_group_scope: The scope (zone or region) of the backend group that handled the connection.
	// backend_subnetwork_name: The name of the subnetwork of the instance that handled the connection.

	// logging_bucket
	// Logging Bucket	An export bucket in Cloud Logging.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// bucket_id: The name of the log bucket.
	// location: The location of the log bucket.
	// source_resource_container: The source resource container (e.g. project, folder, organization) of the log entry that is destined for the log bucket. The format is "projects/project_id"
	// monitored_resource_type: The type field of the monitored resource in the log entry that is destined for the log bucket.

	// logging_exclusion
	// Log Exclusion	An exclusion in Cloud Logging.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// name: The unique name of the exclusion.
	// logging_log
	// Log stream	A Google Cloud Logging log.
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// name: Unique identifier of the log.

	// logging_sink
	// Logging export sink	An export sink in Cloud Logging.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// name: The unique name of the sink.
	// destination: The destination of the sink.

	// managed_service
	// Managed Service	A service managed by Google Service Management.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// service_name: The name of the service.
	// producer_project_id: The id of the project which produces and owns this service.

	// mesh
	// Mesh	A mesh serves as the "key" to deliver configuration to data plane proxy instances.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The location of the control plane
	// mesh: The name of the mesh

	// metastore.googleapis.com/Service
	// Dataproc Metastore Service	A Dataproc Metastore Service.
	//
	// resource_container: The ID of the customer project.
	// location: The region that the service is hosted in.
	// service_id: The service ID.

	// metric
	// Metric Type	A Stackdriver Monitoring metric type.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// name: The name of the metric type, such as "logging.googleapis.com/my-metric-name".

	// ml_job
	// Cloud ML Job	A Cloud Machine Learning job.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// job_id: The job identifier.
	// task_name: The task name.

	// nat_gateway
	// Cloud NAT Gateway	A Cloud NAT Gateway.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// region: The region where the NAT gateway is located.
	// router_id: Identifier of the router under which the NAT gateway is defined.
	// gateway_name: The name of the NAT gateway.

	// network_security_policy
	// Network Security Policy	A network security policy.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// policy_name: The unique user provided name of the security policy.
	// networking.googleapis.com/Location
	// GCP Location	A GCP location: a specific zone or region, or "global".
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: Name of a GCP zone/region, or "global".

	// organization
	// Google Organization	A Google Cloud Platform organization.
	//
	// organization_id: Numeric id of the organization.

	// project
	// Google Project	A Google project.
	//
	// project_id: The identifier of the GCP project associated with this resource (e.g., my-project).

	// pubsub_snapshot
	// Cloud Pub/Sub Snapshot	A snapshot in Google Cloud Pub/Sub.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// snapshot_id: The identifier of the snapshot, such as "my-snapshot".

	// pubsub_subscription
	// Cloud Pub/Sub Subscription	A subscription in Google Cloud Pub/Sub.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// subscription_id: The identifier of the subscription, such as "my-subscription".

	// pubsub_topic
	// Cloud Pub/Sub Topic	A topic in Google Cloud Pub/Sub.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// topic_id: The identifier of the topic, such as "my-topic".
	// recaptchaenterprise.googleapis.com/Key
	// reCAPTCHA Key	Monitoring resource for reCAPTCHA Key.
	// resource_container: The ID of the GCP project associated with this reCAPTCHA Key.
	// location: Location where the reCAPTCHA Key is provisioned.
	// key_id: The ID for this Key.

	// recommender
	// Recommender	A Recommender represents a grouping of similar recommendations.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// recommender_id: The name of the recommender.
	// location: The location of the recommendation.

	// recommender_insight_type
	// InsightType	An InsightType represents a grouping of similar insights.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// insight_type_id: The resource ID of the insight type.
	// location: The location of the insight.

	// redis_instance
	// Cloud Memorystore Redis Instance	A Redis instance hosted on Google Cloud Memorystore.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// region: The Google Cloud region in which the managed instance is running.
	// instance_id: The ID of the managed instance.
	// node_id: The ID of a Redis node within the managed instance.

	// reported_errors
	// Reported Errors	Error data and metadata managed by Stackdriver Error Reporting
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".

	// secretmanager.googleapis.com/Secret
	// Secret Manager Secret	A logical secret whose value and versions can be accessed.
	//
	// resource_container: The identifier of the GCP project associated with this resource.
	// location: Location of secret metadata. Always global.
	// secret_id: The name given to this secret.

	// service_account
	// Service Account	A service account.
	//
	// project_id: The identifier of the GCP project associated with this resource (e.g., my-project).
	// email_id: The service account email id, e.g. "account123@proj123.iam.gserviceaccount.com".
	// unique_id: The unique id of the service account, e.g. "113948692397867021414".

	// service_config
	// Service Configuration	A specific service configuration.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// service_name: The name of the service.
	// service_config_id: The id of the service configuration.

	// service_rollout
	// Service Rollout	A resource type used to describe how a service configuration is deployed to backend systems.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// service_name: The name of the service.
	// rollout_id: The id of the service rollout.

	// servicedirectory_namespace
	// Service Directory Namespace	A namespace in the Service Directory service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The cloud region of the Service Directory namespace.
	// namespace_name: The name of the Service Directory namespace.

	// serviceusage_service
	// Service	A service activated or deactivated by a consumer project.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// service_name: The name of the service.

	// serviceuser_service
	// Service	A service activated or deactivated by a consumer project.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// service_name: The name of the service.

	// spanner_instance
	// Cloud Spanner Instance	A Cloud Spanner instance.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// instance_id: An immutable identifier for an instance.
	// location: Cloud Spanner region.
	// instance_config: Instance config for the instance.

	// storage_transfer_job
	// Cloud Storage Transfer Job	A Google Cloud storage transfer job.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// job_id: A unique name of the storage transfer job.

	// tcp_ssl_proxy_rule
	// Google Cloud TCP/SSL Proxy Rule	A set of definitions for TCP/SSL proxy behavior.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// region: The region on which TCP/SSL proxy is applied, such as 'global' or 'us-central1'. Various other objects are defined per that locality.
	// backend_target_name: The name of the backend target ('backend service', equivalent to 'proxy name').
	// backend_target_type: The type of the backend target. Can only be 'BACKEND_SERVICE' currently.
	// forwarding_rule_name: The name of the forwarding rule.
	// target_proxy_name: The name of the target TCP/SSL proxy.
	// backend_name: The name of the backend group.
	// backend_type: The type of the backend group. Can be 'INSTANCE_GROUP' or 'NETWORK_ENDPOINT_GROUP'.
	// backend_scope: The scope (zone or region) of the backend group.
	// backend_scope_type: The type of the scope of the backend group. Can be either 'ZONE' or 'REGION'.

	// testservice_matrix
	// Test Matrix	A Test Matrix in the Google Cloud Test Lab service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// matrix_id: Unique identifier of the matrix.

	// threat_detector
	// Threat Detector	A detector in the Threat Detection service.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// detector_name: The specific detector that triggered the alert.

	// uptime_url
	// Uptime Check URL	An Uptime Monitoring check against a custom URL.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// host: The hostname or IP address of the check.

	// vmmigration.googleapis.com/MigratingVM
	// Migrate to Virtual Machines Migrating VM	A Migrate to Virtual Machines Migrating VM.
	//
	// resource_container: The identifier of the GCP project associated with this VM resource.
	// location: The GCP location where the VM resource resides.
	// source: The source where the VM resource resides.
	// vm: The VM ID.

	// vmmigration.googleapis.com/Source
	// Migrate to Virtual Machines Source	A Migrate to Virtual Machines Source.
	//
	// resource_container: The identifier of the GCP project associated with this source resource.
	// location: The GCP location where the source resource resides.
	// source: The source ID.

	// vpc_access_connector
	// VPC Access Connector	A connector that can communicate with devices within a VPC.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// location: The region the connector is located in.
	// connector_name: The name of the connector.

	// vpn_gateway
	// Cloud VPN Gateway	A Cloud VPN gateway.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// gateway_id: The VPN Gateway ID.
	// region: The region in which the VPN Gateway is running.

	// vpn_tunnel
	// Cloud VPN Tunnel	A Cloud VPN tunnel.
	//
	// project_id: The identifier of the GCP project associated with this resource, such as "my-project".
	// tunnel_id: The unique numerical identifier of the VPN tunnel.
	// tunnel_name: The unique user provided name of the VPN tunnel.
	// location: Location of the Cloud VPN Tunnel.

	// workflows.googleapis.com/Workflow
	// Workflow	A Workflows specification of steps to execute.
	//
	// resource_container: The identifier of the GCP container associated with the resource.
	// location: The region in which the workflow is deployed.
	// workflow_id: The ID of the workflow.
)

type resource struct {
	pb    *mrpb.MonitoredResource
	attrs detector.ResourceAttributesFetcher
	once  *sync.Once
}

func (r *resource) metadataProjectID() string {
	return r.attrs.Metadata("project/project-id")
}

func (r *resource) metadataZone() string {
	zone := r.attrs.Metadata("instance/zone")
	if zone != "" {
		return zone[strings.LastIndex(zone, "/")+1:]
	}
	return ""
}

func (r *resource) metadataRegion() string {
	region := r.attrs.Metadata("instance/region")
	if region != "" {
		return region[strings.LastIndex(region, "/")+1:]
	}
	return ""
}

// isMetadataActive queries valid response on "/computeMetadata/v1/" URL
func (r *resource) isMetadataActive() bool {
	data := r.attrs.Metadata("")
	return data != ""
}

var resourceDetector = &resource{
	attrs: detector.ResourceAttributes(),
	once:  new(sync.Once),
}

func detectCloudRunResource() *mrpb.MonitoredResource {
	projectID := resourceDetector.metadataProjectID()
	if projectID == "" {
		return nil
	}

	region := resourceDetector.metadataRegion()
	config := resourceDetector.attrs.EnvVar(detector.EnvCloudRunConfig)
	service := resourceDetector.attrs.EnvVar(detector.EnvCloudRunService)
	revision := resourceDetector.attrs.EnvVar(detector.EnvCloudRunRevision)

	return &mrpb.MonitoredResource{
		Type: string(CloudRunRevision),
		Labels: Label{
			"project_id":         projectID,
			"service_name":       service,
			"revision_name":      revision,
			"location":           region,
			"configuration_name": config,
		},
	}
}
