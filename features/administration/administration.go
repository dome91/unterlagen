package administration

import (
	"fmt"
	"runtime"
	"time"
	"unterlagen/features/common"
)

var Version string

type Administration struct {
	*settingsManager
	*users
	taskRepository common.TaskRepository
	startTime      time.Time
}

type RuntimeInfo struct {
	Version         string
	GoVersion       string
	GOOS            string
	GOARCH          string
	NumCPU          int
	NumGoroutine    int
	MemAllocMB      float64
	MemSysMB        float64
	NumGC           uint32
	UptimeSeconds   int64
	UptimeFormatted string
}

func (*Administration) UserRoles() []UserRole {
	return []UserRole{
		UserRoleAdmin,
		UserRoleUser,
	}
}

func (a *Administration) GetAllTasks() ([]common.Task, error) {
	return a.taskRepository.FindAll()
}

func (a *Administration) GetTasksPaginated(page int) ([]common.Task, int, int, error) {
	limit := 10
	offset := (page - 1) * limit
	tasks, total, err := a.taskRepository.FindPaginated(limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := (total + limit - 1) / limit
	return tasks, total, totalPages, nil
}

func (a *Administration) GetRuntimeInfo() RuntimeInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(a.startTime)

	return RuntimeInfo{
		Version:         Version,
		GoVersion:       runtime.Version(),
		GOOS:            runtime.GOOS,
		GOARCH:          runtime.GOARCH,
		NumCPU:          runtime.NumCPU(),
		NumGoroutine:    runtime.NumGoroutine(),
		MemAllocMB:      float64(m.Alloc) / 1024 / 1024,
		MemSysMB:        float64(m.Sys) / 1024 / 1024,
		NumGC:           m.NumGC,
		UptimeSeconds:   int64(uptime.Seconds()),
		UptimeFormatted: formatUptime(uptime),
	}
}

func formatUptime(uptime time.Duration) string {
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}

func New(settingsRepository SettingsRepository, userRepository UserRepository, userMessages UserMessages, taskRepository common.TaskRepository) *Administration {
	return &Administration{
		settingsManager: newSettingsManager(settingsRepository),
		users:           newUsers(userRepository, userMessages),
		taskRepository:  taskRepository,
		startTime:       time.Now(),
	}
}
