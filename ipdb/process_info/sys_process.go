package process_info

import (
	"errors"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"log"
	"regexp"
	"strconv"
)

func processInfos(args []string) []*ProcInfo {
	if len(args) == 0 {
		return allProcessInfo()
	}
	// var res []*ProcInfo
	var res = getProcessList(args)
	return res
}

func getProcessList(args []string) []*ProcInfo {
	var res []*ProcInfo
	for _, processArg := range args {
		var tmp, err = processInfoBySign(processArg)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		res = append(res, &tmp)
	}
	return res
}

// ProcessInfoBySign 支持pid, pname获取进程资源信息
func processInfoBySign(processSign string) (ProcInfo, error) {
	var pidStr = ""
	if isNumber(processSign) {
		pidStr = processSign
	} else {
		//// 获取 进程ID
		pidStr = getPidByName(processSign)
		if pidStr == "" {
			return ProcInfo{}, errors.New("get pidStr error")
		}
	}

	// 获取进程的 CPU、内存和 IO 信息
	processPid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Println("get process pid", "err", err)
	}

	psInfo, err := process.NewProcess(int32(processPid))
	if err != nil {
		log.Println("get process info", "err", err)
		return ProcInfo{}, errors.New("get process info error")
	}
	createTime, err := psInfo.CreateTime()
	if err == nil {
		ProcessMsg.ProcessCreated = createTime
	}

	ProcessMsg.ProcessId = int(psInfo.Pid)
	if pName, errs := psInfo.Name(); errs == nil {
		ProcessMsg.ProcessName = pName
	}

	cwd, _ := psInfo.Exe() //Cwd()
	ProcessMsg.ProcessPath = cwd

	cpuPercent, err := psInfo.CPUPercent()
	if err != nil {
		log.Println("get cpu percent", "err", err)
		ProcessMsg.CpuPercent = 0.0
	} else {
		ProcessMsg.CpuPercent = cpuPercent
		log.Printf("process cpu use percent: %.2f%%\n", cpuPercent)
	}

	memPercent, err := psInfo.MemoryPercent()
	if err != nil {
		log.Println("get cpu percent", "err", err)
		ProcessMsg.MemoryUsedPercent = 0.0
	} else {
		ProcessMsg.MemoryUsedPercent = float64(memPercent)
		log.Printf("process cpu use percent: %.2f%%\n", cpuPercent)
	}

	// 获取系统内存信息，获取总量（GB），剩余量（GB）
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("get memery info", "err", err)
	} else {
		ProcessMsg.MemoryTotal = float64(vmStat.Total) / 1024 / 1024 / 1024 // (GB)
		ProcessMsg.MemoryFree = float64(vmStat.Free) / 1024 / 1024 / 1024   // (GB)
		log.Printf("memery total: %.2fGB, uesd: %.2fGB, free: %.2fGB\n", float64(vmStat.Total)/1024/1024/1024, float64(vmStat.Used)/1024/1024/1024, float64(vmStat.Free)/1024/1024/1024)
	}

	// 获取进程内存使用量（GB）
	memStat, err := psInfo.MemoryInfo()
	if err != nil {
		log.Println("get memery", "err", err)
	} else {
		ProcessMsg.MemoryUsed = float64(memStat.RSS) / 1024 / 1024 / 1024 // (GB)
		log.Printf("memery used: %.2fGB\n", float64(memStat.RSS)/1024/1024/1024)
	}

	ioCounters, err := psInfo.IOCounters()
	if err != nil {
		log.Println("get process io info", "err", err)
	} else {
		ProcessMsg.IORead = float64(ioCounters.ReadCount)
		ProcessMsg.IOWrite = float64(ioCounters.WriteCount)
		log.Printf("process io readCount: %d, ioWriteCount: %d\n", ioCounters.ReadCount, ioCounters.WriteCount)
	}

	return *ProcessMsg, nil
}

func allProcessInfo() []*ProcInfo {
	var res []*ProcInfo
	var pids, err = process.Pids()
	if err != nil {
		res = append(res, ProcessMsg)
		return res
	}
	for _, pid := range pids {
		if pid == 0 {
			continue
		}
		var tmp, _ = processInfoBySign(strconv.Itoa(int(pid)))
		if tmp.MemoryUsedPercent != 0.0 && tmp.MemoryUsed != 0.0 {
			res = append(res, &tmp)
		}
	}
	return res
}

func getPidByName(name string) string {
	var pids, err = process.Pids()
	if err != nil {
		return ""
	}
	var processId = ""
	for _, pid := range pids {
		pn, err := process.NewProcess(pid)
		if err != nil {
			continue
		}
		Pname, err := pn.Name()
		if err != nil {
			continue
		}
		if name == Pname {
			processId = strconv.Itoa(int(pid))
		}
	}
	return processId
}

func isNumber(str string) bool {
	pattern := "^[0-9]+$"
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return match
}
