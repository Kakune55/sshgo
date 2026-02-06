package i18n

// StringKey 定义字符串资源的键
type StringKey string

// 定义所有字符串资源的键
const (
	// 主菜单相关
	InvalidSSHCommand StringKey = "invalid_ssh_command"
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
	ModifyPortAction       StringKey = "modify_port_action"
	NetworkDiagnosticsAction StringKey = "network_diagnostics_action"
	BackAction             StringKey = "back_action"

	// 网络诊断相关
	NetworkDiagnosticsTitle StringKey = "network_diagnostics_title"
	MeasureLatencyAction    StringKey = "measure_latency_action"
	RouteTraceAction        StringKey = "route_trace_action"
	ReturnToMainMenu        StringKey = "return_to_main_menu"
	TestingLatency          StringKey = "testing_latency"
	TestingLatencyShort     StringKey = "testing_latency_short"
	LatencyResult           StringKey = "latency_result"
	LatencySummary          StringKey = "latency_summary"
	TracingRoute            StringKey = "tracing_route"
	TracingRouteShort       StringKey = "tracing_route_short"
	RouteTraceResults       StringKey = "route_trace_results"
	RouteHopInfo            StringKey = "route_hop_info"
	RouteHopTimeout         StringKey = "route_hop_timeout"
	RouteTraceFailed        StringKey = "route_trace_failed"
	PressEscToReturn        StringKey = "press_esc_to_return"

	// 输入提示相关
	EnterNewUsername StringKey = "enter_new_username"
	EnterNewPort     StringKey = "enter_new_port"
	DefaultUsername  StringKey = "default_username"

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

	// 按键帮助文本
	KeySelect  StringKey = "key_select"
	KeyBack    StringKey = "key_back"
	KeyQuit    StringKey = "key_quit"
	KeySearch  StringKey = "key_search"
	KeyConfirm StringKey = "key_confirm"
	KeyCancel  StringKey = "key_cancel"

	// UI 消息
	NoKeyFileConfigured   StringKey = "no_key_file_configured"
	PressAnyKeyToReturn   StringKey = "press_any_key_to_return"
	EnterUsernameForHost  StringKey = "enter_username_for_host"
	NoSSHHostsFound       StringKey = "no_ssh_hosts_found"
	ProgramError          StringKey = "program_error"
	ConfigDeletedReloading StringKey = "config_deleted_reloading"

	// SSH 连接错误
	UsernameNotSet StringKey = "username_not_set"
	Warning        StringKey = "warning"

	// SSH 配置错误
	ParseConfigError      StringKey = "parse_config_error"
	ReadKnownHostsWarning StringKey = "read_known_hosts_warning"
	ParseTestConfigError  StringKey = "parse_test_config_error"
	ReadConfigFileError   StringKey = "read_config_file_error"
	ReadKnownHostsError   StringKey = "read_known_hosts_error"
	ReadConfigFileFailed  StringKey = "read_config_file_failed"
	WriteConfigFileFailed StringKey = "write_config_file_failed"
	ReadKnownHostsFailed  StringKey = "read_known_hosts_failed"
	WriteKnownHostsFailed StringKey = "write_known_hosts_failed"

	// 操作警告
	DeleteHostConfigWarning  StringKey = "delete_host_config_warning"
	DeleteKnownHostsWarning  StringKey = "delete_known_hosts_warning"

	// 其他
	Goodbye      StringKey = "goodbye"
	ConnectingTo StringKey = "connecting_to"
)
