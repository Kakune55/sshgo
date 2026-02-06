package operations

import (
	"fmt"
	"os"

	"sshgo/i18n"
	"sshgo/ssh"
)

// DeleteKeyFile 删除密钥文件（无需确认，确认由 UI 层处理）
func DeleteKeyFile(host ssh.SSHHost) error {
	if host.KeyFile == "" {
		return fmt.Errorf("%s", i18n.T(i18n.FailedToDeleteKey))
	}

	// 删除文件
	err := os.Remove(host.KeyFile)
	if err != nil {
		return fmt.Errorf(i18n.T(i18n.FailedToDeleteKey), err)
	}

	return nil
}

// DeleteHostConfig 删除主机配置（无需确认，确认由 UI 层处理）
func DeleteHostConfig(host ssh.SSHHost) error {
	// 从SSH配置文件中删除主机配置
	err := ssh.RemoveHostFromConfig(host.Host)
	if err != nil {
		fmt.Printf(i18n.T(i18n.DeleteHostConfigWarning)+"\n", err)
	}

	// 从known_hosts文件中删除主机记录
	err = ssh.RemoveHostFromKnownHosts(host.Host)
	if err != nil {
		fmt.Printf(i18n.T(i18n.DeleteKnownHostsWarning)+"\n", err)
	}

	return nil
}

// ModifyUser 修改主机用户（用户名由 UI 层获取）
func ModifyUser(host ssh.SSHHost, newUser string) error {
	if newUser == "" {
		newUser = i18n.T(i18n.DefaultUsername)
	}

	// 保存新用户名到配置文件
	err := ssh.SaveUserToConfig(host.Host, newUser)
	if err != nil {
		return fmt.Errorf(i18n.T(i18n.FailedToModifyUser), err)
	}

	return nil
}

// ModifyPort 修改主机端口（端口由 UI 层获取）
func ModifyPort(host ssh.SSHHost, newPort string) error {
	if newPort == "" {
		newPort = "22"
	}

	// 保存新端口号到配置文件
	err := ssh.SavePortToConfig(host.Host, newPort)
	if err != nil {
		return fmt.Errorf(i18n.T(i18n.FailedToModifyPort), err)
	}

	return nil
}
