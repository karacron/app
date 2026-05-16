package domain

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"
)

type SystemUseCase struct {
	Log     *zap.Logger
	StartAt time.Time
}

func NewSystemUseCase(logger *zap.Logger) *SystemUseCase {
	return &SystemUseCase{
		Log:     logger,
		StartAt: time.Now(),
	}
}

func (u *SystemUseCase) GetOsInfo() OsInfo {
	hostname, _ := os.Hostname()
	info, _ := host.Info()
	release := ""
	if info != nil {
		release = info.KernelVersion
	}
	return OsInfo{
		Platform: runtime.GOOS,
		Release:  release,
		Hostname: hostname,
	}
}

func (u *SystemUseCase) GetCpuInfo() CpuInfo {
	infos, _ := cpu.Info()
	model := "unknown"
	var avgMHz float64
	if len(infos) > 0 {
		model = infos[0].ModelName
		var sum float64
		for _, c := range infos {
			sum += c.Mhz
		}
		avgMHz = sum / float64(len(infos))
	}
	cores, _ := cpu.Counts(true)
	return CpuInfo{
		Model:           model,
		Architecture:    runtime.GOARCH,
		Cores:           cores,
		AverageSpeedMhz: int(math.Round(avgMHz)),
	}
}

func (u *SystemUseCase) GetMemoryInfo() MemoryInfo {
	v, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{TotalBytes: 0, FreeBytes: 0, UsedBytes: 0, UsagePercent: 0}
	}
	usagePct := 0.0
	if v.Total > 0 {
		usagePct = math.Round(float64(v.Used)/float64(v.Total)*10000) / 100
	}
	return MemoryInfo{
		TotalBytes:   v.Total,
		FreeBytes:    v.Free,
		UsedBytes:    v.Used,
		UsagePercent: usagePct,
	}
}

func (u *SystemUseCase) GetUptimeInfo() UptimeInfo {
	seconds, _ := host.Uptime()
	s := float64(seconds)
	return UptimeInfo{
		Seconds: s,
		Minutes: math.Round(s/60*100) / 100,
		Hours:   math.Round(s/3600*100) / 100,
	}
}

func (u *SystemUseCase) GetNetworkInfo() ([]Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		u.Log.Error("GetNetwork", zap.Error(err))
		return nil, err
	}

	result := make([]Interface, 0, len(ifaces))
	for _, i := range ifaces {
		addrs := make([]Address, 0, len(i.Addrs))
		isLoopback := false
		for _, flag := range i.Flags {
			if flag == "loopback" {
				isLoopback = true
				break
			}
		}
		for _, a := range i.Addrs {
			family := "IPv4"
			if len(a.Addr) > 15 {
				family = "IPv6"
			}
			addrs = append(addrs, Address{
				Family:   family,
				Address:  a.Addr,
				Internal: isLoopback,
				Mac:      i.HardwareAddr,
			})
		}
		result = append(result, Interface{Name: i.Name, Addresses: addrs})
	}

	return result, nil
}

func (u *SystemUseCase) GetStorageInfo() StorageInfo {
	partitions, err := disk.Partitions(false)
	if err != nil {
		u.Log.Error("GetStorage partitions", zap.Error(err))
		return StorageInfo{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Disks:     []StorageDisk{},
		}
	}

	disks := make([]StorageDisk, 0, len(partitions))
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		var total, free uint64
		if err == nil {
			total = usage.Total
			free = usage.Free
		}

		writable := false
		f, ferr := os.OpenFile(p.Mountpoint, os.O_RDONLY, 0)
		if ferr == nil {
			f.Close()
			writable = true
		}

		id := fmt.Sprintf("disk-%s", u.sanitizeID(p.Mountpoint))
		label := p.Mountpoint
		if label == "/" {
			label = "Root /"
		}

		disks = append(disks, StorageDisk{
			ID:         id,
			Label:      label,
			Path:       p.Mountpoint,
			Writable:   writable,
			TotalBytes: total,
			FreeBytes:  free,
		})
	}

	return StorageInfo{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Disks:     disks,
	}
}

func (u *SystemUseCase) sanitizeID(s string) string {
	b := []byte(s)
	for i, c := range b {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			b[i] = '-'
		}
	}
	return string(b)
}

