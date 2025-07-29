package ui

import (
	"fmt"

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
	for range 5 {
		result, err := measurer.MeasureLatency(host.Host, "tcp")
		if err != nil {
			return err
		}
		fmt.Printf(i18n.T(i18n.LatencyResult)+"\n", result.Duration)
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
	
	hops, err := tracer.TraceRouteWithCallback(host.Host, callback)
	if err != nil {
		return err
	}
	
	// 如果回调函数为空，则批量显示结果
	if len(hops) > 0 && callback == nil {
		fmt.Println(i18n.T(i18n.RouteTraceResults))
		for _, hop := range hops {
			if hop.IP == nil {
				fmt.Printf(i18n.T(i18n.RouteHopTimeout)+"\n", hop.Index, hop.RTT)
			} else {
				fmt.Printf(i18n.T(i18n.RouteHopInfo)+"\n", hop.Index, hop.IP, hop.RTT)
			}
		}
	}
	
	return nil
}