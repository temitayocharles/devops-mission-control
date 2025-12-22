package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client wraps Git operations
type Client struct{}

// NewClient creates a new Git client
func NewClient() *Client {
	return &Client{}
}

// execGit runs a git command and returns output
func (c *Client) execGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// GetStatus returns repository status
func (c *Client) GetStatus() (string, error) {
	return c.execGit("status")
}

// GetStatusShort returns short repository status
func (c *Client) GetStatusShort() (string, error) {
	return c.execGit("status", "-s")
}

// ListBranches lists all branches
func (c *Client) ListBranches(all bool) (string, error) {
	args := []string{"branch"}
	if all {
		args = append(args, "-a")
	}
	return c.execGit(args...)
}

// GetCurrentBranch returns current branch name
func (c *Client) GetCurrentBranch() (string, error) {
	return c.execGit("rev-parse", "--abbrev-ref", "HEAD")
}

// CreateBranch creates a new branch
func (c *Client) CreateBranch(branchName string) (string, error) {
	return c.execGit("checkout", "-b", branchName)
}

// SwitchBranch switches to a different branch
func (c *Client) SwitchBranch(branchName string) (string, error) {
	return c.execGit("checkout", branchName)
}

// DeleteBranch deletes a branch
func (c *Client) DeleteBranch(branchName string, force bool) (string, error) {
	args := []string{"branch", "-d", branchName}
	if force {
		args = []string{"branch", "-D", branchName}
	}
	return c.execGit(args...)
}

// Add stages changes
func (c *Client) Add(path string) (string, error) {
	return c.execGit("add", path)
}

// AddAll stages all changes
func (c *Client) AddAll() (string, error) {
	return c.execGit("add", "-A")
}

// Commit commits changes
func (c *Client) Commit(message string) (string, error) {
	return c.execGit("commit", "-m", message)
}

// CommitAll commits all changes
func (c *Client) CommitAll(message string) (string, error) {
	return c.execGit("commit", "-am", message)
}

// Push pushes changes to remote
func (c *Client) Push(remote, branch string) (string, error) {
	if remote == "" {
		remote = "origin"
	}
	if branch == "" {
		branch = "HEAD"
	}
	return c.execGit("push", remote, branch)
}

// PushForce pushes with force (--force-with-lease)
func (c *Client) PushForce(remote, branch string) (string, error) {
	if remote == "" {
		remote = "origin"
	}
	if branch == "" {
		branch = "HEAD"
	}
	return c.execGit("push", "--force-with-lease", remote, branch)
}

// Pull pulls changes from remote
func (c *Client) Pull(remote, branch string) (string, error) {
	if remote == "" {
		remote = "origin"
	}
	if branch == "" {
		branch = "HEAD"
	}
	return c.execGit("pull", remote, branch)
}

// PullRebase pulls with rebase
func (c *Client) PullRebase(remote, branch string) (string, error) {
	if remote == "" {
		remote = "origin"
	}
	if branch == "" {
		branch = "HEAD"
	}
	return c.execGit("pull", "--rebase", remote, branch)
}

// GetLog returns commit log
func (c *Client) GetLog(maxCount int) (string, error) {
	args := []string{"log", "--oneline", "--graph"}
	if maxCount > 0 {
		args = append(args, fmt.Sprintf("-n %d", maxCount))
	}
	return c.execGit(args...)
}

// GetLastCommit returns last commit info
func (c *Client) GetLastCommit() (string, error) {
	return c.execGit("log", "-1", "--stat")
}

// GetDiff returns diff of changes
func (c *Client) GetDiff(staged bool) (string, error) {
	args := []string{"diff"}
	if staged {
		args = append(args, "--staged")
	}
	return c.execGit(args...)
}

// GetRemoteList returns list of remotes
func (c *Client) GetRemoteList() (string, error) {
	return c.execGit("remote", "-v")
}

// GetContributors returns list of contributors
func (c *Client) GetContributors() (string, error) {
	return c.execGit("shortlog", "-sn")
}

// ShowFile shows file content from a commit
func (c *Client) ShowFile(ref, file string) (string, error) {
	return c.execGit("show", ref+":"+file)
}

// UndoLastCommit undoes last commit (soft reset)
func (c *Client) UndoLastCommit() (string, error) {
	return c.execGit("reset", "--soft", "HEAD~1")
}

// RevertCommit reverts a commit
func (c *Client) RevertCommit(commitSHA string) error {
	cmd := exec.Command("git", "revert", commitSHA, "--no-edit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// CleanWorktree cleans untracked files
func (c *Client) CleanWorktree(dry bool, directories bool) (string, error) {
	args := []string{"clean"}
	if dry {
		args = append(args, "-n")
	}
	if directories {
		args = append(args, "-d")
	}
	args = append(args, "-f")
	return c.execGit(args...)
}

// GetTags returns list of tags
func (c *Client) GetTags() (string, error) {
	return c.execGit("tag", "-l")
}

// CreateTag creates a new tag
func (c *Client) CreateTag(tagName, message string) (string, error) {
	args := []string{"tag"}
	if message != "" {
		args = append(args, "-a", tagName, "-m", message)
	} else {
		args = append(args, tagName)
	}
	return c.execGit(args...)
}
