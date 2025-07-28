package network

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// 延迟测量结果
type LatencyResult struct {
	Host     string
	Protocol string
	Duration time.Duration
	Error    error
}

// 路由跳点信息
type RouteHop struct {
	Index int
	IP    net.IP
	RTT   time.Duration
}

// 延迟测量器
type LatencyMeasurer struct{}

// 创建新的延迟测量器
func NewLatencyMeasurer() *LatencyMeasurer {
	return &LatencyMeasurer{}
}

// 测量延迟
func (lm *LatencyMeasurer) MeasureLatency(host string, protocol string) (*LatencyResult, error) {
	switch protocol {
	case "tcp":
		return lm.measureTCP(host, 22)
	default:
		return lm.measureTCP(host, 22) // 默认使用TCP
	}
}

// TCP延迟测量
func (lm *LatencyMeasurer) measureTCP(host string, port int) (*LatencyResult, error) {
	start := time.Now()
	
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
	if err != nil {
		return &LatencyResult{
			Host:     host,
			Protocol: "tcp",
			Error:    err,
		}, nil
	}
	defer conn.Close()
	
	duration := time.Since(start)
	return &LatencyResult{
		Host:     host,
		Protocol: "tcp",
		Duration: duration,
	}, nil
}

// 路由追踪器
type RouteTracer struct {
	maxHops int
	timeout time.Duration
}

// 创建新的路由追踪器
func NewRouteTracer() *RouteTracer {
	return &RouteTracer{
		maxHops: 30,
		timeout: 5 * time.Second,
	}
}


// 追踪路由
func (rt *RouteTracer) TraceRoute(host string) ([]RouteHop, error) {
	// 解析目标地址
	destAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return nil, err
	}

	// 创建ICMP连接
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	defer c.Close()

	var hops []RouteHop

	for ttl := 1; ttl <= rt.maxHops; ttl++ {
		// 设置TTL
		p := c.IPv4PacketConn()
		if err := p.SetTTL(ttl); err != nil {
			return nil, err
		}

		// 创建ICMP Echo Request消息
		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: ttl,
				Data: []byte("HELLO-R-U-THERE"),
			},
		}
		wb, err := wm.Marshal(nil)
		if err != nil {
			return nil, err
		}

		// 发送ICMP消息
		start := time.Now()
		if _, err := c.WriteTo(wb, destAddr); err != nil {
			return nil, err
		}

		// 设置读取超时
		if err := c.SetReadDeadline(time.Now().Add(rt.timeout)); err != nil {
			return nil, err
		}

		// 读取响应
		rb := make([]byte, 1500)
		n, peer, err := c.ReadFrom(rb) // peer is net.Addr
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("跳点 %d: * (RTT: %v)\n", ttl, time.Since(start))
				hops = append(hops, RouteHop{Index: ttl})
				continue
			}
			return nil, err
		}
		rtt := time.Since(start)

		// 解析响应
		rm, err := icmp.ParseMessage(ipv4.ICMPType(0).Protocol(), rb[:n])
		if err != nil {
			return nil, err
		}

		// 获取IP地址
		var ip net.IP
		if ipAddr, ok := peer.(*net.IPAddr); ok {
			ip = ipAddr.IP
		} else {
			ip = net.ParseIP(peer.String())
		}

		// 实时打印跳点信息
		fmt.Printf("跳点 %d: %s (RTT: %v)\n", ttl, ip, rtt)
		hops = append(hops, RouteHop{Index: ttl, IP: ip, RTT: rtt})

		// 如果到达目的地，则停止
		if rm.Type == ipv4.ICMPTypeEchoReply {
			break
		}
	}

	return hops, nil
}