package ui

import (
	"fmt"
	"strings"

	"sshgo/i18n"
	"sshgo/operations"
	"sshgo/ssh"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// 状态和类型定义
// ============================================================================

// appState 应用状态
type appState int

const (
	stateHostList appState = iota
	stateActionMenu
	stateHostDetails
	stateConfirmDeleteKey
	stateConfirmDeleteConfig
	stateInputUsername
	stateInputPort
	stateInputConnectUsername
)

// ActionType 操作类型（导出供外部使用）
type ActionType string

const (
	ActionConnect            ActionType = "connect"
	ActionDetails            ActionType = "details"
	ActionDeleteKey          ActionType = "delete_key"
	ActionDeleteConfig       ActionType = "delete_config"
	ActionModifyUser         ActionType = "modify_user"
	ActionModifyPort         ActionType = "modify_port"
	ActionNetworkDiagnostics ActionType = "network_diagnostics"
	ActionBack               ActionType = "back"
	ActionExit               ActionType = "exit"
	ActionNone               ActionType = ""
)

// ============================================================================
// 列表项类型
// ============================================================================

// hostItem 主机列表项
type hostItem struct {
	host        ssh.SSHHost
	displayName string
}

func (i hostItem) Title() string       { return i.displayName }
func (i hostItem) Description() string { return i.host.HostName }
func (i hostItem) FilterValue() string { return i.displayName + " " + i.host.HostName }

// actionItem 操作列表项
type actionItem struct {
	action ActionType
	label  string
}

func (i actionItem) Title() string       { return i.label }
func (i actionItem) Description() string { return "" }
func (i actionItem) FilterValue() string { return i.label }

// ============================================================================
// 键绑定
// ============================================================================

type keyMap struct {
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
	Search key.Binding
	Yes    key.Binding
	No     key.Binding
}

func getKeys() keyMap {
	return keyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", i18n.T(i18n.KeySelect)),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", i18n.T(i18n.KeyBack)),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", i18n.T(i18n.KeyQuit)),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", i18n.T(i18n.KeySearch)),
		),
		Yes: key.NewBinding(
			key.WithKeys("y", "Y"),
			key.WithHelp("y", i18n.T(i18n.KeyConfirm)),
		),
		No: key.NewBinding(
			key.WithKeys("n", "N"),
			key.WithHelp("n", i18n.T(i18n.KeyCancel)),
		),
	}
}

// ============================================================================
// 样式定义
// ============================================================================

var (
	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginLeft(2)

	// 状态信息样式
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)

	// 错误样式
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			MarginLeft(2)

	// 成功样式
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			MarginLeft(2)

	// 详情框样式
	detailsBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			MarginLeft(2).
			MarginTop(1)

	// 帮助样式
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2).
			MarginTop(1)

	// 警告样式
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true).
			MarginLeft(2)

	// 输入框样式
	inputStyle = lipgloss.NewStyle().
			MarginLeft(2)
)

// ============================================================================
// 主应用模型
// ============================================================================

// AppModel 主应用模型
type AppModel struct {
	// 状态
	state appState

	// 数据
	hosts        []ssh.SSHHost
	selectedHost ssh.SSHHost
	configPath   string

	// 列表组件
	hostList   list.Model
	actionList list.Model

	// 输入组件
	textInput textinput.Model

	// 消息显示
	message string
	isError bool

	// 窗口尺寸
	width  int
	height int

	// 退出标志
	quitting bool

	// 返回给调用者的信息
	resultAction ActionType
	resultHost   *ssh.SSHHost
}

