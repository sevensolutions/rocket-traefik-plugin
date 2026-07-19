package pages

import (
	"fmt"
	"os"

	"github.com/sevensolutions/rocket-traefik-plugin/src/utils"
)

const DefaultMaintenanceMessage = "This application is currently undergoing scheduled maintenance. Please check back soon."

const DefaultFallbackMessage = "This application is currently unavailable. Please try again shortly."

// BypassQueryParam is the query/form field name used to submit a maintenance bypass code.
const BypassQueryParam = "rocket_bypass_code"

// rocketLogoSvg is the Rocket wordmark, inlined so the pages render crisply with zero extra
// requests. width/height are intentionally omitted (only viewBox is kept) so the .logo CSS
// class controls the rendered size instead of fighting the SVG's own presentation attributes.
const rocketLogoSvg = `<svg class="logo" viewBox="0 0 249 91" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:2;">
    <g transform="matrix(0.486328,0,0,0.177734,0,0)">
        <clipPath id="_clip1">
            <rect x="0" y="0" width="512" height="512"/>
        </clipPath>
        <g clip-path="url(#_clip1)">
            <g transform="matrix(1.49398,0,0,1.71429,-97.0502,-130.523)">
                <g>
                    <g transform="matrix(0.196165,0,0,0.436585,307.303,76.128)">
                        <path d="M188.9,372L138.5,321.6C157.1,278.9 200.2,183.9 233.6,134.6C304.6,30.1 409,24.6 475.7,36.3C487.4,103 481.9,207.4 377.3,278.3C327.9,311.8 231.8,353.9 188.9,372ZM109,309.2C103.8,321.1 106.5,334.9 115.7,344.1L166.4,394.8C175.5,403.9 189.1,406.7 200.9,401.7C207.4,399 215.2,395.7 223.9,391.9L224,496C224,501.5 226.9,506.7 231.6,509.6C236.3,512.5 242.2,512.8 247.2,510.3L348.7,459.6C370.4,448.8 384.1,426.6 384.1,402.4L384.1,312.2C388.1,309.7 391.8,307.3 395.4,304.9C516.1,222.9 520.1,100.9 506.7,28.1C504.6,16.5 495.5,7.5 483.9,5.3C411.1,-8.1 289.1,-4.1 207.2,116.7C204.8,120.3 202.3,124 199.9,128L109.7,128C85.5,128 63.3,141.7 52.5,163.4L1.7,264.8C-0.8,269.8 -0.5,275.7 2.4,280.4C5.3,285.1 10.5,288 16,288L118.5,288C114.9,296 111.7,303.2 109.1,309.2L109,309.2ZM256,470.1L256,377.6C286.3,363.9 321.4,347.3 352,330.6L352,402.3C352,414.4 345.2,425.5 334.3,430.9L256,470.1ZM109.6,160L181.2,160C164.3,190.7 147.2,225.8 133.1,256L41.9,256L81,177.7C86.4,166.9 97.5,160 109.6,160ZM392,144C392,157.166 381.166,168 368,168C354.834,168 344,157.166 344,144C344,130.834 354.834,120 368,120C381.166,120 392,130.834 392,144ZM368,88C337.279,88 312,113.279 312,144C312,174.721 337.279,200 368,200C398.721,200 424,174.721 424,144C424,113.279 398.721,88 368,88Z" style="fill:rgb(13,109,234);fill-rule:nonzero;"/>
                    </g>
                    <g transform="matrix(2.39726,0,0,3.90131,-379.16,-82.3931)" fill="currentColor">
                        <path d="M197.606,79.037C199.495,79.037 201.068,79.363 202.327,80.014C203.586,80.665 204.495,81.78 205.055,83.358C205.614,84.935 205.762,87.129 205.499,89.939C205.345,91.63 205.062,93.106 204.648,94.368C204.235,95.629 203.65,96.633 202.893,97.379C202.136,98.126 201.155,98.604 199.952,98.814L199.916,99.107C200.357,99.252 200.768,99.54 201.149,99.97C201.531,100.401 201.86,100.988 202.137,101.732C202.415,102.475 202.61,103.421 202.722,104.57L204.114,115.806L199.495,115.806L198.262,104.48C198.137,103.214 197.884,102.352 197.506,101.895C197.127,101.438 196.578,101.21 195.858,101.21C194.254,101.175 192.907,101.151 191.819,101.138C190.73,101.126 189.822,101.108 189.096,101.086C188.369,101.063 187.729,101.042 187.175,101.022L187.621,96.043L196.101,96.043C197.176,96.018 198.038,95.826 198.688,95.466C199.338,95.107 199.845,94.477 200.21,93.576C200.575,92.676 200.834,91.416 200.988,89.796C201.134,88.357 201.126,87.23 200.962,86.415C200.798,85.6 200.442,85.028 199.893,84.698C199.344,84.369 198.558,84.204 197.535,84.204C195.349,84.204 193.503,84.224 191.997,84.264C190.491,84.304 189.472,84.339 188.942,84.37L188.743,79.526C189.7,79.375 190.609,79.265 191.469,79.195C192.329,79.125 193.245,79.081 194.218,79.063C195.191,79.046 196.32,79.037 197.606,79.037ZM193.136,79.526L189.655,115.806L185.262,115.806L188.743,79.526L193.136,79.526Z" style="fill-rule:nonzero;"/>
                        <path d="M217.184,89.34C219.27,89.34 220.934,89.795 222.176,90.704C223.418,91.613 224.258,93.079 224.698,95.102C225.137,97.126 225.19,99.828 224.855,103.209C224.544,106.454 224.022,109.049 223.288,110.995C222.555,112.94 221.527,114.34 220.206,115.193C218.884,116.047 217.181,116.474 215.095,116.474C213.032,116.474 211.374,116.022 210.121,115.118C208.867,114.214 208.013,112.752 207.558,110.734C207.103,108.716 207.034,106.039 207.352,102.703C207.67,99.438 208.194,96.83 208.924,94.88C209.654,92.929 210.687,91.519 212.024,90.647C213.361,89.776 215.081,89.34 217.184,89.34ZM217.143,94.19C216.038,94.19 215.144,94.44 214.462,94.94C213.78,95.44 213.242,96.336 212.848,97.63C212.455,98.923 212.136,100.782 211.894,103.209C211.683,105.394 211.634,107.101 211.748,108.328C211.862,109.556 212.196,110.412 212.75,110.897C213.304,111.381 214.113,111.624 215.178,111.624C216.293,111.624 217.188,111.361 217.863,110.836C218.538,110.311 219.07,109.399 219.461,108.099C219.851,106.798 220.159,105 220.385,102.703C220.603,100.448 220.654,98.71 220.536,97.49C220.419,96.27 220.088,95.414 219.544,94.925C219.001,94.435 218.2,94.19 217.143,94.19Z" style="fill-rule:nonzero;"/>
                        <path d="M237.199,89.34C238.195,89.34 239.278,89.451 240.446,89.672C241.615,89.893 242.651,90.279 243.554,90.83L242.906,94.757C241.975,94.682 241.047,94.632 240.122,94.607C239.197,94.583 238.428,94.57 237.815,94.57C236.605,94.57 235.616,94.803 234.848,95.268C234.081,95.733 233.488,96.567 233.068,97.772C232.649,98.976 232.337,100.715 232.132,102.989C231.921,105.273 231.904,107.004 232.081,108.181C232.258,109.358 232.664,110.16 233.298,110.585C233.933,111.01 234.83,111.222 235.99,111.222C236.4,111.222 236.907,111.207 237.51,111.177C238.114,111.147 238.762,111.099 239.456,111.034C240.151,110.969 240.817,110.872 241.457,110.742L241.765,114.931C240.78,115.463 239.712,115.847 238.559,116.083C237.406,116.319 236.294,116.437 235.223,116.437C233.169,116.437 231.524,115.976 230.287,115.054C229.051,114.133 228.215,112.658 227.778,110.629C227.342,108.601 227.279,105.937 227.59,102.636C227.918,99.406 228.437,96.816 229.146,94.866C229.855,92.915 230.859,91.507 232.158,90.64C233.457,89.773 235.137,89.34 237.199,89.34Z" style="fill-rule:nonzero;"/>
                        <path d="M252.379,79L250.837,95.769C250.734,96.93 250.574,98.03 250.357,99.069C250.14,100.109 249.862,101.196 249.523,102.333C249.571,102.947 249.629,103.607 249.698,104.312C249.766,105.017 249.804,105.665 249.81,106.254L248.837,115.806L244.516,115.806L248.033,79L252.379,79ZM262.51,90.001C261.971,92.652 261.205,94.91 260.214,96.776C259.223,98.642 258.094,100.216 256.828,101.498C255.561,102.78 254.244,103.866 252.876,104.754C251.509,105.643 250.196,106.448 248.938,107.169L248.577,102.629C250.014,101.775 251.387,100.721 252.697,99.47C254.006,98.218 255.136,96.789 256.087,95.182C257.038,93.575 257.692,91.848 258.05,90.001L262.51,90.001ZM255.857,100.949C256.314,101.5 256.705,102.121 257.031,102.813C257.356,103.504 257.629,104.211 257.848,104.933L261.241,115.806L256.689,115.806L252.758,102.666L255.857,100.949Z" style="fill-rule:nonzero;"/>
                        <path d="M273.438,89.34C275.467,89.34 277.036,89.667 278.144,90.321C279.252,90.975 279.997,91.946 280.379,93.233C280.76,94.521 280.843,96.109 280.628,97.998C280.451,99.626 280.092,100.915 279.554,101.866C279.015,102.817 278.269,103.528 277.315,103.998C276.361,104.469 275.156,104.792 273.699,104.967L264.894,106.093L265.258,102.197L273.103,101.07C273.81,100.975 274.386,100.816 274.831,100.594C275.276,100.371 275.626,100.009 275.881,99.507C276.135,99.005 276.312,98.308 276.411,97.418C276.519,96.384 276.47,95.626 276.262,95.144C276.054,94.662 275.695,94.356 275.183,94.226C274.672,94.096 274.004,94.036 273.178,94.046C272.353,94.046 271.653,94.161 271.08,94.391C270.507,94.622 270.028,95.031 269.643,95.618C269.257,96.206 268.94,97.04 268.691,98.12C268.442,99.201 268.22,100.608 268.026,102.344C267.766,104.859 267.739,106.761 267.943,108.051C268.148,109.341 268.577,110.204 269.23,110.642C269.883,111.079 270.761,111.298 271.864,111.298C272.468,111.298 273.164,111.271 273.952,111.219C274.741,111.166 275.538,111.09 276.345,110.99C277.153,110.89 277.89,110.788 278.556,110.682L278.855,114.758C278.187,115.169 277.409,115.501 276.521,115.752C275.632,116.003 274.735,116.18 273.831,116.282C272.926,116.385 272.089,116.437 271.32,116.437C269.103,116.437 267.379,115.924 266.15,114.9C264.92,113.875 264.104,112.32 263.702,110.235C263.301,108.149 263.229,105.542 263.489,102.413C263.714,99.808 264.071,97.655 264.56,95.955C265.05,94.254 265.689,92.926 266.479,91.969C267.269,91.013 268.239,90.336 269.387,89.937C270.536,89.539 271.886,89.34 273.438,89.34Z" style="fill-rule:nonzero;"/>
                        <path d="M291.238,82.628L288.777,108.166C288.672,109.231 288.74,109.97 288.982,110.385C289.224,110.8 289.723,111.007 290.477,111.007L293.027,111.007L293.218,115.504C292.787,115.695 292.282,115.853 291.701,115.978C291.121,116.104 290.56,116.198 290.017,116.261C289.475,116.324 289.037,116.355 288.705,116.355C287.108,116.355 285.927,115.721 285.16,114.454C284.393,113.186 284.121,111.433 284.343,109.195L286.947,82.628L291.238,82.628ZM295.729,90.001L295.298,94.421L282.626,94.421L282.985,90.264L286.441,90.001L295.729,90.001Z" style="fill-rule:nonzero;"/>
                    </g>
                </g>
            </g>
        </g>
    </g>
</svg>`

