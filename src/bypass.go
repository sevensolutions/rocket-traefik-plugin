package src

import (
	"crypto/subtle"
	"net/http"
)

const bypassCookieName = "rocket_bypass"
const bypassCookieMaxAgeSeconds = 24 * 60 * 60 // 24h

// hasValidBypassCookie reports whether the request already carries a cookie matching the
// current bypass code (as returned by Rocket). The cookie value is the code itself, so
// Rocket rotating or clearing the code automatically invalidates any previously issued cookie.
func hasValidBypassCookie(req *http.Request, bypassCode string) bool {
	if bypassCode == "" {
		return false
	}

	cookie, err := req.Cookie(bypassCookieName)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(bypassCode)) == 1
}

func isValidBypassCode(submitted string, bypassCode string) bool {
	if bypassCode == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(submitted), []byte(bypassCode)) == 1
}

func setBypassCookie(rw http.ResponseWriter, req *http.Request, bypassCode string) {
	http.SetCookie(rw, &http.Cookie{
		Name:     bypassCookieName,
		Value:    bypassCode,
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
// code query param, so a page refresh after a successful bypass doesn't resubmit it.
func redirectStrippingBypassParam(rw http.ResponseWriter, req *http.Request) {
	target := req.URL.Path
	if target == "" {
		target = "/"
	}

	http.Redirect(rw, req, target, http.StatusSeeOther)
}
