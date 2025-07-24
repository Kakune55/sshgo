package main

import (
	"fmt"
	"os"
	"strings"

	"sshgo/i18n"
	"sshgo/operations"
	"sshgo/ssh"
	"sshgo/ui"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		// 如果提供了参数，直接连接到指定主机
		hostArg := os.Args[1]
		host := parseHostArgument(hostArg)
		err := ssh.ConnectToHost(host)
		if err != nil {
					fmt.Printf("%v\n", err) // TODO: 需要翻译
				}
		return
	}
	
	// 获取SSH配置文件路径
	configPath := ssh.GetSSHConfigPath()
	
	// 解析SSH配置文件
	hosts, err := ssh.ParseSSHConfig(configPath)
	if err != nil {
			fmt.Printf("%v\n", err) // TODO: 需要翻译
			return
		}
		
		if len(hosts) == 0 {
			fmt.Println("未找到SSH主机配置") // TODO: 需要翻译
			return
		}
	
	for {
		// 显示交互式菜单
		selectedHost, action, err := ui.ShowHostSelectionMenu(hosts)
		if err != nil {
					if err.Error() == "interrupt" {
						fmt.Println("\n" + i18n.T(i18n.Goodbye))
						return
					}
					fmt.Printf("%v\n", err) // TODO: 需要翻译
					return
				}
		
		switch action {
		case "connect":
			// 连接到选中的主机
						err = ssh.ConnectToHost(selectedHost)
						if err != nil {
							fmt.Printf("%v\n", err) // TODO: 需要翻译
						}
		case "details":
			// 显示主机详细信息
			ui.ShowHostDetails(selectedHost)
		case "delete_key":
			// 删除密钥文件
						err = operations.DeleteKeyFile(selectedHost)
						if err != nil {
							fmt.Printf("%v\n", err) // TODO: 需要翻译
						}
		case "delete_config":
			// 删除配置
						err = operations.DeleteHostConfig(selectedHost)
						if err != nil {
							fmt.Printf("%v\n", err) // TODO: 需要翻译
						} else {
							// 重新加载主机列表
							fmt.Println("配置已删除，重新加载主机列表...") // TODO: 需要翻译
							// 重新解析SSH配置文件
							hosts, err = ssh.ParseSSHConfig(configPath)
							if err != nil {
								fmt.Printf("%v\n", err) // TODO: 需要翻译
								return
							}
							if len(hosts) == 0 {
								fmt.Println("未找到SSH主机配置") // TODO: 需要翻译
								return
							}
						}
		case "modify_user":
			// 修改用户
						err = operations.ModifyUser(selectedHost)
						if err != nil {
							fmt.Printf("%v\n", err) // TODO: 需要翻译
						}
		case "modify_port":
			// 修改端口
						err = operations.ModifyPort(selectedHost)
						if err != nil {
							fmt.Printf("%v\n", err) // TODO: 需要翻译
						}
		case "exit":
					fmt.Println(i18n.T(i18n.Goodbye))
					return
		case "back":
			// 返回主菜单
			continue
		}
	}
}

// parseHostArgument 解析命令行参数中的主机信息
func parseHostArgument(hostArg string) ssh.SSHHost {
	host := ssh.SSHHost{
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