// Feather-style line icons (24x24, stroke=currentColor), used inside .icon-badge.
const alertTriangleIconSvg = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`
const toolIconSvg = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76Z"/></svg>`

const pageStyle = `<style>
  :root { color-scheme: light dark; }
  * { box-sizing: border-box; }
  body {
    margin: 0; min-height: 100vh; padding: 1.5rem;
    display: flex; align-items: center; justify-content: center;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    background: linear-gradient(180deg, #f4f6fb 0%, #e9edf5 100%);
  }
  @media (prefers-color-scheme: dark) {
    body { background: linear-gradient(180deg, #0b0d12 0%, #14171f 100%); }
  }
  .page { width: 100%; max-width: 26rem; display: flex; flex-direction: column; align-items: center; }
  .card {
    width: 100%; text-align: center;
    background: #ffffff; border-radius: 20px; padding: 3.5rem 3rem;
    box-shadow: 0 1px 2px rgba(16,24,40,0.06), 0 12px 32px -8px rgba(16,24,40,0.18);
  }
  .logo { width: 116px; height: auto; display: block; margin: 0 0 2rem; color: #344054; }
  @media (prefers-color-scheme: dark) {
    .logo { color: #e4e7ec; }
  }
  .icon-badge {
    width: 56px; height: 56px; border-radius: 16px; margin: 0 auto 1.25rem;
    display: flex; align-items: center; justify-content: center;
    background: rgba(13,109,234,0.1); color: #0d6dea;
  }
  .icon-badge svg { width: 26px; height: 26px; }
  h1 { font-size: 1.375rem; font-weight: 650; letter-spacing: -0.01em; margin: 0 0 0.5rem; color: #101828; }
  .message { color: #667085; font-size: 0.95rem; line-height: 1.55; margin: 0; }
  .bypass { margin-top: 2rem; }
  .bypass-link {
    display: inline-block; list-style: none; font-size: 0.8rem; color: #98a2b3;
    text-decoration: none; border-bottom: 1px dashed #d0d5dd; padding-bottom: 1px; cursor: pointer;
  }
  .bypass-link::-webkit-details-marker { display: none; }
  .bypass-link:hover { color: #667085; border-bottom-color: #98a2b3; }
  .bypass[open] .bypass-link { margin-bottom: 1.25rem; }
  .bypass-row { display: flex; gap: 0.5rem; }
  .bypass-row input[type=text] {
    flex: 1; min-width: 0; padding: 0.6rem 0.85rem; border-radius: 10px;
    border: 1px solid #d0d5dd; font-size: 0.9rem; outline: none; background: #fff; color: #101828;
  }
  .bypass-row input[type=text]:focus { border-color: #0d6dea; box-shadow: 0 0 0 3px rgba(13,109,234,0.15); }
  .btn {
    display: inline-flex; align-items: center; justify-content: center; white-space: nowrap;
    padding: 0.6rem 1.1rem; border-radius: 10px; border: none;
    background: #0d6dea; color: #fff; font-size: 0.9rem; font-weight: 600;
    text-decoration: none; cursor: pointer;
  }
  .btn:hover { background: #0b5cc9; }
  .error-text { color: #d92d20; font-size: 0.8rem; margin-top: 0.6rem; }
</style>`

