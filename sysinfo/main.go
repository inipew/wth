package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	gop "github.com/shirou/gopsutil/net"
)

// Constants
const (
	defaultPort = ":5678"
	apiBasePath = "/api"
	staticDir   = "./frontend/dist"
)

// DiskStats represents disk usage statistics
type DiskStats struct {
	Total string `json:"total"`
	Free  string `json:"free"`
}

// NetworkAdapterStats represents network adapter statistics
type NetworkAdapterStats struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Upload    string `json:"upload"`
	Download  string `json:"download"`
}

// NetworkStats represents overall network statistics
type NetworkStats struct {
	Adapters []NetworkAdapterStats `json:"adapters"`
}

// MemStatsHumanReadable represents memory statistics in human-readable format
type MemStatsHumanReadable struct {
	Total        string  `json:"total_mb"`
	Free         string  `json:"free_mb"`
	Used         string  `json:"used_mb"`
	UsagePercent float64 `json:"usage_percent"`
}

// CPUStats represents CPU statistics
type CPUStats struct {
	UsagePercent float64 `json:"usage_percent"`
	Cores        int     `json:"cores"`
}

// DeviceStats represents overall device statistics
type DeviceStats struct {
	GOOS         string                `json:"goos"`
	GOARCH       string                `json:"goarch"`
	NumCPU       int                   `json:"num_cpu"`
	MemStats     MemStatsHumanReadable `json:"mem_stats"`
	DiskStats    DiskStats             `json:"disk_stats"`
	NetworkStats NetworkStats          `json:"network_stats"`
	LastBoot     string                `json:"last_boot"`
	Uptime       string                `json:"uptime"`
	CPUStats     CPUStats              `json:"cpu_stats"`
	NumProcesses int                   `json:"num_processes"`
}

// formatBytes converts bytes to a human-readable string
func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	}
	return fmt.Sprintf("%d B", bytes)
}

// formatDuration converts a duration to a human-readable string
func formatDuration(duration time.Duration) string {
	seconds := int(duration.Seconds())
	days := seconds / (24 * 3600)
	hours := (seconds % (24 * 3600)) / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60
	return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
}

// getMemStats retrieves memory statistics
func getMemStats() (MemStatsHumanReadable, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return MemStatsHumanReadable{}, fmt.Errorf("failed to get virtual memory info: %w", err)
	}
	return MemStatsHumanReadable{
		Total:        formatBytes(memInfo.Total),
		Free:         formatBytes(memInfo.Available),
		Used:         formatBytes(memInfo.Used),
		UsagePercent: memInfo.UsedPercent,
	}, nil
}

// getDiskStats retrieves disk statistics
func getDiskStats() (DiskStats, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return DiskStats{}, fmt.Errorf("failed to get disk stats: %w", err)
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	return DiskStats{
		Total: formatBytes(total),
		Free:  formatBytes(free),
	}, nil
}

// getNetworkStats retrieves network statistics
func getNetworkStats() (NetworkStats, error) {
	var networkStats NetworkStats
	interfaces, err := net.Interfaces()
	if err != nil {
		return networkStats, fmt.Errorf("failed to get network interfaces: %w", err)
	}
	ioCounters, err := gop.IOCounters(true)
	if err != nil {
		return networkStats, fmt.Errorf("failed to get IO counters: %w", err)
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		ip := ""
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				ip = ipNet.IP.String()
				break
			}
		}
		for _, stat := range ioCounters {
			if stat.Name == iface.Name {
				networkStats.Adapters = append(networkStats.Adapters, NetworkAdapterStats{
					Name:      iface.Name,
					IPAddress: ip,
					Upload:    formatBytes(stat.BytesSent),
					Download:  formatBytes(stat.BytesRecv),
				})
				break
			}
		}
	}
	return networkStats, nil
}

// getUptime retrieves system uptime
func getUptime() (time.Duration, error) {
	uptime, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, fmt.Errorf("failed to read uptime: %w", err)
	}
	var up float64
	if _, err := fmt.Sscanf(string(uptime), "%f", &up); err != nil {
		return 0, fmt.Errorf("failed to parse uptime: %w", err)
	}
	return time.Duration(up) * time.Second, nil
}

// getLastBoot retrieves the last boot time
func getLastBoot() (string, error) {
	bootTime, err := os.Stat("/proc/stat")
	if err != nil {
		return "", fmt.Errorf("failed to get last boot time: %w", err)
	}
	return bootTime.ModTime().Format(time.RFC3339), nil
}

// getCPUStats retrieves CPU statistics
func getCPUStats() (CPUStats, error) {
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return CPUStats{}, fmt.Errorf("failed to get CPU usage: %w", err)
	}
	return CPUStats{
		UsagePercent: cpuUsage[0],
		Cores:        runtime.NumCPU(),
	}, nil
}

