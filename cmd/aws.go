package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	awspkg "github.com/yourusername/ops-tool/pkg/aws"
)

var awsRegion string
var awsProfile string

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "AWS cloud operations",
	Long:  "Manage AWS resources (EC2, S3, RDS, IAM, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var awsWhoCmd = &cobra.Command{
	Use:   "who",
	Short: "Show current AWS identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.GetCallerIdentity()
		if err != nil {
			return fmt.Errorf("failed to get AWS identity: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

// EC2 Commands

var awsEC2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "EC2 instance management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var awsEC2ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List EC2 instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		running, _ := cmd.Flags().GetBool("running")
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.ListEC2Instances(running)
		if err != nil {
			return fmt.Errorf("failed to list EC2 instances: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsEC2StartCmd = &cobra.Command{
	Use:   "start <instance-id>",
	Short: "Start an EC2 instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.StartEC2Instance(args[0])
		if err != nil {
			return fmt.Errorf("failed to start instance: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsEC2StopCmd = &cobra.Command{
	Use:   "stop <instance-id>",
	Short: "Stop an EC2 instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.StopEC2Instance(args[0])
		if err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsEC2TerminateCmd = &cobra.Command{
	Use:   "terminate <instance-id>",
	Short: "Terminate an EC2 instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.TerminateEC2Instance(args[0])
		if err != nil {
			return fmt.Errorf("failed to terminate instance: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsEC2DescribeCmd = &cobra.Command{
	Use:   "describe <instance-id>",
	Short: "Describe an EC2 instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.DescribeEC2Instance(args[0])
		if err != nil {
			return fmt.Errorf("failed to describe instance: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

// S3 Commands

var awsS3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "S3 bucket management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var awsS3ListCmd = &cobra.Command{
	Use:   "list [bucket-name]",
	Short: "List S3 buckets or bucket contents",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)

		if len(args) == 0 {
			output, err := client.ListS3Buckets()
			if err != nil {
				return fmt.Errorf("failed to list buckets: %w", err)
			}
			fmt.Println(output)
			return nil
		}

		recursive, _ := cmd.Flags().GetBool("recursive")
		output, err := client.ListS3BucketContents(args[0], recursive)
		if err != nil {
			return fmt.Errorf("failed to list bucket contents: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsS3SizeCmd = &cobra.Command{
	Use:   "size <bucket-name>",
	Short: "Get S3 bucket size",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.GetS3BucketSize(args[0])
		if err != nil {
			return fmt.Errorf("failed to get bucket size: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

// RDS Commands

var awsRDSCmd = &cobra.Command{
	Use:   "rds",
	Short: "RDS database management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var awsRDSListCmd = &cobra.Command{
	Use:   "list",
	Short: "List RDS instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.ListRDSInstances()
		if err != nil {
			return fmt.Errorf("failed to list RDS instances: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsRDSDescribeCmd = &cobra.Command{
	Use:   "describe <instance-id>",
	Short: "Describe an RDS instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.DescribeRDSInstance(args[0])
		if err != nil {
			return fmt.Errorf("failed to describe RDS instance: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsRDSStartCmd = &cobra.Command{
	Use:   "start <instance-id>",
	Short: "Start an RDS instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.StartRDSInstance(args[0])
		if err != nil {
			return fmt.Errorf("failed to start RDS instance: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsRDSStopCmd = &cobra.Command{
	Use:   "stop <instance-id>",
	Short: "Stop an RDS instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.StopRDSInstance(args[0])
		if err != nil {
			return fmt.Errorf("failed to stop RDS instance: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

// Cost & Optimization Commands

var awsCostCmd = &cobra.Command{
	Use:   "cost",
	Short: "Cost optimization and monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var awsCostOptimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Find unused resources and cost optimization opportunities",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		resources, err := client.FindUnusedResources()
		if err != nil {
			return fmt.Errorf("failed to analyze resources: %w", err)
		}

		fmt.Println("🔍 Cost Optimization Opportunities:")
		fmt.Println("===================================")
		if len(resources) == 0 {
			fmt.Println("✅ No unused resources found!")
		} else {
			for resource, items := range resources {
				fmt.Printf("\n%s (%d):\n", resource, len(items))
				for _, item := range items {
					fmt.Printf("  - %s\n", item)
				}
			}
		}
		return nil
	},
}

// Account & IAM Commands

var awsAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "AWS account operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var awsAccountInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show AWS account information",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.GetAWSAccountInfo()
		if err != nil {
			return fmt.Errorf("failed to get account info: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

// Security Commands

var awsSecurityCmd = &cobra.Command{
	Use:   "security",
	Short: "AWS security operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var awsSecurityGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List security groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.ListSecurityGroups()
		if err != nil {
			return fmt.Errorf("failed to list security groups: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var awsSecurityVpcsCmd = &cobra.Command{
	Use:   "vpcs",
	Short: "List VPCs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := awspkg.NewClient(awsRegion, awsProfile)
		output, err := client.ListVPCs()
		if err != nil {
			return fmt.Errorf("failed to list VPCs: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

func init() {
	// Global AWS flags
	awsCmd.PersistentFlags().StringVarP(&awsRegion, "region", "r", "", "AWS region")
	awsCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "", "AWS profile")

	// Who command
	awsCmd.AddCommand(awsWhoCmd)

	// EC2 subcommands
	awsEC2ListCmd.Flags().BoolP("running", "r", false, "Show only running instances")
	awsEC2Cmd.AddCommand(awsEC2ListCmd)
	awsEC2Cmd.AddCommand(awsEC2StartCmd)
	awsEC2Cmd.AddCommand(awsEC2StopCmd)
	awsEC2Cmd.AddCommand(awsEC2TerminateCmd)
	awsEC2Cmd.AddCommand(awsEC2DescribeCmd)

	// S3 subcommands
	awsS3ListCmd.Flags().BoolP("recursive", "r", false, "Recursive listing")
	awsS3Cmd.AddCommand(awsS3ListCmd)
	awsS3Cmd.AddCommand(awsS3SizeCmd)

	// RDS subcommands
	awsRDSCmd.AddCommand(awsRDSListCmd)
	awsRDSCmd.AddCommand(awsRDSDescribeCmd)
	awsRDSCmd.AddCommand(awsRDSStartCmd)
	awsRDSCmd.AddCommand(awsRDSStopCmd)

	// Cost subcommands
	awsCostCmd.AddCommand(awsCostOptimizeCmd)

	// Account subcommands
	awsAccountCmd.AddCommand(awsAccountInfoCmd)

	// Security subcommands
	awsSecurityCmd.AddCommand(awsSecurityGroupsCmd)
	awsSecurityCmd.AddCommand(awsSecurityVpcsCmd)

	// AWS main subcommands
	awsCmd.AddCommand(awsEC2Cmd)
	awsCmd.AddCommand(awsS3Cmd)
	awsCmd.AddCommand(awsRDSCmd)
	awsCmd.AddCommand(awsCostCmd)
	awsCmd.AddCommand(awsAccountCmd)
	awsCmd.AddCommand(awsSecurityCmd)
}