const DefaultFallbackHtml = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Application Unavailable</title>
` + pageStyle + `
</head>
<body>
<div class="page">
  ` + rocketLogoSvg + `
  <div class="card">
    <div class="icon-badge">` + alertTriangleIconSvg + `</div>
    <h1>Application Unavailable</h1>
    <p class="message">` + DefaultFallbackMessage + `</p>
  </div>
</div>
</body>
</html>
`

const DefaultMaintenanceHtml = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Under Maintenance</title>
` + pageStyle + `
</head>
<body>
<div class="page">
  ` + rocketLogoSvg + `
  <div class="card">
    <div class="icon-badge">` + toolIconSvg + `</div>
    <h1>Under Maintenance</h1>
    <p class="message">{{Message}}</p>
    {{BypassForm}}
  </div>
</div>
</body>
</html>
`

const bypassButtonHtml = `<div class="bypass">
  <a href="?%s=1" class="bypass-link">Continue anyway</a>
</div>`

// bypassFormHtml uses <details>/<summary> so the code entry form stays fully hidden — no JS
// needed — until a visitor deliberately clicks "Have a bypass code?". The %s slots are: the
// "open" attribute (set only after a failed attempt, so the form and error stay visible),
// the query param name, and the invalid-code notice.
const bypassFormHtml = `<details class="bypass"%s>
  <summary class="bypass-link">Have a bypass code?</summary>
  <form method="get" class="bypass-row">
    <input type="text" id="rocket_bypass_code" name="%s" placeholder="Enter code" autocomplete="off">
    <button type="submit" class="btn">Continue</button>
  </form>
  %s
</details>`

const bypassFormInvalidCodeNotice = `<p class="error-text">Invalid bypass code.</p>`

// RenderBypassForm builds the bypass UI shown on the maintenance page. When requiresCode is
// false, bypass is open and a single link grants it. When true, a collapsed disclosure is
// shown instead, expanding into a code-entry form only once clicked (or automatically when
// invalidCode is set, so a failed attempt and its error stay visible).
func RenderBypassForm(requiresCode bool, invalidCode bool) string {
	if !requiresCode {
		return fmt.Sprintf(bypassButtonHtml, BypassQueryParam)
	}

	openAttr := ""
	notice := ""
	if invalidCode {
		openAttr = " open"
		notice = bypassFormInvalidCodeNotice
	}

	return fmt.Sprintf(bypassFormHtml, openAttr, BypassQueryParam, notice)
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