// getNumProcesses retrieves the number of running processes
func getNumProcesses() (int, error) {
	cmd := exec.Command("bash", "-c", "ps aux | wc -l")
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of processes: %w", err)
	}
	processes, err := strconv.Atoi(string(out[:len(out)-1])) // Remove trailing newline
	if err != nil {
		return 0, fmt.Errorf("failed to parse number of processes: %w", err)
	}
	return processes, nil
}

// writeError writes an error response
func writeError(c *fiber.Ctx, err error, statusCode int) error {
	log.Printf("Error: %v", err)
	return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
}

// getDeviceStats handles the request for device statistics
func getDeviceStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	type result struct {
		data interface{}
		err  error
	}

	memStatsCh := make(chan result, 1)
	diskStatsCh := make(chan result, 1)
	networkStatsCh := make(chan result, 1)
	uptimeCh := make(chan result, 1)
	lastBootCh := make(chan result, 1)
	cpuStatsCh := make(chan result, 1)
	numProcessesCh := make(chan result, 1)

	go func() {
		memStats, err := getMemStats()
		memStatsCh <- result{data: memStats, err: err}
	}()

	go func() {
		diskStats, err := getDiskStats()
		diskStatsCh <- result{data: diskStats, err: err}
	}()

	go func() {
		networkStats, err := getNetworkStats()
		networkStatsCh <- result{data: networkStats, err: err}
	}()

	go func() {
		uptime, err := getUptime()
		uptimeCh <- result{data: uptime, err: err}
	}()

	go func() {
		lastBoot, err := getLastBoot()
		lastBootCh <- result{data: lastBoot, err: err}
	}()

	go func() {
		cpuStats, err := getCPUStats()
		cpuStatsCh <- result{data: cpuStats, err: err}
	}()

	go func() {
		numProcesses, err := getNumProcesses()
		numProcessesCh <- result{data: numProcesses, err: err}
	}()

	var stats DeviceStats

	select {
	case r := <-memStatsCh:
		if r.err != nil {
			return writeError(c, r.err, fiber.StatusInternalServerError)
		}
		stats.MemStats = r.data.(MemStatsHumanReadable)
	case <-ctx.Done():
		return writeError(c, ctx.Err(), fiber.StatusRequestTimeout)
	}

	select {
	case r := <-diskStatsCh:
		if r.err != nil {
			return writeError(c, r.err, fiber.StatusInternalServerError)
		}
		stats.DiskStats = r.data.(DiskStats)
	case <-ctx.Done():
		return writeError(c, ctx.Err(), fiber.StatusRequestTimeout)
	}

	select {
	case r := <-networkStatsCh:
		if r.err != nil {
			return writeError(c, r.err, fiber.StatusInternalServerError)
		}
		stats.NetworkStats = r.data.(NetworkStats)
	case <-ctx.Done():
		return writeError(c, ctx.Err(), fiber.StatusRequestTimeout)
	}

	select {
	case r := <-uptimeCh:
		if r.err != nil {
			return writeError(c, r.err, fiber.StatusInternalServerError)
		}
		stats.Uptime = formatDuration(r.data.(time.Duration))
	case <-ctx.Done():
		return writeError(c, ctx.Err(), fiber.StatusRequestTimeout)
	}

	select {
	case r := <-lastBootCh:
		if r.err != nil {
			return writeError(c, r.err, fiber.StatusInternalServerError)
		}
		stats.LastBoot = r.data.(string)
	case <-ctx.Done():
		return writeError(c, ctx.Err(), fiber.StatusRequestTimeout)
	}

	select {
	case r := <-cpuStatsCh:
		if r.err != nil {
			return writeError(c, r.err, fiber.StatusInternalServerError)
		}
		stats.CPUStats = r.data.(CPUStats)
	case <-ctx.Done():
		return writeError(c, ctx.Err(), fiber.StatusRequestTimeout)
	}

	select {
	case r := <-numProcessesCh:
		if r.err != nil {
			return writeError(c, r.err, fiber.StatusInternalServerError)
		}
		stats.NumProcesses = r.data.(int)
	case <-ctx.Done():
		return writeError(c, ctx.Err(), fiber.StatusRequestTimeout)
	}

	stats.GOOS = runtime.GOOS
	stats.GOARCH = runtime.GOARCH
	stats.NumCPU = runtime.NumCPU()

	return c.JSON(stats)
}

func main() {
	app := fiber.New(fiber.Config{
		DisablePreParseMultipartForm: true,
		StreamRequestBody:            true,
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	// API routes
	api := app.Group(apiBasePath)
	api.Get("/device-stats", getDeviceStats)

	// Serve static files
	app.Static("/", staticDir, fiber.Static{
		Compress: true,
	})

	// Fallback route for SPA
	app.Get("/*", func(c *fiber.Ctx) error {
		if _, err := os.Stat(filepath.Join(staticDir, c.Path())); os.IsNotExist(err) {
			return c.SendFile(filepath.Join(staticDir, "index.html"))
		}
		return c.SendFile(filepath.Join(staticDir, c.Path()))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Server running on http://localhost%s\n", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}