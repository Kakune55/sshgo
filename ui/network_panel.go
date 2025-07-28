package ui

import (
	"fmt"
	"sshgo/network"
	"sshgo/ssh"

	"github.com/manifoldco/promptui"
)

// 显示网络诊断菜单
func ShowNetworkDiagnostics(host ssh.SSHHost) error {
	for {
		items := []string{
			"测量延迟",
			"路由追踪",
			"返回主菜单",
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("网络诊断 - %s", host.Host),
			Items: items,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("选择操作失败: %v", err)
		}

		switch result {
		case "测量延迟":
			if err := performLatencyTest(host); err != nil {
				fmt.Printf("延迟测试失败: %v\n", err)
			}
		case "路由追踪":
			if err := performRouteTrace(host); err != nil {
				fmt.Printf("路由追踪失败: %v\n", err)
			}
		case "返回主菜单":
			return nil
		}
	}
}

// 执行延迟测试
func performLatencyTest(host ssh.SSHHost) error {
	fmt.Printf("正在测试到 %s 的延迟...\n", host.Host)
	measurer := network.NewLatencyMeasurer()
	result, err := measurer.MeasureLatency(host.Host, "tcp")
	if err != nil {
		return err
	}
	fmt.Printf("延迟: %v\n", result.Duration)
	return nil
}

// 执行路由追踪
func performRouteTrace(host ssh.SSHHost) error {
	fmt.Printf("正在追踪到 %s 的路由...\n", host.Host)
	tracer := network.NewRouteTracer()
	hops, err := tracer.TraceRoute(host.Host)
	if err != nil {
		return err
	}
	fmt.Println("路由追踪结果:")
	for _, hop := range hops {
		fmt.Printf("跳点 %d: %s (RTT: %v)\n", hop.Index, hop.IP, hop.RTT)
	}
	return nil
}