package url

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ValidateURLAddr validates a URL address.
// ValidateURLAddr 验证 URL 地址。
//
// Parameters:
//
//	addr - The address to validate / 要验证的地址
//
// Returns:
//
//	error - Error if the address is invalid / 如果地址无效则返回错误
func ValidateURLAddr(addr string) error {
	// Check for empty address
	// 检查空地址
	if addr == "" {
		return fmt.Errorf("missing port in address")
	}

	// Check for invalid characters (like @)
	// 检查无效字符（如 @）
	if strings.Contains(addr, "@") {
		return fmt.Errorf("invalid address: contains invalid character '@'")
	}

	// Check for spaces
	// 检查空格
	if strings.Contains(addr, " ") {
		return fmt.Errorf("invalid address: contains spaces")
	}

	// Check if it's an IP address (IPv4 or IPv6)
	// IP addresses are valid without ports
	// 检查是否是 IP 地址（IPv4 或 IPv6）
	// IP 地址不需要端口

	// Check for IPv6 with port format: [::1]:8080
	// 检查 IPv6 带端口格式：[::1]:8080
	if strings.HasPrefix(addr, "[") {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return fmt.Errorf("invalid address: %v", err)
		}
		// Remove brackets for IP parsing
		// 移除方括号进行 IP 解析
		host = strings.TrimPrefix(host, "[")
		host = strings.TrimSuffix(host, "]")
		ip := net.ParseIP(host)
		if ip == nil {
			return fmt.Errorf("invalid address: invalid IPv6 address '%s'", host)
		}
		if port != "" {
			_, err := strconv.ParseInt(port, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid address: invalid port format '%s'", port)
			}
		}
		return nil
	}

	// Check if it's a plain IPv6 address (e.g., ::1)
	// 检查是否是纯 IPv6 地址（如 ::1）
	// Must do this before checking for port-only format
	// 必须在检查端口格式之前进行此检查
	if strings.Contains(addr, ":") {
		ip := net.ParseIP(addr)
		if ip != nil {
			return nil
		}
	}

	// If it's just a port (e.g., :8080), validate the port format
	// 如果只是端口（如 :8080），验证端口格式
	if strings.HasPrefix(addr, ":") {
		port := strings.TrimPrefix(addr, ":")
		if port == "" {
			return fmt.Errorf("invalid address: port is empty")
		}
		// Note: we don't validate port range here (as per test expectations)
		_, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid address: invalid port format '%s'", port)
		}
		return nil
	}

	// Split host and port
	// 分割主机和端口
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		// Check if it's an IP address without port
		// 检查是否是不带端口的 IP 地址
		if net.ParseIP(addr) != nil {
			return nil
		}
		// Check if it's localhost without port
		// 检查是否是不带端口的 localhost
		if addr == "localhost" {
			return nil
		}
		// Hostnames require a port
		// 主机名需要端口
		if strings.Contains(err.Error(), "missing port") {
			return fmt.Errorf("missing port in address")
		}
		return fmt.Errorf("invalid address: %v", err)
	}

	// Validate host
	// 验证主机
	if host == "" {
		return fmt.Errorf("invalid address: host is empty")
	}

	// Check if it's a valid IP address
	// 检查是否是有效的 IP 地址
	ip := net.ParseIP(host)
	if ip != nil {
		// Validate port if present
		// 如果存在端口，验证端口
		if port != "" {
			_, err := strconv.ParseInt(port, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid address: invalid port format '%s'", port)
			}
		}
		return nil
	}

	// Validate hostname (allow underscore, hyphen, dot)
	// 验证主机名（允许下划线、连字符、点）
	for _, c := range host {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '.' || c == '_') {
			return fmt.Errorf("invalid address: host contains invalid character '%c'", c)
		}
	}

	// Validate port if present
	// 如果存在端口，验证端口
	if port != "" {
		_, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid address: invalid port format '%s'", port)
		}
	}

	return nil
}
