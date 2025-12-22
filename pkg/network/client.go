package network

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps network diagnostic operations
type Client struct {
	Timeout string
}

// NewClient creates a new Network client
func NewClient() *Client {
	return &Client{
		Timeout: "5",
	}
}

// execCommand runs a command and returns output
func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Ping
func (c *Client) Ping(host string) (string, error) {
	// -c for Linux/Mac, -n for Windows
	return execCommand("ping", "-c", c.Timeout, host)
}

// Traceroute
func (c *Client) Traceroute(host string) (string, error) {
	// traceroute on Linux/Mac, tracert on Windows
	return execCommand("traceroute", "-m", "30", host)
}

// DNS lookup
func (c *Client) Nslookup(host string) (string, error) {
	return execCommand("nslookup", host)
}

// Dig (DNS information)
func (c *Client) Dig(host string) (string, error) {
	return execCommand("dig", host)
}

// Whois
func (c *Client) Whois(domain string) (string, error) {
	return execCommand("whois", domain)
}

// Check port connectivity
func (c *Client) CheckPort(host string, port string) (string, error) {
	// Using nc (netcat) for port checking
	return execCommand("nc", "-zv", "-w", c.Timeout, host, port)
}

// Get IP address
func (c *Client) GetIP(host string) (string, error) {
	return execCommand("getent", "hosts", host)
}

// Network interfaces
func (c *Client) GetNetworkInterfaces() (string, error) {
	return execCommand("ifconfig")
}

// IP route
func (c *Client) GetRouteTable() (string, error) {
	return execCommand("netstat", "-rn")
}

// DNS resolution with dig
func (c *Client) ReverseLookup(ip string) (string, error) {
	return execCommand("dig", "-x", ip)
}

// Check MTU
func (c *Client) CheckMTU(interface_name string) (string, error) {
	return execCommand("ifconfig", interface_name)
}

// TCP connections
func (c *Client) GetTCPConnections() (string, error) {
	return execCommand("netstat", "-an")
}

// Hostname and domain info
func (c *Client) GetHostname() (string, error) {
	return execCommand("hostname")
}

func (c *Client) GetFQDN() (string, error) {
	return execCommand("hostname", "-f")
}

// DNS servers
func (c *Client) GetDNSServers() (string, error) {
	return execCommand("cat", "/etc/resolv.conf")
}

// Packet loss check
func (c *Client) PacketLoss(host string, count string) (string, error) {
	return execCommand("ping", "-c", count, host)
}

// Bandwidth test (requires iperf3)
func (c *Client) SpeedTest(testServer string) (string, error) {
	// Requires speedtest-cli to be installed
	return execCommand("speedtest", "--simple")
}

// SSL/TLS certificate check
func (c *Client) CheckSSLCertificate(host string, port string) (string, error) {
	if port == "" {
		port = "443"
	}
	return execCommand("openssl", "s_client", "-connect", host+":"+port, "-showcerts")
}

// HTTP response check
func (c *Client) HTTPHead(url string) (string, error) {
	return execCommand("curl", "-I", url)
}

// DNS propagation check
func (c *Client) CheckDNSPropagation(domain string) (string, error) {
	return execCommand("dig", "+short", "NS", domain)
}
