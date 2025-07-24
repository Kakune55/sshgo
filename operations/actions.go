package operations

import (
	"fmt"
	"os"

	"sshgo/i18n"
	"sshgo/ssh"

	"github.com/manifoldco/promptui"
)

// DeleteKeyFile 删除密钥文件
func DeleteKeyFile(host ssh.SSHHost) error {
	if host.KeyFile == "" {
			return fmt.Errorf("%s", i18n.T(i18n.FailedToDeleteKey))
		}
	
	// 确认删除
		prompt := promptui.Prompt{
			Label:     i18n.TWithArgs(i18n.ConfirmDeleteKey, host.KeyFile),
			IsConfirm: true,
		}
	
	_, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("%s", i18n.T(i18n.CancelOperation))
			}
	
	// 删除文件
		err = os.Remove(host.KeyFile)
		if err != nil {
					msg := i18n.TWithArgs(i18n.FailedToDeleteKey, err)
					return fmt.Errorf("%s", msg)
				}
				
				msg := i18n.TWithArgs(i18n.SuccessfullyDeletedKey, host.KeyFile)
				fmt.Printf("%s\n", msg)
	return nil
}

// DeleteHostConfig 删除主机配置
func DeleteHostConfig(host ssh.SSHHost) error {
	// 确认删除
		prompt := promptui.Prompt{
			Label:     i18n.TWithArgs(i18n.ConfirmDeleteConfig, host.Host),
			IsConfirm: true,
		}
	
	_, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("%s", i18n.T(i18n.CancelOperation))
			}
		
		// 从SSH配置文件中删除主机配置
				err = ssh.RemoveHostFromConfig(host.Host)
				if err != nil {
					fmt.Printf("警告: 从配置文件中删除主机配置时出错: %v\n", err) // TODO: 需要翻译
				}
				
				// 从known_hosts文件中删除主机记录
				err = ssh.RemoveHostFromKnownHosts(host.Host)
				if err != nil {
					fmt.Printf("警告: 从known_hosts文件中删除主机记录时出错: %v\n", err) // TODO: 需要翻译
				}
		
		msg := i18n.TWithArgs(i18n.SuccessfullyDeletedConfig, host.Host)
		fmt.Printf("%s\n", msg)
	return nil
}

// ModifyUser 修改主机用户
func ModifyUser(host ssh.SSHHost) error {
	// 询问新用户名
		prompt := promptui.Prompt{
			Label:   i18n.T(i18n.EnterNewUsername),
			Default: host.User,
		}
	
	newUser, err := prompt.Run()
			if err != nil {
				msg := i18n.TWithArgs(i18n.FailedToGetUsername, err)
				return fmt.Errorf("%s", msg)
			}
		
		// 保存新用户名到配置文件
			err = ssh.SaveUserToConfig(host.Host, newUser)
			if err != nil {
				msg := i18n.TWithArgs(i18n.FailedToModifyUser, err)
				return fmt.Errorf("%s", msg)
			}
			
			msg := i18n.TWithArgs(i18n.SuccessfullyModifiedUser, host.Host, newUser)
			fmt.Printf("%s\n", msg)
	return nil
}

// ModifyPort 修改主机端口
func ModifyPort(host ssh.SSHHost) error {
	// 询问新端口号
		prompt := promptui.Prompt{
			Label:   i18n.T(i18n.EnterNewPort),
			Default: host.Port,
		}
	
	newPort, err := prompt.Run()
			if err != nil {
				msg := i18n.TWithArgs(i18n.FailedToGetPort, err)
				return fmt.Errorf("%s", msg)
			}
			
			// 保存新端口号到配置文件
			err = ssh.SavePortToConfig(host.Host, newPort)
			if err != nil {
				msg := i18n.TWithArgs(i18n.FailedToModifyPort, err)
				return fmt.Errorf("%s", msg)
			}
		
		msg := i18n.TWithArgs(i18n.SuccessfullyModifiedPort, host.Host, newPort)
		fmt.Printf("%s\n", msg)
	return nil
}