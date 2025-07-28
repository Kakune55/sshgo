package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"sshgo/i18n"

	"github.com/manifoldco/promptui"
)

// validateSSHCommand 验证SSH命令和参数的有效性
func validateSSHCommand(args []string) error {
	// 检查SSH命令是否存在
	sshPath, err := exec.LookPath("ssh")
	// sshPath变量已被声明但未使用，这里添加一个简单的使用
	_ = sshPath // 忽略未使用的变量警告
	if err != nil {
		return fmt.Errorf("SSH command not found: %w", err)
	}

	// 检查参数是否为空
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}

	// 检查主机名是否为空
	hostName := args[len(args)-1] // 主机名通常是最后一个参数
	if hostName == "" {
		return fmt.Errorf("host name is required")
	}

	// 检查端口号是否有效（如果提供了的话）
	for i, arg := range args {
		if arg == "-p" && i+1 < len(args) {
			port := args[i+1]
			// 简单验证端口号是否为数字
			if port == "" {
				return fmt.Errorf("port number is required when -p flag is used")
			}
		}
	}

	// 检查密钥文件是否存在（如果提供了的话）
	for i, arg := range args {
		if arg == "-i" && i+1 < len(args) {
			keyFile := args[i+1]
			if keyFile == "" {
				return fmt.Errorf("key file path is required when -i flag is used")
			}
			// 检查文件是否存在
			if _, err := os.Stat(keyFile); os.IsNotExist(err) {
				return fmt.Errorf("key file does not exist: %s", keyFile)
			}
		}
	}
	
	return nil
}



// ConnectToHost 连接到指定主机
func ConnectToHost(host SSHHost) error {
	// 如果没有用户名，询问用户
	user := host.User
	if user == "" {
		prompt := promptui.Prompt{
			Label:   i18n.T(i18n.EnterNewUsername),
			Default: i18n.T(i18n.DefaultUsername),
		}

		var err error
		user, err = prompt.Run()
		if err != nil {
			return fmt.Errorf("%s", i18n.TWithArgs(i18n.FailedToGetUsername, err))
		}

		// 保存用户名到配置文件
		err = SaveUserToConfig(host.Host, user)
		if err != nil {
			fmt.Printf("警告: %v\n", err)
		}

		// 更新主机信息
		host.User = user
	}

	// 构建SSH命令
	args := []string{}

	if host.User != "" {
		args = append(args, "-l", host.User)
	}

	if host.Port != "22" {
		args = append(args, "-p", host.Port)
	}

	if host.KeyFile != "" {
		args = append(args, "-i", host.KeyFile)
	}

	hostName := host.HostName
	if hostName == "" {
		hostName = host.Host
	}

	args = append(args, hostName)

	// 预校验SSH命令和参数
	if err := validateSSHCommand(args); err != nil {
		return fmt.Errorf("%s", i18n.TWithArgs(i18n.InvalidSSHCommand, err))
	}

	// 执行SSH命令
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("%s", i18n.TWithArgs(i18n.ConnectingTo, host.User, host.Host)+"\n")
	return cmd.Run()
}
