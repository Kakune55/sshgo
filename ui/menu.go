package ui

import (
	"fmt"
	"io"
	"strings"

	"sshgo/i18n"
	"sshgo/ssh"

	"github.com/manifoldco/promptui"
)

// ShowHostSelectionMenu 显示主机选择菜单
func ShowHostSelectionMenu(hosts []ssh.SSHHost) (ssh.SSHHost, string, error) {
	// 先显示模糊查找提示
		fmt.Println("提示: 可以在选择主机时使用模糊查找功能") // TODO: 需要翻译
	
	var hostNames []string
	hostMap := make(map[string]ssh.SSHHost)
	
	for _, host := range hosts {
		displayName := host.Host
		if host.HostName != "" {
			displayName = fmt.Sprintf("%s (%s)", host.Host, host.HostName)
		}
		hostNames = append(hostNames, displayName)
		hostMap[displayName] = host
	}
	
	// 添加退出搜索选项
		hostNames = append([]string{i18n.T(i18n.SearchHostOption)}, hostNames...)
		hostNames = append([]string{i18n.T(i18n.ExitOption)}, hostNames...)
	
	prompt := promptui.Select{
			Label: i18n.T(i18n.SelectHostLabel),
			Items: hostNames,
			CursorPos: 2,
			Size: 10,
		}
	
	index, result, err := prompt.Run()
		if err != nil {
			if err.Error() == "^C" {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			if err.Error() == "^D" {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			if err == io.EOF {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			return ssh.SSHHost{}, "", fmt.Errorf("选择主机失败: %v", err)
		}
	
	// 如果选择搜索
	if index == 0 {
		// 执行模糊查找
		return performFuzzySearch(hosts)
	}
	
	// 如果选择退出
	if index == len(hostNames)-1 {
		return ssh.SSHHost{}, "exit", nil
	}
	
	// 调整索引以匹配主机列表（因为添加了搜索选项）
	selectedHost := hostMap[result]
	
	// 显示操作菜单
	action, err := ShowActionMenu()
	if err != nil {
		if err.Error() == "^C" {
			return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
		}
		// 如果是返回操作，重新显示主菜单
		if err.Error() == "返回主菜单" {
			return ShowHostSelectionMenu(hosts)
		}
		return ssh.SSHHost{}, "", fmt.Errorf("选择操作失败: %v", err)
	}
	
	// 如果是返回操作，重新显示主菜单
	if action == "back" {
		return ShowHostSelectionMenu(hosts)
	}
	
	return selectedHost, action, nil
}

// performFuzzySearch 执行模糊查找
func performFuzzySearch(hosts []ssh.SSHHost) (ssh.SSHHost, string, error) {
	// 创建主机显示名称到主机对象的映射
	hostMap := make(map[string]ssh.SSHHost)
	var allHostNames []string
	
	for _, host := range hosts {
		displayName := host.Host
		if host.HostName != "" {
			displayName = fmt.Sprintf("%s (%s)", host.Host, host.HostName)
		}
		allHostNames = append(allHostNames, displayName)
		hostMap[displayName] = host
	}
	
	// 询问搜索关键词
		prompt := promptui.Prompt{
			Label: i18n.T(i18n.EnterSearchKeyword),
		}
	
	searchTerm, err := prompt.Run()
		if err != nil {
			if err.Error() == "^C" {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			if err == io.EOF {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			return ssh.SSHHost{}, "", fmt.Errorf("获取搜索关键词失败: %v", err)
		}
	
	// 执行模糊匹配
	var matchedHostNames []string
	for _, hostName := range allHostNames {
		if strings.Contains(strings.ToLower(hostName), strings.ToLower(searchTerm)) {
			matchedHostNames = append(matchedHostNames, hostName)
		}
	}
	
	// 如果没有匹配的主机
		if len(matchedHostNames) == 0 {
			fmt.Println(i18n.T(i18n.NoMatchingHosts))
			return ShowHostSelectionMenu(hosts)
		}
	
	// 如果只有一个匹配的主机，直接选择它
	if len(matchedHostNames) == 1 {
		selectedHost := hostMap[matchedHostNames[0]]
		action, err := ShowActionMenu()
		if err != nil {
			if err.Error() == "^C" {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			// 如果是返回操作，重新显示主菜单
			if err.Error() == "返回主菜单" {
				return ShowHostSelectionMenu(hosts)
			}
			return ssh.SSHHost{}, "", fmt.Errorf("选择操作失败: %v", err)
		}
		
		// 如果是返回操作，重新显示主菜单
		if action == "back" {
			return ShowHostSelectionMenu(hosts)
		}
		
		return selectedHost, action, nil
	}
	
	// 如果有多个匹配的主机，显示选择菜单
		matchedHostNames = append(matchedHostNames, i18n.T(i18n.BackAction))
		
		selectPrompt := promptui.Select{
			Label: i18n.T(i18n.MultipleMatchingHosts),
			Items: matchedHostNames,
			Size: 10,
		}
	
	index, result, err := selectPrompt.Run()
		if err != nil {
			if err.Error() == "^C" {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			if err == io.EOF {
				return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
			}
			return ssh.SSHHost{}, "", fmt.Errorf("选择主机失败: %v", err)
		}
	
	// 如果选择返回
	if index == len(matchedHostNames)-1 {
		return ShowHostSelectionMenu(hosts)
	}
	
	selectedHost := hostMap[result]
	action, err := ShowActionMenu()
	if err != nil {
		if err.Error() == "^C" {
			return ssh.SSHHost{}, "", fmt.Errorf("interrupt")
		}
		// 如果是返回操作，重新显示主菜单
		if err.Error() == "返回主菜单" {
			return ShowHostSelectionMenu(hosts)
		}
		return ssh.SSHHost{}, "", fmt.Errorf("选择操作失败: %v", err)
	}
	
	// 如果是返回操作，重新显示主菜单
	if action == "back" {
		return ShowHostSelectionMenu(hosts)
	}
	
	return selectedHost, action, nil
}

// ShowActionMenu 显示操作菜单
func ShowActionMenu() (string, error) {
	actions := []string{
		i18n.T(i18n.ConnectAction),
		i18n.T(i18n.DetailsAction),
		i18n.T(i18n.DeleteKeyAction),
		i18n.T(i18n.DeleteConfigAction),
		i18n.T(i18n.ModifyUserAction),
		i18n.T(i18n.ModifyPortAction),
		i18n.T(i18n.BackAction),
	}
	
	prompt := promptui.Select{
		Label: i18n.T(i18n.SelectActionLabel),
		Items: actions,
		Size: 10,
	}
	
	_, result, err := prompt.Run()
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interrupt")
			}
			return "", err
		}
	
	switch result {
		case i18n.T(i18n.ConnectAction):
			return "connect", nil
		case i18n.T(i18n.DetailsAction):
			return "details", nil
		case i18n.T(i18n.DeleteKeyAction):
			return "delete_key", nil
		case i18n.T(i18n.DeleteConfigAction):
			return "delete_config", nil
		case i18n.T(i18n.ModifyUserAction):
			return "modify_user", nil
		case i18n.T(i18n.ModifyPortAction):
			return "modify_port", nil
		case i18n.T(i18n.BackAction):
			return "back", fmt.Errorf("返回主菜单") // TODO: 需要翻译
		default:
			return "connect", nil
		}
}

// ShowHostDetails 显示主机详细信息
func ShowHostDetails(host ssh.SSHHost) {
	fmt.Println("\n" + i18n.T(i18n.HostDetailsTitle))
	fmt.Printf(i18n.T(i18n.HostAlias) + "\n", host.Host)
	if host.HostName != "" {
		fmt.Printf(i18n.T(i18n.HostName) + "\n", host.HostName)
	}
	if host.User != "" {
		fmt.Printf(i18n.T(i18n.UserName) + "\n", host.User)
	}
	if host.Port != "22" {
		fmt.Printf(i18n.T(i18n.Port) + "\n", host.Port)
	}
	if host.KeyFile != "" {
		fmt.Printf(i18n.T(i18n.KeyFile) + "\n", host.KeyFile)
	}
	fmt.Print("====================\n\n") // TODO: 需要翻译
}