// NewAppModel 创建新的应用模型
func NewAppModel(hosts []ssh.SSHHost, configPath string) AppModel {
	// 创建主机列表项
	hostItems := make([]list.Item, len(hosts))
	for i, h := range hosts {
		displayName := h.Host
		if h.HostName != "" && h.HostName != h.Host {
			displayName = fmt.Sprintf("%s (%s)", h.Host, h.HostName)
		}
		hostItems[i] = hostItem{host: h, displayName: displayName}
	}

	// 配置主机列表
	hostDelegate := list.NewDefaultDelegate()
	hostDelegate.ShowDescription = true
	hostList := list.New(hostItems, hostDelegate, 0, 0)
	hostList.Title = i18n.T(i18n.SelectHostLabel)
	hostList.SetShowStatusBar(true)
	hostList.SetFilteringEnabled(true)
	hostList.SetShowHelp(true)
	hostList.DisableQuitKeybindings()

	// 创建操作列表项
	actionItems := []list.Item{
		actionItem{action: ActionConnect, label: i18n.T(i18n.ConnectAction)},
		actionItem{action: ActionDetails, label: i18n.T(i18n.DetailsAction)},
		actionItem{action: ActionDeleteKey, label: i18n.T(i18n.DeleteKeyAction)},
		actionItem{action: ActionDeleteConfig, label: i18n.T(i18n.DeleteConfigAction)},
		actionItem{action: ActionModifyUser, label: i18n.T(i18n.ModifyUserAction)},
		actionItem{action: ActionModifyPort, label: i18n.T(i18n.ModifyPortAction)},
		actionItem{action: ActionNetworkDiagnostics, label: i18n.T(i18n.NetworkDiagnosticsAction)},
		actionItem{action: ActionBack, label: i18n.T(i18n.BackAction)},
	}

	// 配置操作列表
	actionDelegate := list.NewDefaultDelegate()
	actionDelegate.ShowDescription = false
	actionList := list.New(actionItems, actionDelegate, 0, 0)
	actionList.Title = i18n.T(i18n.SelectActionLabel)
	actionList.SetShowStatusBar(false)
	actionList.SetFilteringEnabled(false)
	actionList.SetShowHelp(true)
	actionList.DisableQuitKeybindings()

	// 创建文本输入
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 40

	return AppModel{
		state:      stateHostList,
		hosts:      hosts,
		configPath: configPath,
		hostList:   hostList,
		actionList: actionList,
		textInput:  ti,
		width:      80,
		height:     24,
	}
}

// Init 初始化
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update 更新
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 处理通用消息
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h := max(msg.Height, 5)
		m.hostList.SetSize(msg.Width-4, h)
		m.actionList.SetSize(msg.Width-4, h)
		return m, nil

	case tea.KeyMsg:
		// 全局退出（仅在非输入状态）
		if msg.String() == "ctrl+c" {
			m.quitting = true
			m.resultAction = ActionExit
			return m, tea.Quit
		}
	}

	// 根据状态分发处理
	switch m.state {
	case stateHostList:
		return m.updateHostList(msg)
	case stateActionMenu:
		return m.updateActionMenu(msg)
	case stateHostDetails:
		return m.updateHostDetails(msg)
	case stateConfirmDeleteKey:
		return m.updateConfirmDeleteKey(msg)
	case stateConfirmDeleteConfig:
		return m.updateConfirmDeleteConfig(msg)
	case stateInputUsername:
		return m.updateInputUsername(msg)
	case stateInputPort:
		return m.updateInputPort(msg)
	case stateInputConnectUsername:
		return m.updateInputConnectUsername(msg)
	}

	return m, nil
}

// ============================================================================
// 状态更新函数
// ============================================================================

// updateHostList 更新主机列表状态
func (m AppModel) updateHostList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			// 如果列表正在过滤中，不退出
			if m.hostList.FilterState() == list.Filtering {
				break
			}
			m.quitting = true
			m.resultAction = ActionExit
			return m, tea.Quit
		case "enter":
			// 如果列表正在过滤中，让列表处理
			if m.hostList.FilterState() == list.Filtering {
				break
			}
			if item, ok := m.hostList.SelectedItem().(hostItem); ok {
				m.selectedHost = item.host
				m.state = stateActionMenu
				m.message = ""
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.hostList, cmd = m.hostList.Update(msg)
	return m, cmd
}

// updateActionMenu 更新操作菜单状态
func (m AppModel) updateActionMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			m.resultAction = ActionExit
			return m, tea.Quit
		case "esc":
			m.state = stateHostList
			m.message = ""
			return m, nil
		case "enter":
			if item, ok := m.actionList.SelectedItem().(actionItem); ok {
				return m.handleAction(item.action)
			}
		}
	}

	var cmd tea.Cmd
	m.actionList, cmd = m.actionList.Update(msg)
	return m, cmd
}

