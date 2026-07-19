package src

import (
	"net/http"

	"github.com/sevensolutions/rocket-traefik-plugin/src/logging"
)

type RocketTraefikPlugin struct {
	logger *logging.Logger
	next   http.Handler
}

func (toa *RocketTraefikPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
}
