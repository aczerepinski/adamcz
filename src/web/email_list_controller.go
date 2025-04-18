package web

import (
	"net/http"
)

type emailListData struct {
	CDN       string
	Version   string
	MetaTitle string
}

// RenderSubscriptionForm renders the signup form (GET)
func (c *Controller) RenderSubscriptionForm(w http.ResponseWriter, r *http.Request) {
	data := emailListData{
		CDN:       cdnHost,
		Version:   c.version,
		MetaTitle: "adamcz | subscribe",
	}
	c.templates["emailList"].Execute(w, data)
}