// updateHostDetails 更新主机详情状态
func (m AppModel) updateHostDetails(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		m.state = stateActionMenu
		return m, nil
	}
	return m, nil
}

// updateConfirmDeleteKey 更新确认删除密钥状态
func (m AppModel) updateConfirmDeleteKey(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// 执行删除
			err := operations.DeleteKeyFile(m.selectedHost)
			if err != nil {
				m.message = err.Error()
				m.isError = true
			} else {
				m.message = fmt.Sprintf(i18n.T(i18n.SuccessfullyDeletedKey), m.selectedHost.KeyFile)
				m.isError = false
			}
			m.state = stateActionMenu
			return m, nil
		case "n", "N", "esc":
			m.message = i18n.T(i18n.CancelOperation)
			m.isError = false
			m.state = stateActionMenu
			return m, nil
		}
	}
	return m, nil
}

// updateConfirmDeleteConfig 更新确认删除配置状态
func (m AppModel) updateConfirmDeleteConfig(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// 执行删除
			err := operations.DeleteHostConfig(m.selectedHost)
			if err != nil {
				m.message = err.Error()
				m.isError = true
			} else {
				m.message = fmt.Sprintf(i18n.T(i18n.SuccessfullyDeletedConfig), m.selectedHost.Host)
				m.isError = false
				// 需要刷新主机列表
				m.resultAction = ActionDeleteConfig
				m.resultHost = &m.selectedHost
				m.quitting = true
				return m, tea.Quit
			}
			m.state = stateActionMenu
			return m, nil
		case "n", "N", "esc":
			m.message = i18n.T(i18n.CancelOperation)
			m.isError = false
			m.state = stateActionMenu
			return m, nil
		}
	}
	return m, nil
}

