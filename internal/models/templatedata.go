package models

import "github.com/Pdv2323/bread-n-breakfast/internal/forms"

//TemplateData will hold the variables that will be used to store data in the template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
