package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// pageData for teste cases
var pageData = []struct {
	// the name of the page/view
	name string

	// rendering template engine
	renderer string

	// template to run the given test
	template string

	// is an error expected
	errorExpected bool
	errorMessage  string
}{
	// go templates
	{"go_page", "go", "home", false, "error rendering go template"},
	{"go_page_no_template", "go", "no-file", true, "no error rendering non-existent template when one is expected"},

	// jet templates
	{"jet_page", "jet", "home", false, "error rendering jet template"},
	{"jet_page_no_template", "jet", "no-file", true, "no error rendering non-existent jet template when one is expected"},

	// an invalid render engine
	{"invalid_render_engine", "foo", "home", true, "no error rendering with non-existent template engine"},
}

func TestRender_Page(t *testing.T) {

	for _, e := range pageData {
		r, err := http.NewRequest("GET", "/some-url", nil)
		if err != nil {
			t.Error(err)
		}

		// mocks a response writer
		w := httptest.NewRecorder()

		testRenderer.Renderer = e.renderer
		testRenderer.RootPath = "./testdata"

		err = testRenderer.Page(w, r, e.template, nil, nil)
		if e.errorExpected {
			// if expecting error but didn't get one
			if err == nil {
				t.Errorf("%s: %s", e.name, e.errorMessage)
			}
		} else {
			// if not expecting error but get one
			if err != nil {
				t.Errorf("%s: %s: %s", e.name, e.errorMessage, err.Error())
			}
		}
	}
}

func TestRender_GoPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	// test existing page
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestRender_JetPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "jet"

	// test existing page
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error(err)
	}
}
