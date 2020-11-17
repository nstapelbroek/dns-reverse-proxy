package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createRoutes(hostname string) map[string][]string {
	routes := make(map[string][]string)
	routes[hostname] = []string{"8.8.4.4:53"}

	return routes
}

func TestFindProxyRouteSuffixMatch(t *testing.T) {
	routes := createRoutes(".example.com.")

	assert.Len(t, findProxyRoute("subdomain.example.com.", routes), 1)
	assert.Len(t, findProxyRoute("something-else.example.com.", routes), 1)
}

func TestFindProxyRouteSuffixMisMatch(t *testing.T) {
	routes := createRoutes(".example.com.")

	assert.Len(t, findProxyRoute("example.net.", routes), 0)
	assert.Len(t, findProxyRoute("example.com.", routes), 0) // no subdomain here
}

func TestFindProxyRouteAsteriskMatch(t *testing.T) {
	routes := createRoutes("local.*.com.")

	assert.Len(t, findProxyRoute("local.example.com.", routes), 1)
	assert.Len(t, findProxyRoute("local.awesomesite.com.", routes), 1)
}

func TestFindProxyRouteUnnecessaryAsteriskMatch(t *testing.T) {
	routes := createRoutes("*.example.com.")

	assert.Len(t, findProxyRoute("subdomain.example.com.", routes), 1)
	assert.Len(t, findProxyRoute("something-else.example.com.", routes), 1)
}

func TestFindProxyRouteAsteriskMisMatch(t *testing.T) {
	routes := createRoutes("local.*.com.")

	assert.Len(t, findProxyRoute("local.example.co.uk.", routes), 0)
	assert.Len(t, findProxyRoute("local.com.domain.computer.", routes), 0)
}

func TestFindProxyRouteMultipleMatches(t *testing.T) {
	// should return the first match
	routes := make(map[string][]string)
	routes["local.example.com."] = []string{"8.8.4.4:53"}
	routes[".example.com."] = []string{"8.8.8.8:53"}
	routes["api.*.example.com."] = []string{"127.0.0.1:53"}

	assert.Equal(t, []string{"8.8.4.4:53"}, findProxyRoute("local.example.com.", routes))
	assert.Equal(t, []string{"8.8.8.8:53"}, findProxyRoute("subdomain.example.com.", routes))
	assert.Equal(t, []string{"8.8.4.4:53"}, findProxyRoute("api.local.example.com.", routes))
	assert.Equal(t, []string{"8.8.8.8:53"}, findProxyRoute("it-will-never-hit-the-asterisk.example.com.", routes))
}
