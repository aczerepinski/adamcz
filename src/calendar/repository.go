package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Repository struct {
	PastEvents     []*Event
	UpcomingEvents []*Event
}

// InitializeRepository initializes from directory of json files.
// For an example of the expected file format see the tests
// for ParseFile.
func InitializeRepository(root string) (*Repository, error) {
	var r Repository
	now := time.Now()

	if !strings.HasSuffix("/", root) {
		root = fmt.Sprintf("%s/", root)
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return &r, fmt.Errorf("unable to read directory: %v", err)
	}

	for _, f := range files {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", root, f.Name()))
		if err != nil {
			return &r, err
		}
		var events []*Event
		err = json.Unmarshal(data, &events)
		if err != nil {
			return &r, fmt.Errorf("unable to unmarshal file %s: %v", f.Name(), err)
		}
		for _, e := range events {
			if e.Date.Before(now) {
				r.PastEvents = append([]*Event{e}, r.PastEvents...)

			} else {
				r.UpcomingEvents = append(r.UpcomingEvents, e)
			}
		}
	}

	return &r, nil
}
