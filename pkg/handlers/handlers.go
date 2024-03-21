package handlers

import (
	"net/http"

	"github.com/Pdv2323/bread-n-breakfast/pkg/config"
	"github.com/Pdv2323/bread-n-breakfast/pkg/models"
	"github.com/Pdv2323/bread-n-breakfast/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

//Repository is a repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository to the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This is the home page")
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplates(w, "home.page.html", &models.TemplateData{})
}

//About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// Some Logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, Again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//Send data to template
	render.RenderTemplates(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})

}
