package esp32

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// POST /api/v1/esp32/sensor-data
func (h *Handlers) SaveSensorData(c *gin.Context) {
	var req SensorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.service.SaveReading(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reading"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Reading saved successfully",
	})
}

// GET /api/v1/esp32/readings/:device_id?hours=24
func (h *Handlers) GetReadings(c *gin.Context) {
	deviceID := c.Param("device_id")
	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))

	readings, err := h.service.GetRecentReadings(deviceID, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get readings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id": deviceID,
		"hours":     hours,
		"count":     len(readings),
		"readings":  readings,
	})
}

// GET /api/v1/esp32/stats/:device_id?hours=24
func (h *Handlers) GetStats(c *gin.Context) {
	deviceID := c.Param("device_id")
	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))

	stats, err := h.service.GetStats(deviceID, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	timeStats, _ := h.service.GetTimeOfDayStats(deviceID, hours)

	c.JSON(http.StatusOK, gin.H{
		"device_id":   deviceID,
		"hours":       hours,
		"stats":       stats,
		"time_of_day": timeStats,
	})
}
