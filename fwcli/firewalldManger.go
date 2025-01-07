package fwcli

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// FirewallManager 结构体
type FirewallManager struct{}

// NewFirewallManager 返回一个新的 FirewallManager 实例
func NewFirewallManager() *FirewallManager {
	return &FirewallManager{}
}

// ExecuteCommand 执行给定的 shell 命令，并返回输出
func (fm *FirewallManager) ExecuteCommand(args ...string) (string, error) {
	// 在执行命令前，先确认防火墙是否能使用
	if ok, err := fm.isFirewallActive(); !ok {
		return "", fmt.Errorf("cannot execute command, firewall is not active: %v", err)
	}

	command := exec.Command("firewall-cmd", args...)
	// 打印命令及其参数
	fmt.Printf("Executing command:  %v\n", args)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, output: %s", err, output)
	}
	return string(output), nil // 返回命令输出
}

//-----------------rich-rule命令的使用-----------------

// AddIPRule 添加 IP 规则，支持选择协议和端口
func (fm *FirewallManager) AddIPRule(ip string, zone string, protocol string, port string, action Action, permanent bool) error {
	netProto, err := checkIPVersion(ip)
	if err != nil {
		return err
	}

	if action == "" {
		return errors.New("action cannot be empty")
	}

	// 开始构建 rule 字符串
	rule := fmt.Sprintf("rule family=\"%s\"", netProto)

	if ip != "" {
		rule += fmt.Sprintf(" source address=\"%s\"", ip)
	}

	// 添加 protocol 和 port 条件
	if protocol != "" {
		if port != "" {
			rule += fmt.Sprintf(" port port=\"%s\" protocol=\"%s\"", port, protocol)
		} else {
			return errors.New("port cannot be empty")
		}
	} else {
		if port != "" {
			return errors.New("protocol cannot be empty")
		}
	}

	// 最后添加 action
	rule += fmt.Sprintf(" %s", action)

	// 调用添加丰富规则的函数
	return fm.AddRichRule(rule, zone, permanent)
}

// RemoveIPRule 删除 IP 规则，支持选择协议和端口
func (fm *FirewallManager) RemoveIPRule(ip string, zone string, protocol string, port string, action Action, permanent bool) error {
	netProto, err := checkIPVersion(ip)
	if err != nil {
		return err
	}

	if action == "" {
		return errors.New("action cannot be empty")
	}

	// 开始构建 rule 字符串
	rule := fmt.Sprintf("rule family=\"%s\"", netProto)

	if ip != "" {
		rule += fmt.Sprintf(" source address=\"%s\"", ip)
	}

	// 添加 protocol 和 port 条件,
	//且终端输出的格式是：port port="3306" protocol="udp"
	if protocol != "" {
		if port != "" {
			rule += fmt.Sprintf(" port port=\"%s\" protocol=\"%s\"", port, protocol)
		} else {
			return errors.New("port cannot be empty")
		}
	} else {
		if port != "" {
			return errors.New("protocol cannot be empty")
		}
	}

	// 最后添加 action
	rule += fmt.Sprintf(" %s", action)

	// 调用删除丰富规则的函数
	return fm.RemoveRichRule(rule, zone, permanent)
}

// AddRichRule 添加丰富规则，支持协议、源地址、目标地址、端口、服务和动作
func (fm *FirewallManager) AddRichRule(rule string, zone string, permanent bool) error {

	// 检查规则是否已存在
	rules, err := fm.GetRichRules(zone)
	if err != nil {
		return err
	}

	for _, existingRule := range rules {
		if existingRule == rule {
			fmt.Printf("rule already exists:------ %s", rule)
			return nil
		}
	}

	// 构造命令行参数
	args := buildArgs(zone, rule, permanent, true)

	// 执行命令
	if _, err1 := fm.ExecuteCommand(args...); err1 != nil {
		return err1
	}

	err = fm.IsReload(permanent)
	if err != nil {
		return err
	}

	return nil
}

// RemoveRichRule 删除丰富规则
func (fm *FirewallManager) RemoveRichRule(rule string, zone string, permanent bool) error {
	// 检查规则是否存在
	rules, err := fm.GetRichRules(zone)
	if err != nil {
		return err
	}

	exists := false
	for _, existingRule := range rules {
		if existingRule == rule {
			exists = true
			break
		}
	}
	if !exists {
		return errors.New("rule does not exist")
	}

	// 构造命令行参数
	args := buildArgs(zone, rule, permanent, false)

	// 执行命令
	if _, err1 := fm.ExecuteCommand(args...); err1 != nil {
		return err1
	}

	err = fm.IsReload(permanent)
	if err != nil {
		return err
	}

	return nil
}

