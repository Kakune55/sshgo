package ssh

import (
	"strings"
)

// ParseHostArgument 解析命令行参数中的主机信息
func ParseHostArgument(hostArg string) SSHHost {
	host := SSHHost{
		Port: "22", // 默认端口
	}

	// 检查是否包含用户信息 (user@host)
	if strings.Contains(hostArg, "@") {
		parts := strings.Split(hostArg, "@")
		host.User = parts[0]
		hostArg = parts[1]
	}

	// 检查是否包含端口信息 (host:port)
	if strings.Contains(hostArg, ":") {
		parts := strings.Split(hostArg, ":")
		host.HostName = parts[0]
		host.Port = parts[1]
	} else {
		host.HostName = hostArg
	}

	// 设置Host字段为HostName或原始参数
	if host.Host == "" {
		host.Host = host.HostName
	}

	return host
}