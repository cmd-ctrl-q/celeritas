package render

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Render struct {
	// Renderer is the rendering engine
	Renderer string
	RootPath string

	// If Secure is true then the application is
	// running in HTTPS mode
	Secure     bool
	Port       string
	ServerName string
}

type TemplateData struct {
	IsAuthenticated bool

	// IntMap, StringMap, and FloatMap are the types for
	// holding specific data types
	IntMap    map[string]int
	StringMap map[string]string
	FloatMap  map[string]float32

	// Data is the type for holding generic data
	Data map[string]interface{}

	// Cross Site Request Forgery Protection
	CSRFToken  string
	Port       string
	ServerName string
	Secure     bool
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {

	switch strings.ToLower(c.Renderer) {
	// go templates
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
	}

	return nil
}

func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.go.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}
