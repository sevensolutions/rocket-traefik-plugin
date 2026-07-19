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

// maintenanceStatus is the mock's internal, flat representation of a maintenance
// window, also used as the JSON body for the test-control POST/PUT endpoint.
type maintenanceStatus struct {
	MaintenanceEnabled bool   `json:"maintenanceEnabled"`
	Message            string `json:"message"`
	AllowBypass        bool   `json:"allowBypass"`
	BypassCode         string `json:"bypassCode"`
}

// maintenanceMode and statusResponse mirror Rocket's real GET response shape:
// TraefikInstanceStatusResponse { maintenanceMode: TraefikInstanceMaintenanceMode | null }.
type maintenanceMode struct {
	IsEnabled   bool   `json:"isEnabled"`
	Message     string `json:"message"`
	AllowBypass bool   `json:"allowBypass"`
	BypassCode  string `json:"bypassCode"`
}

type statusResponse struct {
	MaintenanceMode *maintenanceMode `json:"maintenanceMode"`
}

type store struct {
	mu         sync.Mutex
	byInstance map[string]maintenanceStatus
}

func (s *store) get(instanceKey string) maintenanceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.byInstance[instanceKey] // zero value = not in maintenance
}

func (s *store) set(instanceKey string, status maintenanceStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.byInstance[instanceKey] = status
}

func main() {
	s := &store{byInstance: map[string]maintenanceStatus{}}

	mux := http.NewServeMux()

	// GET matches Rocket's real endpoint: api/v1/ingress/traefik/instances/{instanceKey}/status.
	// POST/PUT is mock-only test-control, used by taskfile.yml to seed maintenance state.
	mux.HandleFunc("/api/v1/ingress/traefik/instances/", func(rw http.ResponseWriter, req *http.Request) {
		instanceKey, ok := parseInstanceKey(req.URL.Path)
		if !ok {
			http.NotFound(rw, req)
			return
		}

		switch req.Method {
		case http.MethodGet:
			identity := req.Header.Get("X-Rocket-Identity")
			log.Printf("GET status for instance=%q identity=%q", instanceKey, identity)

			status := s.get(instanceKey)

			writeJson(rw, http.StatusOK, statusResponse{
				MaintenanceMode: &maintenanceMode{
					IsEnabled:   status.MaintenanceEnabled,
					Message:     status.Message,
					AllowBypass: status.AllowBypass,
					BypassCode:  status.BypassCode,
				},
			})

		case http.MethodPost, http.MethodPut:
			var status maintenanceStatus
			if err := json.NewDecoder(req.Body).Decode(&status); err != nil {
				http.Error(rw, "invalid JSON body", http.StatusBadRequest)
				return
			}

			s.set(instanceKey, status)
			log.Printf("SET maintenance status for instance=%q -> %+v", instanceKey, status)

			writeJson(rw, http.StatusOK, status)

		default:
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("rocket-mock listening on :80")
	log.Fatal(http.ListenAndServe(":80", mux))
}

// parseInstanceKey extracts {instanceKey} from
// /api/v1/ingress/traefik/instances/{instanceKey}/status.
func parseInstanceKey(path string) (string, bool) {
	const prefix = "/api/v1/ingress/traefik/instances/"
	const suffix = "/status"

	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return "", false
	}

	instanceKey := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	if instanceKey == "" {
		return "", false
	}

	return instanceKey, true
}

func writeJson(rw http.ResponseWriter, statusCode int, body any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(body)
}
