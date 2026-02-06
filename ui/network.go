package ui

import (
	"fmt"
	"strings"
	"time"

	"sshgo/i18n"
	"sshgo/network"
	"sshgo/ssh"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 网络诊断状态
type networkState int

const (
	networkStateMenu networkState = iota
	networkStateLatencyTest
	networkStateRouteTrace
)

// 网络诊断菜单项
type networkMenuItem struct {
	id    string
	label string
}

func (i networkMenuItem) Title() string       { return i.label }
func (i networkMenuItem) Description() string { return "" }
func (i networkMenuItem) FilterValue() string { return i.label }

// 延迟测试结果消息
type latencyResultMsg struct {
	sample   int
	duration time.Duration
	err      error
}

// 延迟测试完成消息
type latencyDoneMsg struct {
	count int
	min   time.Duration
	max   time.Duration
	avg   time.Duration
}

// 路由跳点消息
type routeHopMsg struct {
	hop       network.RouteHop
	isTimeout bool
}

// 路由追踪完成消息
type routeTraceDoneMsg struct{}

// 路由追踪结果消息（包含所有跳点）
type routeTraceResultMsg struct {
	hops []network.RouteHop
}

// 路由跳点接收消息
type routeHopReceivedMsg struct {
	hop      network.RouteHop
	hopChan  <-chan network.RouteHop
	errChan  <-chan error
	doneChan <-chan bool
}

// 路由追踪错误消息
type routeTraceErrorMsg struct {
	err error
}

// NetworkModel 网络诊断模型
type NetworkModel struct {
	state    networkState
	host     ssh.SSHHost
	menu     list.Model
	spinner  spinner.Model
	quitting bool

	// 延迟测试
	latencyResults []string
	latencySummary string

	// 路由追踪
	routeHops []string

	// 窗口尺寸
	width  int
	height int

	// 错误信息
	errorMsg string
}

// NewNetworkModel 创建网络诊断模型
func NewNetworkModel(host ssh.SSHHost) NetworkModel {
	// 创建菜单
	items := []list.Item{
		networkMenuItem{id: "latency", label: i18n.T(i18n.MeasureLatencyAction)},
		networkMenuItem{id: "trace", label: i18n.T(i18n.RouteTraceAction)},
		networkMenuItem{id: "back", label: i18n.T(i18n.ReturnToMainMenu)},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	menu := list.New(items, delegate, 0, 0)
	menu.Title = fmt.Sprintf(i18n.T(i18n.NetworkDiagnosticsTitle), host.Host)
	menu.SetShowStatusBar(false)
	menu.SetFilteringEnabled(false)
	menu.SetShowHelp(true)

	// 创建 spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return NetworkModel{
		state:   networkStateMenu,
		host:    host,
		menu:    menu,
		spinner: s,
		width:   80,
		height:  24,
	}
}

// Init 初始化
func (m NetworkModel) Init() tea.Cmd {
	return nil
}

// Update 更新
func (m NetworkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "q" || msg.String() == "ctrl+c":
			if m.state == networkStateMenu {
				m.quitting = true
				return m, tea.Quit
			}
		case msg.String() == "esc" || msg.String() == "backspace":
			if m.state != networkStateMenu {
				m.state = networkStateMenu
				m.latencyResults = nil
				m.latencySummary = ""
				m.routeHops = nil
				m.errorMsg = ""
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h := msg.Height - 6
		if h < 5 {
			h = 5
		}
		m.menu.SetSize(msg.Width-4, h)
		return m, nil

	case latencyResultMsg:
		if msg.err != nil {
			m.latencyResults = append(m.latencyResults, fmt.Sprintf("  样本 %d: 错误 - %v", msg.sample, msg.err))
		} else {
			m.latencyResults = append(m.latencyResults, fmt.Sprintf("  样本 %d: %v", msg.sample, msg.duration))
		}
		return m, nil

	case latencyDoneMsg:
		m.latencySummary = fmt.Sprintf(i18n.T(i18n.LatencySummary), msg.count, msg.min, msg.max, msg.avg)
		return m, nil

	case routeHopMsg:
		if msg.isTimeout {
			m.routeHops = append(m.routeHops, fmt.Sprintf(i18n.T(i18n.RouteHopTimeout), msg.hop.Index, msg.hop.RTT))
		} else {
			m.routeHops = append(m.routeHops, fmt.Sprintf(i18n.T(i18n.RouteHopInfo), msg.hop.Index, msg.hop.IP, msg.hop.RTT))
		}
		return m, m.spinner.Tick

	case routeTraceResultMsg:
		// 处理所有跳点结果
		for _, hop := range msg.hops {
			if hop.IP == nil {
				m.routeHops = append(m.routeHops, fmt.Sprintf(i18n.T(i18n.RouteHopTimeout), hop.Index, hop.RTT))
			} else {
				m.routeHops = append(m.routeHops, fmt.Sprintf(i18n.T(i18n.RouteHopInfo), hop.Index, hop.IP, hop.RTT))
			}
		}
		return m, nil

	case routeHopReceivedMsg:
		// 实时处理单个跳点
		if msg.hop.IP == nil {
			m.routeHops = append(m.routeHops, fmt.Sprintf(i18n.T(i18n.RouteHopTimeout), msg.hop.Index, msg.hop.RTT))
		} else {
			m.routeHops = append(m.routeHops, fmt.Sprintf(i18n.T(i18n.RouteHopInfo), msg.hop.Index, msg.hop.IP, msg.hop.RTT))
		}
		// 继续等待下一个跳点
		return m, waitForRouteHop(msg.hopChan, msg.errChan, msg.doneChan)

	case routeTraceDoneMsg:
		return m, nil

	case routeTraceErrorMsg:
		m.errorMsg = fmt.Sprintf(i18n.T(i18n.RouteTraceFailed), msg.err)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// 根据状态处理
	switch m.state {
	case networkStateMenu:
		return m.updateMenu(msg)
	case networkStateLatencyTest:
		return m.updateLatencyTest(msg)
	case networkStateRouteTrace:
		return m.updateRouteTrace(msg)
	}

	return m, nil
}

// updateMenu 更新菜单状态
func (m NetworkModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if item, ok := m.menu.SelectedItem().(networkMenuItem); ok {
				switch item.id {
				case "latency":
					m.state = networkStateLatencyTest
					m.latencyResults = nil
					m.latencySummary = ""
					return m, tea.Batch(m.spinner.Tick, m.runLatencyTest())
				case "trace":
					m.state = networkStateRouteTrace
					m.routeHops = nil
					m.errorMsg = ""
					return m, tea.Batch(m.spinner.Tick, m.runRouteTrace())
				case "back":
					m.quitting = true
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

// updateLatencyTest 更新延迟测试状态
func (m NetworkModel) updateLatencyTest(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// updateRouteTrace 更新路由追踪状态
func (m NetworkModel) updateRouteTrace(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// runLatencyTest 运行延迟测试
func (m NetworkModel) runLatencyTest() tea.Cmd {
	host := m.host.HostName
	if host == "" {
		host = m.host.Host
	}
	return func() tea.Msg {
		measurer := network.NewLatencyMeasurer()
		const samples = 5
		var sum time.Duration
		var min, max time.Duration
		var count int

		for range samples {
			result, err := measurer.MeasureLatency(host, "tcp")
			if err != nil {
				continue
			}

			if result.Error != nil {
				continue
			}

			d := result.Duration
			if count == 0 || d < min {
				min = d
			}
			if count == 0 || d > max {
				max = d
			}
			sum += d
			count++
			time.Sleep(200 * time.Millisecond)
		}

		if count > 0 {
			avg := sum / time.Duration(count)
			return latencyDoneMsg{count: count, min: min, max: max, avg: avg}
		}
		return latencyDoneMsg{}
	}
}

// runRouteTrace 运行路由追踪
func (m NetworkModel) runRouteTrace() tea.Cmd {
	host := m.host.HostName
	if host == "" {
		host = m.host.Host
	}
	
	// 创建 channel 用于接收跳点
	hopChan := make(chan network.RouteHop, 1)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)
	
	// 启动后台 goroutine 执行路由追踪
	go func() {
		tracer := network.NewRouteTracer()
		callback := func(hop network.RouteHop, isTimeout bool) {
			hopChan <- hop
		}
		_, err := tracer.TraceRouteWithCallback(host, callback)
		if err != nil {
			errChan <- err
		}
		doneChan <- true
	}()
	
	// 返回一个读取第一个跳点的命令
	return waitForRouteHop(hopChan, errChan, doneChan)
}

// waitForRouteHop 等待路由跳点消息
func waitForRouteHop(hopChan <-chan network.RouteHop, errChan <-chan error, doneChan <-chan bool) tea.Cmd {
	return func() tea.Msg {
		select {
		case hop := <-hopChan:
			return routeHopReceivedMsg{
				hop:      hop,
				hopChan:  hopChan,
				errChan:  errChan,
				doneChan: doneChan,
			}
		case err := <-errChan:
			return routeTraceErrorMsg{err: err}
		case <-doneChan:
			return routeTraceDoneMsg{}
		}
	}
}

// View 渲染
func (m NetworkModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	switch m.state {
	case networkStateMenu:
		s.WriteString(m.menu.View())

	case networkStateLatencyTest:
		s.WriteString(titleStyle.Render(fmt.Sprintf(i18n.T(i18n.TestingLatency), m.host.Host)))
		s.WriteString("\n\n")

		if len(m.latencyResults) == 0 && m.latencySummary == "" {
			s.WriteString("  " + m.spinner.View() + " " + i18n.T(i18n.TestingLatencyShort) + "\n")
		}

		for _, result := range m.latencyResults {
			s.WriteString(result + "\n")
		}

		if m.latencySummary != "" {
			s.WriteString("\n")
			s.WriteString(successStyle.Render(m.latencySummary))
		}

		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render(i18n.T(i18n.PressEscToReturn)))

	case networkStateRouteTrace:
		s.WriteString(titleStyle.Render(fmt.Sprintf(i18n.T(i18n.TracingRoute), m.host.Host)))
		s.WriteString("\n\n")

		if len(m.routeHops) == 0 && m.errorMsg == "" {
			s.WriteString("  " + m.spinner.View() + " " + i18n.T(i18n.TracingRouteShort) + "\n")
		}

		for _, hop := range m.routeHops {
			s.WriteString("  " + hop + "\n")
		}

		if m.errorMsg != "" {
			s.WriteString("\n")
			s.WriteString(errorStyle.Render(m.errorMsg))
		}

		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render(i18n.T(i18n.PressEscToReturn)))
	}

	return s.String()
}

// IsQuitting 是否正在退出
func (m NetworkModel) IsQuitting() bool {
	return m.quitting
}

// RunNetworkDiagnostics 运行网络诊断
func RunNetworkDiagnostics(host ssh.SSHHost) error {
	model := NewNetworkModel(host)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
