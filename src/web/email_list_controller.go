package web

import (
	"net/http"

	email_list "github.com/aczerepinski/adamcz/src/email-list"
)

type emailListFormData struct {
	EmailAddress string
	FullName     string
	Success      bool
	Error        string
	CDN          string
	Version      string
	MetaTitle    string
}

// RenderSubscriptionForm renders the signup form (GET)
func (c *Controller) RenderSubscriptionForm(w http.ResponseWriter, r *http.Request) {
	data := emailListFormData{
		CDN:       cdnHost,
		Version:   c.version,
		MetaTitle: "adamcz | subscribe",
	}
	c.templates["emailList"].Execute(w, data)
}

// ProcessSubscription handles form submission (POST)
func (c *Controller) ProcessSubscription(w http.ResponseWriter, r *http.Request) {
	data := emailListFormData{
		CDN:       cdnHost,
		Version:   c.version,
		MetaTitle: "adamcz | subscribe",
	}
	email := r.FormValue("email")
	data.EmailAddress = email
	if email == "" {
		data.Error = "Email address is required."
	} else {
		// Here you could add logic to save the signup, e.g. to a database or file
		_ = email_list.Signup{EmailAddress: email}
		data.Success = true
		data.EmailAddress = ""
	}
	c.templates["emailList"].Execute(w, data)
}
