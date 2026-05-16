package application

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/system/domain"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type SystemGetUseCase struct {
	*domain.SystemUseCase
}

func NewSystemGetUseCase(logger *zap.Logger) *SystemGetUseCase {
	return &SystemGetUseCase{
		SystemUseCase: domain.NewSystemUseCase(logger),
	}
}

func (u *SystemGetUseCase) SystemGetOverview(c fiber.Ctx) error {
	hostname, _ := os.Hostname()
	return c.JSON(domain.OverviewInfo{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Hostname:  hostname,
		Platform:  u.GetOsInfo(),
		Cpu:       u.GetCpuInfo(),
		Memory:    u.GetMemoryInfo(),
		Uptime:    u.GetUptimeInfo(),
		Links: map[string]string{
			"overview": "/system",
			"os":       "/system/os",
			"cpu":      "/system/cpu",
			"memory":   "/system/memory",
			"network":  "/system/network",
			"storage":  "/system/storage",
			"uptime":   "/system/uptime",
		},
	})
}

func (u *SystemGetUseCase) SystemGetOs(c fiber.Ctx) error {
	return c.JSON(u.GetOsInfo())
}

func (u *SystemGetUseCase) SystemGetCpu(c fiber.Ctx) error {
	return c.JSON(u.GetCpuInfo())
}

func (u *SystemGetUseCase) SystemGetMemory(c fiber.Ctx) error {
	return c.JSON(u.GetMemoryInfo())
}

func (u *SystemGetUseCase) SystemGetNetwork(c fiber.Ctx) error {
	ifaces, err := u.GetNetworkInfo()
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(ifaces)
}

func (u *SystemGetUseCase) SystemGetStorage(c fiber.Ctx) error {
	return c.JSON(u.GetStorageInfo())
}

func (u *SystemGetUseCase) SystemGetUptime(c fiber.Ctx) error {
	return c.JSON(u.GetUptimeInfo())
}

