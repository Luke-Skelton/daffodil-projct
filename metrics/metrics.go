package metrics

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/Luke-Skelton/daffodil-projct/internal/storage" // Import your storage interface
)

// CollectMetrics collects system metrics and returns them as a slice of DataPoints.
func CollectMetrics() ([]storage.DataPoint, error) {
	var dataPoints []storage.DataPoint

	// --- CPU Usage ---
	cpuPercentages, err := cpu.Percent(time.Second, true) // Per-core CPU usage
	if err != nil {
		return nil, fmt.Errorf("error getting CPU usage: %w", err)
	}
	for i, percent := range cpuPercentages {
		dataPoints = append(dataPoints, storage.DataPoint{
			Name:      "cpu_usage",
			Timestamp: time.Now().UnixMilli(),
			Value:     percent,
			Tags: map[string]string{
				"core": fmt.Sprintf("cpu%d", i),
			},
		})
	}

	// --- Memory Usage ---
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting memory usage: %w", err)
	}
	dataPoints = append(dataPoints, storage.DataPoint{
		Name:      "memory_usage",
		Timestamp: time.Now().UnixMilli(),
		Value:     float64(vm.Used), // Convert to float64
		Tags: map[string]string{
			"type": "used",
		},
	})
	dataPoints = append(dataPoints, storage.DataPoint{
		Name:      "memory_total",
		Timestamp: time.Now().UnixMilli(),
		Value:     float64(vm.Total), // Convert to float64
		Tags:      map[string]string{},
	})

	// --- Network I/O ---
	netIO, err := net.IOCounters(true) // Per-interface
	if err != nil {
		return nil, fmt.Errorf("error getting network I/O: %w", err)
	}
	for _, io := range netIO {
		dataPoints = append(dataPoints, storage.DataPoint{
			Name:      "network_bytes_sent",
			Timestamp: time.Now().UnixMilli(),
			Value:     float64(io.BytesSent),
			Tags: map[string]string{
				"interface": io.Name,
			},
		})
		dataPoints = append(dataPoints, storage.DataPoint{
			Name:      "network_bytes_recv",
			Timestamp: time.Now().UnixMilli(),
			Value:     float64(io.BytesRecv),
			Tags: map[string]string{
				"interface": io.Name,
			},
		})
	}

	// --- Disk Usage ---
	parts, err := disk.Partitions(true)
	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			log.Printf("Error getting disk usage for %s: %v", part.Mountpoint, err)
			continue // Skip this partition on error
		}
		dataPoints = append(dataPoints, storage.DataPoint{
			Name:      "disk_usage",
			Timestamp: time.Now().UnixMilli(),
			Value:     float64(usage.Used),
			Tags: map[string]string{
				"mountpoint": part.Mountpoint,
				"fstype":     part.Fstype,
			},
		})
	}

	return dataPoints, nil
}
