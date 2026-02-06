package ssh

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"sshgo/i18n"
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
			return hosts, fmt.Errorf("%s", i18n.TWithArgs(i18n.ParseConfigError, path, err))
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
		fmt.Printf(i18n.T(i18n.ReadKnownHostsWarning)+"\n", err)
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
				return hosts, fmt.Errorf(i18n.T(i18n.ParseTestConfigError), err)
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
		return hosts, fmt.Errorf(i18n.T(i18n.ReadConfigFileError), err)
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
		return hosts, fmt.Errorf(i18n.T(i18n.ReadKnownHostsError), err)
	}

	return hosts, nil
}

// SaveUserToConfig 保存用户名到SSH配置文件
func SaveUserToConfig(host, user string) error {
	return UpdateHostDirective(host, "User", user)
}

// SavePortToConfig 保存端口号到SSH配置文件
func SavePortToConfig(host, port string) error {
	return UpdateHostDirective(host, "Port", port)
}

// UpdateHostDirective 通用更新/新增某个 Host 下的指令 (如 User / Port)
// 逻辑：
// 1. 读取配置 -> 切分为多个块（由 Host 行开始）
// 2. 寻找包含目标 host 的 Host 行（Host 行可能包含多个模式，这里按精确匹配其中一个）
// 3. 在该块内部：若存在同名指令，覆盖；否则在块末尾追加（保持缩进 4 空格）
// 4. 若不存在该 host 块，则在文件末尾新建一个块
func UpdateHostDirective(host, directive, value string) error {
	configPath := GetSSHConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			data = []byte("")
		} else {
			return fmt.Errorf(i18n.T(i18n.ReadConfigFileFailed), err)
		}
	}

	lines := strings.Split(string(data), "\n")

	// 解析为 blocks
	type block struct {
		header string   // 原始 Host 行
		hosts  []string // Host 行里的所有主机模式
		body   []string // 不包含 header 的后续行（直到下一个 Host 行 / 文件结束）
	}

	var blocks []block
	var current *block

	flush := func() {
		if current != nil {
			blocks = append(blocks, *current)
			current = nil
		}
	}

	for _, raw := range lines {
		trim := strings.TrimSpace(raw)
		if strings.HasPrefix(strings.ToLower(trim), "host ") {
			// 新的 block
			flush()
			fields := strings.Fields(trim)
			var hostSpecs []string
			if len(fields) > 1 {
				hostSpecs = fields[1:]
			}
			current = &block{header: raw, hosts: hostSpecs}
			continue
		}
		if current == nil {
			// 文件可能前面有非 Host 行（不标准），我们直接跳过或放入匿名 block
			if raw == "" && len(blocks) == 0 {
				// 顶部空行忽略
				continue
			}
			// 放入一个无 header 的 block（用于保留可能的注释）
			current = &block{header: "", hosts: nil}
		}
		current.body = append(current.body, raw)
	}
	flush()

	// 查找目标 host block
	targetIndex := -1
	for i, b := range blocks {
		for _, h := range b.hosts {
			if h == host { // 精确匹配
				targetIndex = i
				break
			}
		}
		if targetIndex >= 0 { break }
	}

	directiveLower := strings.ToLower(directive)
	updated := false

	if targetIndex >= 0 {
		b := blocks[targetIndex]
		for i, bodyLine := range b.body {
			trim := strings.TrimSpace(bodyLine)
			if strings.HasPrefix(strings.ToLower(trim), strings.ToLower(directive)+" ") {
				b.body[i] = fmt.Sprintf("    %s %s", directive, value)
				updated = true
				break
			}
		}
		if !updated { // append at end (ensure no duplicate trailing blank lines)
			b.body = append(b.body, fmt.Sprintf("    %s %s", directive, value))
		}
		blocks[targetIndex] = b
	} else {
		// 创建新 block
		newBlock := block{
			header: fmt.Sprintf("Host %s", host),
			hosts:  []string{host},
			body:   []string{fmt.Sprintf("    %s %s", directive, value), ""}, // 结尾空行
		}
		blocks = append(blocks, newBlock)
	}

	// 重新拼接
	var out []string
	for _, b := range blocks {
		if b.header != "" {
			out = append(out, b.header)
		}
		out = append(out, b.body...)
	}

	// 清理多余末尾空行
	for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
		out = out[:len(out)-1]
	}
	output := strings.Join(out, "\n") + "\n"

	if err := os.WriteFile(configPath, []byte(output), 0600); err != nil {
		return fmt.Errorf(i18n.T(i18n.WriteConfigFileFailed), err)
	}

	_ = directiveLower // 预留后续需要大小写归一化的扩展
	return nil
}

// RemoveHostFromConfig 从SSH配置文件中删除主机配置
func RemoveHostFromConfig(hostName string) error {
	configPath := GetSSHConfigPath()

	// 读取现有配置文件内容
	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，无需操作
		}
		return fmt.Errorf(i18n.T(i18n.ReadConfigFileFailed), err)
	}

	lines := strings.Split(string(content), "\n")
	newLines := []string{}
	inSectionToDelete := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		isHostLine := strings.HasPrefix(trimmedLine, "Host ")

		if isHostLine {
			// 检查此Host行是否是我们想删除的目标
			fields := strings.Fields(trimmedLine)
			if len(fields) > 1 && fields[1] == hostName {
				inSectionToDelete = true // 进入了待删除的配置区块
			} else {
				inSectionToDelete = false // 遇到了另一个Host，退出删除模式
			}
		}

		if !inSectionToDelete {
			newLines = append(newLines, line)
		}
	}

	// 将修改后的内容写回文件
	output := strings.Join(newLines, "\n")
	err = os.WriteFile(configPath, []byte(output), 0600)
	if err != nil {
		return fmt.Errorf(i18n.T(i18n.WriteConfigFileFailed), err)
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
		return fmt.Errorf(i18n.T(i18n.ReadKnownHostsFailed), err)
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
		return fmt.Errorf(i18n.T(i18n.WriteKnownHostsFailed), err)
	}

	return nil
}
