package schemas

type HealthResponse struct {
	Status         string         `json:"status"`
	Date           string         `json:"date"`
	Timezone       string         `json:"timezone"`
	TimezoneOffset int64          `json:"timezoneOffset"`
	TimezoneName   string         `json:"timezoneName"`
	Version        string         `json:"version"`
	Uptime         float64        `json:"uptime"`
	Memory         map[string]any `json:"memory"`
	OS             string         `json:"os"`
	Arch           string         `json:"arch"`
	Node           string         `json:"node"`
	Args           []string       `json:"args"`
}
