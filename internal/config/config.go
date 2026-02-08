package config

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration from environment.
type Config struct {
	DMARCDir string

	// knownSenderIPs: single IPs and CIDR ranges (e.g. 143.55.232.0/24) for "known sender"
	singleIPs map[string]bool
	networks  []*net.IPNet
}

// IsKnownSenderIP returns true if ip is in KNOWN_SERVER_IPS (exact match or within a CIDR range).
func (c *Config) IsKnownSenderIP(ip string) bool {
	if c == nil || ip == "" {
		return false
	}
	if c.singleIPs != nil && c.singleIPs[ip] {
		return true
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, n := range c.networks {
		if n != nil && n.Contains(parsed) {
			return true
		}
	}
	return false
}

// Load reads .env and builds Config. DMARC_DIR defaults to ./DMARC; KNOWN_SERVER_IPS is optional.
// KNOWN_SERVER_IPS can contain single IPs and/or CIDR ranges (e.g. 143.55.232.0/24), comma-separated.
func Load() (*Config, error) {
	_ = godotenv.Load()

	dmarcDir := os.Getenv("DMARC_DIR")
	if dmarcDir == "" {
		dmarcDir = "DMARC"
	}
	dmarcDir, err := filepath.Abs(dmarcDir)
	if err != nil {
		dmarcDir = os.Getenv("DMARC_DIR")
	}

	singleIPs := make(map[string]bool)
	var networks []*net.IPNet
	if s := os.Getenv("KNOWN_SERVER_IPS"); s != "" {
		for _, entry := range strings.Split(s, ",") {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				continue
			}
			if strings.Contains(entry, "/") {
				_, n, err := net.ParseCIDR(entry)
				if err != nil {
					continue
				}
				networks = append(networks, n)
			} else {
				singleIPs[entry] = true
			}
		}
	}

	return &Config{
		DMARCDir:  dmarcDir,
		singleIPs: singleIPs,
		networks:  networks,
	}, nil
}
