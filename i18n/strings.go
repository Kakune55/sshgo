package i18n

// StringKey 定义字符串资源的键
type StringKey string

// 定义所有字符串资源的键
const (
	// 主菜单相关
	SelectHostLabel  StringKey = "select_host_label"
	SearchHostOption StringKey = "search_host_option"
	ExitOption       StringKey = "exit_option"

	// 搜索相关
	EnterSearchKeyword    StringKey = "enter_search_keyword"
	NoMatchingHosts       StringKey = "no_matching_hosts"
	MultipleMatchingHosts StringKey = "multiple_matching_hosts"

	// 操作菜单相关
	SelectActionLabel  StringKey = "select_action_label"
	ConnectAction      StringKey = "connect_action"
	DetailsAction      StringKey = "details_action"
	DeleteKeyAction    StringKey = "delete_key_action"
	DeleteConfigAction StringKey = "delete_config_action"
	ModifyUserAction   StringKey = "modify_user_action"
	ModifyPortAction   StringKey = "modify_port_action"
	BackAction         StringKey = "back_action"

	// 输入提示相关
	EnterNewUsername StringKey = "enter_new_username"
	EnterNewPort     StringKey = "enter_new_port"

	// 取消操作相关
	CancelOperation StringKey = "cancel_operation"

	// 主机详情相关
	HostDetailsTitle StringKey = "host_details_title"
	HostAlias        StringKey = "host_alias"
	HostName         StringKey = "host_name"
	UserName         StringKey = "user_name"
	Port             StringKey = "port"
	KeyFile          StringKey = "key_file"

	// 确认提示相关
	ConfirmDeleteKey    StringKey = "confirm_delete_key"
	ConfirmDeleteConfig StringKey = "confirm_delete_config"

	// 错误消息相关
	FailedToGetUsername  StringKey = "failed_to_get_username"
	FailedToGetPort      StringKey = "failed_to_get_port"
	FailedToDeleteKey    StringKey = "failed_to_delete_key"
	FailedToDeleteConfig StringKey = "failed_to_delete_config"
	FailedToModifyUser   StringKey = "failed_to_modify_user"
	FailedToModifyPort   StringKey = "failed_to_modify_port"

	// 成功消息相关
	SuccessfullyDeletedKey    StringKey = "successfully_deleted_key"
	SuccessfullyDeletedConfig StringKey = "successfully_deleted_config"
	SuccessfullyModifiedUser  StringKey = "successfully_modified_user"
	SuccessfullyModifiedPort  StringKey = "successfully_modified_port"

	// 其他
	Goodbye      StringKey = "goodbye"
	ConnectingTo StringKey = "connecting_to"
)
