package log

type LogConfig struct {
	Disabled  bool	`json:"disabled"`
	Level     string `json:"level"`
	Output    string `json:"output"`
	Timestamp bool   `json:"timestamp"`
}

func BuildLogConfig() *LogConfig {
	return &LogConfig{
		Disabled:  false,
		Level:     "info",
		Output:    "/etc/wth/log/sing-box.log",
		Timestamp: true,
	}
}