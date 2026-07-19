package config

const (
	ModeFallback    = "fallback"
	ModeMaintenance = "maintenance"
)

type Config struct {
	LogLevel string `json:"log_level"`

	// Mode selects the plugin behavior: "fallback" or "maintenance".
	Mode string `json:"mode"`

	// StatusCode is written when this instance serves its static page.
	StatusCode int `json:"status_code"`

	// FallbackPageContent overrides the built-in fallback HTML. Empty = use the default.
	// Supports the same ${VAR} / ${file:/path} convention as RocketIdentityHeader.
	FallbackPageContent string `json:"fallback_page_content"`

	// --- Maintenance-mode-only fields ---

	AppId                  string `json:"app_id"`
	RocketBaseUrl          string `json:"rocket_base_url"`
	RocketIdentityHeader   string `json:"rocket_identity_header"`
	RocketTimeoutSeconds   int    `json:"rocket_timeout_seconds"`
	CacheTtlSeconds        int    `json:"cache_ttl_seconds"`
	MaintenancePageContent string `json:"maintenance_page_content"`
}
