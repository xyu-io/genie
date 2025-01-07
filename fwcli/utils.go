package fwcli

import (
	"errors"
	"net"
)

// Generic function to check if an item exists in a slice
func itemExists(slice []string, item string) bool {
	for _, existingItem := range slice {
		if existingItem == item {
			return true
		}
	}
	return false
}

// buildArgs rich-rule 的命令参数构造
func buildArgs(zone string, rule string, permanent bool, isAdd bool) []string {
	var args []string
	if zone == "" {
		if isAdd {
			args = []string{"--add-rich-rule", rule}
		} else {
			args = []string{"--remove-rich-rule", rule}
		}
	} else {
		if isAdd {
			args = []string{"--zone", zone, "--add-rich-rule", rule}
		} else {
			args = []string{"--zone", zone, "--remove-rich-rule", rule}
		}
	}
	if permanent {
		args = append(args, "--permanent")
	}
	return args
}

// checkIPVersion 判断IP的类型
func checkIPVersion(ip string) (string, error) {
	if ip == "" {
		return IPv4, nil
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "", errors.New("invalid IP address")
	}
	if parsedIP.To4() != nil {
		return IPv4, nil
	}
	if parsedIP.To16() != nil {
		return IPv6, nil
	}
	return "", errors.New("unknown IP version")
}
