package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	gop "github.com/shirou/gopsutil/net"
)

type DiskStats struct {
    Total string `json:"total"`
    Free  string `json:"free"`
}

type NetworkAdapterStats struct {
    Name       string `json:"name"`
    IPAddress  string `json:"ip_address"`
    Upload     string `json:"upload"`
    Download   string `json:"download"`
}

type NetworkStats struct {
    Adapters []NetworkAdapterStats `json:"adapters"`
}

type MemStatsHumanReadable struct {
    Total        string  `json:"total_mb"`
    Free         string  `json:"free_mb"`
    Used         string  `json:"used_mb"`
    UsagePercent float64 `json:"usage_percent"`
}

type CpuStats struct {
    UsagePercent float64 `json:"usage_percent"`
    Cores        int     `json:"cores"`
}

type DeviceStats struct {
    GOOS         string                  `json:"goos"`
    GOARCH       string                  `json:"goarch"`
    NumCPU       int                     `json:"num_cpu"`
    MemStats     MemStatsHumanReadable   `json:"mem_stats"`
    DiskStats    DiskStats               `json:"disk_stats"`
    NetworkStats NetworkStats            `json:"network_stats"`
    LastBoot     string                  `json:"last_boot"`
    Uptime       string                  `json:"uptime"`
    CpuStats     CpuStats                `json:"cpu_stats"`
    NumProcesses int                     `json:"num_processes"`
}

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

func formatDuration(duration time.Duration) string {
    seconds := int(duration.Seconds())
    days := seconds / (24 * 3600)
    hours := (seconds % (24 * 3600)) / 3600
    minutes := (seconds % 3600) / 60
    seconds = seconds % 60
    return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
}

func getMemStats() (MemStatsHumanReadable, error) {
    memInfo, err := mem.VirtualMemory()
    if err != nil {
        return MemStatsHumanReadable{}, err
    }
    return MemStatsHumanReadable{
        Total:        formatBytes(memInfo.Total),
        Free:         formatBytes(memInfo.Free),
        Used:         formatBytes(memInfo.Used),
        UsagePercent: memInfo.UsedPercent,
    }, nil
}

func getDiskStats() (DiskStats, error) {
    var stat syscall.Statfs_t
    if err := syscall.Statfs("/", &stat); err != nil {
        return DiskStats{}, err
    }
    total := stat.Blocks * uint64(stat.Bsize)
    free := stat.Bfree * uint64(stat.Bsize)
    return DiskStats{
        Total: formatBytes(total),
        Free:  formatBytes(free),
    }, nil
}

func getNetworkStats() (NetworkStats, error) {
    var networkStats NetworkStats
    interfaces, err := net.Interfaces()
    if err != nil {
        return networkStats, err
    }
    ioCounters, err := gop.IOCounters(true)
    if err != nil {
        return networkStats, err
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
            }
        }
    }
    return networkStats, nil
}

func getUptime() (time.Duration, error) {
    uptime, err := os.ReadFile("/proc/uptime")
    if err != nil {
        return 0, err
    }
    var up float64
    _, err = fmt.Sscanf(string(uptime), "%f", &up)
    if err != nil {
        return 0, err
    }
    return time.Duration(up) * time.Second, nil
}

func getLastBoot() (string, error) {
    bootTime, err := os.Stat("/proc/stat")
    if err != nil {
        return "", err
    }
    return bootTime.ModTime().Format(time.RFC3339), nil
}

func getCpuStats() (CpuStats, error) {
    cpuUsage, err := cpu.Percent(0, false)
    if err != nil {
        return CpuStats{}, err
    }
    return CpuStats{
        UsagePercent: cpuUsage[0],
        Cores:        runtime.NumCPU(),
    }, nil
}

func getNumProcesses() (int, error) {
    cmd := exec.Command("bash", "-c", "ps aux | wc -l")
    out, err := cmd.Output()
    if err != nil {
        return 0, err
    }
    processes, err := strconv.Atoi(string(out[:len(out)-1])) // Remove trailing newline
    if err != nil {
        return 0, err
    }
    return processes, nil
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
    http.Error(w, err.Error(), statusCode)
    log.Printf("Error: %v", err)
}

func getDeviceStats(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()

    memStatsCh := make(chan MemStatsHumanReadable, 1)
    diskStatsCh := make(chan DiskStats, 1)
    networkStatsCh := make(chan NetworkStats, 1)
    uptimeCh := make(chan time.Duration, 1)
    lastBootCh := make(chan string, 1)
    cpuStatsCh := make(chan CpuStats, 1)
    numProcessesCh := make(chan int, 1)

    go func() {
        memStats, err := getMemStats()
        if err != nil {
            writeError(w, err, http.StatusInternalServerError)
            return
        }
        memStatsCh <- memStats
    }()

    go func() {
        diskStats, err := getDiskStats()
        if err != nil {
            writeError(w, err, http.StatusInternalServerError)
            return
        }
        diskStatsCh <- diskStats
    }()

    go func() {
        networkStats, err := getNetworkStats()
        if err != nil {
            writeError(w, err, http.StatusInternalServerError)
            return
        }
        networkStatsCh <- networkStats
    }()

    go func() {
        uptime, err := getUptime()
        if err != nil {
            writeError(w, err, http.StatusInternalServerError)
            return
        }
        uptimeCh <- uptime
    }()

    go func() {
        lastBoot, err := getLastBoot()
        if err != nil {
            writeError(w, err, http.StatusInternalServerError)
            return
        }
        lastBootCh <- lastBoot
    }()

    go func() {
        cpuStats, err := getCpuStats()
        if err != nil {
            writeError(w, err, http.StatusInternalServerError)
            return
        }
        cpuStatsCh <- cpuStats
    }()

    go func() {
        numProcesses, err := getNumProcesses()
        if err != nil {
            writeError(w, err, http.StatusInternalServerError)
            return
        }
        numProcessesCh <- numProcesses
    }()

    select {
    case memStats := <-memStatsCh:
        stats := DeviceStats{
            GOOS:         runtime.GOOS,
            GOARCH:       runtime.GOARCH,
            NumCPU:       runtime.NumCPU(),
            MemStats:     memStats,
            DiskStats:    <-diskStatsCh,
            NetworkStats: <-networkStatsCh,
            LastBoot:     <-lastBootCh,
            Uptime:       formatDuration(<-uptimeCh),
            CpuStats:     <-cpuStatsCh,
            NumProcesses: <-numProcessesCh,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(stats)
    case <-ctx.Done():
        writeError(w, ctx.Err(), http.StatusRequestTimeout)
    }
}

func main() {
    http.HandleFunc("/api/device-stats", getDeviceStats)
    log.Fatal(http.ListenAndServe(":5678", nil))
}
