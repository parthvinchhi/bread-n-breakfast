package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/parthvinchhi/bread-n-breakfast/internal/config"
	my_driver "github.com/parthvinchhi/bread-n-breakfast/internal/driver"
	"github.com/parthvinchhi/bread-n-breakfast/internal/forms"
	"github.com/parthvinchhi/bread-n-breakfast/internal/helpers"
	"github.com/parthvinchhi/bread-n-breakfast/internal/models"
	"github.com/parthvinchhi/bread-n-breakfast/internal/render"
	"github.com/parthvinchhi/bread-n-breakfast/internal/repository"
	dbrepo "github.com/parthvinchhi/bread-n-breakfast/internal/repository/db-repo"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is a repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *my_driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository to the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This is the home page")
	// remoteIP := r.RemoteAddr
	// m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Templates(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// Some Logic
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hello, Again"

	// remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	// stringMap["remote_ip"] = remoteIP

	//Send data to template
	render.Templates(w, r, "about.page.html", &models.TemplateData{
		// StringMap: stringMap,
	})

}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Templates(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	err = errors.New("This is an error message")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Templates(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// SingleBed renders the room page
func (m *Repository) SingleBed(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "single-bed.page.html", &models.TemplateData{})
}

// DoubleBed renders the room page.
func (m *Repository) DoubleBed(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "double-bed.page.html", &models.TemplateData{})
}

// SearchAvailability renders the search availibity page.
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostSearchAvailability renders the search availibity page.
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ContactUs renders the contact page.
func (m *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from the session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Templates(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}
