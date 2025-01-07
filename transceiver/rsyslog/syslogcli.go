package rsyslog

import (
	"encoding/json"
	"errors"
	syslog "github.com/RackSec/srslog"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
)

// SyslogClient syslog 客户端
type SyslogClient struct {
	Level logrus.Level
	entry *logrus.Entry
}

type Config struct {
	Network       string           // 协议 udp, tcp, tcp+tls
	Address       string           // 地址 localhost:514
	Cert          string           // tls证书路径
	AppName       string           // log标签
	Formatter     syslog.Formatter // log 格式化模板
	DisableOutput bool             // 是否关闭控制台打印日志, 默认开启
	Level         logrus.Level     // 自定义日志等级，也可以通过SendXX指定其他等级输出
}

// NewSyslogClient 创建 NewSyslogClient 实例
func NewSyslogClient(c Config) *SyslogClient {
	hook, err := CreatHook(c.Network, c.Address, c.AppName, c.Cert, c.Formatter)
	if err != nil {
		return nil
	}

	logger := logrus.New()
	logger.Hooks.Add(hook)
	if c.DisableOutput {
		logger.SetOutput(io.Discard) // 关闭控制台输出
	}
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	// 下面这行打印error等级会覆盖hook中syslog的level
	// logger.Info("this is a syslog Info message")
	// logger.Error("this is a syslog Error message")
	// logger.Fatal("this is a syslog Fatal message")

	return &SyslogClient{Level: c.Level, entry: logrus.NewEntry(logger)}
}

func (s *SyslogClient) Close() {
	// 获取当前 logger 的钩子列表
	hooks := s.entry.Logger.Hooks
	for level, levelHooks := range hooks {
		for i, hook := range levelHooks {
			if syslogHook, ok := hook.(*SyslogTlsHook); ok {
				if err := syslogHook.Writer.Close(); err != nil {
					s.entry.Logger.Errorf("Error closing syslog hook for level %v: %v", level, err)
				}
				// 从钩子列表中移除已关闭的钩子
				hooks[level] = append(hooks[level][:i], hooks[level][i+1:]...)
			}
		}
	}
}

// CreatHook 对外暴露使用，便于支持本库不支持的其他自定义场景
func CreatHook(network, address, appName string, certPath string, formatter syslog.Formatter) (logrus.Hook, error) {
	if certPath != "" {
		hook, err := NewSyslogHookTls(address, appName, certPath, formatter)
		if err != nil {
			return hook, err
		}
		return hook, nil
	}

	hook, err := NewSyslogHook(network, address, appName, formatter)
	if err != nil {
		return hook, err
	}

	return hook, nil
}

func (s *SyslogClient) SendLog(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(s.Level, msg)
}

func (s *SyslogClient) SendInfo(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.InfoLevel, msg)
}

func (s *SyslogClient) SendWarn(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.WarnLevel, msg)
}

func (s *SyslogClient) SendDebug(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.DebugLevel, msg)
}

func (s *SyslogClient) SendError(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.ErrorLevel, msg)
}

func (s *SyslogClient) SendFatal(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.FatalLevel, msg)
}

func (s *SyslogClient) SendTrace(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.TraceLevel, msg)
}

func outMsg(message any) (string, error) {
	var out = ""
	if msg, err := dataType(message); err != nil {
		bytes, _ := json.Marshal(message)
		if json.Valid(bytes) {
			out = string(bytes)
			return out, nil
		}
		out = string(bytes)
	} else {
		out = msg
	}
	//reg := regexp.MustCompile(`[\n\r\t\\/\"]`)
	//out = reg.ReplaceAllString(string(bytes), "'")
	return out, nil
}

func dataType(keyAsAny any) (string, error) {
	switch k := keyAsAny.(type) {
	case string:
		return k, nil
	case bool:
		return strconv.FormatBool(k), nil
	case byte:
		return string(k), nil
	case []byte:
		return string(k), nil
	case int:
		return strconv.Itoa(k), nil
	case int32:
		return strconv.Itoa(int(k)), nil
	case int64:
		return strconv.FormatInt(k, 10), nil
	case uint32:
		return strconv.FormatUint(uint64(k), 10), nil
	case uint64:
		return strconv.FormatUint(k, 10), nil
	case float32:
		return strconv.FormatFloat(float64(k), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(k, 'f', -1, 64), nil
	default:
		return "", errors.New("key type not supported")
	}
}
