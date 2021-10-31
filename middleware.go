package celeritas

import "net/http"

// SessionLoad saves and loads session on every request
func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	return c.Session.LoadAndSave(next)
}
