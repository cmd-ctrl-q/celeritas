package celeritas

import "net/http"

// SessionLoad saves and loads session on every request
func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	c.InfoLog.Println("SessionLoad called")
	return c.Session.LoadAndSave(next)
}
