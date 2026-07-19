package pages

import (
	"fmt"
	"os"

	"github.com/sevensolutions/rocket-traefik-plugin/src/utils"
)

const DefaultMaintenanceMessage = "This application is currently undergoing scheduled maintenance. Please check back soon."

// BypassQueryParam is the query/form field name used to submit a maintenance bypass code.
const BypassQueryParam = "rocket_bypass_code"

const DefaultFallbackHtml = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Application Unavailable</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f5f5f7; color: #1d1d1f; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
  .card { text-align: center; padding: 2.5rem; max-width: 28rem; }
  h1 { font-size: 1.5rem; margin-bottom: 0.5rem; }
  p { color: #6e6e73; }
</style>
</head>
<body>
<div class="card">
  <h1>Application Unavailable</h1>
  <p>This application is currently unavailable. Please try again shortly.</p>
</div>
</body>
</html>
`

const DefaultMaintenanceHtml = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Under Maintenance</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f5f5f7; color: #1d1d1f; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
  .card { text-align: center; padding: 2.5rem; max-width: 28rem; }
  h1 { font-size: 1.5rem; margin-bottom: 0.5rem; }
  p { color: #6e6e73; }
</style>
</head>
<body>
<div class="card">
  <h1>Under Maintenance</h1>
  <p>{{Message}}</p>
  {{BypassForm}}
</div>
</body>
</html>
`

const bypassButtonHtml = `<p style="margin-top: 1.5rem;">
  <a href="?%s=1" style="display:inline-block; padding:0.5rem 1rem; border-radius:0.375rem; background:#0071e3; color:#fff; text-decoration:none;">Continue anyway</a>
</p>`

const bypassFormHtml = `<form method="get" style="margin-top: 1.5rem;">
  <label for="rocket_bypass_code" style="display:block; font-size: 0.85rem; color:#6e6e73; margin-bottom:0.4rem;">Have a bypass code?</label>
  <input type="text" id="rocket_bypass_code" name="%s" placeholder="Bypass code" style="padding:0.5rem; border-radius:0.375rem; border:1px solid #d2d2d7;">
  <button type="submit" style="padding:0.5rem 1rem; border-radius:0.375rem; border:none; background:#0071e3; color:#fff; margin-left:0.5rem;">Continue anyway</button>
  %s
</form>`

const bypassFormInvalidCodeNotice = `<p style="color:#d70015; font-size:0.85rem; margin-top:0.5rem;">Invalid bypass code.</p>`

// RenderBypassForm builds the bypass UI shown on the maintenance page. When requiresCode is
// false, bypass is open and a single link grants it. When true, a code-entry form is shown
// instead, with an "invalid code" notice when invalidCode is set.
func RenderBypassForm(requiresCode bool, invalidCode bool) string {
	if !requiresCode {
		return fmt.Sprintf(bypassButtonHtml, BypassQueryParam)
	}

	notice := ""
	if invalidCode {
		notice = bypassFormInvalidCodeNotice
	}

	return fmt.Sprintf(bypassFormHtml, BypassQueryParam, notice)
}

// ResolveFile returns the contents of the file at path (after ${VAR} expansion) if path is
// set, otherwise defaultContent.
func ResolveFile(path string, defaultContent string) (string, error) {
	if path == "" {
		return defaultContent, nil
	}

	expandedPath := utils.ExpandEnvironmentVariableString(path)

	content, err := os.ReadFile(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %q: %w", expandedPath, err)
	}

	return string(content), nil
}
