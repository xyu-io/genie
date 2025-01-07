package fwcli

// Action 枚举值定义
type Action string

const (
	Accept Action = "accept"
	Reject Action = "reject"
	Drop   Action = "drop"
)

// ItemType 是firewalld基本命令操作的类型
const (
	Interface string = "interface"
	Port      string = "port"
	Source    string = "source"
)

// IP的类型
const (
	IPv4 string = "ipv4"
	IPv6 string = "ipv6"
)
