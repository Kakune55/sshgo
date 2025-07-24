package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"sshgo/i18n"

	"github.com/manifoldco/promptui"
)

// ConnectToHost 连接到指定主机
func ConnectToHost(host SSHHost) error {
	// 如果没有用户名，询问用户
	user := host.User
	if user == "" {
		prompt := promptui.Prompt{
					Label:   i18n.T(i18n.EnterNewUsername),
					Default: "root",
				}
		
		var err error
				user, err = prompt.Run()
				if err != nil {
					return fmt.Errorf("%s", i18n.TWithArgs(i18n.FailedToGetUsername, err))
				}
				
				// 保存用户名到配置文件
				err = SaveUserToConfig(host.Host, user)
				if err != nil {
					fmt.Printf("警告: %v\n", err) // TODO: 需要翻译
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
	
	// 执行SSH命令
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	fmt.Printf("%s", i18n.TWithArgs(i18n.ConnectingTo, host.User, host.Host) + "\n")
	return cmd.Run()
}