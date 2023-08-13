package web

import (
	"fmt"
	"net/http"
)

func getParam(param string, r *http.Request) []string {
	r.ParseForm()
	var vals []string
	if v := r.Form.Get(param); v != "" {
		// todo: parse multiple values
		vals = []string{r.Form.Get(param)}
	}
	return vals
}

func musicFilters(activeInstruments []string) []filter {
	all := []string{
		"trumpet", "flugelhorn", "piano", "rhodes", "bass", "drums", "flute",
	}

	var filters []filter
	for _, instrument := range all {
		active := false
		for _, ai := range activeInstruments {
			if ai == instrument {
				active = true
			}
		}
		path := fmt.Sprintf("/music?instruments=%s", instrument)
		if active {
			path = "/music"
		}
		filters = append(filters,
			filter{
				Name:   instrument,
				Active: active,
				Path:   path,
			})
	}
	return filters
}
