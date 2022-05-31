package web

import (
	"net/http"

	"github.com/aczerepinski/adamcz/src/calendar"
)

type calendarIndex struct {
	CDN            string
	Version        string
	PageTitle      string
	MetaTitle      string
	PastEvents     []*calendar.Event
	UpcomingEvents []*calendar.Event
}

func (c *Controller) calendarIndex(w http.ResponseWriter, r *http.Request) {
	data := calendarIndex{
		CDN:            cdnHost,
		Version:        c.version,
		PageTitle:      "Calendar",
		MetaTitle:      "adamcz | calendar",
		PastEvents:     c.events.PastEvents,
		UpcomingEvents: c.events.UpcomingEvents,
	}

	c.templates["calendarIndex"].Execute(w, data)
}
