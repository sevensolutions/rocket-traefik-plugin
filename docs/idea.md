# Rocket Traefik Plugin

This is a traefik middleware plugin to be used in conjunction with Rocket.
It serves different purposes.

Rocket is a Control-Plane for a cluster, where Traefik is the main reverse proxy.

## Fallback Route

The plugin should be used on a route with priority 1, acting as an underlay to a real route.
The real route is configured dynamically by running an app but if this app is not running, the fallback route should serve a static HTML page, telling the user that the app is currently unavailable.

## Maintenance Mode

Every app route will later register this middleware. The middleware should check the Rocket control plane if the maintenance mode has been enabled for this application and then display a static maintenance page. Here we need to make sure that we add some kind of cache so we don't poll Rocket for every single request.
