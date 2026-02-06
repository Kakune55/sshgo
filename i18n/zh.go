package i18n

// 中文资源
var zhStrings = map[StringKey]string{
	// 主菜单相关
	SelectHostLabel:  "选择要连接的主机",
	SearchHostOption: "搜索主机",
	ExitOption:       "退出",

	// 搜索相关
	EnterSearchKeyword:    "请输入搜索关键词",
	NoMatchingHosts:       "未找到匹配的主机",
	MultipleMatchingHosts: "找到多个匹配的主机，请选择",

	// 操作菜单相关
	SelectActionLabel:  "选择操作",
	ConnectAction:      "连接",
	DetailsAction:      "详细信息",
	DeleteKeyAction:    "删除密钥文件",
	DeleteConfigAction: "删除配置",
	ModifyUserAction:   "修改用户",
	ModifyPortAction:       "修改端口",
	NetworkDiagnosticsAction: "网络诊断",
	BackAction:             "返回",

	// 网络诊断相关
	NetworkDiagnosticsTitle: "网络诊断 - %s",
	MeasureLatencyAction:    "测量延迟",
	RouteTraceAction:        "路由追踪",
	ReturnToMainMenu:        "返回主菜单",
	TestingLatency:          "正在测试到 %s 的延迟...",
	TestingLatencyShort:     "正在测试...",
	LatencyResult:           "延迟: %v",
	LatencySummary:          "样本 %d 个 | 最小: %v 最大: %v 平均: %v",
	TracingRoute:            "正在追踪到 %s 的路由...",
	TracingRouteShort:       "正在追踪路由...",
	RouteTraceResults:       "路由追踪结果:",
	RouteHopInfo:            "跳点 %d: %s (RTT: %v)",
	RouteHopTimeout:         "跳点 %d: * (RTT: %v)",
	RouteTraceFailed:        "路由追踪失败: %v",
	PressEscToReturn:        "esc: 返回",

	// 输入提示相关
	EnterNewUsername: "请输入新用户名",
	EnterNewPort:     "请输入新端口号",
	DefaultUsername:  "root",

	// 取消操作相关
	CancelOperation: "取消删除操作",

	// 主机详情相关
	HostDetailsTitle: "=== 主机详细信息 ===",
	HostAlias:        "主机别名: %s",
	HostName:         "主机地址: %s",
	UserName:         "用户名: %s",
	Port:             "端口: %s",
	KeyFile:          "密钥文件: %s",

	// 确认提示相关
	ConfirmDeleteKey:    "确定要删除密钥文件 '%s' 吗?",
	ConfirmDeleteConfig: "确定要删除主机 '%s' 的配置吗? 这将从config和known_hosts文件中移除相关记录",

	// 错误消息相关
	FailedToGetUsername:  "获取用户名失败: %v",
	FailedToGetPort:      "获取端口号失败: %v",
	FailedToDeleteKey:    "删除密钥文件失败: %v",
	FailedToDeleteConfig: "删除配置时出错: %v",
	FailedToModifyUser:   "修改用户时出错: %v",
	InvalidSSHCommand:    "无效的SSH命令: %v",
	FailedToModifyPort:   "修改端口时出错: %v",

	// 成功消息相关
	SuccessfullyDeletedKey:    "成功删除密钥文件: %s",
	SuccessfullyDeletedConfig: "成功删除主机 '%s' 的配置",
	SuccessfullyModifiedUser:  "成功将主机 '%s' 的用户名修改为 '%s'",
	SuccessfullyModifiedPort:  "成功将主机 '%s' 的端口修改为 '%s'",

	// 按键帮助文本
	KeySelect:  "选择",
	KeyBack:    "返回",
	KeyQuit:    "退出",
	KeySearch:  "搜索",
	KeyConfirm: "确认",
	KeyCancel:  "取消",

	// UI 消息
	NoKeyFileConfigured:    "该主机没有配置密钥文件",
	PressAnyKeyToReturn:    "按任意键返回",
	EnterUsernameForHost:   "输入用户名连接到 %s",
	NoSSHHostsFound:        "未找到SSH主机配置",
	ProgramError:           "程序运行错误: %v",
	ConfigDeletedReloading: "配置已删除，重新加载主机列表...",

	// SSH 连接错误
	UsernameNotSet: "用户名未设置",
	Warning:        "警告: %v",

	// SSH 配置错误
	ParseConfigError:      "解析配置文件 %s 时出错: %v",
	ReadKnownHostsWarning: "警告: 读取known_hosts文件时出错: %v",
	ParseTestConfigError:  "解析测试配置文件时出错: %v",
	ReadConfigFileError:   "读取配置文件时出错: %v",
	ReadKnownHostsError:   "读取known_hosts文件时出错: %v",
	ReadConfigFileFailed:  "读取配置文件失败: %v",
	WriteConfigFileFailed: "写入配置文件失败: %v",
	ReadKnownHostsFailed:  "读取known_hosts文件失败: %v",
	WriteKnownHostsFailed: "写入known_hosts文件失败: %v",

	// 操作警告
	DeleteHostConfigWarning: "警告: 从配置文件中删除主机配置时出错: %v",
	DeleteKnownHostsWarning: "警告: 从known_hosts文件中删除主机记录时出错: %v",

	// 其他
	Goodbye:      "再见!",
	ConnectingTo: "正在连接到 %s@%s...",
}
