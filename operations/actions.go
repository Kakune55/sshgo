package operations

import (
	"fmt"
	"os"

	"sshgo/ssh"

	"github.com/manifoldco/promptui"
)

// DeleteKeyFile 删除密钥文件
func DeleteKeyFile(host ssh.SSHHost) error {
	if host.KeyFile == "" {
		return fmt.Errorf("该主机没有配置密钥文件")
	}
	
	// 确认删除
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("确定要删除密钥文件 '%s' 吗?", host.KeyFile),
		IsConfirm: true,
	}
	
	_, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("取消删除操作")
	}
	
	// 删除文件
	err = os.Remove(host.KeyFile)
	if err != nil {
		return fmt.Errorf("删除密钥文件失败: %v", err)
	}
	
	fmt.Printf("成功删除密钥文件: %s\n", host.KeyFile)
	return nil
}

// DeleteHostConfig 删除主机配置
func DeleteHostConfig(host ssh.SSHHost) error {
	// 确认删除
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("确定要删除主机 '%s' 的配置吗? 这将从config和known_hosts文件中移除相关记录", host.Host),
		IsConfirm: true,
	}
	
	_, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("取消删除操作")
	}
	
	// 从SSH配置文件中删除主机配置
	err = ssh.RemoveHostFromConfig(host.Host)
	if err != nil {
		fmt.Printf("警告: 从配置文件中删除主机配置时出错: %v\n", err)
	}
	
	// 从known_hosts文件中删除主机记录
	err = ssh.RemoveHostFromKnownHosts(host.Host)
	if err != nil {
		fmt.Printf("警告: 从known_hosts文件中删除主机记录时出错: %v\n", err)
	}
	
	fmt.Printf("成功删除主机 '%s' 的配置\n", host.Host)
	return nil
}

// ModifyUser 修改主机用户
func ModifyUser(host ssh.SSHHost) error {
	// 询问新用户名
	prompt := promptui.Prompt{
		Label:   "请输入新用户名",
		Default: host.User,
	}
	
	newUser, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("获取新用户名失败: %v", err)
	}
	
	// 保存新用户名到配置文件
	err = ssh.SaveUserToConfig(host.Host, newUser)
	if err != nil {
		return fmt.Errorf("保存用户名到配置文件失败: %v", err)
	}
	
	fmt.Printf("成功将主机 '%s' 的用户名修改为 '%s'\n", host.Host, newUser)
	return nil
}

// ModifyPort 修改主机端口
func ModifyPort(host ssh.SSHHost) error {
	// 询问新端口号
	prompt := promptui.Prompt{
		Label:   "请输入新端口号",
		Default: host.Port,
	}
	
	newPort, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("获取新端口号失败: %v", err)
	}
	
	// 保存新端口号到配置文件
	err = ssh.SavePortToConfig(host.Host, newPort)
	if err != nil {
		return fmt.Errorf("保存端口号到配置文件失败: %v", err)
	}
	
	fmt.Printf("成功将主机 '%s' 的端口号修改为 '%s'\n", host.Host, newPort)
	return nil
}