// GetRichRules 获取特定区域的rich-rule列表
func (fm *FirewallManager) GetRichRules(zone string) ([]string, error) {

	// 执行命令以获取丰富规则
	output, err := fm.ExecuteCommand("--zone", zone, "--list-rich-rules")
	if err != nil {
		return nil, err
	}

	// 使用正则表达式匹配规则
	re := regexp.MustCompile(`rule\s.*`)
	matches := re.FindAllString(output, -1)

	return matches, nil
}

// -----------------firewalld基础命令使用-----------------

// AddInterface adds an interface to a specified zone
func (fm *FirewallManager) AddInterface(interfaceName string, zone string, permanent bool) error {
	return fm.addItem(Interface, interfaceName, zone, permanent)
}

// RemoveInterface removes an interface from a specified zone
func (fm *FirewallManager) RemoveInterface(interfaceName string, zone string, permanent bool) error {
	return fm.removeItem(Interface, interfaceName, zone, permanent)
}

// AddPort adds a port to a specified zone
func (fm *FirewallManager) AddPort(port string, zone string, permanent bool) error {
	return fm.addItem(Port, port, zone, permanent)
}

// RemovePort removes a port from a specified zone
func (fm *FirewallManager) RemovePort(port string, zone string, permanent bool) error {
	return fm.removeItem(Port, port, zone, permanent)
}

// AddSource adds a source to a specified zone
func (fm *FirewallManager) AddSource(source string, zone string, permanent bool) error {
	return fm.addItem(Source, source, zone, permanent)
}

// RemoveSource removes a source from a specified zone
func (fm *FirewallManager) RemoveSource(source string, zone string, permanent bool) error {
	return fm.removeItem(Source, source, zone, permanent)
}

// Generic function to add an item to a zone
func (fm *FirewallManager) addItem(itemType, item, zone string, permanent bool) error {
	if item == "" {
		return errors.New(itemType + " cannot be empty")
	}

	existingItems, err := fm.GetItemType(zone, itemType)
	if err != nil {
		return err
	}

	if itemExists(existingItems, item) {
		return errors.New(itemType + " already exists")
	}

	// Construct command line arguments
	args := []string{"--zone", zone, "--add-" + itemType, item}
	if permanent {
		args = append(args, "--permanent")
	}

	// 执行命令
	if _, err1 := fm.ExecuteCommand(args...); err1 != nil {
		return err1
	}

	err = fm.IsReload(permanent)
	if err != nil {
		return err
	}

	return nil
}

// Generic function to remove an item from a zone
func (fm *FirewallManager) removeItem(itemType, item, zone string, permanent bool) error {
	if item == "" {
		return errors.New(itemType + " cannot be empty")
	}

	existingItems, err := fm.GetItemType(zone, itemType)
	if err != nil {
		return err
	}

	if !itemExists(existingItems, item) {
		return errors.New(itemType + " does not exist")
	}

	// Construct command line arguments
	args := []string{"--zone", zone, "--remove-" + itemType, item}
	if permanent {
		args = append(args, "--permanent")
	}

	// 执行命令
	if _, err1 := fm.ExecuteCommand(args...); err1 != nil {
		return err1
	}

	err = fm.IsReload(permanent)
	if err != nil {
		return err
	}

	return nil
}

func (fm *FirewallManager) GetItemType(zone string, itemType string) ([]string, error) {

	//执行命令获取item
	output, err := fm.ExecuteCommand("--zone", zone, "--list-"+itemType+"s")
	if err != nil {
		return nil, err
	}
	items := strings.Fields(output)

	return items, nil
}

//--------------------------其他命令的使用------------------------

// Reload 重新加载防火墙
func (fm *FirewallManager) Reload() error {
	_, err := fm.ExecuteCommand("--reload")
	if err != nil {
		return err
	}
	return nil
}

// GetActiveZones 获取活动区域
func (fm *FirewallManager) GetActiveZones() (string, error) {
	output, err := exec.Command("--get-active-zones").Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// ListRules 列出指定区域的规则
func (fm *FirewallManager) ListRules(zone string) (string, error) {
	output, err := exec.Command("--zone", zone, "--list-all").Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

//--------------

// IsReload 是否重载配置
func (fm *FirewallManager) IsReload(setRe bool) error {
	if setRe {
		err := fm.Reload()
		if err != nil {
			return err
		}
	}
	return nil
}

// isFirewallActive 判断防火墙是否存活
func (fm *FirewallManager) isFirewallActive() (bool, error) {
	command := exec.Command("firewall-cmd", "--state")
	output, err := command.CombinedOutput()
	if err != nil {
		return false, err
	}
	if string(output) != "running\n" {
		return false, nil
	}
	return true, nil
}
