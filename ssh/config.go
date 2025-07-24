package ssh

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetAllConfigPaths 获取所有可能的配置文件路径
func GetAllConfigPaths() []string {
	var configPaths []string
	
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		configPaths = append(configPaths, filepath.Join(home, ".ssh", "config"))
		
		// 添加Windows系统中常见的SSH配置文件路径
		configPaths = append(configPaths, `C:\Program Files\Git\etc\ssh\ssh_config`)
		configPaths = append(configPaths, `C:\Program Files (x86)\Git\etc\ssh\ssh_config`)
		configPaths = append(configPaths, filepath.Join(os.Getenv("PROGRAMDATA"), "ssh", "ssh_config"))
	} else {
		home := os.Getenv("HOME")
		configPaths = append(configPaths, filepath.Join(home, ".ssh", "config"))
	}
	
	return configPaths
}

// GetSSHConfigPath 获取SSH配置文件路径
func GetSSHConfigPath() string {
	// 获取所有可能的配置文件路径
	configPaths := GetAllConfigPaths()
	
	// 返回第一个存在的配置文件路径
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// 如果没有找到配置文件，返回默认路径
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		return filepath.Join(home, ".ssh", "config")
	} else {
		home := os.Getenv("HOME")
		return filepath.Join(home, ".ssh", "config")
	}
}

// GetKnownHostsPath 获取known_hosts文件路径
func GetKnownHostsPath() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		return filepath.Join(home, ".ssh", "known_hosts")
	} else {
		home := os.Getenv("HOME")
		return filepath.Join(home, ".ssh", "known_hosts")
	}
}

// ParseSSHConfig 解析SSH配置文件
func ParseSSHConfig(configPath string) ([]SSHHost, error) {
	var hosts []SSHHost
	hostMap := make(map[string]SSHHost) // 用于去重
	
	// 获取所有可能的配置文件路径
	configPaths := GetAllConfigPaths()
	
	// 从所有配置文件中读取主机信息
	for _, path := range configPaths {
		hostsFromFile, err := parseSingleConfigFile(path)
		if err != nil {
			// 如果是文件不存在的错误，跳过继续处理其他文件
			if os.IsNotExist(err) {
				continue
			}
			return hosts, fmt.Errorf("解析配置文件 %s 时出错: %v", path, err)
		}
		
		// 将主机信息添加到结果中，避免重复
		for _, host := range hostsFromFile {
			if _, exists := hostMap[host.Host]; !exists {
				hostMap[host.Host] = host
				hosts = append(hosts, host)
			}
		}
	}
	
	// 从known_hosts文件中读取主机信息
	knownHosts, err := parseKnownHosts()
	if err != nil {
		// 如果读取known_hosts文件出错，只打印警告信息，不中断程序
		fmt.Printf("警告: 读取known_hosts文件时出错: %v\n", err)
	} else {
		// 将known_hosts中的主机信息添加到结果中，避免重复
		for _, host := range knownHosts {
			if _, exists := hostMap[host.Host]; !exists {
				hostMap[host.Host] = host
				hosts = append(hosts, host)
			}
		}
	}
	
	// 如果没有找到任何配置文件，尝试使用测试配置文件
	if len(hosts) == 0 {
		if _, err := os.Stat("test_config"); err == nil {
			hosts, err = parseSingleConfigFile("test_config")
			if err != nil {
				return hosts, fmt.Errorf("解析测试配置文件时出错: %v", err)
			}
		}
	}
	
	return hosts, nil
}

// parseSingleConfigFile 解析单个配置文件
func parseSingleConfigFile(configPath string) ([]SSHHost, error) {
	var hosts []SSHHost
	
	file, err := os.Open(configPath)
	if err != nil {
		return hosts, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	var currentHost *SSHHost
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// 分割键值对
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		
		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")
		
		if key == "host" {
			// 创建新的主机配置
			if currentHost != nil {
				hosts = append(hosts, *currentHost)
			}
			currentHost = &SSHHost{
				Host: value,
				Port: "22", // 默认端口
			}
		} else if currentHost != nil {
			switch key {
			case "hostname":
				currentHost.HostName = value
			case "user":
				currentHost.User = value
			case "port":
				currentHost.Port = value
			case "identityfile":
				// 处理波浪号路径
				if strings.HasPrefix(value, "~") {
					home := os.Getenv("HOME")
					if runtime.GOOS == "windows" {
						home = os.Getenv("USERPROFILE")
					}
					value = filepath.Join(home, value[2:])
				}
				currentHost.KeyFile = value
			}
		}
	}
	
	// 添加最后一个主机
	if currentHost != nil {
		hosts = append(hosts, *currentHost)
	}
	
	if err := scanner.Err(); err != nil {
		return hosts, fmt.Errorf("读取配置文件时出错: %v", err)
	}
	
	return hosts, nil
}

