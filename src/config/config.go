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

	// FallbackPageFile overrides the built-in fallback HTML with the contents of a file on
	// disk. Empty = use the default. Supports ${VAR} expansion, e.g. "${PAGES_DIR}/fallback.html".
	// Only used in "fallback" Mode.
	FallbackPageFile string `json:"fallback_page_file"`

	// --- Rocket / maintenance-check fields, required in both modes ---
	//
	// "maintenance" Mode checks these on every app route and passes through when not in
	// maintenance. "fallback" Mode checks the same thing on the priority-1 underlay route:
	// if the app isn't reachable because Rocket has it in maintenance, the maintenance page
	// (with bypass) is shown; otherwise the plain FallbackPageFile "unavailable" page is shown,
	// since there's no real app behind this route to fall through to.

	AppId                string `json:"app_id"`
	RocketBaseUrl        string `json:"rocket_base_url"`
	RocketIdentityHeader string `json:"rocket_identity_header"`
	RocketTimeoutSeconds int    `json:"rocket_timeout_seconds"`
	CacheTtlSeconds      int    `json:"cache_ttl_seconds"`

	// MaintenancePageFile overrides the built-in maintenance HTML with the contents of a file
	// on disk. Empty = use the default. Supports ${VAR} expansion, same as FallbackPageFile.
	MaintenancePageFile string `json:"maintenance_page_file"`
}
