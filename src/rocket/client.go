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

const identityHeaderName = "X-Rocket-Identity"

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
	BaseUrl        string
	IdentityHeader string
	Timeout        time.Duration
	HttpClient     *http.Client
	Logger         *logging.Logger
}

func NewClient(baseUrl string, identityHeader string, timeout time.Duration, logger *logging.Logger) *Client {
	return &Client{
		BaseUrl:        strings.TrimSuffix(baseUrl, "/"),
		IdentityHeader: identityHeader,
		Timeout:        timeout,
		HttpClient:     &http.Client{},
		Logger:         logger,
	}
}

func (c *Client) CheckMaintenance(ctx context.Context, appId string) (MaintenanceStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	requestUrl := c.BaseUrl + "/api/v1/apps/" + url.PathEscape(appId) + "/maintenance"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return MaintenanceStatus{}, fmt.Errorf("failed to build rocket request: %w", err)
	}
	req.Header.Set(identityHeaderName, c.IdentityHeader)

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
		MaintenanceEnabled bool   `json:"maintenanceEnabled"`
		Message            string `json:"message"`
		AllowBypass        bool   `json:"allowBypass"`
		BypassCode         string `json:"bypassCode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return MaintenanceStatus{}, fmt.Errorf("failed to decode rocket response: %w", err)
	}

	return MaintenanceStatus{
		Enabled:     parsed.MaintenanceEnabled,
		Message:     parsed.Message,
		AllowBypass: parsed.AllowBypass,
		BypassCode:  parsed.BypassCode,
	}, nil
}
