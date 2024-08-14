package hardware

import (
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemInfo struct {
	HostName string
	TotalMem uint64
	UsedMem  uint64
	OS       string
}

type DiskInfo struct {
	TotalDiskSpace uint64
	FreeDiskSpace  uint64
}

type CPUInfo struct {
	ModelName string
	Cores     int
}

func GetSystemInfo() (*SystemInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	hostStat, err := host.Info()
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		HostName: hostStat.Hostname,
		TotalMem: vmStat.Total,
		UsedMem:  vmStat.Used,
		OS:       runtime.GOOS,
	}, nil

}

func GetDiskInfo() (*DiskInfo, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &DiskInfo{
		TotalDiskSpace: diskStat.Total,
		FreeDiskSpace:  diskStat.Free,
	}, nil

}

func GetCPUInfo() (*CPUInfo, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	return &CPUInfo{
		ModelName: cpuStat[0].ModelName,
		Cores:     len(cpuStat),
	}, nil
}
