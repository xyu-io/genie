package process_info

type ProcInfo struct {
	HostIp            string  `json:"ip"`
	SysInfo           string  `json:"sys_info"`
	ProcessId         int     `json:"process_id,omitempty"`
	ProcessName       string  `json:"process_name,omitempty"`
	ProcessCreated    int64   `json:"created_time,omitempty"`
	ProcessPath       string  `json:"process_path,omitempty"`
	CpuPercent        float64 `json:"cpu_percent"`
	CpuLogicNum       int     `json:"cpu_logic_num,omitempty"`
	CpuPhysicalNum    int     `json:"cpu_physical_num,omitempty"`
	MemoryUsedPercent float64 `json:"memory_percent"`
	MemoryUsed        float64 `json:"memory_used"`
	MemoryTotal       float64 `json:"memory_total"`
	MemoryFree        float64 `json:"memory_free"`
	DiskUsedPercent   float64 `json:"disk_percent,omitempty"`
	DiskTotal         float64 `json:"disk_total,omitempty"`
	DiskUsed          float64 `json:"disk_used,omitempty"`
	DiskFree          float64 `json:"disk_free,omitempty"`
	IOWrite           float64 `json:"io_write"`
	IORead            float64 `json:"io_read"`
}

var ProcessMsg = &ProcInfo{}
var SysMsg = &ProcInfo{}