// parseKnownHosts 从known_hosts文件中解析主机信息
func parseKnownHosts() ([]SSHHost, error) {
	var hosts []SSHHost
	hostMap := make(map[string]bool) // 用于去重
	
	knownHostsPath := GetKnownHostsPath()
	
	file, err := os.Open(knownHostsPath)
	if err != nil {
		return hosts, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// 分割字段
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		
		// 第一个字段是主机名或IP地址
		hostName := parts[0]
		
		// 处理逗号分隔的多个主机名
		hostNames := strings.Split(hostName, ",")
		
		for _, name := range hostNames {
			// 去除可能的端口号
			if strings.Contains(name, ":") {
				name = strings.Split(name, ":")[0]
			}
			
			// 如果主机名还没有添加过，则添加到结果中
			if _, exists := hostMap[name]; !exists {
				hostMap[name] = true
				hosts = append(hosts, SSHHost{
					Host:     name,
					HostName: name,
					Port:     "22", // 默认端口
				})
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		return hosts, fmt.Errorf("读取known_hosts文件时出错: %v", err)
	}
	
	return hosts, nil
}

// SaveUserToConfig 保存用户名到SSH配置文件
func SaveUserToConfig(host, user string) error {
	configPath := GetSSHConfigPath()
	
	// 读取现有配置文件内容
	content, err := os.ReadFile(configPath)
	if err != nil {
		// 如果文件不存在，创建新文件
		if os.IsNotExist(err) {
			content = []byte{}
		} else {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
	}
	
	// 将内容转换为字符串
	configStr := string(content)
	
	// 检查是否已存在该主机的配置
	hostSection := fmt.Sprintf("Host %s", host)
	if strings.Contains(configStr, hostSection) {
		// 如果已存在该主机配置，添加User行（如果还没有的话）
		lines := strings.Split(configStr, "\n")
		newLines := []string{}
		inHostSection := false
		userAdded := false
		
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			
			if strings.HasPrefix(trimmedLine, "Host ") {
				inHostSection = (trimmedLine == hostSection)
			}
			
			if inHostSection && strings.HasPrefix(trimmedLine, "User ") {
				// 如果已存在User行，更新用户名
				line = fmt.Sprintf("    User %s", user)
				userAdded = true
			}
			
			newLines = append(newLines, line)
			
			// 如果在主机配置段落中且遇到空行或下一个Host行，添加User行
			if inHostSection && !userAdded && (trimmedLine == "" || strings.HasPrefix(trimmedLine, "Host ")) {
				if trimmedLine != "" {
					// 如果是下一个Host行，先添加User行
					newLines = append(newLines[:len(newLines)-1], fmt.Sprintf("    User %s", user), line)
				} else {
					// 如果是空行，添加User行
					newLines[len(newLines)-1] = fmt.Sprintf("    User %s", user)
					newLines = append(newLines, "")
				}
				userAdded = true
				inHostSection = false
			}
		}
		
		// 如果遍历完所有行仍未添加User行，说明Host段落没有结束，在末尾添加
		if inHostSection && !userAdded {
			newLines = append(newLines, fmt.Sprintf("    User %s", user))
		}
		
		configStr = strings.Join(newLines, "\n")
	} else {
		// 如果不存在该主机配置，添加新的主机配置段落
		if len(configStr) > 0 && !strings.HasSuffix(configStr, "\n") {
			configStr += "\n"
		}
		configStr += fmt.Sprintf("Host %s\n    User %s\n\n", host, user)
	}
	
	// 写入更新后的内容到配置文件
	err = os.WriteFile(configPath, []byte(configStr), 0600)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	return nil
}

// SavePortToConfig 保存端口号到SSH配置文件
func SavePortToConfig(host, port string) error {
	configPath := GetSSHConfigPath()
	
	// 读取现有配置文件内容
	content, err := os.ReadFile(configPath)
	if err != nil {
		// 如果文件不存在，创建新文件
		if os.IsNotExist(err) {
			content = []byte{}
		} else {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
	}
	
	// 将内容转换为字符串
	configStr := string(content)
	
	// 检查是否已存在该主机的配置
	hostSection := fmt.Sprintf("Host %s", host)
	if strings.Contains(configStr, hostSection) {
		// 如果已存在该主机配置，添加Port行（如果还没有的话）
		lines := strings.Split(configStr, "\n")
		newLines := []string{}
		inHostSection := false
		portAdded := false
		
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			
			if strings.HasPrefix(trimmedLine, "Host ") {
				inHostSection = (trimmedLine == hostSection)
			}
			
			if inHostSection && strings.HasPrefix(trimmedLine, "Port ") {
				// 如果已存在Port行，更新端口号
				line = fmt.Sprintf("    Port %s", port)
				portAdded = true
			}
			
			newLines = append(newLines, line)
			
			// 如果在主机配置段落中且遇到空行或下一个Host行，添加Port行
			if inHostSection && !portAdded && (trimmedLine == "" || strings.HasPrefix(trimmedLine, "Host ")) {
				if trimmedLine != "" {
					// 如果是下一个Host行，先添加Port行
					newLines = append(newLines[:len(newLines)-1], fmt.Sprintf("    Port %s", port), line)
				} else {
					// 如果是空行，添加Port行
					newLines[len(newLines)-1] = fmt.Sprintf("    Port %s", port)
					newLines = append(newLines, "")
				}
				portAdded = true
				inHostSection = false
			}
		}
		
		// 如果遍历完所有行仍未添加Port行，说明Host段落没有结束，在末尾添加
		if inHostSection && !portAdded {
			newLines = append(newLines, fmt.Sprintf("    Port %s", port))
		}
		
		configStr = strings.Join(newLines, "\n")
	} else {
		// 如果不存在该主机配置，添加新的主机配置段落
		if len(configStr) > 0 && !strings.HasSuffix(configStr, "\n") {
			configStr += "\n"
		}
		configStr += fmt.Sprintf("Host %s\n    Port %s\n\n", host, port)
	}
	
	// 写入更新后的内容到配置文件
	err = os.WriteFile(configPath, []byte(configStr), 0600)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	return nil
}

// RemoveHostFromConfig 从SSH配置文件中删除主机配置
func RemoveHostFromConfig(hostName string) error {
	configPath := GetSSHConfigPath()
	
	// 读取现有配置文件内容
	content, err := os.ReadFile(configPath)
	if err != nil {
		// 如果文件不存在，直接返回
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("读取配置文件失败: %v", err)
	}
	
	// 将内容转换为字符串并按行分割
	lines := strings.Split(string(content), "\n")
	
	// 查找并删除主机配置段落
	newLines := []string{}
	inHostSection := false
	skipSection := false
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// 检查是否是Host行
		if strings.HasPrefix(trimmedLine, "Host ") {
			// 检查是否是我们要删除的主机
			hostNames := strings.Fields(trimmedLine[5:]) // 去掉"Host "前缀
			skipSection = false
			for _, name := range hostNames {
				if name == hostName {
					skipSection = true
					break
				}
			}
			inHostSection = true
		}
		
		// 如果不在要删除的主机段落中，则保留该行
		if !skipSection {
			newLines = append(newLines, line)
		}
		
		// 检查段落是否结束（遇到空行或其他Host行）
		if inHostSection && (trimmedLine == "" || (strings.HasPrefix(trimmedLine, "Host ") && skipSection)) {
			if skipSection {
				// 如果是结束的Host行，需要保留
				if strings.HasPrefix(trimmedLine, "Host ") && !strings.Contains(trimmedLine, hostName) {
					newLines = append(newLines, line)
				}
				skipSection = false
			}
			inHostSection = strings.HasPrefix(trimmedLine, "Host ")
		}
	}
	
	// 写入更新后的内容到配置文件
	err = os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0600)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	return nil
}

