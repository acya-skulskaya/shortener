package helpers

import (
	"fmt"
	"net"
)

// CheckIPSubnet parses given IP and trusted subnet and checks whether given subnet contains given IP
func CheckIPSubnet(ip string, trustedSubnet string) (bool, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, fmt.Errorf("invalid IP: %s", ip)
	}

	_, ipNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false, fmt.Errorf("invalid trusted subnet %s: %w", ip, err)
	}

	return ipNet.Contains(parsedIP), nil
}
