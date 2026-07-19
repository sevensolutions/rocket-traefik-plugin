package src

import (
	"crypto/subtle"
	"net/http"
)

const bypassCookieName = "rocket_bypass"
const bypassCookieMaxAgeSeconds = 24 * 60 * 60 // 24h

// openBypassCookieValue is the cookie value used when bypass is open (AllowBypass=true, no
// BypassCode configured) — there's no visitor-supplied secret to store instead.
const openBypassCookieValue = "granted"

// expectedBypassCookieValue returns the cookie value that should be issued/expected for the
// given maintenance result. Only meaningful when result.allowBypass is true.
func expectedBypassCookieValue(result maintenanceResult) string {
	if result.bypassCode != "" {
		return result.bypassCode
	}

	return openBypassCookieValue
}

// hasValidBypassCookie reports whether the request already carries a cookie matching what's
// currently required (or open) for this app, as returned by Rocket. Because the expected value
// is derived from the live result, Rocket rotating/clearing a code or disabling bypass entirely
// invalidates any previously issued cookie automatically.
func hasValidBypassCookie(req *http.Request, result maintenanceResult) bool {
	if !result.allowBypass {
		return false
	}

	cookie, err := req.Cookie(bypassCookieName)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(expectedBypassCookieValue(result))) == 1
}

// isValidBypassRequest checks a submitted query value against what's required to grant bypass:
// the exact code when one is configured, or just a non-empty click when bypass is open.
func isValidBypassRequest(submitted string, result maintenanceResult) bool {
	if submitted == "" {
		return false
	}

	if result.bypassCode == "" {
		return true
	}

	return subtle.ConstantTimeCompare([]byte(submitted), []byte(result.bypassCode)) == 1
}

func setBypassCookie(rw http.ResponseWriter, req *http.Request, result maintenanceResult) {
	http.SetCookie(rw, &http.Cookie{
		Name:     bypassCookieName,
		Value:    expectedBypassCookieValue(result),
		Path:     "/",
		MaxAge:   bypassCookieMaxAgeSeconds,
		HttpOnly: true,
		Secure:   req.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

// clearBypassCookieIfPresent removes the bypass cookie once maintenance mode is no longer
// enabled, so a visitor's bypass doesn't silently carry over into a future maintenance window.
func clearBypassCookieIfPresent(rw http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie(bypassCookieName); err != nil {
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     bypassCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   req.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

// redirectStrippingBypassParam sends the browser back to the current path without the bypass
// query param, so a page refresh after a successful bypass doesn't resubmit it.
func redirectStrippingBypassParam(rw http.ResponseWriter, req *http.Request) {
	target := req.URL.Path
	if target == "" {
		target = "/"
	}

	http.Redirect(rw, req, target, http.StatusSeeOther)
}
