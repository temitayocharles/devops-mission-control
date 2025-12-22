package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	gitpkg "github.com/yourusername/ops-tool/pkg/git"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git version control operations",
	Long:  "Manage Git repositories, branches, commits, and workflows",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show repository status",
	RunE: func(cmd *cobra.Command, args []string) error {
		short, _ := cmd.Flags().GetBool("short")
		client := gitpkg.NewClient()

		var output string
		var err error
		if short {
			output, err = client.GetStatusShort()
		} else {
			output, err = client.GetStatus()
		}

		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage branches",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gitBranchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		client := gitpkg.NewClient()
		output, err := client.ListBranches(all)
		if err != nil {
			return fmt.Errorf("failed to list branches: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitBranchCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gitpkg.NewClient()
		branch, err := client.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		fmt.Printf("Current branch: %s\n", branch)
		return nil
	},
}

var gitBranchCreateCmd = &cobra.Command{
	Use:   "create <branch-name>",
	Short: "Create a new branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gitpkg.NewClient()
		output, err := client.CreateBranch(args[0])
		if err != nil {
			return fmt.Errorf("failed to create branch: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitBranchSwitchCmd = &cobra.Command{
	Use:   "switch <branch-name>",
	Short: "Switch to a branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gitpkg.NewClient()
		_, err := client.SwitchBranch(args[0])
		if err != nil {
			return fmt.Errorf("failed to switch branch: %w", err)
		}

		fmt.Printf("✅ Switched to branch: %s\n", args[0])
		return nil
	},
}

var gitBranchDeleteCmd = &cobra.Command{
	Use:   "delete <branch-name>",
	Short: "Delete a branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		client := gitpkg.NewClient()
		output, err := client.DeleteBranch(args[0], force)
		if err != nil {
			return fmt.Errorf("failed to delete branch: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitAddCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Stage changes",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gitpkg.NewClient()

		var output string
		var err error
		if len(args) == 0 {
			output, err = client.AddAll()
		} else {
			output, err = client.Add(args[0])
		}

		if err != nil {
			return fmt.Errorf("failed to add changes: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitCommitCmd = &cobra.Command{
	Use:   "commit <message>",
	Short: "Commit changes",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		client := gitpkg.NewClient()

		message := strings.Join(args, " ")
		var output string
		var err error

		if all {
			output, err = client.CommitAll(message)
		} else {
			output, err = client.Commit(message)
		}

		if err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitPushCmd = &cobra.Command{
	Use:   "push [remote] [branch]",
	Short: "Push changes to remote",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		client := gitpkg.NewClient()

		remote := "origin"
		branch := ""
		if len(args) > 0 {
			remote = args[0]
		}
		if len(args) > 1 {
			branch = args[1]
		}

		var output string
		var err error
		if force {
			output, err = client.PushForce(remote, branch)
		} else {
			output, err = client.Push(remote, branch)
		}

		if err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitPullCmd = &cobra.Command{
	Use:   "pull [remote] [branch]",
	Short: "Pull changes from remote",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		rebase, _ := cmd.Flags().GetBool("rebase")
		client := gitpkg.NewClient()

		remote := "origin"
		branch := ""
		if len(args) > 0 {
			remote = args[0]
		}
		if len(args) > 1 {
			branch = args[1]
		}

		var output string
		var err error
		if rebase {
			output, err = client.PullRebase(remote, branch)
		} else {
			output, err = client.Pull(remote, branch)
		}

		if err != nil {
			return fmt.Errorf("failed to pull: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit log",
	RunE: func(cmd *cobra.Command, args []string) error {
		maxCount, _ := cmd.Flags().GetInt("count")
		client := gitpkg.NewClient()
		output, err := client.GetLog(maxCount)
		if err != nil {
			return fmt.Errorf("failed to get log: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		staged, _ := cmd.Flags().GetBool("staged")
		client := gitpkg.NewClient()
		output, err := client.GetDiff(staged)
		if err != nil {
			return fmt.Errorf("failed to get diff: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitRemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remotes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gitpkg.NewClient()
		output, err := client.GetRemoteList()
		if err != nil {
			return fmt.Errorf("failed to list remotes: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitContributorsCmd = &cobra.Command{
	Use:   "contributors",
	Short: "Show contributors",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gitpkg.NewClient()
		output, err := client.GetContributors()
		if err != nil {
			return fmt.Errorf("failed to get contributors: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gitTagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gitpkg.NewClient()
		output, err := client.GetTags()
		if err != nil {
			return fmt.Errorf("failed to list tags: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var gitTagCreateCmd = &cobra.Command{
	Use:   "create <tag-name>",
	Short: "Create a tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		client := gitpkg.NewClient()
		output, err := client.CreateTag(args[0], message)
		if err != nil {
			return fmt.Errorf("failed to create tag: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

func init() {
	// Status flags
	gitStatusCmd.Flags().BoolP("short", "s", false, "Short format")

	// Branch subcommands
	gitBranchListCmd.Flags().BoolP("all", "a", false, "Show all branches")
	gitBranchCmd.AddCommand(gitBranchListCmd)
	gitBranchCmd.AddCommand(gitBranchCurrentCmd)
	gitBranchCmd.AddCommand(gitBranchCreateCmd)
	gitBranchCmd.AddCommand(gitBranchSwitchCmd)

	gitBranchDeleteCmd.Flags().BoolP("force", "f", false, "Force delete")
	gitBranchCmd.AddCommand(gitBranchDeleteCmd)

	// Add command
	gitCmd.AddCommand(gitAddCmd)

	// Commit command
	gitCommitCmd.Flags().BoolP("all", "a", false, "Add all changes before commit")
	gitCmd.AddCommand(gitCommitCmd)

	// Push command
	gitPushCmd.Flags().BoolP("force", "f", false, "Force push (--force-with-lease)")
	gitCmd.AddCommand(gitPushCmd)

	// Pull command
	gitPullCmd.Flags().BoolP("rebase", "r", false, "Rebase instead of merge")
	gitCmd.AddCommand(gitPullCmd)

	// Log command
	gitLogCmd.Flags().IntP("count", "n", 10, "Number of commits to show")
	gitCmd.AddCommand(gitLogCmd)

	// Diff command
	gitDiffCmd.Flags().BoolP("staged", "s", false, "Show staged changes")
	gitCmd.AddCommand(gitDiffCmd)

	// Tag subcommands
	gitTagCmd.AddCommand(gitTagListCmd)
	gitTagCreateCmd.Flags().StringP("message", "m", "", "Tag message (for annotated tag)")
	gitTagCmd.AddCommand(gitTagCreateCmd)

	// Git main subcommands
	gitCmd.AddCommand(gitStatusCmd)
	gitCmd.AddCommand(gitBranchCmd)
	gitCmd.AddCommand(gitRemoteCmd)
	gitCmd.AddCommand(gitContributorsCmd)
	gitCmd.AddCommand(gitTagCmd)
}
