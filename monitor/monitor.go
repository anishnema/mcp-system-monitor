package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// GetCPUUsage samples CPU usage over a short duration
func GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, fmt.Errorf("no cpu data returned")
}

// GetMemoryUsage fetches virtual memory stats
func GetMemoryUsage() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

// GetDiskUsage fetches disk usage for root partition
func GetDiskUsage() (*disk.UsageStat, error) {
	// Root directory "/" (or "C:\" on Windows)
	return disk.Usage("/")
}

// GetNetworkBandwidth calculates network throughput by taking two snapshots 1 second apart
func GetNetworkBandwidth() (sentPerSec uint64, recvPerSec uint64, err error) {
	net1, err := net.IOCounters(false)
	if err != nil || len(net1) == 0 {
		return 0, 0, err
	}
	bytesSent1 := net1[0].BytesSent
	bytesRecv1 := net1[0].BytesRecv

	time.Sleep(1 * time.Second)

	net2, err := net.IOCounters(false)
	if err != nil || len(net2) == 0 {
		return 0, 0, err
	}
	bytesSent2 := net2[0].BytesSent
	bytesRecv2 := net2[0].BytesRecv

	sentPerSec = bytesSent2 - bytesSent1
	recvPerSec = bytesRecv2 - bytesRecv1
	return sentPerSec, recvPerSec, nil
}