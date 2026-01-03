package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/devops-mission-control/pkg/audit"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Query audit log",
}

var auditListCmd = &cobra.Command{
	Use:   "list",
	Short: "List audit entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		actor, _ := cmd.Flags().GetString("actor")
		action, _ := cmd.Flags().GetString("action")
		sinceStr, _ := cmd.Flags().GetString("since")
		var since time.Time
		if sinceStr != "" {
			t, err := time.Parse(time.RFC3339, sinceStr)
			if err != nil {
				return err
			}
			since = t
		}
		entries, err := audit.ReadEntries("")
		if err != nil {
			return err
		}
		for _, e := range entries {
			if actor != "" && e.Actor != actor {
				continue
			}
			if action != "" && e.Action != action {
				continue
			}
			if !since.IsZero() && e.Timestamp.Before(since) {
				continue
			}
			fmt.Printf("%s %s actor=%s target=%s details=%v\n", e.Timestamp.Format(time.RFC3339), e.Action, e.Actor, e.Target, e.Details)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.AddCommand(auditListCmd)
	auditListCmd.Flags().String("actor", "", "filter by actor")
	auditListCmd.Flags().String("action", "", "filter by action")
	auditListCmd.Flags().String("since", "", "filter entries since RFC3339 timestamp")
}
