package main

import (
	"fmt"
	"os"

	"sshgo/ssh"
	"sshgo/ui"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		// 如果提供了参数，直接连接到指定主机
		hostArg := os.Args[1]
		host := ssh.ParseHostArgument(hostArg)
		
		// 如果没有用户名，使用默认用户名
		if host.User == "" {
			host.User = "root"
		}
		
		err := ssh.ConnectToHost(host)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return
	}

	// 运行主 UI 循环
	ui.RunLoop()
}
