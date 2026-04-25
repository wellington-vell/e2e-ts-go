package routers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"server/internal/models"
)

var startTime = time.Now()

// @Summary Health Check
// @Description Returns the health status of the server along with uptime and memory usage
// @Tags Health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /api/v1/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	_, offsetSeconds := now.Zone()
	resp := models.HealthResponse{
		Status:         "OK",
		Date:           now.Format("02/01/2006 15:04:05"),
		Timezone:       "America/Sao_Paulo",
		TimezoneOffset: int64(offsetSeconds),
		TimezoneName:   "America/Sao_Paulo",
		Version:        "2.0.0",
		Uptime:         time.Since(startTime).Seconds(),
		Memory:         getMemoryUsage(),
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		Node:           runtime.Version(),
		Args:           []string{},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}

func getMemoryUsage() map[string]any {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return map[string]any{
		"heapAlloc": m.HeapAlloc,
		"heapSys":   m.HeapSys,
		"mallocs":   m.Mallocs,
	}
}
