package ssh

// SSHHost 表示一个SSH主机配置
type SSHHost struct {
	Host     string
	HostName string
	User     string
	Port     string
	KeyFile  string
}
