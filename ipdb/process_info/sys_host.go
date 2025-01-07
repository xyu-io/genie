package process_info

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

func hostInfos() *ProcInfo {
	// 获取CPU信息
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Error("get cpu percent", "err", err)
	} else {
		SysMsg.CpuPercent = cpuPercent[0] / 100
	}

	// 获取内存信息
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Error("get memery info", "err", err)
	} else {
		SysMsg.MemoryTotal = float64(vmStat.Total) / 1024 / 1024 / 1024 // (GB)
		SysMsg.MemoryUsed = float64(vmStat.Used) / 1024 / 1024 / 1024   // (GB)
		SysMsg.MemoryFree = float64(vmStat.Free) / 1024 / 1024 / 1024   // (GB)
		SysMsg.MemoryUsedPercent = vmStat.UsedPercent / 100
		// log.Printf("memery total: %.2fGB, used: %.2fGB, free: %.2fGB\n", float64(vmStat.Total)/1024/1024/1024, float64(vmStat.Used)/1024/1024/1024, float64(vmStat.Free)/1024/1024/1024)
	}

	// 获取磁盘IO信息
	diskStat, err := disk.IOCounters()
	if err != nil {
		log.Error("get disk io info", "err", err)
	} else {
		var readCount float64 = 0.0
		var writeCount float64 = 0.0
		for _, stat := range diskStat {
			readCount = readCount + float64(stat.ReadCount)
			writeCount = writeCount + float64(stat.WriteCount)
			// log.Printf("%s disk io readCount: %d，iowriteCount: %d\n", diskName, stat.ReadCount, stat.WriteCount)
		}
		SysMsg.IORead = readCount
		SysMsg.IOWrite = writeCount
	}
	// 磁盘剩余空间大小
	parts, _ := disk.Partitions(true)
	var diskTotal = 0.0
	var diskUsed = 0.0
	var diskFree = 0.0
	for _, path := range parts {
		diskInfo, _ := disk.Usage(path.Mountpoint)
		diskTotal += float64(diskInfo.Total) / 1024 / 1024 / 1024 // (GB)
		diskUsed += float64(diskInfo.Used) / 1024 / 1024 / 1024   // (GB)
		diskFree += float64(diskInfo.Free) / 1024 / 1024 / 1024   // (GB)
	}
	SysMsg.DiskTotal = diskTotal
	SysMsg.DiskUsed = diskUsed
	SysMsg.DiskFree = diskFree
	SysMsg.DiskUsedPercent = diskUsed / diskTotal

	return SysMsg
}
