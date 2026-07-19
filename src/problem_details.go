package src

import (
	"encoding/json"
	"net/http"
)

// problemDetails is an RFC 7807 "Problem Details for HTTP APIs" response body, used instead
// of the HTML maintenance/fallback page when the caller asked for JSON (see wantsJSON).
type problemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func writeProblemDetails(rw http.ResponseWriter, statusCode int, title string, detail string) {
	body, err := json.Marshal(problemDetails{
		Type:   "about:blank",
		Title:  title,
		Status: statusCode,
		Detail: detail,
	})
	if err != nil {
		body = []byte(`{"type":"about:blank"}`)
	}

	rw.Header().Set("Content-Type", "application/problem+json")
	rw.WriteHeader(statusCode)
	rw.Write(body)
}
