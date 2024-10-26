package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"net/http"
)

type SystemStatsHandler struct{}

// NewSystemStatsHandler creates a new instance of the system stats handler
func NewSystemStatsHandler() *SystemStatsHandler {
	return &SystemStatsHandler{}
}

// GetSystemStats retrieves system statistics and returns them in JSON format
func (h *SystemStatsHandler) GetSystemStats(c *gin.Context) {
	// Get CPU statistics
	cpuStats, err := cpu.Percent(0, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get CPU stats"})
		return
	}

	// Get memory statistics
	memStats, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get memory stats"})
		return
	}

	// Get disk statistics
	diskStats, err := disk.Usage("/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get disk stats"})
		return
	}

	// Format and send the response in JSON
	c.JSON(http.StatusOK, gin.H{
		"cpu_usage":    cpuStats[0],
		"memory_total": memStats.Total,
		"memory_used":  memStats.Used,
		"memory_usage": memStats.UsedPercent,
		"disk_total":   diskStats.Total,
		"disk_used":    diskStats.Used,
		"disk_usage":   diskStats.UsedPercent,
	})
}
