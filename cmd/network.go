package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	networkpkg "github.com/yourusername/ops-tool/pkg/network"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network diagnostics and monitoring",
	Long:  "Perform network diagnostics (ping, traceroute, DNS, port checks, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var networkPingCmd = &cobra.Command{
	Use:   "ping <host>",
	Short: "Ping a host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.Ping(args[0])
		if err != nil {
			return fmt.Errorf("failed to ping: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkTracerouteCmd = &cobra.Command{
	Use:   "traceroute <host>",
	Short: "Trace route to host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.Traceroute(args[0])
		if err != nil {
			return fmt.Errorf("failed to traceroute: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkDNSCmd = &cobra.Command{
	Use:   "dns",
	Short: "DNS operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var networkNslookupCmd = &cobra.Command{
	Use:   "nslookup <hostname>",
	Short: "Perform nslookup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.Nslookup(args[0])
		if err != nil {
			return fmt.Errorf("failed to nslookup: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkDigCmd = &cobra.Command{
	Use:   "dig <hostname>",
	Short: "Perform DNS lookup with dig",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.Dig(args[0])
		if err != nil {
			return fmt.Errorf("failed to dig: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkReverseCmd = &cobra.Command{
	Use:   "reverse <ip>",
	Short: "Perform reverse DNS lookup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.ReverseLookup(args[0])
		if err != nil {
			return fmt.Errorf("failed to reverse lookup: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkPortCmd = &cobra.Command{
	Use:   "port <host> <port>",
	Short: "Check if port is open",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.CheckPort(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to check port: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkInterfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "Show network interfaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.GetNetworkInterfaces()
		if err != nil {
			return fmt.Errorf("failed to get interfaces: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkRouteCmd = &cobra.Command{
	Use:   "route",
	Short: "Show routing table",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.GetRouteTable()
		if err != nil {
			return fmt.Errorf("failed to get route table: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkConnectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Show active TCP connections",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.GetTCPConnections()
		if err != nil {
			return fmt.Errorf("failed to get connections: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkHostnameCmd = &cobra.Command{
	Use:   "hostname",
	Short: "Get system hostname",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.GetHostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkFQDNCmd = &cobra.Command{
	Use:   "fqdn",
	Short: "Get fully qualified domain name",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.GetFQDN()
		if err != nil {
			return fmt.Errorf("failed to get FQDN: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkSSLCmd = &cobra.Command{
	Use:   "ssl <host> [port]",
	Short: "Check SSL certificate",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		port := ""
		if len(args) > 1 {
			port = args[1]
		}
		output, err := client.CheckSSLCertificate(args[0], port)
		if err != nil {
			return fmt.Errorf("failed to check SSL: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkMTUCmd = &cobra.Command{
	Use:   "mtu <interface>",
	Short: "Check network MTU",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.CheckMTU(args[0])
		if err != nil {
			return fmt.Errorf("failed to check MTU: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkDNSServersCmd = &cobra.Command{
	Use:   "nameservers",
	Short: "Show DNS servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.GetDNSServers()
		if err != nil {
			return fmt.Errorf("failed to get DNS servers: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var networkWhoisCmd = &cobra.Command{
	Use:   "whois <domain>",
	Short: "Get WHOIS information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := networkpkg.NewClient()
		output, err := client.Whois(args[0])
		if err != nil {
			return fmt.Errorf("failed to whois: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)

	// DNS commands as subgroup
	networkCmd.AddCommand(networkDNSCmd)
	networkDNSCmd.AddCommand(networkNslookupCmd, networkDigCmd, networkReverseCmd)

	// All other commands
	networkCmd.AddCommand(
		networkPingCmd,
		networkTracerouteCmd,
		networkPortCmd,
		networkInterfacesCmd,
		networkRouteCmd,
		networkConnectionsCmd,
		networkHostnameCmd,
		networkFQDNCmd,
		networkSSLCmd,
		networkMTUCmd,
		networkDNSServersCmd,
		networkWhoisCmd,
	)
}
