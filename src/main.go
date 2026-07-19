package src

import (
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sevensolutions/rocket-traefik-plugin/src/config"
	"github.com/sevensolutions/rocket-traefik-plugin/src/logging"
	"github.com/sevensolutions/rocket-traefik-plugin/src/pages"
	"github.com/sevensolutions/rocket-traefik-plugin/src/rocket"
)

type maintenanceCache struct {
	mu          sync.Mutex
	hasValue    bool
	enabled     bool
	message     string
	allowBypass bool
	bypassCode  string
	fetchedAt   time.Time
}

func (c *maintenanceCache) toResult() maintenanceResult {
	return maintenanceResult{
		enabled:     c.enabled,
		message:     c.message,
		allowBypass: c.allowBypass,
		bypassCode:  c.bypassCode,
	}
}

type RocketTraefikPlugin struct {
	logger *logging.Logger
	next   http.Handler

	mode       string
	statusCode int

	fallbackHtml string

	appId           string
	maintenanceHtml string
	cacheTtl        time.Duration
	rocketClient    *rocket.Client
	cache           maintenanceCache
}

func (p *RocketTraefikPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if p.mode == config.ModeFallback {
		p.writeHtmlPage(rw, p.fallbackHtml)
		return
	}

	result := p.checkMaintenance()

	if !result.enabled {
		clearBypassCookieIfPresent(rw, req)
		p.next.ServeHTTP(rw, req)
		return
	}

	if hasValidBypassCookie(req, result) {
		p.next.ServeHTTP(rw, req)
		return
	}

	if result.allowBypass {
		if submitted := req.URL.Query().Get(pages.BypassQueryParam); submitted != "" {
			if isValidBypassRequest(submitted, result) {
				setBypassCookie(rw, req, result)
				redirectStrippingBypassParam(rw, req)
				return
			}

			p.writeMaintenancePage(rw, result, true)
			return
		}
	}

	p.writeMaintenancePage(rw, result, false)
}

func (p *RocketTraefikPlugin) writeHtmlPage(rw http.ResponseWriter, htmlContent string) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(p.statusCode)
	rw.Write([]byte(htmlContent))
}

func (p *RocketTraefikPlugin) writeMaintenancePage(rw http.ResponseWriter, result maintenanceResult, invalidBypassCode bool) {
	message := result.message
	if message == "" {
		message = pages.DefaultMaintenanceMessage
	}

	bypassForm := ""
	if result.allowBypass {
		bypassForm = pages.RenderBypassForm(result.bypassCode != "", invalidBypassCode)
	}

	body := strings.ReplaceAll(p.maintenanceHtml, "{{Message}}", html.EscapeString(message))
	body = strings.ReplaceAll(body, "{{BypassForm}}", bypassForm)

	p.writeHtmlPage(rw, body)
}
