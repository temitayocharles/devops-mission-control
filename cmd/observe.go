package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	dashboardpkg "github.com/yourusername/ops-tool/pkg/dashboard"
	metricspkg "github.com/yourusername/ops-tool/pkg/metrics"
	slackpkg "github.com/yourusername/ops-tool/pkg/slack"
)

var (
	dashboardAddr  string
	metricsStore   *metricspkg.MetricsStore
	dashboardInst  *dashboardpkg.Dashboard
)

var observabilityCmd = &cobra.Command{
	Use:   "observe",
	Short: "Observability and monitoring",
	Long:  "Manage observability features (dashboard, metrics, alerts, Slack)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var dashboardStartCmd = &cobra.Command{
	Use:   "dashboard start",
	Short: "Start the observability dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		dashboardInst = dashboardpkg.NewDashboard(dashboardAddr, metricsStore)
		go dashboardInst.Start()

		fmt.Printf("✅ Dashboard started at http://%s\n", dashboardAddr)
		fmt.Println("Press Ctrl+C to stop")

		// Keep running
		select {}
	},
}

var dashboardStopCmd = &cobra.Command{
	Use:   "dashboard stop",
	Short: "Stop the observability dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dashboardInst == nil {
			return fmt.Errorf("dashboard not running")
		}

		dashboardInst.Stop()
		fmt.Println("✅ Dashboard stopped")
		return nil
	},
}

var metricsListCmd = &cobra.Command{
	Use:   "metrics list",
	Short: "List all metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		metrics := metricsStore.GetMetrics()
		if len(metrics) == 0 {
			fmt.Println("No metrics recorded")
			return nil
		}

		fmt.Println("Metrics:")
		for _, m := range metrics {
			fmt.Printf("  %s: %.2f %s (tags: %v)\n", m.Name, m.Value, m.Unit, m.Tags)
		}
		return nil
	},
}

var metricsRecordCmd = &cobra.Command{
	Use:   "metrics record <name> <value> <unit>",
	Short: "Record a metric",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		var value float64
		_, err := fmt.Sscanf(args[1], "%f", &value)
		if err != nil {
			return fmt.Errorf("invalid metric value: %w", err)
		}

		metricsStore.RecordMetric(args[0], value, args[2], nil)
		fmt.Printf("✅ Recorded metric: %s = %.2f %s\n", args[0], value, args[2])
		return nil
	},
}

var metricsStatsCmd = &cobra.Command{
	Use:   "metrics stats",
	Short: "Show metrics statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		stats := metricsStore.GetStats()
		fmt.Println("Metrics Statistics:")
		for key, value := range stats {
			fmt.Printf("  %s: %v\n", key, value)
		}
		return nil
	},
}

var alertsListCmd = &cobra.Command{
	Use:   "alerts list",
	Short: "List all alerts",
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		alerts := metricsStore.GetAlerts()
		if len(alerts) == 0 {
			fmt.Println("No alerts")
			return nil
		}

		fmt.Println("Alerts:")
		for _, a := range alerts {
			status := "Active"
			if a.Resolved {
				status = "Resolved"
			}
			fmt.Printf("  [%s] %s: %s (%s)\n", a.Severity, a.Name, a.Message, status)
		}
		return nil
	},
}

var alertsActiveCmd = &cobra.Command{
	Use:   "alerts active",
	Short: "Show active alerts",
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		alerts := metricsStore.GetActiveAlerts()
		if len(alerts) == 0 {
			fmt.Println("✅ No active alerts")
			return nil
		}

		fmt.Println("Active Alerts:")
		for _, a := range alerts {
			fmt.Printf("  [%s] %s: %s\n", a.Severity, a.Name, a.Message)
		}
		return nil
	},
}

var alertsCreateCmd = &cobra.Command{
	Use:   "alerts create <name> <severity> <message>",
	Short: "Create an alert",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		metricsStore.CreateAlert(args[0], args[1], args[2], nil)
		fmt.Printf("✅ Alert created: %s (%s)\n", args[0], args[1])
		return nil
	},
}

var eventsListCmd = &cobra.Command{
	Use:   "events list",
	Short: "List all events",
	RunE: func(cmd *cobra.Command, args []string) error {
		if metricsStore == nil {
			metricsStore = metricspkg.NewMetricsStore(10000)
		}

		events := metricsStore.GetEvents()
		if len(events) == 0 {
			fmt.Println("No events")
			return nil
		}

		fmt.Println("Events:")
		for _, e := range events {
			fmt.Printf("  [%s] %s: %s (%s)\n", e.Type, e.Message, e.Status, e.Duration)
		}
		return nil
	},
}

var slackSendCmd = &cobra.Command{
	Use:   "slack send <webhook-url> <channel> <message>",
	Short: "Send message to Slack",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := slackpkg.NewClient(args[0], args[1])
		err := client.SendMessage(args[2])
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		fmt.Println("✅ Message sent to Slack")
		return nil
	},
}

var slackAlertCmd = &cobra.Command{
	Use:   "slack alert <webhook-url> <channel> <alert-name> <severity> <message>",
	Short: "Send alert to Slack",
	Args:  cobra.ExactArgs(5),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := slackpkg.NewClient(args[0], args[1])
		err := client.SendAlert(args[2], args[3], args[4], nil)
		if err != nil {
			return fmt.Errorf("failed to send alert: %w", err)
		}
		fmt.Println("✅ Alert sent to Slack")
		return nil
	},
}

var slackDeployCmd = &cobra.Command{
	Use:   "slack deploy <webhook-url> <channel> <app> <status> <version>",
	Short: "Send deployment notification to Slack",
	Args:  cobra.ExactArgs(5),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := slackpkg.NewClient(args[0], args[1])
		err := client.SendDeployment(args[2], args[3], args[4], nil)
		if err != nil {
			return fmt.Errorf("failed to send deployment: %w", err)
		}
		fmt.Println("✅ Deployment notification sent to Slack")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(observabilityCmd)

	// Dashboard commands
	dashboardCmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Dashboard operations",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
			}
		},
	}
	observabilityCmd.AddCommand(dashboardCmd)
	dashboardCmd.AddCommand(dashboardStartCmd, dashboardStopCmd)
	dashboardCmd.PersistentFlags().StringVarP(&dashboardAddr, "addr", "a", "localhost:8080", "Dashboard address")

	// Metrics commands
	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Metrics operations",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
			}
		},
	}
	observabilityCmd.AddCommand(metricsCmd)
	metricsCmd.AddCommand(metricsListCmd, metricsRecordCmd, metricsStatsCmd)

	// Alerts commands
	alertsCmd := &cobra.Command{
		Use:   "alerts",
		Short: "Alert operations",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
			}
		},
	}
	observabilityCmd.AddCommand(alertsCmd)
	alertsCmd.AddCommand(alertsListCmd, alertsActiveCmd, alertsCreateCmd)

	// Events commands
	eventsCmd := &cobra.Command{
		Use:   "events",
		Short: "Event operations",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
			}
		},
	}
	observabilityCmd.AddCommand(eventsCmd)
	eventsCmd.AddCommand(eventsListCmd)

	// Slack commands
	slackCmd := &cobra.Command{
		Use:   "slack",
		Short: "Slack integration",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
			}
		},
	}
	observabilityCmd.AddCommand(slackCmd)
	slackCmd.AddCommand(slackSendCmd, slackAlertCmd, slackDeployCmd)
}