// updateInputUsername 更新用户名输入状态
func (m AppModel) updateInputUsername(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			username := m.textInput.Value()
			if username == "" {
				username = i18n.T(i18n.DefaultUsername)
			}
			err := operations.ModifyUser(m.selectedHost, username)
			if err != nil {
				m.message = err.Error()
				m.isError = true
			} else {
				m.message = fmt.Sprintf(i18n.T(i18n.SuccessfullyModifiedUser), m.selectedHost.Host, username)
				m.isError = false
			}
			m.state = stateActionMenu
			return m, nil
		case "esc":
			m.message = i18n.T(i18n.CancelOperation)
			m.isError = false
			m.state = stateActionMenu
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// updateInputPort 更新端口输入状态
func (m AppModel) updateInputPort(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			port := m.textInput.Value()
			if port == "" {
				port = "22"
			}
			err := operations.ModifyPort(m.selectedHost, port)
			if err != nil {
				m.message = err.Error()
				m.isError = true
			} else {
				m.message = fmt.Sprintf(i18n.T(i18n.SuccessfullyModifiedPort), m.selectedHost.Host, port)
				m.isError = false
			}
			m.state = stateActionMenu
			return m, nil
		case "esc":
			m.message = i18n.T(i18n.CancelOperation)
			m.isError = false
			m.state = stateActionMenu
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// updateInputConnectUsername 更新连接用户名输入状态
func (m AppModel) updateInputConnectUsername(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			username := m.textInput.Value()
			if username == "" {
				username = i18n.T(i18n.DefaultUsername)
			}
			// 设置用户名并准备连接
			m.selectedHost.User = username
			// 保存用户名
			_ = ssh.SaveUserToConfig(m.selectedHost.Host, username)
			// 返回连接操作
			m.resultAction = ActionConnect
			m.resultHost = &m.selectedHost
			m.quitting = true
			return m, tea.Quit
		case "esc":
			m.message = i18n.T(i18n.CancelOperation)
			m.isError = false
			m.state = stateActionMenu
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// ============================================================================
// 操作处理
// ============================================================================

// handleAction 处理操作
func (m AppModel) handleAction(action ActionType) (tea.Model, tea.Cmd) {
	m.message = ""
	m.isError = false

	switch action {
	case ActionConnect:
		// 检查是否需要用户名
		if ssh.NeedsUsername(m.selectedHost) {
			m.textInput.SetValue("")
			m.textInput.Placeholder = i18n.T(i18n.DefaultUsername)
			m.textInput.Focus()
			m.state = stateInputConnectUsername
			return m, textinput.Blink
		}
		// 直接连接
		m.resultAction = ActionConnect
		m.resultHost = &m.selectedHost
		m.quitting = true
		return m, tea.Quit

	case ActionDetails:
		m.state = stateHostDetails
		return m, nil

	case ActionDeleteKey:
		if m.selectedHost.KeyFile == "" {
			m.message = i18n.T(i18n.NoKeyFileConfigured)
			m.isError = true
			return m, nil
		}
		m.state = stateConfirmDeleteKey
		return m, nil

	case ActionDeleteConfig:
		m.state = stateConfirmDeleteConfig
		return m, nil

	case ActionModifyUser:
		m.textInput.SetValue(m.selectedHost.User)
		m.textInput.Placeholder = i18n.T(i18n.DefaultUsername)
		m.textInput.Focus()
		m.state = stateInputUsername
		return m, textinput.Blink

	case ActionModifyPort:
		port := m.selectedHost.Port
		if port == "" {
			port = "22"
		}
		m.textInput.SetValue(port)
		m.textInput.Placeholder = "22"
		m.textInput.Focus()
		m.state = stateInputPort
		return m, textinput.Blink

	case ActionNetworkDiagnostics:
		m.resultAction = ActionNetworkDiagnostics
		m.resultHost = &m.selectedHost
		m.quitting = true
		return m, tea.Quit

	case ActionBack:
		m.state = stateHostList
		return m, nil
	}

	return m, nil
}

// ============================================================================
// 视图渲染
// ============================================================================

// View 渲染
func (m AppModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	switch m.state {
	case stateHostList:
		s.WriteString(m.hostList.View())

	case stateActionMenu:
		m.actionList.Title = fmt.Sprintf("%s: %s", i18n.T(i18n.SelectActionLabel), m.selectedHost.Host)
		s.WriteString(m.actionList.View())

	case stateHostDetails:
		s.WriteString(m.renderHostDetails())

	case stateConfirmDeleteKey:
		s.WriteString(m.renderConfirmDeleteKey())

	case stateConfirmDeleteConfig:
		s.WriteString(m.renderConfirmDeleteConfig())

	case stateInputUsername, stateInputConnectUsername:
		s.WriteString(m.renderInputUsername())

	case stateInputPort:
		s.WriteString(m.renderInputPort())
	}

	// 显示消息
	if m.message != "" {
		s.WriteString("\n\n")
		if m.isError {
			s.WriteString(errorStyle.Render("✗ " + m.message))
		} else {
			s.WriteString(successStyle.Render("✓ " + m.message))
		}
	}

	return s.String()
}

// renderHostDetails 渲染主机详情
func (m AppModel) renderHostDetails() string {
	var details strings.Builder

	details.WriteString(fmt.Sprintf(i18n.T(i18n.HostAlias), m.selectedHost.Host))
	if m.selectedHost.HostName != "" {
		details.WriteString("\n")
		details.WriteString(fmt.Sprintf(i18n.T(i18n.HostName), m.selectedHost.HostName))
	}
	if m.selectedHost.User != "" {
		details.WriteString("\n")
		details.WriteString(fmt.Sprintf(i18n.T(i18n.UserName), m.selectedHost.User))
	}
	port := m.selectedHost.Port
	if port == "" {
		port = "22"
	}
	details.WriteString("\n")
	fmt.Fprintf(&details, i18n.T(i18n.Port), port)
	if m.selectedHost.KeyFile != "" {
		details.WriteString("\n")
		details.WriteString(fmt.Sprintf(i18n.T(i18n.KeyFile), m.selectedHost.KeyFile))
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render(i18n.T(i18n.HostDetailsTitle)))
	s.WriteString("\n")
	s.WriteString(detailsBoxStyle.Render(details.String()))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render(i18n.T(i18n.PressAnyKeyToReturn)))

	return s.String()
}

// renderConfirmDeleteKey 渲染确认删除密钥
func (m AppModel) renderConfirmDeleteKey() string {
	var s strings.Builder

	s.WriteString(warningStyle.Render("⚠ " + i18n.T(i18n.KeyConfirm)))
	s.WriteString("\n\n")
	s.WriteString(statusStyle.Render(fmt.Sprintf(i18n.T(i18n.ConfirmDeleteKey), m.selectedHost.KeyFile)))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("y: " + i18n.T(i18n.KeyConfirm) + " • n/esc: " + i18n.T(i18n.KeyCancel)))

	return s.String()
}

// renderConfirmDeleteConfig 渲染确认删除配置
func (m AppModel) renderConfirmDeleteConfig() string {
	var s strings.Builder

	s.WriteString(warningStyle.Render("⚠ " + i18n.T(i18n.KeyConfirm)))
	s.WriteString("\n\n")
	s.WriteString(statusStyle.Render(fmt.Sprintf(i18n.T(i18n.ConfirmDeleteConfig), m.selectedHost.Host)))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("y: " + i18n.T(i18n.KeyConfirm) + " • n/esc: " + i18n.T(i18n.KeyCancel)))

	return s.String()
}

