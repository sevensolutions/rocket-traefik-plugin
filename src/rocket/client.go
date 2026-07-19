package rocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sevensolutions/rocket-traefik-plugin/src/logging"
)

const AuthTokenHeaderName = "X-Rocket-Identity"

type MaintenanceStatus struct {
	Enabled bool
	Message string
	// AllowBypass controls whether the maintenance page offers a way to skip maintenance mode.
	AllowBypass bool
	// BypassCode is optional even when AllowBypass is true: if set, visitors must submit this
	// exact code; if empty, bypass is open and a single click grants it. Rocket controls both
	// values entirely.
	BypassCode string
}

type Client struct {
	BaseUrl    string
	AuthToken  string
	Timeout    time.Duration
	HttpClient *http.Client
	Logger     *logging.Logger
}

func NewClient(baseUrl string, authToken string, timeout time.Duration, logger *logging.Logger) *Client {
	return &Client{
		BaseUrl:    strings.TrimSuffix(baseUrl, "/"),
		AuthToken:  authToken,
		Timeout:    timeout,
		HttpClient: &http.Client{},
		Logger:     logger,
	}
}

func (c *Client) CheckMaintenance(ctx context.Context, instanceKey string) (MaintenanceStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	requestUrl := c.BaseUrl + "/api/v1/ingress/traefik/instances/" + url.PathEscape(instanceKey) + "/status"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return MaintenanceStatus{}, fmt.Errorf("failed to build rocket request: %w", err)
	}
	req.Header.Set(AuthTokenHeaderName, c.AuthToken)

	c.Logger.Log(logging.LevelDebug, "Checking maintenance status with Rocket: %s", requestUrl)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return MaintenanceStatus{}, fmt.Errorf("rocket request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MaintenanceStatus{}, fmt.Errorf("rocket returned unexpected status %d", resp.StatusCode)
	}

	var parsed struct {
		MaintenanceMode *struct {
			IsEnabled   bool   `json:"isEnabled"`
			Message     string `json:"message"`
			AllowBypass bool   `json:"allowBypass"`
			BypassCode  string `json:"bypassCode"`
		} `json:"maintenanceMode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return MaintenanceStatus{}, fmt.Errorf("failed to decode rocket response: %w", err)
	}

	// A nil MaintenanceMode means Rocket has no maintenance mode configured for this
	// instance at all, which is equivalent to it being disabled.
	if parsed.MaintenanceMode == nil {
		return MaintenanceStatus{}, nil
	}

	return MaintenanceStatus{
		Enabled:     parsed.MaintenanceMode.IsEnabled,
		Message:     parsed.MaintenanceMode.Message,
		AllowBypass: parsed.MaintenanceMode.AllowBypass,
		BypassCode:  parsed.MaintenanceMode.BypassCode,
	}, nil
}
