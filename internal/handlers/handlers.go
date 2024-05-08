package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Pdv2323/bread-n-breakfast/internal/config"
	"github.com/Pdv2323/bread-n-breakfast/internal/models"
	"github.com/Pdv2323/bread-n-breakfast/internal/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is a repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository to the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This is the home page")
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplates(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// Some Logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, Again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//Send data to template
	render.RenderTemplates(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})

}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "make-reservation.page.html", &models.TemplateData{})
}

// SingleBed renders the room page
func (m *Repository) SingleBed(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "single-bed.page.html", &models.TemplateData{})
}

// DoubleBed renders the room page.
func (m *Repository) DoubleBed(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "double-bed.page.html", &models.TemplateData{})
}

// SearchAvailability renders the search availibity page.
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostSearchAvailability renders the search availibity page.
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ContactUs renders the contact page.
func (m *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "contact.page.html", &models.TemplateData{})
}
