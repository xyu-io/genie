package fwcli

import (
	"fmt"
	"testing"
)

// TestFirewall01 测试rich-rule
func TestFirewall01(t *testing.T) {
	manager := NewFirewallManager()

	// Prepare parameters
	ip := "10.52.3.1"
	zone := ""
	protocol := "tcp"
	port := "3306"
	action := Accept
	permanent := true

	//添加rich-rule
	err := manager.AddIPRule(ip, zone, protocol, port, action, permanent)
	if err != nil {
		t.Fatalf("Error adding IP rule: %v", err)
	}

	//删除rich-rule
	err1 := manager.RemoveIPRule(ip, zone, protocol, port, action, permanent)
	if err1 != nil {
		t.Fatalf("Error removing IP rule: %v", err1)
	}

}

// TestFirewall02 测试添加接口
func TestFirewall02(t *testing.T) {
	manager := NewFirewallManager()
	err := manager.AddInterface("ens35", "", false)
	if err != nil {
		fmt.Println("-------------", err)
	}
	err1 := manager.RemoveInterface("ens34", "public", false)
	if err1 != nil {
		fmt.Println("-------------", err1)
	}
}

// TestFirewall03 测试添加端口
func TestFirewall03(t *testing.T) {
	manager := NewFirewallManager()
	err := manager.AddPort("3306/tcp", "", false)
	if err != nil {
		fmt.Println("-------------", err)
	}
	err1 := manager.RemovePort("3306/tcp", "", false)
	if err1 != nil {
		fmt.Println("-------------", err1)
	}
}

// TestFirewall04 测试添加IP
func TestFirewall04(t *testing.T) {
	manager := NewFirewallManager()
	err := manager.AddSource("192.168.80.1", "", true)
	if err != nil {
		fmt.Println("---", err)
	}

	err1 := manager.RemoveSource("192.168.80.1", "", true)
	if err1 != nil {
		fmt.Println("---", err1)
	}
}

// TestFirewall05 测试获取interfaces
func TestFirewall05(t *testing.T) {
	manager := NewFirewallManager()
	interfaces, err := manager.GetItemType("", Interface)
	if err != nil {
		fmt.Println("---", err)
	}
	fmt.Println(interfaces)
}
