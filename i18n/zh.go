package i18n

// 中文资源
var zhStrings = map[StringKey]string{
	// 主菜单相关
	SelectHostLabel:      "选择要连接的主机",
	SearchHostOption:     "搜索主机",
	ExitOption:           "退出",
	
	// 搜索相关
	EnterSearchKeyword:   "请输入搜索关键词",
	NoMatchingHosts:      "未找到匹配的主机",
	MultipleMatchingHosts: "找到多个匹配的主机，请选择",
	
	// 操作菜单相关
	SelectActionLabel:    "选择操作",
	ConnectAction:        "连接",
	DetailsAction:        "详细信息",
	DeleteKeyAction:      "删除密钥文件",
	DeleteConfigAction:   "删除配置",
	ModifyUserAction:     "修改用户",
	ModifyPortAction:     "修改端口",
	BackAction:           "返回",
	
	// 输入提示相关
	EnterNewUsername:     "请输入新用户名",
	EnterNewPort:         "请输入新端口号",
	
	// 取消操作相关
	CancelOperation:      "取消删除操作",
	
	// 主机详情相关
	HostDetailsTitle:     "=== 主机详细信息 ===",
	HostAlias:            "主机别名: %s",
	HostName:             "主机地址: %s",
	UserName:             "用户名: %s",
	Port:                 "端口: %s",
	KeyFile:              "密钥文件: %s",
	
	// 确认提示相关
	ConfirmDeleteKey:     "确定要删除密钥文件 '%s' 吗?",
	ConfirmDeleteConfig:  "确定要删除主机 '%s' 的配置吗? 这将从config和known_hosts文件中移除相关记录",
	
	// 错误消息相关
	FailedToGetUsername:  "获取用户名失败: %v",
	FailedToGetPort:      "获取端口号失败: %v",
	FailedToDeleteKey:    "删除密钥文件失败: %v",
	FailedToDeleteConfig: "删除配置时出错: %v",
	FailedToModifyUser:   "修改用户时出错: %v",
	FailedToModifyPort:   "修改端口时出错: %v",
	
	// 成功消息相关
	SuccessfullyDeletedKey:    "成功删除密钥文件: %s",
	SuccessfullyDeletedConfig: "成功删除主机 '%s' 的配置",
	SuccessfullyModifiedUser:  "成功将主机 '%s' 的用户名修改为 '%s'",
	SuccessfullyModifiedPort:  "成功将主机 '%s' 的端口修改为 '%s'",
	
	// 其他
	Goodbye:              "再见!",
	ConnectingTo:         "正在连接到 %s@%s...",
}