package gcp

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps Google Cloud operations
type Client struct {
	Project string
	Region  string
}

// NewClient creates a new GCP client
func NewClient(project, region string) *Client {
	if project == "" {
		project = getDefaultProject()
	}
	if region == "" {
		region = "us-central1"
	}
	return &Client{
		Project: project,
		Region:  region,
	}
}

// execGcloud runs a gcloud command and returns output
func (c *Client) execGcloud(args ...string) (string, error) {
	// Add project and region to command if not already present
	fullArgs := []string{"--project=" + c.Project}
	fullArgs = append(fullArgs, args...)

	cmd := exec.Command("gcloud", fullArgs...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("gcloud error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// getDefaultProject gets the default GCP project
func getDefaultProject() string {
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}

// Compute Engine operations
func (c *Client) ListInstances() (string, error) {
	return c.execGcloud("compute", "instances", "list", "--format=table(name,zone,machineType,status)")
}

func (c *Client) DescribeInstance(name string) (string, error) {
	return c.execGcloud("compute", "instances", "describe", name, "--zone="+c.Region)
}

func (c *Client) StartInstance(name string) (string, error) {
	return c.execGcloud("compute", "instances", "start", name, "--zone="+c.Region)
}

func (c *Client) StopInstance(name string) (string, error) {
	return c.execGcloud("compute", "instances", "stop", name, "--zone="+c.Region)
}

func (c *Client) DeleteInstance(name string) (string, error) {
	return c.execGcloud("compute", "instances", "delete", name, "--zone="+c.Region, "--quiet")
}

func (c *Client) ListImages() (string, error) {
	return c.execGcloud("compute", "images", "list", "--format=table(name,status,deprecated)")
}

// Cloud Storage operations
func (c *Client) ListBuckets() (string, error) {
	return c.execGcloud("storage", "buckets", "list", "--format=table(name,location,storageClass)")
}

func (c *Client) ListBucketContents(bucket string) (string, error) {
	return c.execGcloud("storage", "objects", "list", "gs://"+bucket, "--format=table(name,size,timeCreated)")
}

func (c *Client) CreateBucket(name, location string) (string, error) {
	if location == "" {
		location = "us-central1"
	}
	return c.execGcloud("storage", "buckets", "create", "gs://"+name, "--location="+location)
}

func (c *Client) DeleteBucket(name string) (string, error) {
	return c.execGcloud("storage", "buckets", "delete", "gs://"+name, "--quiet")
}

func (c *Client) GetBucketSize(bucket string) (string, error) {
	return c.execGcloud("storage", "du", "gs://"+bucket, "--summarize")
}

// Cloud SQL operations
func (c *Client) ListSQLInstances() (string, error) {
	return c.execGcloud("sql", "instances", "list", "--format=table(name,databaseVersion,state,ipAddresses[0].ipAddress)")
}

func (c *Client) DescribeSQLInstance(name string) (string, error) {
	return c.execGcloud("sql", "instances", "describe", name)
}

func (c *Client) StartSQLInstance(name string) (string, error) {
	return c.execGcloud("sql", "instances", "patch", name, "--activation-policy=ALWAYS")
}

func (c *Client) StopSQLInstance(name string) (string, error) {
	return c.execGcloud("sql", "instances", "patch", name, "--activation-policy=NEVER")
}

// Cloud Run operations
func (c *Client) ListCloudRunServices() (string, error) {
	return c.execGcloud("run", "services", "list", "--region="+c.Region, "--format=table(metadata.name,status.url,status.conditions[0].lastTransitionTime)")
}

func (c *Client) DescribeCloudRunService(name string) (string, error) {
	return c.execGcloud("run", "services", "describe", name, "--region="+c.Region)
}

// BigQuery operations
func (c *Client) ListDatasets() (string, error) {
	return c.execGcloud("bq", "ls", "--format=prettyjson")
}

func (c *Client) GetDatasetSize(dataset string) (string, error) {
	return c.execGcloud("bq", "show", "--format=prettyjson", dataset)
}

// IAM operations
func (c *Client) ListServiceAccounts() (string, error) {
	return c.execGcloud("iam", "service-accounts", "list", "--format=table(email,displayName,disabled)")
}

func (c *Client) ListIAMBindings() (string, error) {
	return c.execGcloud("projects", "get-iam-policy", c.Project, "--format=json")
}

// Project and configuration
func (c *Client) GetProjectInfo() (string, error) {
	return c.execGcloud("projects", "describe", c.Project, "--format=json")
}

func (c *Client) GetActiveServices() (string, error) {
	return c.execGcloud("services", "list", "--enabled", "--format=table(config.name,state)")
}

func (c *Client) GetBillingInfo() (string, error) {
	return c.execGcloud("billing", "budgets", "list", "--format=table(displayName,budgetAmount.currencyCode,budgetAmount.nanos,thresholdRule[0].thresholdPercent)")
}
