package ntp

type NTPConfig struct {
	Interval    string `json:"interval"`
	Server      string `json:"server"`
	ServerPort  int    `json:"server_port"`
	Detour      string `json:"detour"`
}

func BuildNTPConfig() *NTPConfig {
	return &NTPConfig{
		Interval:   "5m0s",
		Server:     "time.apple.com",
		ServerPort: 123,
		Detour:     "direct",
	}
}