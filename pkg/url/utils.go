package url

import (
	"fmt"
	"net"
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
	// 如果是本地地址，例如 "localhost" 或 "127.0.0.1"
	if strings.HasPrefix(addr, "localhost") || strings.HasPrefix(addr, "127.0.0.1") {
		return nil
	}

	// 使用 net.ParseIP 检查是否是有效的 IP 地址
	ip := net.ParseIP(addr)
	if ip != nil {
		return nil
	}

	// 如果是端口检查 (例如 :8080)
	if strings.HasPrefix(addr, ":") {
		return nil
	}

	// 如果是主机名和端口号的组合，例如 localhost:8080
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %v", err)
	}

	return nil
}
