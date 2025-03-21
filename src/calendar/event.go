package calendar

import "time"

var layout string = "1/2/06, 3:04pm"

// Event represents a performance/concert/whatever
type Event struct {
	BandName    string
	Instruments []string
	Date        time.Time
	Venue       Venue
	Performers  []string
	Featured    bool
}

func (e *Event) DateAndTime() string {
	if e.Date.IsZero() {
		return ""
	}

	return e.Date.Format(layout)
}
