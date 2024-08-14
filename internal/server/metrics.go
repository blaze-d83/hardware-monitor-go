package server

type Metrics struct {
	HostName       string
	TotalMemory    uint64
	UsedMemory     uint64
	OS             string
	TotalDiskSpace uint64
	FreeDiskSpace  uint64
	CPUModelName   string
	Cores          uint8
}
