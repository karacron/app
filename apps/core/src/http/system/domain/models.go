package domain

type OsInfo struct {
	Platform string `json:"platform"`
	Release  string `json:"release"`
	Hostname string `json:"hostname"`
}

type CpuInfo struct {
	Model           string `json:"model"`
	Architecture    string `json:"architecture"`
	Cores           int    `json:"cores"`
	AverageSpeedMhz int    `json:"averageSpeedMhz"`
}

type MemoryInfo struct {
	TotalBytes   uint64  `json:"totalBytes"`
	FreeBytes    uint64  `json:"freeBytes"`
	UsedBytes    uint64  `json:"usedBytes"`
	UsagePercent float64 `json:"usagePercent"`
}

type UptimeInfo struct {
	Seconds float64 `json:"seconds"`
	Minutes float64 `json:"minutes"`
	Hours   float64 `json:"hours"`
}

type Address struct {
	Family   string `json:"family"`
	Address  string `json:"address"`
	Internal bool   `json:"internal"`
	Mac      string `json:"mac"`
}

type Interface struct {
	Name      string    `json:"name"`
	Addresses []Address `json:"addresses"`
}

type StorageDisk struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Path       string `json:"path"`
	Writable   bool   `json:"writable"`
	TotalBytes uint64 `json:"totalBytes"`
	FreeBytes  uint64 `json:"freeBytes"`
}

type StorageInfo struct {
	Timestamp string        `json:"timestamp"`
	Disks     []StorageDisk `json:"disks"`
}

type OverviewInfo struct {
	Timestamp string      `json:"timestamp"`
	Hostname  string      `json:"hostname"`
	Platform  OsInfo      `json:"platform"`
	Cpu       CpuInfo     `json:"cpu"`
	Memory    MemoryInfo  `json:"memory"`
	Uptime    UptimeInfo  `json:"uptime"`
	Links     map[string]string `json:"links"`
}

