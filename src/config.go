package src

import (
	"context"
	"net/http"

	"github.com/sevensolutions/rocket-traefik-plugin/src/config"
	"github.com/sevensolutions/rocket-traefik-plugin/src/logging"
	"github.com/sevensolutions/rocket-traefik-plugin/src/utils"
)

func CreateConfig() *config.Config {
	return &config.Config{
		LogLevel: logging.LevelWarn,
	}
}

// Will be called by traefik
func New(uctx context.Context, next http.Handler, cfg *config.Config, name string) (http.Handler, error) {
	cfg.LogLevel = utils.ExpandEnvironmentVariableString(cfg.LogLevel)

	logger := logging.CreateLogger(cfg.LogLevel)

	logger.Log(logging.LevelInfo, "Loading Configuration...")

	logger.Log(logging.LevelInfo, "Configuration loaded successfully, starting middleware...")

	return &RocketTraefikPlugin{
		logger: logger,
		next:   next,
	}, nil
}
