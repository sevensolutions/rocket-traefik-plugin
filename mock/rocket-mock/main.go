// rocket-mock is a tiny standalone HTTP server that mimics the Rocket control-plane
// maintenance-status endpoint, for local development against the plugin's Maintenance
// Mode. It is a regular Go program built/run with the standard toolchain — it is not
// part of the Traefik plugin (package src) and is never interpreted by Yaegi.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

type maintenanceStatus struct {
	MaintenanceEnabled bool   `json:"maintenanceEnabled"`
	Message            string `json:"message"`
	BypassCode         string `json:"bypassCode"`
}

type store struct {
	mu    sync.Mutex
	byApp map[string]maintenanceStatus
}

func (s *store) get(appId string) maintenanceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.byApp[appId] // zero value = not in maintenance
}

func (s *store) set(appId string, status maintenanceStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.byApp[appId] = status
}

func main() {
	s := &store{byApp: map[string]maintenanceStatus{}}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/apps/", func(rw http.ResponseWriter, req *http.Request) {
		appId, ok := parseAppId(req.URL.Path)
		if !ok {
			http.NotFound(rw, req)
			return
		}

		switch req.Method {
		case http.MethodGet:
			identity := req.Header.Get("X-Rocket-Identity")
			log.Printf("GET maintenance status for app=%q identity=%q", appId, identity)

			writeJson(rw, http.StatusOK, s.get(appId))

		case http.MethodPost, http.MethodPut:
			var status maintenanceStatus
			if err := json.NewDecoder(req.Body).Decode(&status); err != nil {
				http.Error(rw, "invalid JSON body", http.StatusBadRequest)
				return
			}

			s.set(appId, status)
			log.Printf("SET maintenance status for app=%q -> %+v", appId, status)

			writeJson(rw, http.StatusOK, status)

		default:
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("rocket-mock listening on :80")
	log.Fatal(http.ListenAndServe(":80", mux))
}

// parseAppId extracts {appId} from /api/v1/apps/{appId}/maintenance.
func parseAppId(path string) (string, bool) {
	const prefix = "/api/v1/apps/"
	const suffix = "/maintenance"

	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return "", false
	}

	appId := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	if appId == "" {
		return "", false
	}

	return appId, true
}

func writeJson(rw http.ResponseWriter, statusCode int, body any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(body)
}
