package src

import (
	"context"
	"time"

	"github.com/sevensolutions/rocket-traefik-plugin/src/logging"
)

type maintenanceResult struct {
	enabled     bool
	message     string
	allowBypass bool
	bypassCode  string
}

// checkMaintenance returns the current maintenance status for this plugin instance's app,
// refreshing from Rocket when the cached value is missing or has expired.
//
// On a Rocket error, a stale cached value is served rather than failing open, so a transient
// Rocket outage doesn't briefly disable maintenance mode. Only when there has never been a
// successful check does this fail open (treat as not-in-maintenance).
func (p *RocketTraefikPlugin) checkMaintenance() maintenanceResult {
	p.cache.mu.Lock()
	defer p.cache.mu.Unlock()

	if p.cache.hasValue && time.Since(p.cache.fetchedAt) < p.cacheTtl {
		return p.cache.toResult()
	}

	status, err := p.rocketClient.CheckMaintenance(context.Background(), p.instanceKey)

	if err != nil {
		p.logger.Log(logging.LevelWarn, "Failed to fetch maintenance status from Rocket for instance %q: %s", p.instanceKey, err)

		if p.cache.hasValue {
			p.logger.Log(logging.LevelDebug, "Serving stale cached maintenance status for instance %q", p.instanceKey)
			return p.cache.toResult()
		}

		p.logger.Log(logging.LevelWarn, "No cached maintenance status for instance %q yet, failing open", p.instanceKey)
		return maintenanceResult{}
	}

	p.cache.hasValue = true
	p.cache.enabled = status.Enabled
	p.cache.message = status.Message
	p.cache.allowBypass = status.AllowBypass
	p.cache.bypassCode = status.BypassCode
	p.cache.fetchedAt = time.Now()

	return p.cache.toResult()
}
