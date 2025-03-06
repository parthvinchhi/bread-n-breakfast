package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/parthvinchhi/bread-n-breakfast/internal/config"
	"github.com/parthvinchhi/bread-n-breakfast/internal/models"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
}

var app *config.AppConfig

var pathToTemplates = "./templates"

// Iterate returns a slice of ints, starting at 1 going to count
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 1; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in DD-MM-YYYY format
func HumanDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "success")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Templates renders templates using html/templates
func Templates(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		//get template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writting template to browser", err)
		return err
	}

	return nil

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// pages, err := filepath.Glob("/home/vinchhi-parth/Desktop/Git - Golang/learn-web-dev/templates/*.page.html")
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, fmt.Errorf("error finding page templates: %w", err)
	}

	if len(pages) == 0 {
		return myCache, fmt.Errorf("no page templates found in path: %s", pathToTemplates)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, fmt.Errorf("error parsing page template %s: %w", page, err)
		}

		// matches, err := filepath.Glob("/home/vinchhi-parth/Desktop/Git - Golang/learn-web-dev/templates/*.layout.html")
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, fmt.Errorf("error finding layout templates: %w", err)
		}

		if len(matches) > 0 {
			// ts, err = ts.ParseGlob("/home/vinchhi-parth/Desktop/Git - Golang/learn-web-dev/templates/*.layout.html")
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, fmt.Errorf("error parsing layout templates: %w", err)
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
