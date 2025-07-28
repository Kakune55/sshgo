package i18n

// 英文资源
var enStrings = map[StringKey]string{
	// 主菜单相关
	SelectHostLabel:  "Select a host to connect to",
	SearchHostOption: "Search Host",
	ExitOption:       "Exit",

	// 搜索相关
	EnterSearchKeyword:    "Please enter a search keyword",
	NoMatchingHosts:       "No matching hosts found",
	MultipleMatchingHosts: "Found multiple matching hosts, please select",

	// 操作菜单相关
	SelectActionLabel:  "Select an action",
	ConnectAction:      "Connect",
	DetailsAction:      "Details",
	DeleteKeyAction:    "Delete Key File",
	DeleteConfigAction: "Delete Configuration",
	ModifyUserAction:   "Modify User",
	ModifyPortAction:   "Modify Port",
	BackAction:         "Back",

	// 输入提示相关
	EnterNewUsername: "Please enter a new username",
	EnterNewPort:     "Please enter a new port number",
	DefaultUsername:  "root",

	// 取消操作相关
	CancelOperation: "Cancel operation",

	// 主机详情相关
	HostDetailsTitle: "=== Host Details ===",
	HostAlias:        "Host Alias: %s",
	HostName:         "Host Name: %s",
	UserName:         "Username: %s",
	Port:             "Port: %s",
	KeyFile:          "Key File: %s",

	// 确认提示相关
	ConfirmDeleteKey:    "Are you sure you want to delete the key file '%s'?",
	ConfirmDeleteConfig: "Are you sure you want to delete the configuration for host '%s'? This will remove records from config and known_hosts files.",

	// 错误消息相关
	FailedToGetUsername:  "Failed to get username: %v",
	FailedToGetPort:      "Failed to get port number: %v",
	FailedToDeleteKey:    "Failed to delete key file: %v",
	FailedToDeleteConfig: "Failed to delete configuration: %v",
	FailedToModifyUser:   "Failed to modify user: %v",
	InvalidSSHCommand:    "Invalid SSH command: %v",
	FailedToModifyPort:   "Failed to modify port: %v",

	// 成功消息相关
	SuccessfullyDeletedKey:    "Successfully deleted key file: %s",
	SuccessfullyDeletedConfig: "Successfully deleted configuration for host '%s'",
	SuccessfullyModifiedUser:  "Successfully modified username for host '%s' to '%s'",
	SuccessfullyModifiedPort:  "Successfully modified port for host '%s' to '%s'",

	// 其他
	Goodbye:      "Goodbye!",
	ConnectingTo: "Connecting to %s@%s...",
}
