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
	mu         sync.Mutex
	hasValue   bool
	enabled    bool
	message    string
	bypassCode string
	fetchedAt  time.Time
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

	enabled, message, bypassCode := p.checkMaintenance()

	if !enabled {
		clearBypassCookieIfPresent(rw, req)
		p.next.ServeHTTP(rw, req)
		return
	}

	if hasValidBypassCookie(req, bypassCode) {
		p.next.ServeHTTP(rw, req)
		return
	}

	bypassAvailable := bypassCode != ""

	if bypassAvailable {
		if code := req.URL.Query().Get(pages.BypassQueryParam); code != "" {
			if isValidBypassCode(code, bypassCode) {
				setBypassCookie(rw, req, bypassCode)
				redirectStrippingBypassParam(rw, req)
				return
			}

			p.writeMaintenancePage(rw, message, true, bypassAvailable)
			return
		}
	}

	p.writeMaintenancePage(rw, message, false, bypassAvailable)
}

func (p *RocketTraefikPlugin) writeHtmlPage(rw http.ResponseWriter, htmlContent string) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(p.statusCode)
	rw.Write([]byte(htmlContent))
}

func (p *RocketTraefikPlugin) writeMaintenancePage(rw http.ResponseWriter, message string, invalidBypassCode bool, bypassAvailable bool) {
	if message == "" {
		message = pages.DefaultMaintenanceMessage
	}

	bypassForm := ""
	if bypassAvailable {
		bypassForm = pages.RenderBypassForm(invalidBypassCode)
	}

	body := strings.ReplaceAll(p.maintenanceHtml, "{{Message}}", html.EscapeString(message))
	body = strings.ReplaceAll(body, "{{BypassForm}}", bypassForm)

	p.writeHtmlPage(rw, body)
}
