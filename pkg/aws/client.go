package aws

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps AWS CLI operations
type Client struct {
	Region  string
	Profile string
}

// NewClient creates a new AWS client
func NewClient(region, profile string) *Client {
	return &Client{
		Region:  region,
		Profile: profile,
	}
}

// execAWS runs an AWS CLI command and returns output
func (c *Client) execAWS(service string, args ...string) (string, error) {
	cmdArgs := []string{service}
	if c.Profile != "" {
		cmdArgs = append(cmdArgs, "--profile", c.Profile)
	}
	if c.Region != "" {
		cmdArgs = append(cmdArgs, "--region", c.Region)
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("aws", cmdArgs...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("aws error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// GetCallerIdentity returns current AWS account info
func (c *Client) GetCallerIdentity() (string, error) {
	return c.execAWS("sts", "get-caller-identity", "--output", "table")
}

// EC2 Operations

// ListEC2Instances lists all EC2 instances
func (c *Client) ListEC2Instances(running bool) (string, error) {
	args := []string{"describe-instances", "--query", "Reservations[*].Instances[*].[InstanceId,InstanceType,State.Name,PrivateIpAddress,PublicIpAddress]", "--output", "table"}
	if running {
		args = []string{"describe-instances", "--filters", "Name=instance-state-name,Values=running", "--query", "Reservations[*].Instances[*].[InstanceId,InstanceType,State.Name,PrivateIpAddress,PublicIpAddress]", "--output", "table"}
	}
	return c.execAWS("ec2", args...)
}

// StopEC2Instance stops an EC2 instance
func (c *Client) StopEC2Instance(instanceID string) (string, error) {
	return c.execAWS("ec2", "stop-instances", "--instance-ids", instanceID, "--output", "table")
}

// StartEC2Instance starts an EC2 instance
func (c *Client) StartEC2Instance(instanceID string) (string, error) {
	return c.execAWS("ec2", "start-instances", "--instance-ids", instanceID, "--output", "table")
}

// TerminateEC2Instance terminates an EC2 instance
func (c *Client) TerminateEC2Instance(instanceID string) (string, error) {
	return c.execAWS("ec2", "terminate-instances", "--instance-ids", instanceID, "--output", "table")
}

// DescribeEC2Instance gets details about an EC2 instance
func (c *Client) DescribeEC2Instance(instanceID string) (string, error) {
	return c.execAWS("ec2", "describe-instances", "--instance-ids", instanceID, "--output", "table")
}

// S3 Operations

// ListS3Buckets lists all S3 buckets
func (c *Client) ListS3Buckets() (string, error) {
	return c.execAWS("s3", "ls")
}

// ListS3BucketContents lists contents of an S3 bucket
func (c *Client) ListS3BucketContents(bucket string, recursive bool) (string, error) {
	args := []string{"ls", "s3://" + bucket}
	if recursive {
		args = append(args, "--recursive")
	}
	return c.execAWS("s3", args...)
}

// GetS3BucketSize gets the size of an S3 bucket
func (c *Client) GetS3BucketSize(bucket string) (string, error) {
	return c.execAWS("s3", "ls", "s3://"+bucket, "--recursive", "--summarize")
}

// DeleteS3Object deletes an object from S3
func (c *Client) DeleteS3Object(bucket, key string) (string, error) {
	return c.execAWS("s3", "rm", "s3://"+bucket+"/"+key)
}

// SyncS3 syncs a local directory to S3
func (c *Client) SyncS3(localPath, s3Path string, delete bool) (string, error) {
	args := []string{"sync", localPath, s3Path}
	if delete {
		args = append(args, "--delete")
	}
	return c.execAWS("s3", args...)
}

// RDS Operations

// ListRDSInstances lists all RDS instances
func (c *Client) ListRDSInstances() (string, error) {
	return c.execAWS("rds", "describe-db-instances", "--query", "DBInstances[*].[DBInstanceIdentifier,DBInstanceClass,Engine,DBInstanceStatus]", "--output", "table")
}

// DescribeRDSInstance gets details about an RDS instance
func (c *Client) DescribeRDSInstance(instanceID string) (string, error) {
	return c.execAWS("rds", "describe-db-instances", "--db-instance-identifier", instanceID, "--output", "table")
}

// StartRDSInstance starts an RDS instance
func (c *Client) StartRDSInstance(instanceID string) (string, error) {
	return c.execAWS("rds", "start-db-instance", "--db-instance-identifier", instanceID, "--output", "table")
}

// StopRDSInstance stops an RDS instance
func (c *Client) StopRDSInstance(instanceID string) (string, error) {
	return c.execAWS("rds", "stop-db-instance", "--db-instance-identifier", instanceID, "--output", "table")
}

// Cost Operations

// EstimateCosts provides a basic cost estimate (placeholder - requires Cost Explorer API)
func (c *Client) EstimateCosts(days int) (string, error) {
	return fmt.Sprintf("Cost estimation for last %d days (requires Cost Explorer setup)\n", days), nil
}

// FindUnusedResources identifies potentially unused resources
func (c *Client) FindUnusedResources() (map[string][]string, error) {
	resources := make(map[string][]string)

	// Check for unattached EBS volumes
	ebs, err := c.execAWS("ec2", "describe-volumes", "--filters", "Name=status,Values=available", "--query", "Volumes[*].VolumeId", "--output", "text")
	if err == nil && ebs != "" {
		resources["unused_ebs_volumes"] = strings.Fields(ebs)
	}

	// Check for unattached elastic IPs
	eip, err := c.execAWS("ec2", "describe-addresses", "--filters", "Name=instance-id,Values=''", "--query", "Addresses[*].PublicIp", "--output", "text")
	if err == nil && eip != "" {
		resources["unattached_elastic_ips"] = strings.Fields(eip)
	}

	// Check for stopped instances
	stopped, err := c.execAWS("ec2", "describe-instances", "--filters", "Name=instance-state-name,Values=stopped", "--query", "Reservations[*].Instances[*].InstanceId", "--output", "text")
	if err == nil && stopped != "" {
		resources["stopped_instances"] = strings.Fields(stopped)
	}

	return resources, nil
}

// SecurityOperations

// ListSecurityGroups lists all security groups
func (c *Client) ListSecurityGroups() (string, error) {
	return c.execAWS("ec2", "describe-security-groups", "--query", "SecurityGroups[*].[GroupId,GroupName,VpcId]", "--output", "table")
}

// ListVPCs lists all VPCs
func (c *Client) ListVPCs() (string, error) {
	return c.execAWS("ec2", "describe-vpcs", "--query", "Vpcs[*].[VpcId,CidrBlock,State]", "--output", "table")
}

// GetAWSAccountInfo gets general account information
func (c *Client) GetAWSAccountInfo() (string, error) {
	return c.execAWS("iam", "list-account-aliases", "--output", "table")
}
