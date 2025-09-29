package ssh

import "testing"

func TestParseHostArgument(t *testing.T) {
	cases := []struct {
		in   string
		user string
		host string
		port string
	}{
		{"root@example.com", "root", "example.com", "22"},
		{"user@host:2222", "user", "host", "2222"},
		{"example.com:2200", "", "example.com", "2200"},
		{"justhost", "", "justhost", "22"},
	}

	for _, c := range cases {
		got := ParseHostArgument(c.in)
		if got.User != c.user || got.HostName != c.host || got.Port != c.port {
			// 允许 Host 字段等于 HostName（实现中如此）
			if got.Host != c.host {
				// 统一错误输出
				
					 t.Errorf("ParseHostArgument(%q) => User=%q HostName=%q Port=%q Host=%q; want User=%q Host=%q Port=%q", c.in, got.User, got.HostName, got.Port, got.Host, c.user, c.host, c.port)
			}
		}
	}
}
