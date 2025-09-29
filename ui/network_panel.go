package ui

import (
	"fmt"
	"time"
	"sshgo/i18n"
	"sshgo/network"
	"sshgo/ssh"
	"github.com/manifoldco/promptui"
)

// 显示网络诊断菜单
func ShowNetworkDiagnostics(host ssh.SSHHost) error {
	for {
		items := []string{
			i18n.T(i18n.MeasureLatencyAction),
			i18n.T(i18n.RouteTraceAction),
			i18n.T(i18n.ReturnToMainMenu),
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf(i18n.T(i18n.NetworkDiagnosticsTitle), host.Host),
			Items: items,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("选择操作失败: %v", err)
		}

		switch result {
		case i18n.T(i18n.MeasureLatencyAction):
			if err := performLatencyTest(host); err != nil {
				fmt.Printf("延迟测试失败: %v\n", err)
			}
		case i18n.T(i18n.RouteTraceAction):
			if err := performRouteTrace(host); err != nil {
				fmt.Printf("路由追踪失败: %v\n", err)
			}
		case i18n.T(i18n.ReturnToMainMenu):
			return nil
		}
	}
}

// 执行延迟测试
func performLatencyTest(host ssh.SSHHost) error {
	fmt.Printf(i18n.T(i18n.TestingLatency)+"\n", host.Host)
	measurer := network.NewLatencyMeasurer()
	const samples = 5
	var sum time.Duration
	var min time.Duration
	var max time.Duration
	var count int
	for i := 0; i < samples; i++ {
		result, err := measurer.MeasureLatency(host.Host, "tcp")
		if err != nil {
			return err
		}
		if result.Error != nil {
			fmt.Printf(i18n.T(i18n.LatencyResult)+"\n", result.Error)
			continue
		}
		d := result.Duration
		if count == 0 || d < min { min = d }
		if count == 0 || d > max { max = d }
		sum += d
		count++
		fmt.Printf(i18n.T(i18n.LatencyResult)+"\n", d)
		time.Sleep(200 * time.Millisecond)
	}
	if count > 0 {
		avg := sum / time.Duration(count)
		fmt.Printf(i18n.T(i18n.LatencySummary)+"\n", count, min, max, avg)
	}
	return nil
}

// 执行路由追踪
func performRouteTrace(host ssh.SSHHost) error {
	fmt.Printf(i18n.T(i18n.TracingRoute)+"\n", host.Host)
	tracer := network.NewRouteTracer()
	
	// 使用回调函数实现实时更新
	callback := func(hop network.RouteHop, isTimeout bool) {
		if isTimeout {
			fmt.Printf(i18n.T(i18n.RouteHopTimeout)+"\n", hop.Index, hop.RTT)
		} else {
			fmt.Printf(i18n.T(i18n.RouteHopInfo)+"\n", hop.Index, hop.IP, hop.RTT)
		}
	}
	
	_, err := tracer.TraceRouteWithCallback(host.Host, callback)
	if err != nil {
		return err
	}
	
	// 当前使用实时回调显示；如果未来支持无回调模式，可在此补充批量输出
	
	return nil
}