// RemoveHostFromKnownHosts 从known_hosts文件中删除主机记录
func RemoveHostFromKnownHosts(hostName string) error {
	knownHostsPath := GetKnownHostsPath()
	
	// 读取现有known_hosts文件内容
	content, err := os.ReadFile(knownHostsPath)
	if err != nil {
		// 如果文件不存在，直接返回
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("读取known_hosts文件失败: %v", err)
	}
	
	// 将内容转换为字符串并按行分割
	lines := strings.Split(string(content), "\n")
	
	// 查找并删除主机记录
	newLines := []string{}
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// 跳过空行和注释
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			newLines = append(newLines, line)
			continue
		}
		
		// 分割字段
		parts := strings.Fields(trimmedLine)
		if len(parts) < 2 {
			newLines = append(newLines, line)
			continue
		}
		
		// 第一个字段是主机名或IP地址
		hostField := parts[0]
		
		// 处理逗号分隔的多个主机名
		hostNames := strings.Split(hostField, ",")
		keepLine := true
		
		for _, name := range hostNames {
			// 去除可能的端口号
			if strings.Contains(name, ":") {
				name = strings.Split(name, ":")[0]
			}
			
			// 如果是我们要删除的主机，标记为不保留
			if name == hostName {
				keepLine = false
				break
			}
		}
		
		// 如果不是要删除的主机，保留该行
		if keepLine {
			newLines = append(newLines, line)
		}
	}
	
	// 写入更新后的内容到known_hosts文件
	err = os.WriteFile(knownHostsPath, []byte(strings.Join(newLines, "\n")), 0600)
	if err != nil {
		return fmt.Errorf("写入known_hosts文件失败: %v", err)
	}
	
	return nil
}