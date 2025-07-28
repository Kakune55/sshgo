package main

import (
	"fmt"
	"os"

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
		host := ssh.ParseHostArgument(hostArg)
		err := ssh.ConnectToHost(host)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return
	}

	// 获取SSH配置文件路径
	configPath := ssh.GetSSHConfigPath()

	// 解析SSH配置文件
	hosts, err := ssh.ParseSSHConfig(configPath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if len(hosts) == 0 {
		fmt.Println("未找到SSH主机配置")
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
			fmt.Printf("%v\n", err)
			return
		}

		switch action {
		case "connect":
			// 连接到选中的主机
			err = ssh.ConnectToHost(selectedHost)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		case "details":
			// 显示主机详细信息
			ui.ShowHostDetails(selectedHost)
		case "delete_key":
			// 删除密钥文件
			err = operations.DeleteKeyFile(selectedHost)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		case "delete_config":
			// 删除配置
			err = operations.DeleteHostConfig(selectedHost)
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				// 重新加载主机列表
				fmt.Println("配置已删除，重新加载主机列表...")
				// 重新解析SSH配置文件
				hosts, err = ssh.ParseSSHConfig(configPath)
				if err != nil {
					fmt.Printf("%v\n", err)
					return
				}
				if len(hosts) == 0 {
					fmt.Println("未找到SSH主机配置")
					return
				}
			}
		case "modify_user":
			// 修改用户
			err = operations.ModifyUser(selectedHost)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		case "modify_port":
			// 修改端口
			err = operations.ModifyPort(selectedHost)
			if err != nil {
				fmt.Printf("%v\n", err)
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
