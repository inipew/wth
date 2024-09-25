package experimental

type ExperimentalConfig struct {
	CacheFile CacheFileConfig `json:"cache_file"`
	ClashAPI  ClashAPIConfig  `json:"clash_api"`
}

type CacheFileConfig struct {
	Enabled   bool   `json:"enabled"`
	Path      string `json:"path"`
	CacheID   string `json:"cache_id"`
	StoreRDRC bool   `json:"store_rdrc"`
}

type ClashAPIConfig struct {
	ExternalController     string `json:"external_controller"`
	ExternalUI             string `json:"external_ui"`
	ExternalUIDownloadURL  string `json:"external_ui_download_url"`
	ExternalUIDownloadDetour string `json:"external_ui_download_detour"`
	Secret                 string `json:"secret"`
}

func BuildExperimentalConfig() *ExperimentalConfig {
	return &ExperimentalConfig{
		CacheFile: CacheFileConfig{
			Enabled:   true,
			Path:      "caches.db",
			CacheID:   "sing",
			StoreRDRC: true,
		},
		ClashAPI: ClashAPIConfig{
			ExternalController:       "[::]:9090",
			ExternalUI:               "dashboard",
			ExternalUIDownloadURL:    "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip",
			ExternalUIDownloadDetour: "direct",
			Secret:                   "qwe12345",
		},
	}
}