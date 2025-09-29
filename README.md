# SSHGo - 简单的SSH管理工具

SSHGo是一个简单易用的SSH管理工具，支持跨平台使用。它可以从SSH配置文件中自动读取主机信息，并提供交互式界面选择要连接的主机。

## 功能特性

- 从SSH配置文件自动读取主机列表
- 交互式终端界面，使用上下键选择主机
- 支持直接命令行参数指定主机（如 `sshgo root@192.168.1.100`）
- 支持主机详细信息查看
- 支持删除密钥文件功能
- 支持修改主机用户和端口
- 支持模糊查找主机功能
- 支持基础网络诊断（TCP延迟、路由追踪）
- 自动语言检测（中/英），可通过环境变量覆盖 `SSHGO_LANG=zh|en`
- 跨平台支持（Windows、Linux、macOS）

## 安装

### 方法1：使用Go安装
```bash
go install github.com/yourusername/sshgo@latest
```

### 方法2：从源码构建
```bash
git clone https://github.com/yourusername/sshgo.git
cd sshgo
go build -o sshgo main.go
```

## 使用方法

### 交互式使用
直接运行程序，将显示SSH配置文件中的主机列表：
```bash
./sshgo
```

使用上下键选择主机，按回车确认选择，然后选择操作：
- 连接：直接SSH连接到选中的主机
- 详细信息：查看主机的详细配置信息
- 删除密钥文件：删除该主机关联的密钥文件
- 删除配置：从config和known_hosts文件中删除主机配置
- 修改用户：修改主机的用户名配置
- 修改端口：修改主机的端口配置
- 网络诊断：对主机进行 TCP 延迟测量与路由追踪（需要管理员/适当权限进行 ICMP 操作）
- 返回：返回主机选择菜单

### 模糊查找功能
在主机选择菜单中，第一行提供了模糊查找功能。选择"搜索主机 (模糊查找)"选项，
然后输入关键词即可搜索匹配的主机。

### 命令行直接连接
可以直接指定主机信息进行连接：
```bash
./sshgo root@192.168.1.100
./sshgo user@host:port
./sshgo hostname
```

## SSH配置文件

SSHGo会自动读取默认的SSH配置文件：
- Linux/macOS: `~/.ssh/config`
- Windows: `%USERPROFILE%\.ssh\config`

在Windows系统中，SSHGo还会尝试读取以下位置的配置文件：
- `C:\Program Files\Git\etc\ssh\ssh_config`
- `C:\Program Files (x86)\Git\etc\ssh\ssh_config`
- `%PROGRAMDATA%\ssh\ssh_config`

配置文件格式示例：
```
Host server1
    HostName 192.168.1.100
    User ubuntu
    Port 22
    IdentityFile ~/.ssh/server1_key

Host server2
    HostName example.com
    User root
    Port 2222
```

## 跨平台兼容性

SSHGo支持以下平台：
- Windows (需要安装SSH客户端)
- Linux (通常预装SSH客户端)
- macOS (通常预装SSH客户端)

## 许可证

MIT License（见 LICENSE 文件）。

## 语言 / 本地化

程序会自动检测系统语言（优先级：显式环境变量 > LANG/LC_ALL > 默认英文）。

手动指定：
```pwsh
$env:SSHGO_LANG = "zh"  # 或 "en"
./sshgo
```

## 网络诊断说明

网络诊断包含：
- TCP 延迟：尝试与目标主机 22 端口建立连接测量时间；失败会显示错误。
- 路由追踪：基于 ICMP Echo。部分系统需要管理员权限；若权限不足可能失败或无结果。

后续可扩展建议：
- 增加并发多次延迟统计（平均/最小/最大/抖动）。
- 支持自定义端口或协议（如 ICMP / UDP）。
- 导出结果为 JSON 便于脚本化集成。