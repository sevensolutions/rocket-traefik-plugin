package src

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sevensolutions/rocket-traefik-plugin/src/config"
	"github.com/sevensolutions/rocket-traefik-plugin/src/logging"
	"github.com/sevensolutions/rocket-traefik-plugin/src/pages"
	"github.com/sevensolutions/rocket-traefik-plugin/src/rocket"
	"github.com/sevensolutions/rocket-traefik-plugin/src/utils"
)

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

	plugin := &RocketTraefikPlugin{
		logger:     logger,
		next:       next,
		mode:       mode,
		statusCode: statusCode,
	}

	if mode == config.ModeFallback {
		plugin.fallbackHtml = pages.Resolve(cfg.FallbackPageContent, pages.DefaultFallbackHtml)
	} else {
		appId := utils.ExpandEnvironmentVariableString(strings.TrimSpace(cfg.AppId))
		rocketBaseUrl := utils.ExpandEnvironmentVariableString(strings.TrimSpace(cfg.RocketBaseUrl))
		identityHeader := utils.ExpandEnvironmentVariableString(cfg.RocketIdentityHeader)

		if appId == "" {
			return nil, fmt.Errorf("AppId is required when Mode is %q", config.ModeMaintenance)
		}
		if rocketBaseUrl == "" {
			return nil, fmt.Errorf("RocketBaseUrl is required when Mode is %q", config.ModeMaintenance)
		}

		timeoutSeconds := cfg.RocketTimeoutSeconds
		if timeoutSeconds <= 0 {
			timeoutSeconds = 5
		}

		cacheTtlSeconds := cfg.CacheTtlSeconds
		if cacheTtlSeconds <= 0 {
			cacheTtlSeconds = 30
		}

		plugin.appId = appId
		plugin.cacheTtl = time.Duration(cacheTtlSeconds) * time.Second
		plugin.rocketClient = rocket.NewClient(rocketBaseUrl, identityHeader, time.Duration(timeoutSeconds)*time.Second, logger)
		plugin.maintenanceHtml = pages.Resolve(cfg.MaintenancePageContent, pages.DefaultMaintenanceHtml)
	}

	logger.Log(logging.LevelInfo, "Configuration loaded successfully, starting middleware in %q mode...", mode)

	return plugin, nil
}
