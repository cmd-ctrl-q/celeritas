package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
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
	JetViews   *jet.Set
	Session    *scs.SessionManager
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
	Error      string
	// Flash is a message that gets sent to the session only once
	Flash string
}

// defaultData modifies data from TemplateData based on data from Request
func (c *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	// pass data to request
	td.Secure = c.Secure
	td.ServerName = c.ServerName
	td.CSRFToken = nosurf.Token(r)
	td.Port = c.Port
	if c.Session.Exists(r.Context(), "userID") {
		// user is authenticated
		td.IsAuthenticated = true
	}
	td.Error = c.Session.PopString(r.Context(), "error")
	td.Flash = c.Session.PopString(r.Context(), "flash")
	return td
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {

	switch strings.ToLower(c.Renderer) {
	// go templates
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, variables, data)
	default:

	}

	return errors.New("no rendering engine specified")
}

// GoPage renders a standard Go template
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

// JetPage renders a template using the Jet templating engine
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap
	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	// add to default data
	td = c.defaultData(td, r)

	// render template
	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
