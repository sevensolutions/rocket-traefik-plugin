package src

import (
	"net/http"
	"strconv"
	"strings"
)

// wantsJSON reports whether the caller expects a machine-readable (JSON) error response
// instead of the HTML maintenance/fallback page — e.g. a background fetch/XHR call made by a
// single-page app or API client, rather than a browser navigating to the page directly.
func wantsJSON(req *http.Request) bool {
	// Sent by jQuery and most other AJAX helpers on every request, but not by top-level
	// browser navigations.
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}

	return acceptPrefersJSON(req.Header.Get("Accept"))
}

// acceptPrefersJSON reports whether the Accept header explicitly favors JSON over HTML. A
// plain browser navigation always lists text/html (usually with q=1), so this only trips for
// callers that either omit text/html entirely or rank JSON above it.
func acceptPrefersJSON(accept string) bool {
	if accept == "" {
		return false
	}

	htmlQ, htmlOk := acceptQuality(accept, "text", "html")
	jsonQ, jsonOk := acceptQuality(accept, "application", "json")

	if !jsonOk {
		return false
	}
	if !htmlOk {
		return true
	}

	return jsonQ > htmlQ
}

// acceptQuality scans a comma-separated Accept header for the entry (exact, subtype wildcard,
// or full wildcard) that best matches typ/subtype, returning its q value (defaulting to 1 when
// unspecified). matched is false if nothing in the header applies to typ/subtype at all.
func acceptQuality(accept string, typ string, subtype string) (quality float64, matched bool) {
	best := -1.0

	for _, entry := range strings.Split(accept, ",") {
		segments := strings.Split(entry, ";")
		mediaRange := strings.TrimSpace(segments[0])
		if mediaRange == "" {
			continue
		}

		if !acceptRangeMatches(mediaRange, typ, subtype) {
			continue
		}

		q := 1.0
		for _, param := range segments[1:] {
			param = strings.TrimSpace(param)
			value, ok := strings.CutPrefix(param, "q=")
			if !ok {
				continue
			}
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				q = parsed
			}
		}

		if q > best {
			best = q
		}
	}

	if best < 0 {
		return 0, false
	}

	return best, true
}

// acceptRangeMatches reports whether an Accept media-range entry (e.g. "*/*", "application/*",
// "application/json") covers typ/subtype.
func acceptRangeMatches(mediaRange string, typ string, subtype string) bool {
	if mediaRange == "*/*" {
		return true
	}

	rangeType, rangeSubtype, ok := strings.Cut(mediaRange, "/")
	if !ok {
		return false
	}

	if rangeType != "*" && !strings.EqualFold(rangeType, typ) {
		return false
	}
	if rangeSubtype != "*" && !strings.EqualFold(rangeSubtype, subtype) {
		return false
	}

	return true
}