// renderInputUsername 渲染用户名输入
func (m AppModel) renderInputUsername() string {
	var s strings.Builder

	title := i18n.T(i18n.EnterNewUsername)
	if m.state == stateInputConnectUsername {
		title = fmt.Sprintf(i18n.T(i18n.EnterUsernameForHost), m.selectedHost.Host)
	}

	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")
	s.WriteString(inputStyle.Render(m.textInput.View()))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("enter: " + i18n.T(i18n.KeyConfirm) + " • esc: " + i18n.T(i18n.KeyCancel)))

	return s.String()
}

// renderInputPort 渲染端口输入
func (m AppModel) renderInputPort() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(i18n.T(i18n.EnterNewPort)))
	s.WriteString("\n\n")
	s.WriteString(inputStyle.Render(m.textInput.View()))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("enter: " + i18n.T(i18n.KeyConfirm) + " • esc: " + i18n.T(i18n.KeyCancel)))

	return s.String()
}

// ============================================================================
// 导出方法
// ============================================================================

// GetResultAction 获取结果操作
func (m AppModel) GetResultAction() ActionType {
	return m.resultAction
}

// GetResultHost 获取结果主机
func (m AppModel) GetResultHost() *ssh.SSHHost {
	return m.resultHost
}

// IsQuitting 是否正在退出
func (m AppModel) IsQuitting() bool {
	return m.quitting
}

// ============================================================================
// 入口函数
// ============================================================================

// Run 运行主应用
func Run(hosts []ssh.SSHHost, configPath string) (ActionType, *ssh.SSHHost, error) {
	model := NewAppModel(hosts, configPath)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return ActionNone, nil, err
	}

	m, ok := finalModel.(AppModel)
	if !ok {
		return ActionNone, nil, fmt.Errorf("unexpected model type")
	}

	return m.GetResultAction(), m.GetResultHost(), nil
}

// RunLoop 运行主循环
func RunLoop() {
	// 获取SSH配置文件路径
	configPath := ssh.GetSSHConfigPath()

	// 解析SSH配置文件
	hosts, err := ssh.ParseSSHConfig(configPath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if len(hosts) == 0 {
		fmt.Println(i18n.T(i18n.NoSSHHostsFound))
		return
	}

	for {
		action, host, err := Run(hosts, configPath)
		if err != nil {
			fmt.Printf(i18n.T(i18n.ProgramError)+"\n", err)
			return
		}

		switch action {
		case ActionExit, ActionNone:
			fmt.Println(i18n.T(i18n.Goodbye))
			return

		case ActionConnect:
			if host != nil {
				err = ssh.ConnectToHost(*host)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}

		case ActionNetworkDiagnostics:
			if host != nil {
				err = RunNetworkDiagnostics(*host)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}

		case ActionDeleteConfig:
			// 重新加载主机列表
			fmt.Println(i18n.T(i18n.ConfigDeletedReloading))
			hosts, err = ssh.ParseSSHConfig(configPath)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			if len(hosts) == 0 {
				fmt.Println(i18n.T(i18n.NoSSHHostsFound))
				return
			}
		}
	}
}