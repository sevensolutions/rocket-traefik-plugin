package src

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sevensolutions/rocket-traefik-plugin/src/config"
	"github.com/sevensolutions/rocket-traefik-plugin/src/logging"
	"github.com/sevensolutions/rocket-traefik-plugin/src/pages"
	"github.com/sevensolutions/rocket-traefik-plugin/src/rocket"
	"github.com/sevensolutions/rocket-traefik-plugin/src/utils"
)

// Rocket's base URL and auth token are deliberately not part of the dynamic (per-middleware)
// config: they're deployment-wide secrets/endpoints, not something that should vary between
// routers or be visible in Traefik's dynamic configuration.
const rocketUrlEnvVar = "ROCKET_URL"
const rocketTokenEnvVar = "ROCKET_TOKEN"

func CreateConfig() *config.Config {
	return &config.Config{
		LogLevel:             logging.LevelWarn,
		Mode:                 config.ModeFallback,
		StatusCode:           http.StatusServiceUnavailable,
		RocketTimeoutSeconds: 5,
		CacheTtlSeconds:      30,
	}
}

// Will be called by traefik
func New(uctx context.Context, next http.Handler, cfg *config.Config, name string) (http.Handler, error) {
	cfg.LogLevel = utils.ExpandEnvironmentVariableString(cfg.LogLevel)

	logger := logging.CreateLogger(cfg.LogLevel)

	logger.Log(logging.LevelInfo, "Loading Configuration...")

	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))

	if mode != config.ModeFallback && mode != config.ModeMaintenance {
		return nil, fmt.Errorf("invalid Mode %q, must be %q or %q", cfg.Mode, config.ModeFallback, config.ModeMaintenance)
	}

	statusCode := cfg.StatusCode
	if statusCode <= 0 {
		statusCode = http.StatusServiceUnavailable
	}

	// Both modes check Rocket for maintenance status: "maintenance" mode does so on the app's
	// real route (passing through when not in maintenance), "fallback" mode does so on the
	// priority-1 underlay route (showing the plain fallback page when not in maintenance,
	// since there's no working app route to fall through to in that case).
	instanceKey := utils.ExpandEnvironmentVariableString(strings.TrimSpace(cfg.InstanceKey))
	rocketBaseUrl := strings.TrimSpace(os.Getenv(rocketUrlEnvVar))
	rocketToken := os.Getenv(rocketTokenEnvVar)

	if instanceKey == "" {
		return nil, fmt.Errorf("InstanceKey is required")
	}
	if rocketBaseUrl == "" {
		return nil, fmt.Errorf("%s environment variable is required", rocketUrlEnvVar)
	}
	if rocketToken == "" {
		return nil, fmt.Errorf("%s environment variable is required", rocketTokenEnvVar)
	}

	timeoutSeconds := cfg.RocketTimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = 5
	}

	cacheTtlSeconds := cfg.CacheTtlSeconds
	if cacheTtlSeconds <= 0 {
		cacheTtlSeconds = 30
	}

	maintenanceHtml, err := pages.ResolveFile(cfg.MaintenancePageFile, pages.DefaultMaintenanceHtml)
	if err != nil {
		return nil, fmt.Errorf("failed to load MaintenancePageFile: %w", err)
	}

	plugin := &RocketTraefikPlugin{
		logger:          logger,
		next:            next,
		mode:            mode,
		statusCode:      statusCode,
		instanceKey:     instanceKey,
		cacheTtl:        time.Duration(cacheTtlSeconds) * time.Second,
		rocketClient:    rocket.NewClient(rocketBaseUrl, rocketToken, time.Duration(timeoutSeconds)*time.Second, logger),
		maintenanceHtml: maintenanceHtml,
	}

	if mode == config.ModeFallback {
		fallbackHtml, err := pages.ResolveFile(cfg.FallbackPageFile, pages.DefaultFallbackHtml)
		if err != nil {
			return nil, fmt.Errorf("failed to load FallbackPageFile: %w", err)
		}
		plugin.fallbackHtml = fallbackHtml
	}

	logger.Log(logging.LevelInfo, "Configuration loaded successfully, starting middleware in %q mode...", mode)

	return plugin, nil
}
