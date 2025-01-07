package process_info

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/xyu-io/genie/iper"
)

func init() {
	// 系统信息，版本
	var sysInfo = ""
	hostInfo, err := host.Info()
	if err == nil {
		sysInfo = hostInfo.Platform + "@" + hostInfo.KernelVersion
	}
	// 获取CPU核心数
	logicNum, _ := cpu.Counts(true)
	physicalNum, _ := cpu.Counts(false)
	SysMsg = &ProcInfo{
		HostIp:         iper.LocalIP(),
		SysInfo:        sysInfo,
		CpuLogicNum:    logicNum,
		CpuPhysicalNum: physicalNum,
		CpuPercent:     0,
		MemoryUsed:     0,
		MemoryFree:     0,
		IOWrite:        0,
		IORead:         0,
	}
	ProcessMsg = &ProcInfo{
		HostIp:         iper.LocalIP(),
		SysInfo:        sysInfo,
		ProcessName:    "NULL",
		ProcessCreated: 0,
		ProcessPath:    "NULL",
		ProcessId:      999999,
		CpuPercent:     0,
		MemoryUsed:     0,
		MemoryFree:     0,
		IOWrite:        0,
		IORead:         0,
	}
}

func HostInfo() *ProcInfo {
	return hostInfos()
}

func ProcessInfos(args []string) []*ProcInfo {
	return processInfos(args)
}
